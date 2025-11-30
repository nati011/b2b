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

// Cache to prevent unnecessary API calls
let plansCache: PricingPlan[] | null = null;
let cacheTimestamp: number = 0;
const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

const usePlanstore = create<PricingPlanStore>((set) => ({
    success: null,
    isLoggedIn: true,
    plans: [],
    loading: false,
    error: null,

    fetchPricingPlan: async () => {
        // Return cached data if available and not expired
        const now = Date.now();
        if (plansCache && (now - cacheTimestamp) < CACHE_DURATION) {
            set({ plans: plansCache, loading: false });
            return;
        }

        set({ loading: true, error: null });
        try {
            const response = await fetchPricingPlans();
            plansCache = response;
            cacheTimestamp = now;
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

