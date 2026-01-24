import axiosIns from '../axios';
import { Catalogue } from '../types';

export const GetAllCatalogues = async (): Promise<Catalogue[]> => {
  try {
    const response = await axiosIns.get('/api/v1/catalogue');
    return response.data.body.products || [];
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || 'An error occurred while fetching products';
    }
    throw 'An error occurred while fetching products';
  }
};

