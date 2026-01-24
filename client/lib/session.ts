"use client";
import { useSession } from "next-auth/react";

/**
 * Client-side session helper
 * Use this in components to get the current session
 */
export function useClientSession() {
  return useSession();
}

/**
 * Get access token from session
 * This can be used to manually add auth headers
 */
export function getAccessToken(session: any): string | null {
  return session?.data?.accessToken || null;
}

