import axios from "axios";

type OrderItem = {
    ProductId: number
    Quantity: number
}

type PlaceOrderRequest = {
    RetailerId: number
    Items: OrderItem[]
}


const API_URL = "/api/v1/order";

export const placeOrder = async (orderReqest: any): Promise<number> => {
    try {
        const response = await axios.post(API_URL, orderReqest);
        console.log(response.body)
        return response.data.body.order.id || 0;
    } catch (error) {
        console.error("Error fetching categories:", error);
        throw error;
    }
};
