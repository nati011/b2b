import axiosIns from '../axios';

export interface CreateCustomerRequest {
  full_name: string;
  first_name?: string;
  last_name?: string;
  email?: string;
  phone?: string;
  phone_number?: string;
  city?: string;
  region?: string;
  woreda?: string;
  status?: string;
  is_active?: boolean;
  password?: string;
}

export interface UpdateCustomerRequest {
  full_name: string;
  first_name?: string;
  last_name?: string;
  email?: string;
  phone?: string;
  phone_number?: string;
  city?: string;
  region?: string;
  woreda?: string;
  status?: string;
  is_active?: boolean;
}

export interface CustomerResponse {
  id: number;
  full_name: string;
  status: string;
  city?: string;
  region?: string;
  woreda?: string;
  phone_number?: string;
  email?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateCustomerResponse {
  message: string;
  customer: CustomerResponse;
}

export interface CustomerListResponse {
  items: CustomerResponse[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export const CreateCustomer = async (customer: CreateCustomerRequest): Promise<CreateCustomerResponse> => {
  try {
    const response = await axiosIns.post('/customer', customer);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while creating customer';
    }
    throw 'An error occurred while creating customer';
  }
};

export const GetCustomer = async (id: number): Promise<CustomerResponse> => {
  try {
    const response = await axiosIns.get(`/customer/${id}`);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while fetching customer';
    }
    throw 'An error occurred while fetching customer';
  }
};

export const GetCustomerByEmail = async (email: string): Promise<CustomerResponse | null> => {
  try {
    const response = await axiosIns.get(`/customer?email=${encodeURIComponent(email)}`);
    // Handle response wrapped in body or direct
    const data = response.data?.body || response.data;
    
    // If the API returns a list, get the first item
    if (data?.items && Array.isArray(data.items) && data.items.length > 0) {
      return data.items[0];
    }
    // If it returns a single customer directly
    if (data?.id) {
      return data;
    }
    return null;
  } catch (error: any) {
    // If customer not found, return null instead of throwing
    if (error.response?.status === 404) {
      return null;
    }
    console.error('Error fetching customer by email:', error);
    return null;
  }
};

export const UpdateCustomer = async (id: number, customer: UpdateCustomerRequest): Promise<CustomerResponse> => {
  try {
    const response = await axiosIns.put(`/customer/${id}`, customer);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while updating customer';
    }
    throw 'An error occurred while updating customer';
  }
};

export const DeleteCustomer = async (id: number): Promise<void> => {
  try {
    await axiosIns.delete(`/customer/${id}`);
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || error.response.data.error || 'An error occurred while deleting customer';
    }
    throw 'An error occurred while deleting customer';
  }
};

export const ListCustomers = async (params?: {
  page?: number;
  limit?: number;
}): Promise<CustomerListResponse> => {
  try {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const url = `/customer?${queryParams.toString()}`;
    console.log('📡 Fetching customers from:', url);
    
    const response = await axiosIns.get(url);
    
    console.log('✅ Customers response received:', {
      status: response.status,
      dataKeys: Object.keys(response.data || {}),
      itemsCount: response.data?.items?.length || 0,
    });
    
    // Ensure response has the expected structure
    if (!response.data) {
      throw new Error('Invalid response: no data received');
    }
    
    // Check if response structure is valid
    const data = response.data;
    if (!data || !Array.isArray(data.items)) {
      console.error('Invalid customer list response structure:', data);
      throw new Error('Invalid response structure: expected items array');
    }
    
    return data;
  } catch (error: any) {
    console.error('❌ Customers fetch error:', {
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
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching customers';
      throw new Error(`${errorMessage} (Status: ${error.response.status})`);
    }
    if (error.request) {
      throw new Error('Network error: Unable to reach the server. Please check if the backend is running.');
    }
    throw error;
  }
};




