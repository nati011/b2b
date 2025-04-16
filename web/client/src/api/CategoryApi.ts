import axios from "axios";

interface Category {
  id: number;
  name: string;
}

const API_URL = "/api/v1/category";

export const fetchCategories = async (): Promise<Category[]> => {
  try {
    const response = await axios.get(API_URL);
    return response.data.body.category.categories || [];
  } catch (error) {
    console.error("Error fetching categories:", error);
    throw error;
  }
};
