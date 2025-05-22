import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Invoice, Order, Transaction } from '@/app/libs/types';

interface OrdersStore {
    orders: Order[];
    order: Order,
    transactions: Transaction[]
    invoice: Invoice
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchOrders: (url?: string) => Promise<void>;
    fetchOrder: (id: number) => Promise<void>;
    createOrders: (OrdersData: Partial<Order>) => Promise<void>;
    fetchTransactions: (url?: string) => Promise<void>
    fetchInvoice: (order_id?: number) => Promise<void>
}

const useOrdersStore = create<OrdersStore>((set) => ({
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
    transactions: [],
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

    fetchOrders: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get('/api/v1/order');
            set({
                orders: response.data.body.List,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch order', loading: false });
        }
    },
    fetchOrder: async (order_id?: number) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get(`/api/v1/order?id=${order_id}`);
            set({
                order: response.data.body.order,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch transactions', loading: false });
        }
    },
    createOrders: async (OrdersData: Partial<Order>) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post('/api/v1/order/', OrdersData);
            set(state => ({
                Orders: [...state.orders, response.data.detail],
                loading: false
            }));
        } catch (error) {

            set({ error: 'Failed to create order', loading: false });
        }
    },
    fetchTransactions: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get('/api/transaction');
            set({
                transactions: response.data.body.transactions,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch transactions', loading: false });
        }
    },
    fetchInvoice: async (order_id?: number) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get(`/api/invoice?order_id=${order_id}`);
            set({
                invoice: response.data.body.invoice.List[0],
                loading: false
            });
            console.log(response.data)
        } catch (error) {
            set({ error: 'Failed to fetch transactions', loading: false });
        }
    },
}));

export default useOrdersStore;