import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Distributor, DistributorRequest, UserAccount } from '@/app/libs/types';

interface DistributorsStore {
    distributors: Distributor[];
    distributor: Distributor;
    distributorUser: UserAccount
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchDistributors: (url?: string) => Promise<void>;
    createDistributors: (DistributorsData: DistributorRequest) => Promise<void>;
    fetchDistributorDetail: (id: number) => Promise<void>
    fetchDistributorUser: (id: number) => Promise<void>
}

const useDistributorsStore = create<DistributorsStore>((set) => ({
    distributors: [],
    distributor: {
        id: 0,
        name: '',
        tin: '',
        latitude: 0,
        longitude: 0,
        general_zone: '',
        region: '',
        woreda: '',
        user: []
    },
    distributorUser: {
        id: 0,
        first_name: '',
        last_name: '',
        email: '',
        phone: '',
        username: '',
        dob: '',
        is_active: false,
        external_id: ''
    },
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
            const response = await axiosIns.get(`/api/distributor?id=${id}`);
            set({
                distributor: response.data.body.distributor,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch distributor', loading: false });
        }
    },
    fetchDistributorUser: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get(`/api/user?id=${id}`);
            set({
                distributorUser: response.data.body.user,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch distributor', loading: false });
        }
    }
}));

export default useDistributorsStore;