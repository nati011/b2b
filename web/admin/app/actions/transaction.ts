'use server'
import axiosIns from "@/app/libs/axios";

export async function fetchTransactions() {
    try {
        const response = await axiosIns.get('/api/transaction');
        console.log(response.data)
        return response.data.body.transactions;
    } catch (error: any) {
        if (error.response) {
            throw Error(error.response.message)
        }
        throw new Error('Failed to create order');
    }
}
