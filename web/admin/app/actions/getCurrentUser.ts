"use server"
import { getUserIdentityFromSession } from "./getUserIdentity"

export default async function getCurrentUser() {
    try {
        return await getUserIdentityFromSession()
    } catch (error: any) {
        console.error("Error in getCurrentUser:", error.message)
        return null
    }
}