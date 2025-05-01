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
  user: UserAccount;
};

export type Product = {
  Id: number;
  Name: string;
  Desc: string;
  ExternalID: string;
  Images: string;
  Price: number;
  Attributes: any;
  DistributorId: number;
  CategoryId: number;
  Stock: number;
  IsActive: number;
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