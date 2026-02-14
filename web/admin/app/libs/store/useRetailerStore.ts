import { create } from 'zustand'
import { Retailer } from '@/app/libs/types';
import { Create, GetAll, GetById } from '@/app/actions/retailer';

interface RetailersStore {
  success: string | null
  retailers: Retailer[];
  retailer: Retailer;
  loading: boolean;
  error: string | null;
  next: string | null;
  previous: string | null;
  total: number | null;

  fetchRetailers: (page: number) => Promise<void>;
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
  success: null,
  error: null,
  next: null,
  previous: null,
  total: 0,

  fetchRetailers: async (page: number) => {
    set({ loading: true, error: null });
    try {
      const response = await GetAll(page);
      set({
        retailers: response.list,
        total: response.total_count,
        loading: false
      });
    } catch (error: any) {
      set({ error: error, loading: false });
    }
  },
  fetchRetailer: async (id?: number) => {
    console.log(id, "IDDDDDDDDDDDDDD")
    set({ loading: true, error: null });
    try {
      const response = await GetById(id);
      console.log(response, "Retailer Response")
      set({
        retailer: response,
        loading: false
      });
    } catch (error: any) {
      console.log(error)
      set({ error: error, loading: false });
    }
  },

  createRetailers: async (RetailersData: Partial<Retailer>) => {
    set({ loading: true, error: null });
    try {
      const response = await Create(RetailersData);
      set(state => ({
        success: response,
        loading: false
      }));
    } catch (error: any) {
      set({ error: error, loading: false });
    }
  },

}));

export default useRetailersStore;