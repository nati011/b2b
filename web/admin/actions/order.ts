'use server'
import axiosIns from "@/app/libs/axios";

export async function fetchOrders() {
    try {
        const response = await axiosIns.get('http://localhost:3000/api/order');
        console.log(response.data)
        return response.data.body.List;
    } catch (error) {
        console.log(error)
        throw new Error('Failed to fetch orders');
    }
}

export async function getOrderById(orderId: string) {
    try {
        const response = await axiosIns.get(`/api/order?id=${orderId}`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch order');
    }
}

export async function updateOrderStatus(orderId: string, status: string) {
    try {
        const response = await axiosIns.put(`/api/order/${orderId}`, { status });
        return response.data;
    } catch (error) {
        throw new Error('Failed to update order status');
    }
}

export async function createOrder(orderData: any) {
    try {
        const response = await axiosIns.post('/api/order', orderData);
        return response.data;
    } catch (error) {
        throw new Error('Failed to create order');
    }
}