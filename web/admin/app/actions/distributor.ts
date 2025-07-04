'use server'
import axiosIns from '@/app/libs/axios';
import { DistributorRequest } from '@/app/libs/types';


export const GetAll = async (status?: string) => {
    try {
        let requestUrl = '/distributor'
        if (status != "ALL"){
            requestUrl = `/distributor?verdict=${status}`
        }
        const response = await axiosIns.get(requestUrl);
        console.log(response.data)
        return response.data.body
    } catch (error) {
        throw error
    }
}


export const Create = async (DistributorsData: DistributorRequest) => {
    try {
        const response = await axiosIns.post('/distributor/', DistributorsData);
        console.log(response)
        return response.data.detail
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message
        }


    }
}
export const GetById = async (id: number) => {
    try {
        const response = await axiosIns.get(`/distributor?id=${id}`);
        console.log(response.data, "HERE___________________")
        return response.data.body.distributor
    } catch (error) {
        throw error
    }
}
export const GetDistributorUser = async () => {
    try {
        const response = await axiosIns.get(`/distributor/user`);
        console.log(response.data)
        return response.data.body.users
    } catch (error) {
        console.log(error)
        throw error
    }
}


export const DistributorOnBoardingReview = async (id: number, command: string, comment?: string) => {
    try {
        let url = `/distributor/${id}/onboarding_review?command=${command}`;
        if (comment && command === 'reject') {
            url += `&comment=${encodeURIComponent(comment)}`;
        }
        const response = await axiosIns.patch(url);
        console.log(response.data)
        return response.data.body.detail
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error has occured while processing your request"
        }
        throw "An error has occured while processing your request"
    }
}



export const UpdateDistributorStatus = async (id: number, command: string) => {
    try {
        const response = await axiosIns.patch(`/distributor/${id}/status?command=${command}`);
        console.log(response.data)
        return response.data.body.detail
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "An error has occured while processing your request"
        }
        throw "An error has occured while processing your request"
    }
}

