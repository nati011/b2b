'use server'
import axiosIns from "@/app/libs/axios";
import { withErrorHandling } from "@/app/libs/error-handling";

export async function fetchOrders(offset: number) {
    return withErrorHandling(async () => {
        const response = await axiosIns.get(`/order?limit=10&offset=${offset}`);
        console.log(response.data)
        return response.data.body;
      });

}

export async function getOrderById(orderId: number) {
    try {
        console.log(orderId)
        const response = await axiosIns.get(`/order?id=${orderId}`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch order');
    }
}

export async function updateOrderStatus(orderId: number, status: string) {
    try {
        const response = await axiosIns.patch(`/order?id=${orderId}&command=${status}`);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw error.response.data.message || "Failed to update order status";
        }
        throw "Failed to update order status";
    }
        
}

export async function createOrder(orderData: any) {
    try {
        const response = await axiosIns.post('/order', orderData);
        return response.data;
    } catch (error: any) {
        if (error.response) {
            throw Error(error.response.message)
        }
        throw new Error('Failed to create order');
    }
}