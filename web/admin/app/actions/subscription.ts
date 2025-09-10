'use server'
import axiosIns from '@/app/libs/axios'


export const Get = async () => {
    try {
        const response = await axiosIns.get("/subscription/plan");
        console.log(response.data)
        return response.data.body.plan
    } catch (error) {
        throw error
    }
}
