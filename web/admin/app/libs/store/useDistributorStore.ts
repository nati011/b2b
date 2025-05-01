import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Distributor } from '@/app/libs/types';

interface DistributorsStore {
    distributors: Distributor[];
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchDistributors: (url?: string) => Promise<void>;
    createDistributors: (DistributorsData: Partial<Distributor>) => Promise<void>;
}

const useDistributorsStore = create<DistributorsStore>((set) => ({
    distributors: [],
    loading: false,
    error: null,
    next: null,
    previous: null,

    fetchDistributors: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get('/api/distributor');
            set({
                distributors: response.data.body.distributors.list,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch distributor', loading: false });
        }
    },

    createDistributors: async (DistributorsData: Partial<Distributor>) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post('/api/distributor/', DistributorsData);
            set(state => ({
                Distributors: [...state.distributors, response.data.detail],
                loading: false
            }));
        } catch (error) {
            set({ error: 'Failed to create distributor', loading: false });
        }
    },

}));

export default useDistributorsStore;