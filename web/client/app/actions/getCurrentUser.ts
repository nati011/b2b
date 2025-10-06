"use server"
import axiosIns from "@/lib/axios";
import { UserIdentity } from "@/lib/types";

export default async function getCurrentUser(): Promise<UserIdentity | null> {
  try {
    const { data } = await axiosIns.get<{ body:{ user: UserIdentity} }>(
      "/identity/user"
    );
    console.log(data)

    if (!data?.body.user) {
      console.error("No user data received");
      return null;
    }

    return data.body.user;
  } catch (error: any) {
    console.error(
      "Error fetching current user:",
      error.response?.data || error.message
    );
    return null;
  } 
}
