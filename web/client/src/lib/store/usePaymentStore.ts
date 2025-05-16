import { create } from "zustand";
import axiosIns from "@/lib/axios";
import { Partner, CheckoutRequest } from "@/lib/types";


interface ProductsStore {
    success: string | null;
    partners: Partner[];
    loading: boolean;
    error: string | null;
    isLoggedIn: boolean

    fetchPaymentPartners: () => Promise<void>;
    checkout: (request: CheckoutRequest) => Promise<void>;
}

const usePartnerStore = create<ProductsStore>((set) => ({
    success: null,
    isLoggedIn: true,
    partners: [],
    loading: false,
    error: null,

    fetchPaymentPartners: async () => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get("/api/v1/payment_option");

            set({
                partners: response.data.body.payment_options,
                loading: false,
            });
        } catch (error) {
            set({ error: "Failed to fetch payment partner", loading: false });
        }
    },
    checkout: async (request: CheckoutRequest) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post("/api/v1/order", request);
            if (response.status == 202) {
                console.log(response.data)
                window.location.href = response.data.body.order.checkout_url
            }
        } catch (error) {
            set({ error: "Failed to checkout", loading: false });
        }
    }

}));

export default usePartnerStore;
