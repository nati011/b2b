// Note: Category API endpoint is not available in the backend
// Categories are referenced via category_ids in products
// This file is kept for backward compatibility but will need to be updated
// when category endpoints are implemented in the backend

import axiosIns from '../axios';

export interface Category {
  id: number;
  name: string;
}

export const GetAllCategories = async (): Promise<Category[]> => {
  // TODO: Implement when category endpoint is available in backend
  // For now, return empty array or fetch from products
  console.warn('Category API endpoint not available in backend. Returning empty array.');
  return [];
  
  // Uncomment when category endpoint is implemented:
  // try {
  //   const response = await axiosIns.get('/category');
  //   return response.data.items || [];
  // } catch (error: any) {
  //   if (error.response) {
  //     throw error.response.data.message || error.response.data.error || 'An error occurred while fetching categories';
  //   }
  //   throw 'An error occurred while fetching categories';
  // }
};

