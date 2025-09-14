'use server'
import axiosIns from '@/app/libs/axios'
import { withErrorHandling } from "@/app/libs/error-handling";


export const Get = async () => {
        return withErrorHandling(async () => {
        const response = await axiosIns.get("/subscription/plan");
        console.log(response.data)
        return response.data.body.plan
        })
  
}
