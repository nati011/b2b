import axiosIns from "@/lib/axios";

export interface Customer {
    id: number;
    full_name: string;
    email: string;
    phone_number?: string;
    city?: string;
    region?: string;
    woreda?: string;
    status: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export async function getCustomerByEmail(email: string): Promise<Customer | null> {
  try {
    const response = await axiosIns.get(`/customer?email=${encodeURIComponent(email)}`);
    // Handle response wrapped in body or direct (matching mobile app pattern)
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
}

export interface CreateCustomerRequest {
  full_name?: string;
  first_name?: string;
  last_name?: string;
  email: string;
  phone_number?: string;
  phone?: string;
  city?: string;
  region?: string;
  woreda?: string;
  status?: string;
}

export async function createCustomer(customerData: CreateCustomerRequest): Promise<Customer> {
  try {
    // Prepare request payload - backend accepts full_name or first_name+last_name
    const payload: any = {
      email: customerData.email,
    };

    if (customerData.full_name) {
      payload.full_name = customerData.full_name;
    } else if (customerData.first_name || customerData.last_name) {
      payload.first_name = customerData.first_name || '';
      payload.last_name = customerData.last_name || '';
    }

    if (customerData.phone_number) {
      payload.phone_number = customerData.phone_number;
    } else if (customerData.phone) {
      payload.phone = customerData.phone;
    }

    if (customerData.city) payload.city = customerData.city;
    if (customerData.region) payload.region = customerData.region;
    if (customerData.woreda) payload.woreda = customerData.woreda;
    if (customerData.status) payload.status = customerData.status;

    const response = await axiosIns.post('/customer', payload);
    
    // Handle response wrapped in body or direct
    const data = response.data?.body || response.data;
    const customer = data?.customer || data;
    
    if (customer?.id) {
      return customer;
    }
    
    throw new Error('Invalid response from customer creation');
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || 
                          error.response.data?.error || 
                          'An error occurred while creating customer';
      throw new Error(errorMessage);
    }
    throw new Error(error.message || 'An error occurred while creating customer');
  }
}

