'use server'
import axiosIns from "@/lib/axios";
import { CheckoutRequest } from "@/lib/types";

export async function fetchOrders(limit: number, offset: number, status: string) {
    try {
        console.log(status, "STatus")
        let req_url = `/orders/retailer?limit=${limit}&offset=${limit * offset}`
        if (status != "ALL") {
            req_url = `/orders/retailer?limit=${limit}&offset=${limit * offset}&status=${status}`
        }
        const response = await axiosIns.get(req_url);
        console.log(response.data)
        return response.data.body.orders;
    } catch (error) {
        console.log(error)
        throw new Error('Failed to fetch orders');
    }
}

export async function getOrderById(orderId: number) {
    try {
        const response = await axiosIns.get(`/order?id=${orderId}`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch order');
    }
}


export async function getInvoice(orderId: number) {
    try {
        const response = await axiosIns.get(`/invoice?order_id=${orderId}`);
        return response.data.body.invoice.List[0];
    } catch (error) {
        throw new Error('Failed to fetch invoice');
    }
}


export async function updateOrderStatus(orderId: string, status: string) {
    try {
        const response = await axiosIns.patch(`/order?id=${orderId}&command=${status}`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to update order status');
    }
}

export async function createOrder(orderData: CheckoutRequest) {
    try {
        console.log(orderData)
        const response = await axiosIns.post('/order', orderData);
        return response.data.body
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message
        }
    }
}


export async function completePayment(orderId: number) {
    try {
        const response = await axiosIns.post(`/order/init_payment?id=${orderId}`);
        return response.data.body
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message
        }
    }
}

export async function verifyPayment(tx_ref: string) {
    try {
        console.log(tx_ref, "Tx Ref")
        const response = await axiosIns.post(`/payment/verify?tx_ref=${tx_ref}`);
        console.log(response.data, "Response")
        return response.data.body
    } catch (error: any) {
        console.log(error)
        if (error.response) {
            throw error.response.data.message
        }
    }
}

