export type AuthModel = {
    email: string
    password: string
}

export type Register = {
    first_name: string
    last_name: string
    email: string
    phone: string
    password: string
    confirm_password: string
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

export type Retailer = {
    id: number
    name: string
    tin: string
    latitude: string
    longitude: string
    general_zone: string
    region: string
    woreda: string
    user: User
}

export type Partner = {
    id: number;
    name: string;
    icon: string;
    base_url: string
}

export type Item = {
    id: number
    quantity: number
}

export type CheckoutRequest = {
    retailer_id: number;
    items: Item[]
    payment_partner_id: number
}