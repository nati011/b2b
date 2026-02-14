import axiosIns from '../axios';
import { Product } from '../types';

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

export const GetProduct = async (id: number): Promise<ProductResponse> => {
  try {
    const url = `/product?id=${id}`;
    console.log('📡 Fetching product from:', url);
    
    const response = await axiosIns.get(url);
    
    console.log('✅ Product response received:', {
      status: response.status,
      data: response.data,
    });
    
    return response.data;
  } catch (error: any) {
    console.error('❌ Product fetch error:', {
      id,
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
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching product';
      throw new Error(`${errorMessage} (Status: ${error.response.status})`);
    }
    if (error.request) {
      throw new Error(`Network error: Unable to reach the server. Please check if the backend is running.`);
    }
    throw new Error(error.message || 'An error occurred while fetching product');
  }
};

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

    const response = await axiosIns.get(`/products?${queryParams.toString()}`);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while fetching products';
    }
    throw 'An error occurred while fetching products';
  }
};

export const CreateProduct = async (product: {
  name: string;
  description?: string;
  external_id?: string;
  attributes?: any;
  unit?: string;
  is_active: boolean;
  supplier_id: number;
  price?: number;
  total_quantity?: number;
  reserved_quantity?: number;
  category_ids?: number[];
}): Promise<ProductResponse> => {
  try {
    const response = await axiosIns.post('/product', product);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while creating product';
    }
    throw 'An error occurred while creating product';
  }
};

export const UpdateProduct = async (id: number, product: {
  name: string;
  description?: string;
  external_id?: string;
  attributes?: any;
  unit?: string;
  is_active: boolean;
  price?: number;
  total_quantity?: number;
  reserved_quantity?: number;
  category_ids?: number[];
}): Promise<ProductResponse> => {
  try {
    const response = await axiosIns.put(`/product?id=${id}`, product);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while updating product';
    }
    throw 'An error occurred while updating product';
  }
};

export const DeleteProduct = async (id: number): Promise<void> => {
  try {
    await axiosIns.delete(`/product?id=${id}`);
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while deleting product';
    }
    throw 'An error occurred while deleting product';
  }
};

