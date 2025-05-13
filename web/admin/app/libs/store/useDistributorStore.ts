import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Distributor, DistributorRequest } from '@/app/libs/types';

interface DistributorsStore {
    distributors: Distributor[];
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchDistributors: (url?: string) => Promise<void>;
    createDistributors: (DistributorsData: DistributorRequest) => Promise<void>;
    fetchDistributorDetail: (id: number) => Promise<void>
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
                distributors: response.data.body.distributors,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch distributor', loading: false });
        }
    },

    createDistributors: async (DistributorsData: DistributorRequest) => {
        set({ loading: true, error: null });
        try {
            console.log(DistributorsData)
            const response = await axiosIns.post('/api/distributor/', DistributorsData);
            set(state => ({
                Distributors: [...state.distributors, response.data.detail],
                loading: false
            }));
        } catch (error) {
            set({ error: 'Failed to create distributor', loading: false });
        }
    },
    fetchDistributorDetail: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get(`/api/distributor/${id}`);
            set({
                distributors: response.data.body.distributors,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch distributor', loading: false });
        }
    }
}));

export default useDistributorsStore;