import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Product, ConfigurableProduct, ProductForm } from "@/app/libs/types";

interface ProductsStore {
  success: string;
  products: Product[];
  configurable_products: ConfigurableProduct[];
  configurable_product: ConfigurableProduct;
  product: Product;
  loading: boolean;
  error: string | null;
  next: string | null;
  previous: string | null;

  fetchProducts: (url?: string) => Promise<void>;
  fetchConfigurableProducts: (url?: string) => Promise<void>;
  updateProduct: (ProductsData: Partial<ProductForm>) => Promise<void>;
  createProduct: (productData: any) => Promise<void>;
  createConfigurableProduct: (productData: any) => Promise<void>
  addStock: (stock: number, id: number) => Promise<void>;
  depleteStock: (stock: number, id: number) => Promise<void>;
  updateProductStatus: (id: number, productStatus: boolean) => Promise<void>;
  fetchProductDetail: (id: number) => Promise<void>;
  fetchConfigurableProductDetail: (id: number) => Promise<void>;
  updateConfigurableProductStatus: (id: number, productStatus: boolean) => Promise<void>;
}

const useProductsStore = create<ProductsStore>((set) => ({
  success: "",
  products: [],
  product: {
    Id: 0,
    Name: "",
    Desc: "",
    ExternalID: "",
    Images: [],
    Price: 0,
    Attributes: [],
    DistributorId: 0,
    CategoryId: [],
    Stock: 0,
    AvailableStock: 0,
    ReservedStock: 0,
    IsActive: 0
  },
  configurable_product: {
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
  configurable_products: [],
  loading: false,
  error: null,
  next: null,
  previous: null,
  categories: [],
  categoriesLoading: false,
  categoriesError: null,


  fetchProducts: async (url?: string) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.get("/api/product");
      console.log(response.data)
      set({
        products: response.data.body.List,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
  fetchConfigurableProducts: async (url?: string) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.get("/api/configurable_product");
      console.log(response.data)
      set({
        configurable_products: response.data.body.configurable_products.List,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },

  updateProduct: async (ProductsData: Partial<ProductForm>) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.put("/api/product/", ProductsData);
      set((state) => ({
        loading: false,
      }));
      await useProductsStore.getState().fetchProducts();
    } catch (error) {
      set({ error: "Failed to create product", loading: false });
    }
  },

  createProduct: async (productData: any) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.post("/api/product", productData);
      await useProductsStore.getState().fetchProducts();
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
      await useProductsStore.getState().fetchProducts();
      if (response.status = 200) {
        window.location.href = '/products/configurable'
      }
      set({ loading: false });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  addStock: async (stock: number, id: number) => {
    try {
      const response = await axiosIns.patch(`/api/product/${id}/stock?amount=${stock}&command=receive`);
      console.log(response.data)
      await useProductsStore.getState().fetchProducts();
      set({ loading: false });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  depleteStock: async (stock: number, id: number) => {
    try {
      const response = await axiosIns.patch(`/api/product/${id}/stock?amount=${stock}&command=deplete`, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
      await useProductsStore.getState().fetchProducts();
      set({ loading: false });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  updateProductStatus: async (id: number, productStatus: boolean) => {
    try {
      const command = productStatus ? "deactivate" : "activate"
      const response = await axiosIns.patch(`/api/product/${id}/status?command=${command}`, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
      await useProductsStore.getState().fetchProducts();
      console.log(response.data)
      set({ loading: false, success: response.data.message, });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  fetchProductDetail: async (id: number) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.get(`/api/product?id=${id}`);
      response.data.body.Attributes = [response.data.body.Attributes]
      set({
        product: response.data.body,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
  fetchConfigurableProductDetail: async (id: number) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.get(`/api/configurable_product?id=${id}`);
      console.log(response.data.body)
      // Note: TEMP
      response.data.body.Attributes = [response.data.body.Attributes]
      set({
        product: response.data.body,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
  updateConfigurableProductStatus: async (id: number, productStatus: boolean) => {
    try {
      const command = productStatus ? "deactivate" : "activate"
      const response = await axiosIns.patch(`/api/configurable_product/${id}/status?command=${command}`);
      await useProductsStore.getState().fetchConfigurableProducts();
      console.log(response.data)
      set({ loading: false, success: response.data.message, });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
}));

export default useProductsStore;
