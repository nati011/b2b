'use server'
import axiosIns from "@/app/libs/axios";
import { ProductForm } from "@/app/libs/types";
import { withErrorHandling } from "@/app/libs/error-handling";

export const fetchProducts = async (url?: string) => {
    return withErrorHandling(async () => {
        let requestUrl = '/product'
        console.log(requestUrl)
        const response = await axiosIns.get(requestUrl);
        console.log(response.data)
        return response.data.body
      });
}
export const fetchConfigurableProducts = async (url?: string) => {
    return withErrorHandling(async () => {
        let requestUrl = '/configurable_product'
        console.log(requestUrl)
        const response = await axiosIns.get(requestUrl);
        console.log(response.data)
        return response.data
    });
    }

export const updateProduct = async (ProductsData: Partial<ProductForm>) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.put("/product/", ProductsData);
        return response.data
    });
}


export const createProduct = async (productData: any) => {
    return withErrorHandling(async () => {
        console.log(productData)
        const response = await axiosIns.post("/product", productData);
        return response.data.detail
    });
}

export const createConfigurableProduct = async (productData: any) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.post('/configurable_product', productData);
        return response.data
    });
}

export const addStock = async (stock: number, id: number) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.patch(`/product/${id}/stock?amount=${stock}&command=receive`);
        return response.data
    });
};

export const depleteStock = async (stock: number, id: number) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.patch(`/product/${id}/stock?amount=${stock}&command=deplete`)
        return response.data
    });
}

export const updateProductStatus = async (id: number, command: string) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.patch(`/product/${id}/status?command=${command}`);
        return response.data
    });
};

export const fetchProductDetail = async (id: number) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.get(`/product?id=${id}`);
        const attrs = response.data.body.Attributes;
        if (attrs && typeof attrs === "object" && !Array.isArray(attrs)) {
            response.data.body.Attributes = Object.entries(attrs).map(([key, value]) => ({ key, value }));
        } else if (Array.isArray(attrs)) {
            response.data.body.Attributes = attrs;
        } else {
            response.data.body.Attributes = [];
        }
        return response.data.body
    });
};


export const fetchConfigurableProductDetail = async (id: number) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.get(`/configurable_product?id=${id}`);
        console.log(response.data.body)
        const attrs = response.data.body.configurable_product.Attributes;
        if (attrs && typeof attrs === "object" && !Array.isArray(attrs)) {
            response.data.body.configurable_product.Attributes = Object.entries(attrs).map(([key, value]) => ({ key, value }));
        } else if (Array.isArray(attrs)) {
            response.data.body.configurable_product.Attributes = attrs;
        } else {
            response.data.body.configurable_product.Attributes = [];
        }
        console.log(response.data.body)
        return response.data.body.configurable_product
    });
};


export const updateConfigurableProductStatus = async (id: number, command: string) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.patch(`/configurable_product/${id}/status?command=${command}`);
        return response.data
    });
};

export const updateConfigurableProduct = async (productData: any) => {
    return withErrorHandling(async () => {
        const response = await axiosIns.put('/configurable_product', productData);
        return response.data;
    });
};