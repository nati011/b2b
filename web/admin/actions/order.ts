'use server'
import axiosIns from "@/app/libs/axios";

export async function fetchOrders() {
    try {
        const response = await axiosIns.get('/order');
        console.log(response.data)
        return response.data.body.List;
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

export async function updateOrderStatus(orderId: string, status: string) {
    try {
        const response = await axiosIns.put(`/order/${orderId}`, { status });
        return response.data;
    } catch (error) {
        throw new Error('Failed to update order status');
    }
}

export async function createOrder(orderData: any) {
    try {
        const response = await axiosIns.post('/order', orderData);
        return response.data;
    } catch (error) {
        throw new Error('Failed to create order');
    }
}