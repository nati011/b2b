import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { ConfigurableProduct } from "@/app/libs/types";

interface ProductsStore {
    success: string;
    configurableProducts: ConfigurableProduct[];
    product: ConfigurableProduct;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchConfigurableProduct: (url?: string) => Promise<void>;
    updateConfigurableProduct: (ProductsData: Partial<ConfigurableProduct>) => Promise<void>;
    createConfigurabileProduct: (productData: any) => Promise<void>;
    createConfigurableProduct: (productData: any) => Promise<void>
    addStock: (stock: number, id: number) => Promise<void>;
    depleteStock: (stock: number, id: number) => Promise<void>;
    updateProductStatus: (id: number, productStatus: boolean) => Promise<void>;
    fetchConfigurableProductDetail: (id: number) => Promise<void>;
}

const useConfigurableProductStore = create<ProductsStore>((set) => ({
    success: "",
    configurableProducts: [],
    product: {
        Id: 0,
        Name: "",
        Desc: "",
        ExternalId: "",
        Images: [],
        Attributes: [],
        DistributorId: 0,
        CategoryId: 0,
        PriceRange: {
            min: 0,
            max: 0
        },
        IsAvailable: false,
        Products: []
    },
    loading: false,
    error: null,
    next: null,
    previous: null,


    fetchConfigurableProduct: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get("/api/configurable_product");
            console.log(response.data)
            set({
                configurableProducts: response.data.body.List,
                loading: false,
            });
        } catch (error) {
            set({ error: "Failed to fetch configurableProducts", loading: false });
        }
    },

    updateConfigurableProduct: async (ProductsData: Partial<ConfigurableProduct>) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.put("/api/configurable_product", ProductsData);
            set((state) => ({
                loading: false,
            }));
            await useConfigurableProductStore.getState().fetchConfigurableProduct();
        } catch (error) {
            set({ error: "Failed to create product", loading: false });
        }
    },


    createConfigurabileProduct: async (productData: any) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post("/api/configurable_product", productData, {
                headers: {
                    "Content-Type": "application/json",
                },
            });
            await useConfigurableProductStore.getState().fetchConfigurableProduct();
            set({ loading: false, success: response.data });
        } catch (error: any) {
            set({ loading: false, error: error.message });
        }
    },
    createConfigurableProduct: async (productData: any) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post('/api/configurable_product', productData, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            await useConfigurableProductStore.getState().fetchConfigurableProduct();
            set({ loading: false });
        } catch (error: any) {
            set({ loading: false, error: error.message });
        }
    },
    addStock: async (stock: number, id: number) => {
        try {
            const response = await axiosIns.patch(`/api/configurable_product${id}/stock?amount=${stock}&command=receive`);
            console.log(response.data)
            await useConfigurableProductStore.getState().fetchConfigurableProduct();
            set({ loading: false });
        } catch (error: any) {
            set({ loading: false, error: error.message });
        }
    },
    depleteStock: async (stock: number, id: number) => {
        try {
            const response = await axiosIns.patch(`/api/configurable_product${id}/stock?amount=${stock}&command=deplete`, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            await useConfigurableProductStore.getState().fetchConfigurableProduct();
            set({ loading: false });
        } catch (error: any) {
            set({ loading: false, error: error.message });
        }
    },
    updateProductStatus: async (id: number, productStatus: boolean) => {
        try {
            const command = productStatus ? "deactivate" : "activate"
            const response = await axiosIns.patch(`/api/configurable_product${id}/status?command=${command}`, {
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            await useConfigurableProductStore.getState().fetchConfigurableProduct();
            console.log(response.data)
            set({ loading: false, success: response.data.message, });
        } catch (error: any) {
            set({ loading: false, error: error.message });
        }
    },
    fetchConfigurableProductDetail: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get(`/api/configurable_product?id=${id}`);
            set({
                product: response.data.body.configurable_product,
                loading: false,
            });
        } catch (error) {
            set({ error: "Failed to fetch configurableProducts", loading: false });
        }
    },
}));

export default useConfigurableProductStore;
