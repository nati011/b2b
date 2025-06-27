'use server'
import axiosIns from '@/app/libs/axios'
import { Role } from '@/app/libs/types';



export const GetAll = async () => {
    try {
        const response = await axiosIns.get("/role");
        return response.data.body.roles
    } catch (error) {
        throw error
    }
}

export const Create = async (data: Partial<Role>) => {
    try {
        const response = await axiosIns.post(
            "/role",
            data,
        );
        return response.data.detail
    } catch (error: any) {
        throw error
    }
}

export const Delete = async (id: number) => {
    try {
        await axiosIns.delete(`/role?id=${id}`);
    } catch (error: any) {
        throw error
    }
}

export const Update = async (data: Partial<Role>) => {
    try {
        const response = await axiosIns.put(`/role`,data);
        return response.data.message
    } catch (error: any) {
        throw error
    }
}