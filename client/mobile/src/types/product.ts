export type ProductCategory = 'Equipment' | 'Tools' | 'Accessories' | 'Footwear';

export interface Product {
  id: string;
  name: string;
  brand: string;
  price: number;
  originalPrice?: number;
  category: ProductCategory;
  images: string[];
  description: string;
  details: string[];
  fit: string;
  sizes: string[];
  colors: string[];
  isNew?: boolean;
  isTrending?: boolean;
}

export interface CartItem {
  product: Product;
  quantity: number;
  size: string;
  color: string;
}

export interface Brand {
  id: string;
  name: string;
  logo: string;
}
