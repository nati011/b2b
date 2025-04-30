import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Product } from '@/app/libs/types';

interface ProductsStore {
  products: Product[];
  loading: boolean;
  error: string | null;
  next: string | null;
  previous: string | null;

  fetchProducts: (url?: string) => Promise<void>;
  createProducts: (ProductsData: Partial<Product>) => Promise<void>;
}

const useProductsStore = create<ProductsStore>((set) => ({
  products: [],
  loading: false,
  error: null,
  next: null,
  previous: null,

  fetchProducts: async (url?: string) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.get('/product');
      set({
        products: response.data.body.Products,
        loading: false
      });
    } catch (error) {
      set({ error: 'Failed to fetch Loan', loading: false });
    }
  },

  createProducts: async (ProductsData: Partial<Product>) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.post('/product/', ProductsData);
      set(state => ({
        Products: [...state.products, response.data.detail],
        loading: false
      }));
    } catch (error) {
      set({ error: 'Failed to create product', loading: false });
    }
  },

}));

export default useProductsStore;