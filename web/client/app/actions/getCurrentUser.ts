import { getSession } from "@/app/actions/getSession";

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