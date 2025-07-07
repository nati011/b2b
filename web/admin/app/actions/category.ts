'use server'
import axiosIns from "@/app/libs/axios";
import { withErrorHandling } from "@/app/libs/error-handling";

export const GetAll = async () => {
    try {
        const response = await axiosIns.get("/category");
        return response.data.body.categories
    } catch (error) {
        throw error
    }
}

export const Create = async (name: string) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.post(
            "/category",
            { name: name },
        );
        return response.data.detail
      });
    
}

export const Delete = async (id: number) => {
    return withErrorHandling(async () => {
        await axiosIns.delete(`/category?id=${id}`);
      });
}

export const Update = async (id: number, name: string) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.patch(`/category/${id}`, {
            "name": name
        });
        return response.data.message
      });
}