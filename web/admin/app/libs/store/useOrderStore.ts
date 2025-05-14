import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Order } from '@/app/libs/types';

interface OrdersStore {
    orders: Order[];
    transactions: []
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;

    fetchOrders: (url?: string) => Promise<void>;
    createOrders: (OrdersData: Partial<Order>) => Promise<void>;
    fetchTransactions: (url?: string) => Promise<void>
}

const useOrdersStore = create<OrdersStore>((set) => ({
    orders: [],
    transactions: [],
    loading: false,
    error: null,
    next: null,
    previous: null,

    fetchOrders: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.get('/api/order');
            set({
                orders: response.data.body.Orders,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch order', loading: false });
        }
    },

    createOrders: async (OrdersData: Partial<Order>) => {
        set({ loading: true, error: null });
        try {
            const response = await axiosIns.post('/api/order/', OrdersData);
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
            const response = await axiosIns.get('/api/order');
            set({
                orders: response.data.body.Orders,
                loading: false
            });
        } catch (error) {
            set({ error: 'Failed to fetch order', loading: false });
        }
    },
}));

export default useOrdersStore;