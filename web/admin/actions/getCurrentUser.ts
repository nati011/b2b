import { getSession } from "@/actions/getSession";
import axiosIns from "@/app/libs/axios";

export default async function getCurrentUser() {
    try {
        const session = await getSession();
        // const resp = await axiosIns.get("/users/me")
        const currentUser = session?.user

        if (!currentUser) {
            return null;
        }

        return currentUser
    } catch (error: any) {
        return null;
    }
}