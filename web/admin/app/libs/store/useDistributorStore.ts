import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Distributor, DistributorRequest, UserAccount } from '@/app/libs/types';
import { Create, GetAll, GetDistributorUser } from '@/actions/distributor';
import { GetById } from '@/actions/retailer';

interface DistributorsStore {
    success: string | null
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
        latitude: "",
        longitude: "",
        general_zone: '',
        region: '',
        woreda: '',
        user: [],
        is_active: false
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
    success: null,
    loading: false,
    error: null,
    next: null,
    previous: null,

    fetchDistributors: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await GetAll()
            set({
                distributors: response,
                loading: false
            });
        } catch (error: any) {
            console.log(error)
            set({ error: error.response.data, loading: false });
        }
    },

    createDistributors: async (DistributorsData: DistributorRequest) => {
        set({ loading: true, error: null });
        try {
            console.log(DistributorsData)
            const response = await Create(DistributorsData);
            console.log(response)
            set(state => ({
                success: response,
                loading: false
            }));
        } catch (error: any) {
            console.log(error)
            set({ error: error.message, loading: false });
        }
    },
    fetchDistributorDetail: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await GetById(id);
            set({
                distributor: response.data.body.distributor,
                loading: false
            });
        } catch (error: any) {
            set({ error: error.response.data || "An error has occured", loading: false });
        }
    },
    fetchDistributorUser: async (id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await GetDistributorUser(id);
            set({
                distributorUser: response.data.body.user,
                loading: false
            });
        } catch (error: any) {
            set({ error: error, loading: false });
        }
    }
}));

export default useDistributorsStore;