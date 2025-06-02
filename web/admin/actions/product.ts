'use server'
import axiosIns from "@/app/libs/axios";
import { ProductForm } from "@/app/libs/types";

export const fetchProducts = async (url?: string) => {
    try {
        const response = await axiosIns.get("/product");
        return response.data
    } catch (error) {
        throw error
    }
}
export const fetchConfigurableProducts = async (url?: string) => {
    try {
        const response = await axiosIns.get("/configurable_product");
        return response.data
    } catch (error) {
        throw error
    }
}

export const updateProduct = async (ProductsData: Partial<ProductForm>) => {
    try {
        const response = await axiosIns.put("/product/", ProductsData);
        return response.data
    } catch (error) {
        throw error
    }
}


export const createProduct = async (productData: any) => {
    try {
        const response = await axiosIns.post("/product", productData);
        return response.data.detail
    } catch (error: any) {

        if (error.response) {
            throw error.response.data.message || "An error has occured while creating the product"
        }

        throw "An error has occured while creating the product"
    }
}

export const createConfigurableProduct = async (productData: any) => {
    try {
        const response = await axiosIns.post('/configurable_product', productData);
        return response.data
    } catch (error) {
        throw error
    }
}

export const addStock = async (stock: number, id: number) => {
    try {
        const response = await axiosIns.patch(`/product/${id}/stock?amount=${stock}&command=receive`);
        return response.data
    } catch (error) {
        throw error
    }
}

export const depleteStock = async (stock: number, id: number) => {
    try {
        const response = await axiosIns.patch(`/product/${id}/stock?amount=${stock}&command=deplete`)
        return response.data
    } catch (error) {
        throw error
    }
}

export const updateProductStatus = async (id: number, command: string) => {
    try {
        const response = await axiosIns.patch(`/product/${id}/status?command=${command}`);
        return response.data
    } catch (error) {
        throw error
    }
}

export const fetchProductDetail = async (id: number) => {
    try {
        const response = await axiosIns.get(`/product?id=${id}`);
        response.data.body.Attributes = [response.data.body.Attributes]
        return response.data.body
    } catch (error) {
        throw error
    }
}

export const fetchConfigurableProductDetail = async (id: number) => {
    try {
        const response = await axiosIns.get(`/configurable_product?id=${id}`);
        // Note: TEMP
        response.data.body.Attributes = [response.data.body.Attributes]
        return response.data.body
    } catch (error) {
        throw error
    }
}

export const updateConfigurableProductStatus = async (id: number, command: string) => {
    try {
        const response = await axiosIns.patch(`/configurable_product/${id}/status?command=${command}`);
        return response.data
    } catch (error) {
        throw error
    }
}