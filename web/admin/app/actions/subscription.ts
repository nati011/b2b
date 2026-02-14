'use server'
import axiosIns from '@/app/libs/axios'
import { withErrorHandling } from "@/app/libs/error-handling";
import { getSession } from "@/app/actions/getSession";


export const Get = async () => {
        return withErrorHandling(async () => {
        const session = await getSession() as any
        const roles: string[] = session?.user?.roles || []
        if (roles.includes('superadmin') || roles.includes('admin')) {
            return null
        }

        const response = await axiosIns.get("/subscription/plan");
        console.log(response.data)
        return response.data.body.plan
        })
  
}
