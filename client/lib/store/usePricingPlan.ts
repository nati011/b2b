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
            console.log('Using cached pricing plans, cache age:', Math.round((now - cacheTimestamp) / 1000), 'seconds');
            set({ plans: plansCache, loading: false });
            return;
        }

        console.log('Fetching pricing plans from backend API...');
        set({ loading: true, error: null });
        try {
            const response = await fetchPricingPlans();
            console.log('Pricing plans fetched successfully from backend:', response);
            console.log('Response type:', typeof response, 'Is array:', Array.isArray(response), 'Length:', response?.length);
            
            if (response && Array.isArray(response) && response.length > 0) {
                plansCache = response;
                cacheTimestamp = now;
                set({
                    plans: response,
                    loading: false,
                    error: null,
                });
            } else if (response && Array.isArray(response) && response.length === 0) {
                // Empty array - no plans available, but this is a valid response
                console.warn('Backend returned empty plans array - no pricing plans available');
                set({ 
                    plans: [], 
                    loading: false, 
                    error: "No pricing plans available" 
                });
            } else {
                // Unexpected response format
                console.error('Unexpected response format from fetchPricingPlans:', response);
                set({ 
                    plans: [], 
                    loading: false, 
                    error: "Unable to parse pricing plans from server response" 
                });
            }
        } catch (error: any) {
            console.error('Error fetching pricing plans from backend:', error);
            console.error('Error details:', {
                message: error?.message,
                code: error?.code,
                response: error?.response?.data,
                status: error?.response?.status
            });
            // Use the error message from the thrown error, which should be user-friendly
            const errorMessage = error?.message || "Unable to fetch supplier pricing. Please try again later.";
            set({ 
                error: errorMessage, 
                loading: false,
                plans: []
            });
        }
    }

}));

export default usePlanstore;

