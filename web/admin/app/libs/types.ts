
export type UserAccount = {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  username: string;
  dob: string;
  is_active: boolean;
  external_id: string;
};

export type Retailer = {
  id: number;
  name: string;
  tin: string;
  latitude: string;
  longitude: string;
  general_zone: string;
  region: string;
  woreda: string;
  user: UserAccount;
};

export type Distributor = {
  id: number;
  name: string;
  tin: string;
  latitude: string;
  longitude: string;
  general_zone: string;
  region: string;
  woreda: string;
  user: number[];
};

export type Image = {
  ImageUrl: string
  BlurHash: string
}

export type Product = {
  Id: number;
  Name: string;
  Desc: string;
  ExternalID: string;
  Images: Image[];
  Price: number;
  Attributes: Record<string, string>[];
  DistributorId: number;
  CategoryId: number[];
  Stock: number;
  AvailableStock: number;
  ReservedStock: number;
  IsActive: number;
};



export type ProductForm = {
  Id: number
  Name: string
  ExternalID: string
  Price: GLfloat
  Attributes: any[];
  Desc: string
  Images: string[]
  CategoryId: number[]
}

export type PriceRange = {
  min: number
  max: number
}

export type ConfigurableProduct = {
  Id: number;
  Name: string;
  Desc: string;
  ExternalId: string;
  Images: Image[];
  Attributes: string[];
  DistributorId: number;
  CategoryId: number;
  PriceRange: PriceRange;
  IsAvailable: boolean;
  Products: number[];
};


export type ConfigurableProductForm = {
  id?: number;
  name: string;
  desc: string;
  external_id: string;
  images: string[];
  attribute_keys: string[];
  products: number[];
};

export type Category = {
  id: number;
  name: string;
};

export type DistributorRequest = {
  name: string;
  tin: string;
  latitude: string;
  longitude: string
  general_zone: string;
  region: string;
  woreda: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  username: string;
  dob: string;
  external_id: string;
}


export type Item = {
  ProductId: number
  ProductName: string
  Quantity: number
}


export type Order = {
  Id: number;
  RetailerId: number;
  Items: Item[]
  Total: number
  Status: string
  DeliveryStatus: string
  PaymentStatus: string
  CreatedAt: string
}

export type Profile = {
  id: number
  first_name: string
  last_name: string
  email: string
  phone: string
  username: string
  dob: string
}


export type Role = {
  id: number;
  name: string;
  desc: string;
};

export type Resource = {
  id: number;
  name: string;
  action: string;
};

export type Transaction = {
  id: number
  date: string
  amount: GLfloat
  partner_id: number
  tx_ref: string
  status: string
}


export type InvoiceItem = {
  ProductId: number
  ProductName: string
  ProductQuantity: number
  ProductPrice: GLfloat
}

export type Invoice = {
  Id: number
  Created_Date: string
  ExternalId: string
  Status: string
  OrderId: number
  SubTotal: GLfloat
  LineItems: InvoiceItem[]
  TaxAmount: GLfloat
}