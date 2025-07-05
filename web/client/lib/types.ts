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
    DistributorId: number;
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
    distributor_id: number;
    categories: number[];
    is_active: boolean


}

export type Catalogue = {
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
    retailer_id: number;
    items: CheckoutRequestItem[]
    payment_partner_id: number
    payment_method: string
}

export type Partner = {
    id: number;
    name: string;
    icon: string;
    base_url: string;
    payment_method: string;
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

export type Transaction = {
    id: number
    date: string
    amount: GLfloat
    partner_id: number
    tx_ref: string
    status: string
}


export type Order = {
    Id: number;
    RetailerId: number;
    RetailerName: string;
    Items: Item[]
    Total: number
    Status: string
    DeliveryStatus: string
    PaymentStatus: string
    ConfirmationStatus: string
    CreatedAt: string
    ExpiresAt: string
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