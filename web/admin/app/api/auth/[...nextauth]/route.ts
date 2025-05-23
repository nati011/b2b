import NextAuth, { AuthOptions } from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"
import { jwtDecode } from "jwt-decode"
import axios from "axios"
import { NextResponse } from "next/server";

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

async function refreshAccessToken(refreshToken: string) {
  try {
    console.log(refreshToken)
    const response = await axios.post(`${baseURL}/api/v1/auth/refresh`, {
      refresh_token: refreshToken,
    });
    console.log(response.data)
    const decoded = jwtDecode<KeycloakJWT>(response.data.body.jwt.access_token)
    return {
      accessToken: response.data.body.jwt.access_token,
      refreshToken: response.data.body.jwt.refresh_token,
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
    NextResponse.redirect("/auth/signin")
    throw error
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
        //@ts-expect-error
        return await refreshAccessToken(token.refreshToken)
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