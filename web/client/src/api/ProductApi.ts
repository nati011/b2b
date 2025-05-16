import axios from "axios";

const API_URL = "/api/v1/catalogue";

export const fetchProducts = async () => {
  try {
    const response = await axios.get(API_URL);
    const products = response.data.body.products.map((product) => ({
      name: product.name,
      price:
        product.configurables.length > 1
          ? Math.min(...product.configurables.map((cfg) => cfg.price))
          : product.configurables[0].price,
      description: product.desc,
      images: product.images,
      configurable_attributes: product.configurable_attributes,
      configurables: product.configurables,
    }));
    return products;
  } catch (error) {
    console.error("Error fetching products:", error);
    throw error;
  }
};
