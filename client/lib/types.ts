export type Image = {
    ImageUrl: string
    BlurHash: string
}

export type Item = {
    ProductId: number
    ProductName: string
    ProductPrice: GLfloat
    Quantity: number
}

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
    price: GLfloat;
    stock: number;
    external_id: string
    attributes: Record<string, string>;
    images: Image[];
    supplier_id: number;
    categories: number[];
    is_active: boolean


}

export type Catalogue = {
    id?: number; // Product ID from backend
    name: string;
    desc: string;
    price?: GLfloat
    is_active: boolean;
    images: Image[];
    configurable_attributes: any;
    configurables: Configurables[]
}

export type CheckoutRequestItem = {
    id: number
    quantity: number
}
export type CheckoutRequest = {
    customer_id: number;
    items: CheckoutRequestItem[]
}

export type RegisterRequest = {
    first_name: string
    last_name: string
    email: string
    phone: string
    password: string
}

export type User = {
    id?: number
    first_name: string
    last_name: string
    email: string
    phone: string
    username?: string
    dob?: string
    external_id?: string
}

// Backend Order Response (snake_case from API)
export type BackendOrder = {
    id: number;
    customer_id: number;
    status: string;
    payment_status?: string;
    delivery_status?: string;
    confirmation_status?: string;
    total?: number;
    customer_snapshot?: string | object;
    cart_snapshot: string | object; // Always included, can be empty array []
    items?: OrderItem[];
    created_at: string;
    updated_at: string;
}

// Backend Order Item Response
export type OrderItem = {
    product_id: number;
    quantity: number;
    price?: number;
}

// Frontend Order (camelCase for component usage)
export type Order = {
    Id: number;
    CustomerId: number;
    CustomerName: string;
    Items: Item[]
    Total: number
    Status: string
    DeliveryStatus: string
    PaymentStatus: string
    ConfirmationStatus: string
    CreatedAt: string
    ExpiresAt: string
    CartSnapshot?: string | object; // New field from backend
    CustomerSnapshot?: string | object; // New field from backend
}

export interface CartItem {
    id: number;
    name: string;
    price: number;
    image: string;
    quantity: number;
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



export type Permissions = {
    Id: number
    Name: string
    Action: string
}
export type UserIdentity = {
    id: number
    first_name: string
    last_name: string
    email: string
    phone: string
    username: string
    dob: string
    is_active: boolean
    permissions: Permissions[]
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
  }

  export type SupplierRequest = {
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
    licence_url: string;
    password: string;
    confirm_password: string
  }


export type PricingPlan = {
    id: number;
    name: string;
    price: number;
    term_in_month: number;
    desc: string;
}