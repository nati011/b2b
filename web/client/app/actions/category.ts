'use server'
import axiosIns from "@/lib/axios";

export const GetAllCategories = async () => {
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
            "/category",
            { name: name },
        );
        return response.data.detail
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error has occured while creating the product"
        }
        throw "An error has occured while creating the product"
    }
}

export const Delete = async (id: number) => {
    try {
        const response = await axiosIns.delete(`/category?id=${id}`);
        return response.data.detail
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error has occured while creating the product"
        }
        throw "An error has occured while creating the product"
    }
}

export const Update = async (id: number, name: string) => {
    try {
        const response = await axiosIns.patch(`/category/${id}`, {
            "name": name
        });
        return response.data.detail
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error has occured while creating the product"
        }
        throw "An error has occured while creating the product"
    }
}