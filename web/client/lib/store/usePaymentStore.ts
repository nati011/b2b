import { create } from "zustand";
import axiosIns from "@/lib/axios";
import { Partner } from "@/lib/types";
import { fetchPaymentPartners } from "@/app/actions/paymentpartner";


interface PartnerStore {
    success: string | null;
    partners: Partner[];
    loading: boolean;
    error: string | null;


    fetchPaymentPartners: () => Promise<void>;
}

const usePartnerStore = create<PartnerStore>((set) => ({
    success: null,
    isLoggedIn: true,
    partners: [],
    loading: false,
    error: null,

    fetchPaymentPartners: async () => {
        set({ loading: true, error: null });
        try {
            const response = await fetchPaymentPartners()

            set({
                partners: response,
                loading: false,
            });
        } catch (error) {
            set({ error: "Failed to fetch payment partner", loading: false });
        }
    }

}));

export default usePartnerStore;

