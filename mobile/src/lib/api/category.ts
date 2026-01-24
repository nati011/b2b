import axiosIns from '../axios';
import { Category } from '../types';

export const GetAllCategories = async (): Promise<Category[]> => {
  try {
    const response = await axiosIns.get('/api/v1/category');
    return response.data.body || [];
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || 'An error occurred while fetching categories';
    }
    throw 'An error occurred while fetching categories';
  }
};

