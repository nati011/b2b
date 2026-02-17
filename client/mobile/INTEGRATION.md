# Mobile Client – API Integration Points

This document lists all backend APIs used by the mobile app, including base URL, authentication, and per-endpoint details.

## Base configuration

| Setting | Source | Default |
|--------|--------|---------|
| Base URL | `VITE_API_BASE_URL` env | `http://localhost:8090` |
| Timeout | Axios instance | 15s |
| Content-Type | Default header | `application/json` |

**Authentication:** Requests use either **Basic Auth** (email + password from `localStorage`: `user_email`, `user_password`) or **Bearer** token (`auth_token`). No auth is sent if neither is present.

**Axios instance:** `src/lib/axios.ts`

---

## 1. Auth

**Module:** `src/lib/api/auth.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Customer login |
| POST | `/api/v1/auth/sso` | Google SSO login |

### POST `/api/v1/auth/login`

- **Request body:** `{ email: string, password: string }`
- **Success:** `202`; response body has `body.access_token` and optional `body.refresh_token`.
- **Side effects:** Stores `auth_token`, `refresh_token`, `user_id`, `user_email`, `user_name` in `localStorage`.

### POST `/api/v1/auth/sso`

- **Request body:** `{ token: string, first_name: string, last_name: string, email: string }` (Google OAuth token + profile).
- **Success:** `202`; same token/refresh shape as login.
- **Side effects:** Same as login.

**External call (Google):** Login page uses `fetch('https://www.googleapis.com/oauth2/v2/userinfo')` with the Google access token to get user profile before calling `/api/v1/auth/sso`.

---

## 2. Orders

**Module:** `src/lib/api/order.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/order` | Create order |
| GET | `/order?id={id}` | Get order by ID |
| PATCH | `/order?id={id}&command={command}` | Update order status |
| GET | `/orders/customer?customer_id=&status=&limit=&offset=` | List customer orders |

### POST `/order`

- **Request body:** `CreateOrderRequest`
  - `customer_id: number`
  - `status?`, `payment_status?`, `delivery_status?`, `confirmation_status?`: string
  - `total?`: number
  - `customer_snapshot?`, `shipping_address_snapshot?`, `billing_address_snapshot?`: any
  - `items: OrderItemRequest[]` where each item: `product_id?`, `id?`, `quantity: number`, `price?`
- **Response:** Single order object (`OrderResponse`).

### GET `/order?id={id}`

- **Response:** Single order (`OrderResponse`), including `cart_snapshot`, `customer_snapshot`, `items`, etc.

### PATCH `/order?id={id}&command={command}`

- **Query:** `id`, `command` (e.g. status update command).
- **Response:** Updated order.

### GET `/orders/customer`

- **Query:** `customer_id` (required), `status`, `limit`, `offset`.
- **Response:** Handles both `response.data.body` and `response.data`; expects `orders` array and optional `total`, `limit`, `offset`.

**Used in:** Checkout (create), Orders (list, get, update status).

---

## 3. Products

**Module:** `src/lib/api/product.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/product?id={id}` | Get single product |
| GET | `/products?supplier_id=&category_id=&is_active=&limit=&offset=` | List products |
| POST | `/product` | Create product |
| PUT | `/product?id={id}` | Update product |
| DELETE | `/product?id={id}` | Delete product |

### GET `/product?id={id}`

- **Response:** Single product (`ProductResponse`: id, name, description, external_id, attributes, unit, is_active, supplier_id, price, total_quantity, reserved_quantity, available_quantity, category_ids, created_at, updated_at).

### GET `/products`

- **Query:** `supplier_id`, `category_id`, `is_active`, `limit`, `offset`.
- **Response:** `{ products: ProductResponse[], total, limit, offset }`.

### POST `/product`

- **Request body:** Product fields (name, description, external_id, attributes, unit, is_active, supplier_id, price, total_quantity, reserved_quantity, category_ids).
- **Response:** Created product.

### PUT `/product?id={id}`

- **Request body:** Same product fields (no supplier_id on update).
- **Response:** Updated product.

### DELETE `/product?id={id}`

- **Response:** No body on success.

**Used in:** ProductDetail (get), SupplierDetail (list by supplier).

---

## 4. Catalogue

**Module:** `src/lib/api/catalogue.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/catalogue?supplier_id=&category_id=&is_active=&limit=&offset=` | Get catalogue (product list) |

### GET `/catalogue`

- **Query:** `supplier_id`, `category_id`, `is_active`, `limit`, `offset`.
- **Response:** `{ products: Product[], total, limit, offset }` (products in app `Product` type).

**Used in:** Home, HeroSection (catalogue + categories/suppliers).

---

## 5. Customer

**Module:** `src/lib/api/customer.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/customer` | Create customer |
| GET | `/customer/{id}` | Get customer by ID |
| GET | `/customer?email={email}` | Get customer by email |
| GET | `/customer?page=&limit=` | List customers (paginated) |
| PUT | `/customer/{id}` | Update customer |
| DELETE | `/customer/{id}` | Delete customer |

### POST `/customer`

- **Request body:** `CreateCustomerRequest` (full_name, first_name, last_name, email, phone, phone_number, city, region, woreda, status, is_active, password).
- **Response:** `{ message, customer: CustomerResponse }`.

### GET `/customer/{id}`

- **Response:** Single `CustomerResponse` (id, full_name, status, city, region, woreda, phone_number, email, is_active, created_at, updated_at).

### GET `/customer?email={email}`

- **Response:** Handles `body` wrapper; if `items` array and non-empty, returns first item; if single object with `id`, returns it; else `null`. Used to resolve customer by logged-in user email.

### GET `/customer?page=&limit=`

- **Response:** Expects `items` array; returns full response (items, total, page, limit, total_pages, has_next, has_prev).

### PUT `/customer/{id}`

- **Request body:** `UpdateCustomerRequest` (same fields as create except password).
- **Response:** Updated customer.

### DELETE `/customer/{id}`

- **Response:** No body on success.

**Used in:** AuthContext (GetCustomerByEmail), Profile (ListCustomers), Checkout (GetCustomerByEmail), Orders (GetCustomerByEmail).

---

## 6. Supplier

**Module:** `src/lib/api/supplier.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/supplier?page=&limit=` | List suppliers |
| GET | `/supplier/{id}` or `/supplier?id={id}` | Get supplier by ID |
| POST | `/supplier` | Create supplier |
| PUT | `/supplier/{id}` | Update supplier |
| DELETE | `/supplier/{id}` | Delete supplier |

### GET `/supplier`

- **Query:** `page`, `limit`.
- **Response:** `SupplierListResponse` (items, total, page, limit, total_pages, has_next, has_prev). Each item: id, business_name, status, support_email, support_phone, is_active, created_at, updated_at.

### GET `/supplier/{id}` or `/supplier?id={id}`

- **Response:** Tries path first; on 404/400 falls back to query. Unwraps `response.data.supplier` or uses `response.data`; returns single supplier.

### POST `/supplier`

- **Request body:** `CreateSupplierRequest` (business_name, status, support_email, support_phone, is_active).
- **Response:** `{ message, supplier: SupplierResponse }`.

### PUT `/supplier/{id}`

- **Request body:** `UpdateSupplierRequest`.
- **Response:** Updated supplier.

### DELETE `/supplier/{id}`

- **Response:** No body on success.

**Used in:** Home, HeroSection (GetAllSuppliers), Suppliers (GetAllSuppliers), SupplierDetail (GetSupplier, ListProducts).

---

## 7. Category

**Module:** `src/lib/api/category.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| (none) | — | Categories not implemented on backend |

- `GetAllCategories()` currently returns `[]` and does not call the backend. Commented placeholder: `GET /category` to return `response.data.items`.

**Used in:** Home, HeroSection, Categories.

---

## 8. Pricing

**Module:** `src/lib/api/pricing.ts`

| Method | Endpoint | Description |
|--------|----------|-------------|
| (none) | — | Pricing plans not implemented on backend |

- `GetAllPricingPlans()` returns `[]` and does not call the backend. Commented placeholder: `GET /plan`.

**Used in:** (Reserved for future use.)

---

## Summary table

| Domain | Base path | Endpoints used by mobile |
|--------|-----------|---------------------------|
| Auth | `/api/v1/auth` | `POST /login`, `POST /sso` |
| Orders | `/order`, `/orders` | `POST /order`, `GET /order?id=`, `PATCH /order?id=&command=`, `GET /orders/customer?…` |
| Products | `/product`, `/products` | `GET /product?id=`, `GET /products?…`, `POST /product`, `PUT /product?id=`, `DELETE /product?id=` |
| Catalogue | `/catalogue` | `GET /catalogue?…` |
| Customer | `/customer` | `POST /customer`, `GET /customer/{id}`, `GET /customer?email=`, `GET /customer?page=&limit=`, `PUT /customer/{id}`, `DELETE /customer/{id}` |
| Supplier | `/supplier` | `GET /supplier?…`, `GET /supplier/{id}` or `?id=`, `POST /supplier`, `PUT /supplier/{id}`, `DELETE /supplier/{id}` |
| Category | `/category` | Not called (stub returns `[]`) |
| Pricing | `/plan` | Not called (stub returns `[]`) |

---

## Order item payload (configurable products)

The mobile app currently sends order items as `{ product_id, quantity, price }` and does not send `selected_attributes`. The backend supports optional `selected_attributes` per item for configurable products. To support configurable products on mobile, add `selected_attributes?: Record<string, string>` to each item in the checkout payload when the user has selected variant options.
