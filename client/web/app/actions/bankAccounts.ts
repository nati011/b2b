import axiosIns from "@/lib/axios";

export interface BankAccountRequest {
  supplier_id: number;
  bank_name: string;
  account_number: string;
  account_holder_name: string;
  branch_name?: string;
  account_type?: string;
  is_primary: boolean;
}

export interface BankAccountResponse {
  id: number;
  supplier_id: number;
  bank_name: string;
  account_number: string;
  account_holder_name: string;
  branch_name?: string;
  account_type: string;
  is_primary: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export async function listBankAccounts(supplierId: number): Promise<BankAccountResponse[]> {
  try {
    const response = await axiosIns.get(`/supplier/bank-account?supplier_id=${supplierId}`);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while fetching bank accounts';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while fetching bank accounts');
  }
}

export async function createBankAccount(data: BankAccountRequest): Promise<BankAccountResponse> {
  try {
    const response = await axiosIns.post('/supplier/bank-account', data);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while creating bank account';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while creating bank account');
  }
}

export async function updateBankAccount(id: number, data: BankAccountRequest): Promise<BankAccountResponse> {
  try {
    const response = await axiosIns.put(`/supplier/bank-account/${id}`, data);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while updating bank account';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while updating bank account');
  }
}

export async function deleteBankAccount(id: number): Promise<void> {
  try {
    await axiosIns.delete(`/supplier/bank-account/${id}`);
  } catch (error: any) {
    if (error.response) {
      const errorMessage = error.response.data?.message || error.response.data?.error || 'An error occurred while deleting bank account';
      throw new Error(errorMessage);
    }
    throw new Error('An error occurred while deleting bank account');
  }
}

