export type Image = {
  ImageUrl: string;
  BlurHash: string;
};

export type Item = {
  ProductId: number;
  ProductName: string;
  ProductPrice: number;
  Quantity: number;
};

export type Product = {
  Id: number;
  Name: string;
  Desc: string;
  ExternalID: string;
  Images: Image[];
  Price: number;
  Attributes: Record<string, string>[];
  SupplierId: number;
  CategoryId: number[];
  Stock: number;
  AvailableStock: number;
  ReservedStock: number;
  IsActive: number;
};

type Configurables = {
  id: number;
  name: string;
  desc: string;
  price: number;
  stock: number;
  external_id: string;
  attributes: Record<string, string>;
  images: Image[];
  supplier_id: number;
  categories: number[];
  is_active: boolean;
};

export type Catalogue = {
  name: string;
  desc: string;
  price?: number;
  is_active: boolean;
  images: Image[];
  configurable_attributes: any;
  configurables: Configurables[];
};

export type Category = {
  id: number;
  name: string;
};

export type PricingPlan = {
  id: number;
  name: string;
  price: number;
  term_in_month: number;
  desc: string;
};

export interface CartItem {
  id: number;
  name: string;
  price: number;
  image: string;
  quantity: number;
}

export type Supplier = {
  id: number;
  name: string;
  tin: string;
  latitude: string;
  longitude: string;
  general_zone: string;
  region: string;
  woreda: string;
  user: number[];
  is_active: boolean;
};

