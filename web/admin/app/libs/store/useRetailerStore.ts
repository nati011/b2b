import {create} from 'zustand'
import axiosIns from '@/app/libs/axios'
import {Retailer} from '@/app/libs/types';

interface RetailersStore {
  retailers: Retailer[];
  loading: boolean;
  error: string | null;
  next: string | null;
  previous: string | null;

  fetchRetailers: (url?: string) => Promise<void>;
  createRetailers: (RetailersData: Partial<Retailer>) => Promise<void>;
}

const useRetailersStore = create<RetailersStore>((set) => ({
  retailers: [],
  loading: false,
  error: null,
  next: null,
  previous: null,

  fetchRetailers: async (url?: string) => {
    set({ loading: true, error: null });
    try {
        const response = await axiosIns.get('/api/retailer');
        set({ 
          retailers: response.data.body.Retailers,
          loading: false 
        });
    } catch (error) {
        set({ error: 'Failed to fetch Loan', loading: false });
    }
  },

  createRetailers: async (RetailersData: Partial<Retailer>) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.post('/retailer/', RetailersData);
      set(state => ({
        Retailers: [...state.retailers, response.data.detail],
        loading: false
      }));
    } catch (error) {
      set({ error: 'Failed to create retailer', loading: false });
    }
  },

}));

export default useRetailersStore;