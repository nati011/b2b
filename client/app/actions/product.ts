import axiosIns from "@/lib/axios";

export interface ProductResponse {
  id: number;
  name: string;
  description?: string;
  external_id?: string;
  attributes?: any;
  unit?: string;
  is_active: boolean;
  supplier_id: number;
  price?: number;
  total_quantity: number;
  reserved_quantity: number;
  available_quantity: number;
  category_ids?: number[];
  created_at: string;
  updated_at: string;
}

export interface ProductListResponse {
  products: ProductResponse[];
  total: number;
  limit: number;
  offset: number;
}

export const ListProducts = async (params?: {
  supplier_id?: number;
  category_id?: number;
  is_active?: boolean;
  limit?: number;
  offset?: number;
}): Promise<ProductListResponse> => {
  try {
    const queryParams = new URLSearchParams();
    if (params?.supplier_id) queryParams.append('supplier_id', params.supplier_id.toString());
    if (params?.category_id) queryParams.append('category_id', params.category_id.toString());
    if (params?.is_active !== undefined) queryParams.append('is_active', params.is_active.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const url = `/products?${queryParams.toString()}`;
    console.log('📡 Fetching products from:', url);
    
    const response = await axiosIns.get(url);
    
    console.log('✅ Products response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      productsCount: response.data?.products?.length || 0,
    });
    
    return response.data;
  } catch (error: any) {
    console.error('❌ Products fetch error:', {
      message: error.message,
      code: error.code,
      response: error.response ? {
        status: error.response.status,
        statusText: error.response.statusText,
        data: error.response.data,
      } : null,
      request: error.request ? {
        url: error.config?.url,
        baseURL: error.config?.baseURL,
      } : null,
    });
    
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching products';
      throw new Error(`${errorMessage} (Status: ${error.response.status})`);
    }
    if (error.request) {
      throw new Error(`Network error: Unable to reach the server. Please check if the backend is running.`);
    }
    throw new Error(error.message || 'An error occurred while fetching products');
  }
};


