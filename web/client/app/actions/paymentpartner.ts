"use server";
import axiosIns from "@/lib/axios";

export async function fetchPaymentPartners() {
  try {
    const response = await axiosIns.get("/payment_option/active");
    console.log(response.data)
    return response.data.body.payment_options;
  } catch (error) {
    console.log(error);
    throw new Error("Failed to fetch orders");
  }
}
