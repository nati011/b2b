import { create } from 'zustand'
import { Distributor, DistributorRequest, DistributorUserRequest, UserDetail } from '@/app/libs/types';
import { Create, GetAll, GetDistributorUser, GetById, DistributorOnBoardingReview, UpdateDistributorStatus, CreateDistributorUser } from '@/app/actions/distributor';
import { toast } from 'sonner';

interface DistributorsStore {
    success: string | null
    distributors: Distributor[];
    distributor: Distributor |  null;
    distributorUser: UserDetail[];
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    totalCount: number | null;
    userCount: number |null;

    fetchDistributors: (status?: string) => Promise<void>;
    createDistributors: (DistributorsData: DistributorRequest) => Promise<void>;
    fetchDistributorDetail: (id: number) => Promise<void>
    fetchDistributorUser: () => Promise<void>
    approveDistributor: (id: number) => Promise<void>
    rejectDistributor: (id: number, comment: string) => Promise<void>
    activateDistributor: (id: number) => Promise<void>
    deactivateDistributor: (id: number) => Promise<void>
    createDistirbutorUser: (user: DistributorUserRequest) => Promise<void>
}

const useDistributorsStore = create<DistributorsStore>((set) => ({
    distributors: [],
    distributor: null,
    distributorUser: [],
    success: null,
    loading: false,
    error: null,
    next: null,
    previous: null,
    totalCount: null,
    userCount: null,

    fetchDistributors: async (status?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await GetAll(status)
            console.log(response)
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
    fetchDistributorUser: async () => {
        set({ loading: true, error: null });
        try {
            const response = await GetDistributorUser();
            set({
                distributorUser: response.List,
                userCount: response.TotalCount,
                loading: false
            });
        } catch (error: any) {
            toast.error(error.message)
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
    },

    createDistirbutorUser: async(user: DistributorUserRequest) =>{
        set({ loading: true, error: null })
        try {
            const response = await CreateDistributorUser(user)
            set({ success: response, loading: false })
            await useDistributorsStore.getState().fetchDistributorUser()
        } catch (error: any) {
            set({ error: error.message || "Error occured while creating distributor user.", loading: false });
        }
    }
}));

export default useDistributorsStore;