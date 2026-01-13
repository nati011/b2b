"use client";
import axiosIns from "@/lib/axios";

// Debug function - can be called from browser console: globalThis.testPricingAPI()
if (globalThis !== undefined && globalThis.window !== undefined) {
  (globalThis as any).testPricingAPI = async () => {
    console.log('🧪 Testing pricing API directly...');
    console.log('Base URL:', axiosIns.defaults.baseURL);
    console.log('Full URL:', axiosIns.defaults.baseURL + '/api/v1/plan');
    
    try {
      const response = await fetch(axiosIns.defaults.baseURL + '/api/v1/plan', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      console.log('✅ Fetch response status:', response.status);
      console.log('✅ Fetch response headers:', Object.fromEntries(response.headers.entries()));
      const text = await response.text();
      console.log('✅ Fetch response text:', text);
      try {
        const json = JSON.parse(text);
        console.log('✅ Fetch response JSON:', json);
      } catch (e) {
        console.warn('⚠️ Response is not valid JSON');
      }
    } catch (error) {
      console.error('❌ Fetch error:', error);
    }
  };
}

export async function fetchPricingPlans() {
  try {
    const apiUrl = axiosIns.defaults.baseURL + '/api/v1/plan';
    console.log('Making API call to fetch pricing plans from:', apiUrl);
    // Use a longer timeout for pricing plans (15 seconds)
    const response = await axiosIns.get("/api/v1/plan", {
      timeout: 15000,
    });
    console.log('✅ API call successful!');
    console.log('API response status:', response.status);
    console.log('API response headers:', response.headers);
    console.log('API response data (full):', response.data);
    console.log('API response data (stringified):', JSON.stringify(response.data, null, 2));
    
    // Handle the response structure: {body: {plans: {list: [...]}}}
    // This is the expected structure from OperationSuccessResponse
    if (response.data?.body?.plans?.list && Array.isArray(response.data.body.plans.list)) {
      console.log('Fetched pricing plans from backend (body.plans.list):', response.data.body.plans.list);
      return response.data.body.plans.list;
    }
    
    // Alternative: {body: {plans: [...]}} (if plans is directly an array)
    if (response.data?.body?.plans && Array.isArray(response.data.body.plans)) {
      console.log('Fetched pricing plans from backend (body.plans):', response.data.body.plans);
      return response.data.body.plans;
    }
    
    // Alternative response structure: {plans: {list: [...]}}
    if (response.data?.plans?.list && Array.isArray(response.data.plans.list)) {
      console.log('Fetched pricing plans from backend (plans.list):', response.data.plans.list);
      return response.data.plans.list;
    }
    
    // Alternative response structure: {plans: [...]}
    if (response.data?.plans && Array.isArray(response.data.plans)) {
      console.log('Fetched pricing plans from backend (plans):', response.data.plans);
      return response.data.plans;
    }
    
    // Fallback: if response is directly an array
    if (response.data && Array.isArray(response.data)) {
      console.log('Fetched pricing plans from backend (direct array):', response.data);
      return response.data;
    }
    
    // If response has a different structure, log it and try to extract data
    console.warn('Unexpected response structure for pricing plans. Full response:', JSON.stringify(response.data, null, 2));
    console.warn('Response structure keys:', Object.keys(response.data || {}));
    
    // Try to find any array in the response that might be the plans
    const findArrayInObject = (obj: any, depth = 0): any[] | null => {
      if (depth > 5) return null; // Prevent infinite recursion
      if (Array.isArray(obj)) {
        // Check if this looks like a plans array (has objects with id, name, price, etc.)
        if (obj.length > 0 && obj[0] && typeof obj[0] === 'object') {
          const firstItem = obj[0];
          if (firstItem.id !== undefined && (firstItem.name !== undefined || firstItem.price !== undefined)) {
            return obj;
          }
        }
      }
      if (obj && typeof obj === 'object') {
        for (const key in obj) {
          const result = findArrayInObject(obj[key], depth + 1);
          if (result) return result;
        }
      }
      return null;
    };
    
    const foundArray = findArrayInObject(response.data);
    if (foundArray) {
      console.log('Found plans array in unexpected location:', foundArray);
      return foundArray;
    }
    
    if (response.data?.body) {
      console.warn('Body keys:', Object.keys(response.data.body));
      if (response.data.body.plans) {
        console.warn('Plans type:', typeof response.data.body.plans, 'Is array:', Array.isArray(response.data.body.plans));
        console.warn('Plans value:', response.data.body.plans);
        // If plans exists but isn't in the expected format, try to extract it
        if (response.data.body.plans && typeof response.data.body.plans === 'object') {
          // Maybe plans is an object with a list property we missed
          if (response.data.body.plans.list && Array.isArray(response.data.body.plans.list)) {
            console.log('Found plans.list in body.plans:', response.data.body.plans.list);
            return response.data.body.plans.list;
          }
        }
      }
    }
    
    // Last resort: return empty array (don't throw error for successful responses)
    console.error('Could not parse pricing plans from response. Returning empty array.');
    return [];
  } catch (error: any) {
    // Log the full error for debugging
    console.error('Failed to fetch pricing plans from backend:', error);
    
    // Handle different error types
    if (error.code === 'ECONNABORTED' || error.message?.includes('timeout')) {
      const timeoutError = new Error('Request timed out. Please check your connection and try again.');
      timeoutError.name = 'TimeoutError';
      throw timeoutError;
    }
    
    if (error.code === 'ERR_NETWORK' || error.message === 'Network Error') {
      const networkError = new Error('Unable to fetch distributor pricing. Please check your internet connection and try again.');
      networkError.name = 'NetworkError';
      throw networkError;
    }
    
    if (error.response) {
      console.error('Error response:', error.response.data);
      console.error('Error status:', error.response.status);
      const statusError = new Error(
        error.response.data?.message || 
        `Failed to fetch pricing plans (${error.response.status}). Please try again later.`
      );
      statusError.name = 'APIError';
      throw statusError;
    }
    
    // Generic error
    const genericError = new Error(error?.message || 'Unable to fetch distributor pricing. Please try again later.');
    genericError.name = 'FetchError';
    throw genericError;
  }
}
