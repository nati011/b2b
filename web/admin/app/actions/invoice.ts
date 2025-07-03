'use server'
import axiosIns from "@/app/libs/axios";

export async function fetchInvoice(order_id?: number) {
    try {
        const response = await axiosIns.get(`/invoice?order_id=${order_id}`);
        console.log(response.data)
        return response.data.body.invoice.List[0]
    } catch (error: any) {
        if (error.response) {
            throw Error(error.response.message)
        }
        throw new Error('Failed to fetch invoice');
    }
}
