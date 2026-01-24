import axiosIns from '../axios';
import { Supplier } from '../types';

export const GetAllSuppliers = async (): Promise<Supplier[]> => {
  try {
    const response = await axiosIns.get('/api/v1/supplier');
    return response.data.body || [];
  } catch (error: any) {
    if (error.response) {
      throw error.response.data.message || 'An error occurred while fetching suppliers';
    }
    throw 'An error occurred while fetching suppliers';
  }
};

