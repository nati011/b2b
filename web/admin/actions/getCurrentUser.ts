import { getSession } from "@/actions/getSession";
import axiosIns from "@/app/libs/axios";

export default async function getCurrentUser() {
    try {
        const session = await getSession();
        // @ts-ignore
        if (!session?.user?.employeeNumber) {
            return null;
        }

        const resp = await axiosIns.get("/users/me")
        const currentUser = resp.data

        if (!currentUser) {
            return null;
        }

        return currentUser
    } catch (error: any) {
        return null;
    }
}