'use server'
import axiosIns from "@/app/libs/axios";

export const GetAll = async () => {
    try {
        const response = await axiosIns.get("/category");
        return response.data.body.categories
    } catch (error) {
        throw error
    }
}

export const Create = async (name: string) => {
    try {
        const response = await axiosIns.post(
            "/api/category",
            { name: name },
        );
        return response.data.detail
    } catch (error: any) {
        throw error
    }
}

export const Delete = async (id: number) => {
    try {
        await axiosIns.delete(`/category?id=${id}`);
    } catch (error: any) {
        throw error
    }
}

export const Update = async (id: number, name: string) => {
    try {
        await axiosIns.patch(`/category/${id}`, {
            "name": name
        });
    } catch (error: any) {
        throw error
    }
}