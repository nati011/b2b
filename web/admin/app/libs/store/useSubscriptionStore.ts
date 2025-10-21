import {Get} from "@/app/actions/subscription"
import {Subscription} from "@/app/libs/types"
import { toast } from "sonner";
import { create } from "zustand";


interface SubscriptionStore {
    success: string | null
    subscription: Subscription | null,
    subscriptionLoading: boolean | null;
    subscriptionError: string | null;
   
    fetchSubscription: () => Promise<void>;
}

export const useSubscriptionStore = create<SubscriptionStore>((set) => ({
    success: null,
   subscription: null,
   subscriptionLoading: null,
   subscriptionError: null,
   fetchSubscription: async () => {
        set({ subscriptionLoading: true, subscriptionError: null });
        try {
            const response = await Get();
            set({
                subscription: response,
                subscriptionLoading: false,
            });
        } catch (error: any) {
            // set({
            //     subscriptionError: error.message || "Failed to fetch roles",
            //     subscriptionLoading: false,
            // });
            toast.error(error.message || "Failed to fetch roles");
        }
   }
}))

