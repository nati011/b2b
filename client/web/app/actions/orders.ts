import axiosIns from "@/lib/axios";

// Order types matching mobile app and backend
export interface OrderItemRequest {
  id?: number;
  product_id?: number;
  quantity: number;
  price?: number;
}

export interface OrderItemResponse {
  product_id: number;
  quantity: number;
  price?: number;
}

export interface CreateOrderRequest {
  customer_id: number;
  status?: string;
  payment_status?: string;
  delivery_status?: string;
  confirmation_status?: string;
  total?: number;
  referral_code?: string;
  delivery_address: string;
  customer_snapshot?: any;
  shipping_address_snapshot?: any;
  billing_address_snapshot?: any;
  items: OrderItemRequest[];
}

export interface OrderResponse {
  id: number;
  customer_id: number;
  status: string;
  payment_status?: string;
  delivery_status?: string;
  confirmation_status?: string;
  total?: number;
  referral_code?: string;
  delivery_address?: string;
  customer_snapshot?: any;
  cart_snapshot?: any;
  items?: OrderItemResponse[];
  created_at: string;
  updated_at: string;
}

export interface OrderListResponse {
  orders: OrderResponse[];
  total: number;
  limit: number;
  offset: number;
}

// Map backend order to frontend Order type (for backward compatibility)
function mapOrderToFrontendFormat(order: OrderResponse): any {
  // Parse cart_snapshot if it's a string
  let items: any[] = [];
  if (order.cart_snapshot) {
    try {
      const cartSnapshot = typeof order.cart_snapshot === 'string' 
        ? JSON.parse(order.cart_snapshot) 
        : order.cart_snapshot;
      
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
  if (order.customer_snapshot) {
    try {
      const customerSnapshot = typeof order.customer_snapshot === 'string'
        ? JSON.parse(order.customer_snapshot)
        : order.customer_snapshot;
      customerName = customerSnapshot?.full_name || customerSnapshot?.name || '';
    } catch (e) {
      console.warn('Failed to parse customer_snapshot:', e);
    }
  }

  return {
    Id: order.id,
    CustomerId: order.customer_id,
    CustomerName: customerName,
    Items: items,
    Total: order.total || 0,
    Status: order.status || '',
    DeliveryStatus: order.delivery_status || '',
    PaymentStatus: order.payment_status || '',
    ConfirmationStatus: order.confirmation_status || '',
    ReferralCode: order.referral_code || '',
    DeliveryAddress: order.delivery_address || '',
    CreatedAt: order.created_at,
    ExpiresAt: '',
    CartSnapshot: order.cart_snapshot,
    CustomerSnapshot: order.customer_snapshot
  };
}

export async function ListCustomerOrders(params: {
  customer_id: number;
  status?: string;
  limit?: number;
  offset?: number;
}): Promise<OrderListResponse> {
  try {
    const queryParams = new URLSearchParams();
    queryParams.append('customer_id', params.customer_id.toString());
    if (params.status) queryParams.append('status', params.status);
    if (params.limit) queryParams.append('limit', params.limit.toString());
    if (params.offset) queryParams.append('offset', params.offset.toString());

    const url = `/orders/customer?${queryParams.toString()}`;
    console.log('📡 Fetching orders from:', url);
    
    const response = await axiosIns.get(url);
    console.log('✅ Orders response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      hasBody: !!response.data?.body,
      ordersCount: response.data?.body?.orders?.length || response.data?.orders?.length || 0,
    });
    
    // Handle response wrapped in body or direct
    const data = response.data?.body || response.data;
    
    // Ensure we have the expected structure
    if (data && (data.orders || Array.isArray(data))) {
      const result = {
        orders: data.orders || data,
        total: data.total || (data.orders ? data.orders.length : 0),
        limit: data.limit || params.limit || 10,
        offset: data.offset || params.offset || 0,
      };
      console.log('📦 Processed orders result:', {
        ordersCount: result.orders.length,
        total: result.total,
      });
      return result;
    }
    
    console.warn('⚠️ Unexpected orders response structure:', data);
    return data;
  } catch (error: any) {
    console.error('Error fetching orders:', error);
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching orders';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while fetching orders');
  }
}

// ListSupplierOrders fetches orders for a supplier
// Uses the dedicated /order/supplier endpoint (server automatically filters by authenticated supplier)
export async function ListSupplierOrders(params?: {
  supplier_id?: number; // Not needed - server gets it from authenticated user
  status?: string;
  limit?: number;
  offset?: number;
}): Promise<OrderListResponse> {
  try {
    const queryParams = new URLSearchParams();
    // supplier_id is not needed - server automatically gets it from authenticated user
    if (params?.status && params.status !== 'ALL') {
      queryParams.append('status', params.status);
    }
    if (params?.limit) {
      queryParams.append('limit', params.limit.toString());
    }
    if (params?.offset !== undefined) {
      queryParams.append('offset', params.offset.toString());
    }

    const url = `/order/supplier?${queryParams.toString()}`;
    console.log('📡 Fetching supplier orders from:', url);
    
    const response = await axiosIns.get(url);
    console.log('✅ Supplier orders response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      ordersCount: response.data?.orders?.length || 0,
    });
    
    // Backend returns: { orders: [], total: 0, limit: 10, offset: 0 } directly
    const orderData: OrderListResponse = response.data;
    
    return orderData;
  } catch (error: any) {
    console.error('❌ Supplier orders fetch error:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status,
    });
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching supplier orders';
      throw new Error(errorMessage);
    }
    throw new Error(error.message || 'An error occurred while fetching supplier orders');
  }
}

export async function fetchOrders(limit: number, offset: number, status: string, customerId?: number, isSupplier?: boolean, supplierId?: number) {
  try {
    const queryParams = new URLSearchParams();
    
    // Use dedicated supplier endpoint if user is a supplier
    if (isSupplier) {
      // Use /order/supplier endpoint - server automatically filters by authenticated supplier
      if (status && status !== 'ALL') queryParams.append('status', status);
      queryParams.append('limit', limit.toString());
      queryParams.append('offset', (limit * offset).toString());
      
      const url = `/order/supplier?${queryParams.toString()}`;
      console.log('📡 Fetching supplier orders from:', url);
      
      const response = await axiosIns.get(url);
      console.log('✅ Supplier orders response received:', {
        status: response.status,
        dataKeys: Object.keys(response.data || {}),
        ordersCount: response.data?.orders?.length || 0,
      });
      
      const orderData: OrderListResponse = response.data;
      const mappedOrders = (orderData.orders || []).map((order: OrderResponse) => 
        mapOrderToFrontendFormat(order)
      );
      
      return {
        List: mappedOrders,
        TotalCount: orderData.total || 0,
        Limit: orderData.limit || limit,
        Offset: orderData.offset || offset
      };
    }
    
    // For non-suppliers, use customer endpoint
    if (customerId) queryParams.append('customer_id', customerId.toString());
    if (status && status !== 'ALL') queryParams.append('status', status);
    queryParams.append('limit', limit.toString());
    queryParams.append('offset', (limit * offset).toString());

    const url = `/orders/customer?${queryParams.toString()}`;
    console.log('📡 Fetching orders from:', url);
    
    const response = await axiosIns.get(url);
    console.log('✅ Orders response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      ordersCount: response.data?.orders?.length || 0,
    });
    
    // Backend returns: { orders: [], total: 0, limit: 10, offset: 0 } directly
    const orderData: OrderListResponse = response.data;
    
    // Map backend orders to frontend format for backward compatibility
    const mappedOrders = (orderData.orders || []).map((order: OrderResponse) => 
      mapOrderToFrontendFormat(order)
    );
    
    return {
      List: mappedOrders,
      TotalCount: orderData.total || 0,
      Limit: orderData.limit || limit,
      Offset: orderData.offset || offset
    };
  } catch (error: any) {
    console.error('❌ Error fetching orders:', error);
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching orders';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while fetching orders');
  }
}

export async function getOrderById(orderId: number): Promise<OrderResponse> {
  try {
    const response = await axiosIns.get(`/order?id=${orderId}`);
    // Backend returns OrderResponse directly (not wrapped in body)
    return response.data;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching order';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while fetching order');
  }
}


export async function getInvoice(orderId: number) {
  try {
    const response = await axiosIns.get(`/invoice?order_id=${orderId}`);
    return response.data.body.invoice.List[0];
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching invoice';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while fetching invoice');
  }
}


export async function updateOrderStatus(orderId: string, command: string): Promise<OrderResponse> {
  try {
    const response = await axiosIns.patch(`/order?id=${orderId}&command=${encodeURIComponent(command)}`);
    // Backend returns OrderResponse directly (not wrapped in body)
    return response.data;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while updating order status';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while updating order status');
  }
}

export async function updateOrderPaymentStatus(orderId: string, paymentStatus: string): Promise<OrderResponse> {
  try {
    const response = await axiosIns.patch(`/order?id=${orderId}&payment_status=${encodeURIComponent(paymentStatus)}`);
    // Backend returns OrderResponse directly (not wrapped in body)
    return response.data;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while updating payment status';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while updating payment status');
  }
}

// Helper to convert CheckoutRequest to CreateOrderRequest
export function convertCheckoutToCreateOrder(checkout: { customer_id: number; referral_code?: string; delivery_address?: string; items: Array<{ id: number; quantity: number }> }): CreateOrderRequest {
  return {
    customer_id: checkout.customer_id,
    referral_code: checkout.referral_code,
    delivery_address: checkout.delivery_address?.trim() || '',
    items: checkout.items.map(item => ({
      product_id: item.id,
      quantity: item.quantity
    }))
  };
}

export async function createOrder(orderData: CreateOrderRequest | { customer_id: number; items: Array<{ id: number; quantity: number }> }): Promise<OrderResponse> {
  try {
    // Convert CheckoutRequest format if needed
    let request: CreateOrderRequest;
    
    if ('items' in orderData && orderData.items.length > 0 && 'product_id' in orderData.items[0]) {
      request = orderData as CreateOrderRequest;
    } else {
      request = convertCheckoutToCreateOrder(orderData as { customer_id: number; items: Array<{ id: number; quantity: number }> });
    }
    
    console.log('📤 Creating order:', request);
    const response = await axiosIns.post('/order', request);
    
    // Handle both response formats: direct or wrapped in body
    const orderResponse = response.data?.body || response.data;
    
    // Validate response has required fields
    if (!orderResponse || !orderResponse.id) {
      console.error('Invalid order response:', orderResponse);
      throw new Error('Invalid response from server');
    }
    
    return orderResponse;
  } catch (error: any) {
    console.error('Error creating order:', error);
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || error.response.data?.body?.message || 'An error occurred while creating order';
      throw new Error(errorMessage);
    }
    if (error.message) {
      throw error;
    }
    throw new Error('An error occurred while creating order');
  }
}
