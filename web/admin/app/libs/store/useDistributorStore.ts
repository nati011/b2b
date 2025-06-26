import { create } from 'zustand'
import { Distributor, DistributorRequest, UserAccount } from '@/app/libs/types';
import { Create, GetAll, GetDistributorUser, GetById, DistributorOnBoardingReview, UpdateDistributorStatus } from '@/app/actions/distributor';

interface DistributorsStore {
    success: string | null
    distributors: Distributor[];
    distributor: Distributor |  null;
    distributorUser: UserAccount
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    totalCount: number | null;

    fetchDistributors: (status?: string) => Promise<void>;
    createDistributors: (DistributorsData: DistributorRequest) => Promise<void>;
    fetchDistributorDetail: (id: number) => Promise<void>
    fetchDistributorUser: (id: number) => Promise<void>
    approveDistributor: (id: number) => Promise<void>
    rejectDistributor: (id: number, comment: string) => Promise<void>
    activateDistributor: (id: number) => Promise<void>
    deactivateDistributor: (id: number) => Promise<void>
}

const useDistributorsStore = create<DistributorsStore>((set) => ({
    distributors: [],
    distributor: null,
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
    totalCount: null,

    fetchDistributors: async (status?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await GetAll(status)
            set({
                distributors: response.distributors,
                totalCount: response.total_count,
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
            set({
                success: response,
                loading: false
            });
        } catch (error: any) {
            console.log(error)
            set({ error: error.message, loading: false });
        }
    },
    fetchDistributorDetail: async (id: number) => {
        set({ loading: true, error: null, distributor: null });
        try {
            const response = await GetById(id);
            console.log(response, "Response")
            set({
                distributor: response,
                loading: false
            });
        } catch (error: any) {
            console.log(error)
            set({ error: error.message || "An error has occured", loading: false });
        }
    },
    fetchDistributorUser: async (id: number) => {
        console.log(id)
        set({ loading: true, error: null });
        try {
            console.log("EZIII", id)
            const response = await GetDistributorUser(id);
            set({
                // @ts-ignore
                distributorUser: response,
                loading: false
            });
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    },
    activateDistributor: async (id: number) => {
        set({ loading: true, error: null })
        try {
            const response = await UpdateDistributorStatus(id, "activate")
            set({ success: response, loading: false })
            await useDistributorsStore.getState().fetchDistributorDetail(id);
        } catch (error: any) {
            set({ error: error.message || "Error occured while approving distributor.", loading: false });
        }
    },
    deactivateDistributor: async (id: number) => {
        set({ loading: true, error: null })
        try {
            const response = await UpdateDistributorStatus(id, "deactivate")
            set({ success: response, loading: false })
            await useDistributorsStore.getState().fetchDistributorDetail(id);
        } catch (error: any) {
            set({ error: error.message || "Error occured while approving distributor.", loading: false });
        }
    },
    approveDistributor: async (id: number) => {
        set({ loading: true, error: null })
        try {
            const response = await DistributorOnBoardingReview(id, "approve")
            set({ success: response, loading: false })
            await useDistributorsStore.getState().fetchDistributorDetail(id);
        } catch (error: any) {
            set({ error: error.message || "Error occured while approving distributor.", loading: false });
        }
    },

    rejectDistributor: async (id: number, comment: string) => {
        set({ loading: true, error: null })
        try {
            const response = await DistributorOnBoardingReview(id, "reject", comment)
            set({ success: response, loading: false })
            await useDistributorsStore.getState().fetchDistributorDetail(id);
        } catch (error: any) {
            set({ error: error.message || "Error occured while rejecting distributor.", loading: false });
        }
    }
}));

export default useDistributorsStore;