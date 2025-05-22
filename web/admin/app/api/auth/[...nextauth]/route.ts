import type { NextApiRequest, NextApiResponse } from "next"
import NextAuth, { AuthOptions } from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"
import { jwtDecode } from "jwt-decode"
import axios from "axios"

const baseURL = process.env.NEXT_BASE_URL || "https://b2b-67gk.onrender.com"

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

async function refreshAccessToken() {
  const refreshToken = localStorage.getItem('refreshToken');
  if (!refreshToken) {
    throw new Error('No refresh token available');
  }

  try {
    const response = await axios.post(`${baseURL}/api/auth/refresh/`, {
      refresh: refreshToken,
    });
    const decoded = jwtDecode<KeycloakJWT>(response.data.body.access_token)
    return {
      // @ts-expect-error
      accessToken: user.body.access_token,
      // @ts-expect-error
      refreshToken: user.body.refresh_token,
      accessTokenExpires: decoded.exp * 1000,
      user: {
        id: decoded.sub,
        name: decoded.name,
        email: decoded.email,
        username: decoded.preferred_username,
        roles: decoded.realm_access?.roles || []
      }
    }
  } catch (error) {
    window.location.href = '/login';
    throw error;
  }
}


export const authOptions: AuthOptions = {
  providers: [
    CredentialsProvider({
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        try {
          const res = await fetch(`${baseURL}/api/v1/auth/login`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              "Access-Control-Allow-Origin": "*",
            },
            body: JSON.stringify(credentials)
          })

          if (!res.ok) {
            const error = await res.json()
            throw new Error(error.message || "Authentication failed")
          }

          return await res.json()
        } catch (error) {
          console.error("Authentication error:", error)
          return null
        }
      }
    })
  ],
  secret: process.env.NEXTAUTH_SECRET,
  pages: {
    signIn: "/auth/signin",
    error: "/error"
  },
  callbacks: {
    async jwt({ token, user, trigger, session }) {
      if (user) {
        // @ts-expect-error
        const decoded = jwtDecode<KeycloakJWT>(user.body.access_token)
        return {
          ...token,
          // @ts-expect-error
          accessToken: user.body.access_token,
          // @ts-expect-error
          refreshToken: user.body.refresh_token,
          accessTokenExpires: decoded.exp * 1000,
          user: {
            id: decoded.sub,
            name: decoded.name,
            email: decoded.email,
            username: decoded.preferred_username,
            roles: decoded.realm_access?.roles || []
          }
        }
      }

      // @ts-expect-error
      if (Date.now() > token.accessTokenExpires) {
        return await refreshAccessToken()
      }

      return token
    },
    async session({ session, token }) {
      console.log(session)
      console.log(token)
      return {
        ...session,
        error: token.error,
        accessToken: token.accessToken,
        user: {
          ...session.user,
          // @ts-expect-error
          ...token.user,
          // @ts-expect-error
          roles: token.user?.roles || []
        }
      }
    }
  },
  session: {
    strategy: "jwt",
    maxAge: 4 * 60 * 60
  },
  debug: process.env.NODE_ENV !== 'production'
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST }