import { create } from 'zustand'
import axiosIns from '@/lib/axios'
import { Invoice, Order } from '@/lib/types';
import { getInvoice, verifyPayment } from '@/app/actions/orders';

interface InvoiceStore {
    invoice: Invoice
    verified: boolean;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchInvoice: (order_id: number) => Promise<void>
    paymentVerification: () => Promise<void>
}

const useInvoiceStore = create<InvoiceStore>((set) => ({
    invoice: {
        Id: 0,
        Created_Date: '',
        ExternalId: '',
        Status: '',
        OrderId: 0,
        SubTotal: 0,
        LineItems: [],
        TaxAmount: 0
    },
    verified: false,
    loading: false,
    error: null,
    next: null,
    previous: null,


    fetchInvoice: async (order_id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await getInvoice(order_id)
            set({
                invoice: response,
                loading: false
            });
        } catch (error: any) {

            set({ error: error.message, loading: false });
        }
    },

    paymentVerification: async () => {
        set({ loading: true, error: null });
        try {
            const tx_ref = localStorage.getItem("tx_ref")
            const response = await verifyPayment(tx_ref || "")
            set({
                verified: response,
                loading: false
            });
        } catch (error: any) {
            set({ error: error.message, loading: false });
        } finally {
            localStorage.removeItem("tx_ref")
        }
    },
}));

export default useInvoiceStore;