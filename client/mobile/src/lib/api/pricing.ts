// Note: Pricing plan API endpoint is not available in the backend
// This file is kept for backward compatibility but will need to be updated
// when pricing plan endpoints are implemented in the backend

import axiosIns from '../axios';

export interface PricingPlan {
  id: number;
  name: string;
  price: number;
  term_in_month: number;
  desc: string;
}

export const GetAllPricingPlans = async (): Promise<PricingPlan[]> => {
  // TODO: Implement when pricing plan endpoint is available in backend
  // For now, return empty array
  console.warn('Pricing plan API endpoint not available in backend. Returning empty array.');
  return [];
  
  // Uncomment when pricing plan endpoint is implemented:
  // try {
  //   const response = await axiosIns.get('/plan');
  //   return response.data.items || [];
  // } catch (error: any) {
  //   if (error.response) {
  //     throw error.response.data.message || error.response.data.error || 'An error occurred while fetching pricing plans';
  //   }
  //   throw 'An error occurred while fetching pricing plans';
  // }
};

