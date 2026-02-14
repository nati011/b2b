// Image type for product images (if needed in future)
export type Image = {
  ImageUrl: string;
  BlurHash: string;
};

// Order item type
export type Item = {
  ProductId: number;
  ProductName: string;
  ProductPrice: number;
  Quantity: number;
};

// Product type matching backend response
export type Product = {
  id: number;
  name: string;
  description?: string;
  external_id?: string;
  attributes?: any;
  unit?: string;
  is_active: boolean;
  supplier_id: number;
  price?: number;
  total_quantity: number;
  reserved_quantity: number;
  available_quantity: number;
  category_ids?: number[];
  created_at: string;
  updated_at: string;
};

// Legacy Catalogue type (kept for backward compatibility)
// Use Product type for new code
export type Catalogue = {
  id?: number; // Product ID from backend
  name: string;
  desc: string;
  price?: number;
  is_active: boolean;
  images: Image[];
  configurable_attributes: any;
  configurables: any[];
};

// Category type (when category endpoint is implemented)
export type Category = {
  id: number;
  name: string;
};

// Pricing plan type (when pricing endpoint is implemented)
export type PricingPlan = {
  id: number;
  name: string;
  price: number;
  term_in_month: number;
  desc: string;
};

// Cart item type
export interface CartItem {
  product: Product;
  quantity: number;
  size: string;
  color: string;
}

// Supplier type matching backend response
export type Supplier = {
  id: number;
  business_name: string;
  status: string;
  support_email?: string;
  support_phone?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

// Customer type matching backend response
export type Customer = {
  id: number;
  full_name: string;
  status: string;
  city?: string;
  region?: string;
  woreda?: string;
  phone_number?: string;
  email?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

// Order type matching backend response
export type Order = {
  id: number;
  customer_id: number;
  status: string;
  payment_status?: string;
  delivery_status?: string;
  confirmation_status?: string;
  total?: number;
  customer_snapshot?: any;
  shipping_address_snapshot?: any;
  billing_address_snapshot?: any;
  items?: Array<{
    product_id: number;
    quantity: number;
    price?: number;
  }>;
  created_at: string;
  updated_at: string;
};

