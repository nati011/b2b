import axiosIns from '../axios';

export interface SupplierResponse {
  id: number;
  business_name: string;
  status: string;
  support_email?: string;
  support_phone?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateSupplierRequest {
  business_name: string;
  status?: string;
  support_email?: string;
  support_phone?: string;
  is_active?: boolean;
}

export interface UpdateSupplierRequest {
  business_name: string;
  status?: string;
  support_email?: string;
  support_phone?: string;
  is_active?: boolean;
}

export interface CreateSupplierResponse {
  message: string;
  supplier: SupplierResponse;
}

export interface SupplierListResponse {
  items: SupplierResponse[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export const GetAllSuppliers = async (params?: {
  page?: number;
  limit?: number;
}): Promise<SupplierListResponse> => {
  try {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const url = `/supplier?${queryParams.toString()}`;
    console.log('📡 Fetching suppliers from:', url);
    
    const response = await axiosIns.get(url);
    
    console.log('✅ Suppliers response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      itemsCount: response.data?.items?.length || 0,
    });
    
    return response.data;
  } catch (error: any) {
    console.error('❌ Suppliers fetch error:', {
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
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching suppliers';
      throw new Error(`${errorMessage} (Status: ${error.response.status})`);
    }
    if (error.request) {
      throw new Error(`Network error: Unable to reach the server. Please check if the backend is running.`);
    }
    throw new Error(error.message || 'An error occurred while fetching suppliers');
  }
};

export const GetSupplier = async (id: number): Promise<SupplierResponse> => {
  try {
    // Try path parameter first
    let url = `/supplier/${id}`;
    console.log('📡 Fetching supplier from:', url);
    
    let response;
    try {
      response = await axiosIns.get(url);
    } catch (pathError: any) {
      // If path parameter fails, try query parameter as fallback
      if (pathError.response?.status === 404 || pathError.response?.status === 400) {
        console.log('⚠️ Path parameter failed, trying query parameter format');
        url = `/supplier?id=${id}`;
        response = await axiosIns.get(url);
      } else {
        throw pathError;
      }
    }
    
    console.log('✅ Supplier response received:', {
      status: response.status,
      data: response.data,
    });
    
    // Handle both direct response and wrapped response
    const supplierData = response.data?.supplier || response.data;
    
    if (!supplierData || !supplierData.id) {
      throw new Error('Invalid supplier data received from server');
    }
    
    return supplierData;
  } catch (error: any) {
    console.error('❌ Supplier fetch error:', {
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
      const status = error.response.status;
      if (status === 404) {
        throw new Error('Supplier not found');
      }
      const errorMessage = error.response.data?.message || error.response.data?.error || `An error occurred while fetching supplier (Status: ${status})`;
      throw new Error(errorMessage);
    }
    if (error.request) {
      throw new Error(`Network error: Unable to reach the server. Please check if the backend is running.`);
    }
    throw new Error(error.message || 'An error occurred while fetching supplier');
  }
};

export const CreateSupplier = async (supplier: CreateSupplierRequest): Promise<CreateSupplierResponse> => {
  try {
    const response = await axiosIns.post('/supplier', supplier);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while creating supplier';
    }
    throw 'An error occurred while creating supplier';
  }
};

export const UpdateSupplier = async (id: number, supplier: UpdateSupplierRequest): Promise<SupplierResponse> => {
  try {
    const response = await axiosIns.put(`/supplier/${id}`, supplier);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while updating supplier';
    }
    throw 'An error occurred while updating supplier';
  }
};

export const DeleteSupplier = async (id: number): Promise<void> => {
  try {
    await axiosIns.delete(`/supplier/${id}`);
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while deleting supplier';
    }
    throw 'An error occurred while deleting supplier';
  }
};

