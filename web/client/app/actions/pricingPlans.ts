"use server";
import axios from "axios";

export async function fetchPricingPlans() {
  try {
    // For server-side requests, use internal Docker network URL
    // Fallback to external URL if internal doesn't work
    const baseURL = process.env.NEXT_PUBLIC_BASE_URL || 'http://backend:8080';
    
    const axiosInstance = axios.create({
      baseURL: baseURL,
      timeout: 10000,
    });

    const response = await axiosInstance.get("/api/v1/plan");
    
    // Handle the response structure: {body: {plans: {list: [...]}}}
    if (response.data?.body?.plans?.list) {
      return response.data.body.plans.list;
    }
    
    // Alternative response structure: {plans: {list: [...]}}
    if (response.data?.plans?.list) {
      return response.data.plans.list;
    }
    
    // Alternative response structure: {plans: [...]}
    if (response.data?.plans && Array.isArray(response.data.plans)) {
      return response.data.plans;
    }
    
    // Fallback: if response is directly an array
    if (response.data && Array.isArray(response.data)) {
      return response.data;
    }
    
    // If response has a different structure, return empty array
    return [];
  } catch (error: any) {
    // If it's a network error, try with external URL
    if (error.code === 'ECONNREFUSED' || error.code === 'ENOTFOUND') {
      try {
        const externalURL = 'http://localhost:8082';
        const axiosInstance = axios.create({
          baseURL: externalURL,
          timeout: 10000,
        });
        const response = await axiosInstance.get("/api/v1/plan");
        
        if (response.data?.body?.plans?.list) {
          return response.data.body.plans.list;
        }
        if (response.data?.plans?.list) {
          return response.data.plans.list;
        }
      } catch (retryError) {
        // Retry failed, continue to return empty array
      }
    }
    
    // For other errors, return empty array to prevent UI crash
    return [];
  }
}
