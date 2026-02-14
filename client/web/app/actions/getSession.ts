"use client";
import { getSession as getNextAuthSession } from "next-auth/react";

/**
 * Client-side session getter
 * This replaces the server-side getSession
 */
export async function getSession() {
  return await getNextAuthSession();
}
