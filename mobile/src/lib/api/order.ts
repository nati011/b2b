import axiosIns from '../axios';

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
  customer_snapshot?: any;
  cart_snapshot?: any;
  shipping_address_snapshot?: any;
  billing_address_snapshot?: any;
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

export const CreateOrder = async (order: CreateOrderRequest): Promise<OrderResponse> => {
  try {
    const response = await axiosIns.post('/order', order);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while creating order';
    }
    throw 'An error occurred while creating order';
  }
};

export const GetOrder = async (id: number): Promise<OrderResponse> => {
  try {
    const response = await axiosIns.get(`/order?id=${id}`);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while fetching order';
    }
    throw 'An error occurred while fetching order';
  }
};

export const UpdateOrderStatus = async (id: number, command: string): Promise<OrderResponse> => {
  try {
    const response = await axiosIns.patch(`/order?id=${id}&command=${encodeURIComponent(command)}`);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while updating order status';
    }
    throw 'An error occurred while updating order status';
  }
};

export const ListCustomerOrders = async (params: {
  customer_id: number;
  status?: string;
  limit?: number;
  offset?: number;
}): Promise<OrderListResponse> => {
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
};




