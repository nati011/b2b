'use server'
import axiosIns from '@/app/libs/axios';
import { DistributorRequest } from '@/app/libs/types';


export const GetAll = async (url?: string) => {
    try {
        const response = await axiosIns.get('/distributor');
        return response.data.body.distributors
    } catch (error) {
        throw error
    }
}


export const Create = async (DistributorsData: DistributorRequest) => {
    try {
        const response = await axiosIns.post('/distributor/', DistributorsData);
        return response.data.detail
    } catch (error) {
        throw error
    }
}
export const GetById = async (id: number) => {
    try {
        const response = await axiosIns.get(`/distributor?id=${id}`);
        return response.data.body.Distributor
    } catch (error) {
        throw error
    }
}
export const GetDistributorUser = async (id: number) => {
    try {
        const response = await axiosIns.get(`/user?id=${id}`);
        return response.data.body.user
    } catch (error) {
        throw error
    }
}
