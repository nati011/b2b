import { create } from "zustand";
import { SupplierRequest } from "@/lib/types";
import { RegisterSupplier } from "@/app/actions/auth";
import { toast } from "sonner";

interface SupplierStore {
    success: string | null;
    loading: boolean;
    error: string | null;
    supplier_id: number | null;

    register: (profile: SupplierRequest) => Promise<void>;
}

const useSupplierStore = create<SupplierStore>((set) => ({
    success: null,
    loading: false,
    error: null,
    supplier_id: null,

    register: async (profile: SupplierRequest) => {
        set({ loading: true, error: null, success: null, supplier_id: null });
        try {
            const supplier_id = await RegisterSupplier(profile);
            // Only set success if we got a valid supplier_id
            if (supplier_id) {
                set({ loading: false, success: "Supplier registered successfully", supplier_id: supplier_id, error: null });
            } else {
                const errMsg = "Registration succeeded but no supplier ID was returned";
                toast.error(errMsg);
                set({ loading: false, error: errMsg, success: null, supplier_id: null });
                throw new Error(errMsg);
            }
        } catch (error: any) {
            const errMsg = typeof error === "string" ? error : error?.message || "Failed to register supplier";
            toast.error(errMsg);
            set({ loading: false, error: errMsg, success: null, supplier_id: null });
            throw error; // Re-throw to allow component to handle it
        }
    },
}));

export default useSupplierStore;
