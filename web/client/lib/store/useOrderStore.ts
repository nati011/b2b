import { create } from 'zustand'
import axiosIns from '@/lib/axios'
import { CheckoutRequest, Invoice, Order } from '@/lib/types';
import { completePayment, createOrder, fetchOrders, getInvoice, getOrderById } from '@/app/actions/orders';

interface OrdersStore {
    totalOrder: number;
    orders: Order[];
    order: Order,
    invoice: Invoice
    invoiceloading: boolean;
    invoiceError: string | null;
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchOrders: (page: number, status: string) => Promise<void>;
    fetchOrder: (id: number) => Promise<void>;
    fetchInvoice: (order_id: number) => Promise<void>;
    checkout: (request: CheckoutRequest) => Promise<void>;
    completePayment: (id: number) => Promise<void>;
}

const useOrdersStore = create<OrdersStore>((set) => ({
    totalOrder: 0,
    orders: [],
    order: {
        Id: 0,
        RetailerId: 0,
        RetailerName: "",
        Items: [],
        Total: 0,
        Status: '',
        DeliveryStatus: '',
        PaymentStatus: '',
        CreatedAt: ''
    },

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
    invoiceloading: false,
    invoiceError: null,
    loading: false,
    error: null,
    next: null,
    previous: null,

    fetchOrders: async (page: number, status: string) => {
        set({ loading: true, error: null });
        try {
            const response = await fetchOrders(10, page, status);
            set({
                orders: response.List,
                totalOrder: response.TotalCount,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch order', loading: false });
        }
    },
    fetchOrder: async (order_id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await getOrderById(order_id);
            console.log(response)
            set({
                order: response.body.order,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch transactions', loading: false });
        }
    },

    fetchInvoice: async (order_id: number) => {
        set({ invoiceloading: true, error: null });
        try {
            const response = await getInvoice(order_id)
            set({
                invoice: response,
                invoiceloading: false
            });
            console.log(response.data)
        } catch (error: any) {
            set({ invoiceError: error.message, loading: false });
        }
    },
    checkout: async (request: CheckoutRequest) => {
        set({ loading: true, error: null });
        try {
            const response = await createOrder(request)
            console.log(response)
            localStorage.setItem("tx_ref", response.tx_ref)
            if (response.checkout_url) {
                window.location.href = response.checkout_url
            }
        } catch (error: any) {
            set({ error: error.message || "An error occured while checkingout", loading: false });
        }
    },
    completePayment: async (order_id: number) => {
        set({ loading: true, error: null });
        try {
            const response = await completePayment(order_id)
            console.log(response)
            localStorage.setItem("tx_ref", response)
            if (response) {
                window.location.href = response.checkout_url
            }
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    },
}));

export default useOrdersStore;