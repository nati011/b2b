import { AuthOptions, Session, TokenSet } from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"
import GoogleProvider from "next-auth/providers/google"

import { jwtDecode } from "jwt-decode"
import axios, { AxiosError } from "axios"
import { getUserIdentityWithToken } from "@/app/actions/getUserIdentity"
import { Permissions } from "@/app/libs/types"

const API_BASE_URL = process.env.NEXT_BASE_URL || "https://b2b-67gk.onrender.com"
const NEXTAUTH_SECRET = process.env.NEXTAUTH_SECRET
const GOOGLE_CLIENT_ID = process.env.GOOGLE_CLIENT_ID
const GOOGLE_CLIENT_SECRET = process.env.GOOGLE_CLIENT_SECRET

interface AuthResponse {
    body: {
      access_token: string
      refresh_token: string
    }
    message?: string
}

interface KeycloakJWT {
    exp: number
    iat: number
    sub: string
    email: string
    name: string
    preferred_username: string
    given_name: string
    family_name: string
    realm_access: {
      roles: string[]
    }
}

interface UserToken {
    id: string
    name: string
    email: string
    username: string
    roles: string[]
}


interface UserIdentity {
    id: number
    first_name: string
    last_name: string
    email: string
    phone: string
    username: string
    dob: string
    is_active: boolean
    permissions: {
      List: Permissions[]
    }
}

interface AppToken extends TokenSet {
    accessToken: string
    refreshToken: string
    accessTokenExpires: number
    user: UserToken & {
        permissions?: {
            List: Permissions[]
        }
    }
    userIdentity?: UserIdentity
    error?: "RefreshAccessTokenError" | "TokenExpiredError"
    refreshAttempts?: number
}

function decodeToken(token: string): KeycloakJWT {
    try {
      return jwtDecode<KeycloakJWT>(token)
    } catch (error) {
      throw new Error("Invalid token format")
    }
}

function createUserFromToken(decoded: KeycloakJWT): UserToken {
    return {
      id: decoded.sub,
      name: decoded.name || `${decoded.given_name || ''} ${decoded.family_name || ''}`.trim(),
      email: decoded.email,
      username: decoded.preferred_username,
      roles: decoded.realm_access?.roles || []
    }
}

async function refreshAccessToken(token: AppToken): Promise<AppToken> {
    const MAX_REFRESH_ATTEMPTS = 3
    const currentAttempts = token.refreshAttempts || 0
    
    if (currentAttempts >= MAX_REFRESH_ATTEMPTS) {
      console.error("Maximum refresh attempts exceeded")
      return {
        ...token,
        error: "RefreshAccessTokenError",
        accessToken: "",
        refreshToken: "",
        accessTokenExpires: 0,
        user: {
          id: "",
          name: "",
          email: "",
          username: "",
          roles: []
        },
        refreshAttempts: currentAttempts
      }
    }

    try {
      const response = await axios.post<AuthResponse>(`${API_BASE_URL}/api/v1/auth/refresh`, {
        refresh_token: token.refreshToken,
      }, {
        headers: {
          'Content-Type': 'application/json',
        }
      })

      if (!response.data.body?.access_token) {
        throw new Error("Invalid refresh response")
      }

      const decoded = decodeToken(response.data.body.access_token)
      const user = createUserFromToken(decoded)

      const userIdentity = await getUserIdentityWithToken(response.data.body.access_token)

      return {
        ...token,
        accessToken: response.data.body.access_token,
        refreshToken: response.data.body.refresh_token || token.refreshToken,
        accessTokenExpires: decoded.exp * 1000,
        user: {
          ...user,
          // @ts-expect-error
          permissions: userIdentity?.permissions || []
        },
        error: undefined,
        refreshAttempts: 0
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        const axiosError = error as AxiosError<{ message?: string }>
        
        if (axiosError.response?.status === 400 || axiosError.response?.status === 401) {
          return {
            ...token,
            error: "RefreshAccessTokenError",
            accessToken: "",
            refreshToken: "",
            accessTokenExpires: 0,
            user: {
              id: "",
              name: "",
              email: "",
              username: "",
              roles: []
            },
            refreshAttempts: currentAttempts + 1
          }
        }
        
        if (axiosError.code === 'ECONNABORTED' || axiosError.code === 'NETWORK_ERROR') {
          console.error("Network error during token refresh")
        }
      }

      return {
        ...token,
        error: "RefreshAccessTokenError",
        refreshAttempts: currentAttempts + 1
      }
    }
}

export const authOptions: AuthOptions = {
    providers: [
      ...(GOOGLE_CLIENT_ID && GOOGLE_CLIENT_SECRET ? [
        GoogleProvider({
          clientId: GOOGLE_CLIENT_ID,
          clientSecret: GOOGLE_CLIENT_SECRET,
        })
      ] : []),
      
      // Credentials provider
      CredentialsProvider({
        name: "Credentials",
        credentials: {
          email: { 
            label: "Email", 
            type: "email",
            placeholder: "Enter your email"
          },
          password: { 
            label: "Password", 
            type: "password",
            placeholder: "Enter your password"
          }
        },
        // @ts-ignore
        async authorize(credentials) {
          if (!credentials?.email || !credentials?.password) {
            throw new Error("Email and password are required")
          }

          try {
            const response = await axios.post<AuthResponse>(
              `${API_BASE_URL}/api/v1/auth/login`,
              {
                email: credentials.email,
                password: credentials.password
              },
              {
                timeout: 10000,
              }
            )

            if (response.status !== 202) {
              console.log(response.data)
              throw new Error(response.data?.message || "Authentication failed")
            }

            if (!response.data.body?.access_token) {
              throw new Error("Invalid authentication response")
            }

            const decoded = decodeToken(response.data.body.access_token)
            const user = createUserFromToken(decoded)
            console.log('User created from token:', user)

            const userIdentity = await getUserIdentityWithToken(response.data.body.access_token)
            console.log('User identity:', userIdentity)

            return {
              id: user.id,
              ...response.data,
              user,
              userIdentity
            }
          } catch (error) {
            if (axios.isAxiosError(error)) {
              const axiosError = error as AxiosError<{ message?: string }>
              const errorMessage = axiosError.response?.data?.message || "Authentication failed"
              console.log(axiosError.response?.data)
              console.error("Authentication 11:", axiosError.response?.data, axiosError.response?.status)
              throw new Error(errorMessage)
            }
            
            console.error("Authentication 22:", error)
            throw new Error("Authentication failed")
          }
        }
      })
    ],
    
    secret: NEXTAUTH_SECRET,
    
    pages: {
      signIn: "/auth/signin",
      error: "/auth/signin",
      signOut: "/auth/signin"
    },
    
    callbacks: {
      async redirect({ url, baseUrl }) {
        if (url.startsWith("/")) return `${baseUrl}${url}`
        if (new URL(url).origin === baseUrl) return url
        return baseUrl
      },
      
      async jwt({ token, user, account }) {
        if (user && account) {
          const userData = user as any
          const decoded = decodeToken(userData.body.access_token)
          console.log('JWT callback - userData:', userData)
          
          return {
            accessToken: userData.body.access_token,
            refreshToken: userData.body.refresh_token,
            accessTokenExpires: decoded.exp * 1000,
            user: {
              id: userData.user.id,
              name: userData.user.name,
              email: userData.user.email,
              permissions: userData.userIdentity?.permissions
            },
            refreshAttempts: 0,
          }
        }

        const appToken = token as AppToken
        if (appToken.accessTokenExpires && Date.now() < appToken.accessTokenExpires) {
          return token
        }

        if (appToken.refreshToken) {
          const refreshedToken = await refreshAccessToken(appToken)
        
          if (refreshedToken.error === "RefreshAccessTokenError" && !refreshedToken.accessToken) {
            return {
              ...refreshedToken,
              user: {
                id: "",
                name: "",
                email: "",
                username: "",
                roles: []
              }
            }
          }
          
          return refreshedToken
        }

        // No refresh token available
        return {
          ...token,
          error: "TokenExpiredError"
        }
      },
      
      async session({ session, token }) {
        const appToken = token as any
        
        if (appToken.error === "RefreshAccessTokenError" && !appToken.accessToken) {
          session.user = undefined
          session.expires = new Date(0).toISOString()
          return session
        }
        
        session.user = {
          name: appToken.user.name,
          email: appToken.user.email,
          image: null
        }
        
        const customSession = session as any
        customSession.user = {
          ...session.user,
          id: appToken.user.id,
          first_name: appToken?.userIdentity?.first_name,
          last_name: appToken?.userIdentity?.last_name,
          permissions: appToken.user?.permissions || appToken.userIdentity?.permissions
        }
        customSession.accessToken = appToken.accessToken
        customSession.error = appToken.error
        return customSession
      }
    },
    
    session: {
      strategy: "jwt",
      maxAge: 72 * 60 * 60,
    },
    
    jwt: {
      maxAge: 72 * 60 * 60, 
    },
    
    debug: process.env.NODE_ENV === 'development',

    useSecureCookies: process.env.NODE_ENV === 'production',
    cookies: {
      sessionToken: {
        name: process.env.NODE_ENV === 'production' ? '__Secure-next-auth.session-token' : 'next-auth.session-token',
        options: {
          httpOnly: true,
          sameSite: 'lax',
          path: '/',
          secure: process.env.NODE_ENV === 'production',
          maxAge: 72 * 60 * 60 
        }
      }
    }
}