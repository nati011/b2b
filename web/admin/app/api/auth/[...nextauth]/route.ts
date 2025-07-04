import NextAuth, { AuthOptions, Session, TokenSet } from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"
import GoogleProvider from "next-auth/providers/google"
import { jwtDecode } from "jwt-decode"
import axios, { AxiosError } from "axios"
import { getUserIdentityWithToken } from "@/app/actions/getUserIdentity"

// Environment variables
const API_BASE_URL = process.env.NEXT_BASE_URL || "https://b2b-67gk.onrender.com"
const NEXTAUTH_SECRET = process.env.NEXTAUTH_SECRET
const GOOGLE_CLIENT_ID = process.env.GOOGLE_CLIENT_ID
const GOOGLE_CLIENT_SECRET = process.env.GOOGLE_CLIENT_SECRET

// Types
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
    Id: number
    Name: string
    Action: string
  }[]
}

interface AppToken extends TokenSet {
  accessToken: string
  refreshToken: string
  accessTokenExpires: number
  user: UserToken
  userIdentity?: UserIdentity
  error?: "RefreshAccessTokenError" | "TokenExpiredError"
  refreshAttempts?: number
}

interface AuthResponse {
  body: {
    access_token: string
    refresh_token: string
  }
  message?: string
}

// Utility functions
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

// Token refresh function with improved error handling
async function refreshAccessToken(token: AppToken): Promise<AppToken> {
  const MAX_REFRESH_ATTEMPTS = 3
  const currentAttempts = token.refreshAttempts || 0
  
  // Check if we've exceeded the maximum refresh attempts
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
      timeout: 10000, // 10 second timeout
      headers: {
        'Content-Type': 'application/json',
      }
    })

    if (!response.data.body?.access_token) {
      throw new Error("Invalid refresh response")
    }

    const decoded = decodeToken(response.data.body.access_token)
    const user = createUserFromToken(decoded)

    // Fetch user identity with the new token using the dedicated function
    const userIdentity = await getUserIdentityWithToken(response.data.body.access_token)

    return {
      ...token,
      accessToken: response.data.body.access_token,
      refreshToken: response.data.body.refresh_token || token.refreshToken,
      accessTokenExpires: decoded.exp * 1000,
      user,
      userIdentity: userIdentity || token.userIdentity,
      error: undefined,
      refreshAttempts: 0
    }
  } catch (error) {
    console.error("Token refresh failed:", error)

    // Handle specific error cases
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError<{ message?: string }>
      
      // Handle 400/401 errors (invalid refresh token)
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
      
      // Handle network errors
      if (axiosError.code === 'ECONNABORTED' || axiosError.code === 'NETWORK_ERROR') {
        console.error("Network error during token refresh")
      }
    }

    // Return token with error for any other failure
    return {
      ...token,
      error: "RefreshAccessTokenError",
      refreshAttempts: currentAttempts + 1
    }
  }
}

// Main auth configuration
export const authOptions: AuthOptions = {
  providers: [
    // Google OAuth provider
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
              headers: {
                'Content-Type': 'application/json',
              }
            }
          )

          if (response.status !== 202) {
            throw new Error(response.data?.message || "Authentication failed")
          }

          if (!response.data.body?.access_token) {
            throw new Error("Invalid authentication response")
          }

          const decoded = decodeToken(response.data.body.access_token)
          const user = createUserFromToken(decoded)

          // Fetch user identity after successful authentication using the dedicated function
          const userIdentity = await getUserIdentityWithToken(response.data.body.access_token)

          return {
            id: user.id,
            name: user.name,
            email: user.email,
            image: null,
            ...response.data,
            user,
            userIdentity
          }
        } catch (error) {
          if (axios.isAxiosError(error)) {
            const axiosError = error as AxiosError<{ message?: string }>
            const errorMessage = axiosError.response?.data?.message || "Authentication failed"
            console.error("Authentication error:", errorMessage)
            throw new Error(errorMessage)
          }
          
          console.error("Authentication error:", error)
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
      // Allow relative URLs
      if (url.startsWith("/")) return `${baseUrl}${url}`
      
      // Allow URLs from the same origin
      if (new URL(url).origin === baseUrl) return url
      
      // Default to base URL for external URLs
      return baseUrl
    },
    
    async jwt({ token, user, account }) {
      // Initial sign in
      if (user && account) {
        const userData = user as any
        const decoded = decodeToken(userData.body.access_token)
        
        return {
          accessToken: userData.body.access_token,
          refreshToken: userData.body.refresh_token,
          accessTokenExpires: decoded.exp * 1000,
          user: userData.user,
          userIdentity: userData.userIdentity,
          refreshAttempts: 0
        }
      }

      // Return previous token if the access token has not expired
      const appToken = token as AppToken
      if (appToken.accessTokenExpires && Date.now() < appToken.accessTokenExpires) {
        return token
      }

      // Access token has expired, try to refresh it
      if (appToken.refreshToken) {
        const refreshedToken = await refreshAccessToken(appToken)
        
        // If refresh failed and we have no access token, clear the session
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
      const appToken = token as AppToken
      
      // Handle token errors
      if (appToken.error === "RefreshAccessTokenError" && !appToken.accessToken) {
        session.user = undefined
        session.expires = new Date(0).toISOString()
        return session
      }
      
      // Set session data
      session.user = appToken.user
      ;(session as any).accessToken = appToken.accessToken
      ;(session as any).userIdentity = appToken.userIdentity
      ;(session as any).error = appToken.error
      
      return session
    }
  },
  
  session: {
    strategy: "jwt",
    maxAge: 72 * 60 * 60, // 72 hours
  },
  
  jwt: {
    maxAge: 72 * 60 * 60, // 72 hours
  },
  
  debug: process.env.NODE_ENV === 'development',
  
  // Security settings
  useSecureCookies: process.env.NODE_ENV === 'production',
  cookies: {
    sessionToken: {
      name: process.env.NODE_ENV === 'production' ? '__Secure-next-auth.session-token' : 'next-auth.session-token',
      options: {
        httpOnly: true,
        sameSite: 'lax',
        path: '/',
        secure: process.env.NODE_ENV === 'production'
      }
    }
  }
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST }