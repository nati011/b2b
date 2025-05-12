import { create } from "zustand";
import axiosIns from "@/app/libs/axios";
import { Product, Category } from "@/app/libs/types";

interface ProductsStore {
  success: string;
  products: Product[];
  product: Product;
  loading: boolean;
  error: string | null;
  next: string | null;
  previous: string | null;
  categories: Category[];
  categoriesLoading: boolean;
  categoriesError: string | null;

  fetchProducts: (url?: string) => Promise<void>;
  // createProducts: (ProductsData: Partial<Product>) => Promise<void>;
  fetchCategories: () => Promise<void>;
  createCategory: (name: string) => Promise<void>;
  deleteCategory: (id: number) => Promise<void>;
  createProduct: (productData: any) => Promise<void>;
  createConfigurableProduct: (productData: any) => Promise<void>
  addStock: (stock: number, id: number) => Promise<void>;
  depleteStock: (stock: number, id: number) => Promise<void>;
  updateProductStatus: (id: number, productStatus: boolean) => Promise<void>;
  fetchProductDetail: (id: number) => Promise<void>;
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
    CategoryId: 0,
    Stock: 0,
    AvailableStock: 0,
    ReservedStock: 0,
    IsActive: 0
  },
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

  // createProducts: async (ProductsData: Partial<Product>) => {
  //   set({ loading: true, error: null });
  //   try {
  //     const response = await axiosIns.post("/product/", ProductsData);
  //     set((state) => ({
  //       products: [...state.products, response.data.detail],
  //       loading: false,
  //     }));
  //   } catch (error) {
  //     set({ error: "Failed to create product", loading: false });
  //   }
  // },

  fetchCategories: async () => {
    set({ categoriesLoading: true, categoriesError: null });
    try {
      const response = await axiosIns.get("/api/category");
      set({
        categories: response.data.body.category.categories,
        categoriesLoading: false,
      });
    } catch (error) {
      set({
        categoriesError: "Failed to fetch categories",
        categoriesLoading: false,
      });
    }
  },
  createCategory: async (name: string) => {
    set({ categoriesLoading: true, categoriesError: null });
    try {
      const response = await axiosIns.post(
        "/api/category",
        { name: name },
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      await useProductsStore.getState().fetchCategories();
      set({ categoriesLoading: false });
    } catch (error: any) {
      set({ categoriesLoading: false, categoriesError: error.message });
    }
  },
  deleteCategory: async (id: number) => {
    set({ categoriesLoading: true, categoriesError: null });
    try {
      await axiosIns.delete(`/api/category?id=${id}`);
      await useProductsStore.getState().fetchCategories();
      set({ categoriesLoading: false });
      await useProductsStore.getState().fetchCategories();
    } catch (error: any) {
      set({ categoriesLoading: false, categoriesError: error.message });
    }
  },
  createProduct: async (productData: any) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.post("/api/product", productData, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      console.log(productData)
      console.log(response.data)
      await useProductsStore.getState().fetchProducts();
      set({ loading: false });
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
      const response = await axiosIns.get(`/api/product/${id}`);
      console.log(response.data)
      set({
        product: response.data.body,
        loading: false,
      });
    } catch (error) {
      set({ error: "Failed to fetch products", loading: false });
    }
  },
}));

export default useProductsStore;
