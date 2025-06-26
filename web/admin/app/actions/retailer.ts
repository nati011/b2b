'use server'
import axiosIns from '@/app/libs/axios'
import { Retailer } from '@/app/libs/types';

export async function GetAll(page: number) {
  try {
    const response = await axiosIns.get(`/retailer?limit=10&offset=${page}`);
    console.log(response.data)
    return response.data.retailers
  } catch (error: any) {
    console.log(error)
    throw error.data.message || "An error has occured while processing your request"
  }
}
export async function GetById(id?: number) {
  try {
    const response = await axiosIns.get(`/retailer?id=${id}`);
    return response.data.retailer
  } catch (error: any) {
    throw error.message || "An error has occured while processing your request"
  }
}

export async function Create(RetailersData: Partial<Retailer>) {
  try {
    const response = await axiosIns.post('/retailer/', RetailersData);
    return response.data.detail
  } catch (error: any) {
    throw error.data.message || "An error has occured while processing your request"
  }
}