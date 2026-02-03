import { create } from 'zustand'
import { Invoice, Order } from '@/lib/types';
import { createOrder, fetchOrders, getInvoice, getOrderById, updateOrderStatus, CreateOrderRequest, OrderResponse } from '@/app/actions/orders';
import { toast } from 'sonner';

// Map OrderResponse to frontend Order format
function mapOrderResponseToOrder(orderResponse: OrderResponse): Order {
  // Parse cart_snapshot if it's a string
  let items: any[] = [];
  if (orderResponse.cart_snapshot) {
    try {
      const cartSnapshot = typeof orderResponse.cart_snapshot === 'string' 
        ? JSON.parse(orderResponse.cart_snapshot) 
        : orderResponse.cart_snapshot;
      
      if (Array.isArray(cartSnapshot)) {
        items = cartSnapshot.map((item: any) => ({
          ProductId: item.product_id || item.ProductId || 0,
          ProductName: item.product_name || item.ProductName || '',
          ProductPrice: item.price || item.ProductPrice || 0,
          Quantity: item.quantity || item.Quantity || 0
        }));
      }
    } catch (e) {
      console.warn('Failed to parse cart_snapshot:', e);
    }
  }
  
  // Parse customer_snapshot to get customer name
  let customerName = '';
  if (orderResponse.customer_snapshot) {
    try {
      const customerSnapshot = typeof orderResponse.customer_snapshot === 'string'
        ? JSON.parse(orderResponse.customer_snapshot)
        : orderResponse.customer_snapshot;
      customerName = customerSnapshot?.full_name || customerSnapshot?.name || '';
    } catch (e) {
      console.warn('Failed to parse customer_snapshot:', e);
    }
  }

  return {
    Id: orderResponse.id,
    CustomerId: orderResponse.customer_id,
    CustomerName: customerName,
    Items: items,
    Total: orderResponse.total || 0,
    Status: orderResponse.status || '',
    DeliveryStatus: orderResponse.delivery_status || '',
    PaymentStatus: orderResponse.payment_status || '',
    ConfirmationStatus: orderResponse.confirmation_status || '',
    CreatedAt: orderResponse.created_at,
    ExpiresAt: '',
    CartSnapshot: orderResponse.cart_snapshot,
    CustomerSnapshot: orderResponse.customer_snapshot
  };
}

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

    fetchOrders: (page: number, status: string, customerId?: number) => Promise<void>;
    fetchOrder: (id: number) => Promise<void>;
    fetchInvoice: (order_id: number) => Promise<void>;
    checkout: (request: CreateOrderRequest) => Promise<void>;
    cancelOrder: (id: number) => Promise<void>;
}

const useOrdersStore = create<OrdersStore>((set) => ({
    totalOrder: 0,
    orders: [],
    order: {
        Id: 0,
        CustomerId: 0,
        CustomerName: "",
        Items: [],
        Total: 0,
        Status: '',
        DeliveryStatus: '',
        PaymentStatus: '',
        CreatedAt: '',
        ConfirmationStatus: '',
        ExpiresAt: ''
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

    fetchOrders: async (page: number, status: string, customerId?: number) => {
        set({ loading: true, error: null });
        try {
            // Check if user is a supplier from localStorage
            let isSupplier = false;
            if (typeof window !== 'undefined') {
                const userRolesStr = localStorage.getItem('user_roles');
                if (userRolesStr) {
                    try {
                        const userRoles: string[] = JSON.parse(userRolesStr);
                        isSupplier = userRoles.includes('supplier');
                    } catch (e) {
                        // Ignore parse errors
                    }
                }
            }
            
            const response = await fetchOrders(10, page, status, customerId, isSupplier);
            set({
                orders: response.List,
                totalOrder: response.TotalCount,
                loading: false
            });
        } catch (error: any) {
            set({ error: error.message, loading: false });
        }
    },
    fetchOrder: async (order_id: number) => {
        set({ loading: true, error: null });
        try {
            const orderResponse = await getOrderById(order_id);
            // Map OrderResponse to frontend Order format
            const order = mapOrderResponseToOrder(orderResponse);
            set({
                order: order,
                loading: false
            });
        } catch (error: any) {
            set({ loading: false, error: error.message });
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
        } catch (error: any) {
            set({ loading: false, error: error.message });
            toast.error(error.message)
        }
    },
    checkout: async (request: CreateOrderRequest) => {
        set({ loading: true, error: null });
        try {
            await createOrder(request)
            set({ loading: false });
        } catch (error: any) {
            const er = error.message || "An error occured while checkingout"
            toast.error(er)
            set({ loading: false });
            throw error;
        }
    },
    cancelOrder: async (order_id: number) => {
        set({ loading: true, error: null });
        try {
            await updateOrderStatus(order_id.toString(), "CANCELLED");
            set({ loading: false });
            // Check if user is a supplier from localStorage
            let isSupplier = false;
            if (typeof window !== 'undefined') {
                const userRolesStr = localStorage.getItem('user_roles');
                if (userRolesStr) {
                    try {
                        const userRoles: string[] = JSON.parse(userRolesStr);
                        isSupplier = userRoles.includes('supplier');
                    } catch (e) {
                        // Ignore parse errors
                    }
                }
            }
            await fetchOrders(10, 0, "PENDING", undefined, isSupplier)
        } catch (error: any) {
            toast.error(error.message)  
            set({ loading: false });
        }
    },
}));

export default useOrdersStore;