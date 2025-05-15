import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Retailer } from '@/app/libs/types';

interface RetailersStore {
  retailers: Retailer[];
  retailer: Retailer;
  loading: boolean;
  error: string | null;
  next: string | null;
  previous: string | null;

  fetchRetailers: (url?: string) => Promise<void>;
  fetchRetailer: (id: number) => Promise<void>;
  createRetailers: (RetailersData: Partial<Retailer>) => Promise<void>;
}

const useRetailersStore = create<RetailersStore>((set) => ({
  retailers: [],
  retailer: {
    id: 0,
    name: '',
    tin: '',
    latitude: '',
    longitude: '',
    general_zone: '',
    region: '',
    woreda: '',
    user: {
      id: 0,
      first_name: '',
      last_name: '',
      email: '',
      phone: '',
      username: '',
      dob: '',
      is_active: false,
      external_id: ''
    }
  },
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
  fetchRetailer: async (id?: number) => {
    set({ loading: true, error: null });
    try {
      const response = await axiosIns.get(`/api/retailer?id=${id}`);
      set({
        retailer: response.data.body.retailer,
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