import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Product, ConfigurableProduct, ProductForm } from "@/app/libs/types";
import { addStock, createConfigurableProduct, createProduct, depleteStock, fetchConfigurableProductDetail, fetchConfigurableProducts, fetchProductDetail, fetchProducts, updateConfigurableProductStatus, updateProduct, updateProductStatus, updateConfigurableProduct } from "@/app/actions/product";

interface ProductsStore {
  success: string | null;
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
  updateConfigurableProduct: (productData: any) => Promise<void>;
}

const useProductsStore = create<ProductsStore>((set) => ({
  success: null,
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

  fetchProducts: async (url?: string) => {
    set({ loading: true, error: null });
    try {
      const resp = await fetchProducts()
      console.log(resp)
      set({
        products: resp.body.List,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
  fetchConfigurableProducts: async (url?: string) => {
    set({ loading: true, error: null });
    try {
      const resp = await fetchConfigurableProducts()
      set({
        configurable_products: resp.body.configurable_products.List,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },

  updateProduct: async (ProductsData: Partial<ProductForm>) => {
    set({ loading: true, error: null });
    try {
      await updateProduct(ProductsData)
      set({
        loading: false,
      });
      await useProductsStore.getState().fetchProducts();
    } catch (error) {
      set({ error: "Failed to create product", loading: false });
    }
  },

  createProduct: async (productData: any) => {
    set({ loading: true, error: null });
    try {
      const response = await createProduct(productData);
      console.log(productData)
      await useProductsStore.getState().fetchProducts();
      set({ loading: false, success: response });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  createConfigurableProduct: async (productData: any) => {
    set({ loading: true, error: null });
    try {
      const response = await createConfigurableProduct(productData);
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
      await addStock(stock, id);
      await useProductsStore.getState().fetchProducts();
      set({ loading: false });
    } catch (error: any) {
      console.log("Error here____________________________")
      set({ loading: false, error: error.message || "An error has occured while updating stock." });
    }
  },
  depleteStock: async (stock: number, id: number) => {
    try {
      await depleteStock(stock, id);
      await useProductsStore.getState().fetchProducts();
      set({ loading: false });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  updateProductStatus: async (id: number, productStatus: boolean) => {
    try {
      const command = productStatus ? "deactivate" : "activate"
      const response = await updateProductStatus(id, command);
      await useProductsStore.getState().fetchProducts();
      set({ loading: false, success: response.message, });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  fetchProductDetail: async (id: number) => {
    set({ loading: true, error: null });
    try {
      const detail = await fetchProductDetail(id);
      set({
        product: detail,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
  fetchConfigurableProductDetail: async (id: number) => {
    set({ loading: true, error: null });
    try {
      const detail = await fetchConfigurableProductDetail(id);
      console.log(detail)
      set({
        product: detail,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
  updateConfigurableProductStatus: async (id: number, productStatus: boolean) => {
    try {
      const command = productStatus ? "deactivate" : "activate"
      const response = await updateConfigurableProductStatus(id, command);
      await useProductsStore.getState().fetchConfigurableProducts();
      set({ loading: false, success: response.message, });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
  updateConfigurableProduct: async (productData: any) => {
    set({ loading: true, error: null });
    try {
      const response = await updateConfigurableProduct(productData);
      await useProductsStore.getState().fetchConfigurableProducts();
      set({ loading: false, success: response });
    } catch (error: any) {
      set({ loading: false, error: error.message });
    }
  },
}));

export default useProductsStore;
