"use server";
import axiosIns from "@/lib/axios";

export const GetAllCatalogues = async (url?: string) => {
  try {
    const response = await axiosIns.get("/api/v1/catalogue");
    return response.data.body.products;
  } catch (error: any) {
    if (error.response) {
      throw (
        error.response.data.message ||
        "An error has occured while creating the product"
      );
    }
    throw "An error has occured while creating the product";
  }
};
