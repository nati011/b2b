'use server'
import axiosIns from "@/app/libs/axios";

export const InitResetPassword = async (email: string) => {
    try {
        console.log(email)
        const response = await axiosIns.post("/user/init_reset", {
            email: email
        });
        return response.data.message
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error has occured while reseting the password"
        }
        throw "An error has occured while reseting the password"
    }
}


export const ResetPassword = async (token: string, password: string) => {
    try {
        const response = await axiosIns.post(`/auth/reset/${token}`, {
            password: password
        });
        return response.data.message
    } catch (error: any) {
        console.log(error)
        if (error.response) {
            throw error.response.data.message || "An error has occured while creating the product"
        }
        throw "An error has occured while creating the product"
    }
}