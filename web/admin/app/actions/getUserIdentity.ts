"use server"
import axios from 'axios'
import { getServerSession } from "next-auth/next"
import { authOptions } from "@/app/api/auth/[...nextauth]/route"
import { UserIdentity } from '@/app/libs/types'

// Environment variables
const API_BASE_URL = process.env.NEXT_BASE_URL || "https://b2b-67gk.onrender.com"

// Extended session type to include userIdentity
interface ExtendedSession {
  user?: {
    id: string
    name: string
    email: string
    username: string
    roles: string[]
  } | null
  accessToken?: string
  userIdentity?: UserIdentity
  error?: string
  expires: string
}

type IdentityResponse  = {
  user: UserIdentity
}

/**
 * Fetch user identity using a provided access token
 * This function can be used in NextAuth callbacks where session context is not available
 */
export async function getUserIdentityWithToken(accessToken: string): Promise<UserIdentity | null> {
  try {
    console.log(accessToken)
    const response = await axios.get<{ body: IdentityResponse }>(
      `${API_BASE_URL}/api/v1/identity/user`,
      {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
        timeout: 10000,
      }
    )

    if (!response.data?.body.user) {
      console.error("No user data received from API")
      return null
    }

    return response.data.body.user
  } catch (error: any) {
    console.error("Error fetching user identity with token:", error.response?.data || error.message)
    return null
  }
}

/**
 * Fetch user identity using the current session
 * This function first checks the session for cached user identity, then falls back to API call
 */
export async function getUserIdentityFromSession(): Promise<UserIdentity | null> {
  try {
    // First, try to get user identity from the session
    const session = await getServerSession(authOptions) as ExtendedSession
    
    if (session?.userIdentity) {
      return session.userIdentity
    }

    // If not in session, fetch from API using the access token
    if (session?.accessToken) {
      return await getUserIdentityWithToken(session.accessToken)
    }

    console.error("No session or access token available")
    return null
  } catch (error: any) {
    console.error("Error fetching user identity from session:", error.message)
    return null
  }
} 