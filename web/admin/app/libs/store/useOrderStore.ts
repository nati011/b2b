import { create } from 'zustand'
import axiosIns from '@/app/libs/axios'
import { Invoice, Order, Transaction } from '@/app/libs/types';
import { createOrder, fetchOrders, getOrderById,updateOrderStatus } from '@/app/actions/order';
import { fetchTransactions } from '@/app/actions/transaction';
import { fetchInvoice } from '@/app/actions/invoice';
import { toast } from 'sonner';

interface OrdersStore {
    orders: Order[];
    order: Order,
    transactions: Transaction[]
    invoice: Invoice
    loading: boolean;
    error: string | null;
    next: string | null;
    previous: string | null;
    total: number | null;

    fetchOrders: (page: number) => Promise<void>;
    fetchOrder: (id: number) => Promise<void>;
    createOrders: (OrdersData: Partial<Order>) => Promise<void>;
    fetchTransactions: (url?: string) => Promise<void>
    fetchInvoice: (order_id?: number) => Promise<void>
    updateOrderStatus: (status: string, id: number) => Promise<void>
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
        ConfirmationStatus:'',
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
    total: 0,

    fetchOrders: async (page: number) => {
        set({ loading: true, error: null });
        try {
            const response = await fetchOrders(page);
            set({
                orders: response.List ? response.List : [],
                total: response.TotalCount,
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
    createOrders: async (OrdersData: Partial<Order>) => {
        set({ loading: true, error: null });
        try {
            const response = await createOrder(OrdersData)
            set({
                loading: false
            });
            await useOrdersStore.getState().fetchOrders(0)
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    },
    fetchTransactions: async (url?: string) => {
        set({ loading: true, error: null });
        try {
            const response = await fetchTransactions()
            set({
                transactions: response.data.body.transactions,
                loading: false
            });
        } catch (error: any) {
            set({ error: 'Failed to fetch transactions', loading: false });
            toast.error(error.message || "Failed to fetch transactions");

        }
    },
    fetchInvoice: async (order_id?: number) => {
        set({ loading: true, error: null });
        try {
            const response = await fetchInvoice(order_id)
            set({
                invoice: response,
                loading: false
            });
            console.log(response.data)
        } catch (error: any) {
            set({ error: error.message, loading: false });
            toast.error(error.message || "Failed to fetch invoice");

        }
    },
    updateOrderStatus: async(status: string, order_id: number) =>{
        set({ loading: true, error: null });
        try {
            const response = await updateOrderStatus(order_id, status)
            
            await useOrdersStore.getState().fetchOrder(order_id)
        } catch (error: any) {
            set({ error: error.message, loading: false });
            toast.error(error.message || "Failed to update order status");

        }
    }
}));

export default useOrdersStore;