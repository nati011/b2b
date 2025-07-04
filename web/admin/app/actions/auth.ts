'use server'
import axiosIns from "@/app/libs/axios";
import { UserIdentity } from "../libs/types";


export const FetchUserDetail = async (id: number) => {
    try {
        const response = await axiosIns.get(`/user?id=${id}`);
        console.log(response.data)
        return response.data.body.user
    } catch (error) {
        console.log(error)
        throw error
    }
}


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
            throw error.response.data.message || "An error has occured while reseting your password"
        }
        throw "An error has occured while reseting your password"
    }
}

export const UpdateProfile  = async (data: Partial<UserIdentity>) =>{
    try{
        const response = await axiosIns.patch("/user", data)
        console.log(response.data)
        return response.data.message
    } catch (error: any) {
        console.log(error)
        if (error.response) {
            throw error.response.data.message || "An error has occured while updating your profile"
        }
        throw "An error has occured while updating your profile"
    }
}