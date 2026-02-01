import axiosIns from '../axios';
import { Product } from '../types';

export interface CatalogueResponse {
  products: Product[];
  total: number;
  limit: number;
  offset: number;
}

export const GetAllCatalogues = async (params?: {
  supplier_id?: number;
  category_id?: number;
  is_active?: boolean;
  limit?: number;
  offset?: number;
}): Promise<CatalogueResponse> => {
  try {
    const queryParams = new URLSearchParams();
    if (params?.supplier_id) queryParams.append('supplier_id', params.supplier_id.toString());
    if (params?.category_id) queryParams.append('category_id', params.category_id.toString());
    if (params?.is_active !== undefined) queryParams.append('is_active', params.is_active.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const url = `/catalogue?${queryParams.toString()}`;
    console.log('📡 Fetching catalogue from:', url);
    
    const response = await axiosIns.get(url);
    
    console.log('✅ Catalogue response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      productsCount: response.data?.products?.length || 0,
    });
    
    return response.data;
  } catch (error: any) {
    console.error('❌ Catalogue fetch error:', {
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
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching catalogue';
      throw new Error(`${errorMessage} (Status: ${error.response.status})`);
    }
    if (error.request) {
      throw new Error(`Network error: Unable to reach the server. Please check if the backend is running.`);
    }
    throw new Error(error.message || 'An error occurred while fetching catalogue');
  }
};

