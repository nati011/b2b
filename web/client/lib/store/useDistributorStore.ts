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
        set({ loading: true, error: null, success: null, distributor_id: null });
        try {
            const distributor_id = await RegisterDistributor(profile);
            // Only set success if we got a valid distributor_id
            if (distributor_id) {
                set({ loading: false, success: "Distributor registered successfully", distributor_id: distributor_id, error: null });
            } else {
                const errMsg = "Registration succeeded but no distributor ID was returned";
                toast.error(errMsg);
                set({ loading: false, error: errMsg, success: null, distributor_id: null });
                throw new Error(errMsg);
            }
        } catch (error: any) {
            const errMsg = typeof error === "string" ? error : error?.message || "Failed to register distributor";
            toast.error(errMsg);
            set({ loading: false, error: errMsg, success: null, distributor_id: null });
            throw error; // Re-throw to allow component to handle it
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
            set({ loading: false, error: null });
        } catch (error: any) {
            const errMsg = typeof error === "string" ? error : error?.message || "Failed to buy subscription";
            toast.error(errMsg)
            set({ loading: false, error: errMsg });
        }
    },
}));

export default useDistributorStore;
