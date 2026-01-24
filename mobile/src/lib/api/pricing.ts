import axiosIns from '../axios';
import { PricingPlan } from '../types';

export const GetAllPricingPlans = async (): Promise<PricingPlan[]> => {
  try {
    const response = await axiosIns.get('/api/v1/plan');
    return response.data.body || [];
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || 'An error occurred while fetching pricing plans';
    }
    throw 'An error occurred while fetching pricing plans';
  }
};

