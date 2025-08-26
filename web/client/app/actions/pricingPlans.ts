"use server";
import axiosIns from "@/lib/axios";

export async function fetchPricingPlans() {
  try {
    const response = await axiosIns.get("/plan");
    console.log(response.data)
    return response.data.body.plans.list;
  } catch (error) {
    console.log(error);
    throw new Error("Failed to fetch pricing plans");
  }
}
