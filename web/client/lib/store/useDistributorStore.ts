import { create } from "zustand";
import { BuySubscriptionRequest, DistributorRequest } from "@/lib/types";
import { BuySubscription, RegisterDistributor } from "@/app/actions/auth";
import { toast } from "sonner";

interface DistributorStore {
    success: string | null;
    loading: boolean;
    error: string | null;
    distributor_id: number | null;

    register: (profile: DistributorRequest) => Promise<void>;
    buySubscription: (profile: BuySubscriptionRequest) => Promise<void>;
}

const useDistributorStore = create<DistributorStore>((set) => ({
    success: null,
    loading: false,
    error: null,
    distributor_id: null,

    register: async (profile: DistributorRequest) => {
        set({ loading: true, error: null, success: null });
        try {
            const distributor_id = await RegisterDistributor(profile);
            set({ loading: false, success: "Distributor registered successfully", distributor_id: distributor_id });
        } catch (error: any) {
            const errMsg = typeof error === "string" ? error : error?.message || "Failed to register distributor";
            toast.error(errMsg)
            set({ loading: false });
        }
    },
    buySubscription: async (profile: BuySubscriptionRequest) => {
        set({ loading: true, error: null, success: null });
        try {
            const response = await BuySubscription(profile);
            localStorage.setItem("tx_ref", response.tx_ref)

            if (response) {
                window.location.href = response.checkout_url
            }
            set({ loading: false});
        } catch (error: any) {
            const errMsg = typeof error === "string" ? error : error?.message || "Failed to buy subscription";
            toast.error(errMsg)
            set({ loading: false});
        }
    },
}));

export default useDistributorStore;
