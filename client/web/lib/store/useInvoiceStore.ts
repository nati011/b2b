import { create } from 'zustand'
import { Invoice } from '@/lib/types';
import { getInvoice } from '@/app/actions/orders';

interface InvoiceStore {
    invoice: Invoice
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchInvoice: (order_id: number) => Promise<void>
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
}));

export default useInvoiceStore;