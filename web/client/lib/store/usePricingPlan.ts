import { create } from "zustand";
import { PricingPlan } from "@/lib/types";
import { fetchPricingPlans } from "@/app/actions/pricingPlans";


interface PricingPlanStore {
    success: string | null;
    plans: PricingPlan[];
    loading: boolean;
    error: string | null;

    fetchPricingPlan: () => Promise<void>;
}

const usePlanstore = create<PricingPlanStore>((set) => ({
    success: null,
    isLoggedIn: true,
    plans: [],
    loading: false,
    error: null,

    fetchPricingPlan: async () => {
        set({ loading: true, error: null });
        try {
            const response = await fetchPricingPlans()
            set({
                plans: response,
                loading: false,
            });
        } catch (error) {
            set({ error: "Failed to fetch pricing plan", loading: false });
        }
    }

}));

export default usePlanstore;

