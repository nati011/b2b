import axios from "axios";

interface Category {
  Name: string;
  // Add other properties if needed
}

const API_URL = "/api/v1/category";

export const fetchCategories = async (): Promise<Category[]> => {
  try {
    const response = await axios.get(API_URL);
    return response.data.categories?.List || [];
  } catch (error) {
    console.error("Error fetching categories:", error);
    throw error; // Re-throw the error for handling in the component
  }
};
