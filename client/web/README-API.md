# Web Client – API Reference

All HTTP APIs used by the web client. The backend base URL is set via `NEXT_PUBLIC_BASE_URL` (default `http://localhost:8090` in the browser). Requests use the axios instance from `lib/axios.ts`, which adds auth and base URL.

---

## Backend APIs (same origin / proxy)

### Auth (`/api/v1` and legacy paths)

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| POST | `/api/v1/auth/login` | Email/password login. Returns `access_token`, `refresh_token`. | `lib/auth.ts` (Credentials), `app/actions/auth.ts` (Login) |
| POST | `/api/v1/auth/refresh` | Refresh access token. Body: `{ refresh_token }`. | `lib/auth.ts` |
| POST | `/api/v1/auth/sso` | Google SSO. Body: `{ token, first_name, last_name, email }`. | `lib/auth.ts` |
| POST | `/api/v1/customer` | Register customer. | `app/actions/auth.ts` (RegisterCustomer) |
| POST | `/api/v1/supplier` | Register supplier. | `app/actions/auth.ts` (RegisterSupplier) |
| POST | `/user/init_reset` | Initiate password reset. Body: `{ email }`. | `app/actions/auth.ts` (InitResetPassword) |
| POST | `/auth/reset/:token` | Reset password with token. Body: `{ password }`. | `app/actions/auth.ts` (ResetPassword) |
| PATCH | `/user` | Update current user profile. | `app/actions/auth.ts` (UpdateProfile) |

### Identity

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/identity/user` | Current authenticated user. | `app/actions/getCurrentUser.ts` |

### Admin – Users & RBAC (`/api/v1`)

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/api/v1/users` | List users. Query: `page`, `limit`. | `app/actions/admin.ts` (getUsers) |
| GET | `/api/v1/users/:id` | Get user by id. | `app/actions/admin.ts` (getUserById) |
| POST | `/api/v1/users` | Create user. | `app/actions/admin.ts` (createUser) |
| PUT | `/api/v1/users/:id` | Update user. | `app/actions/admin.ts` (updateUser) |
| DELETE | `/api/v1/users/:id` | Delete user. | `app/actions/admin.ts` (deleteUser) |
| POST | `/api/v1/users/:userId/roles` | Assign roles. Body: `{ role_ids }`. | `app/actions/admin.ts` (assignRolesToUser) |
| POST | `/api/v1/users/:id/activate` | Activate user. | `app/actions/admin.ts` (activateUser) |
| POST | `/api/v1/users/:id/deactivate` | Deactivate user. | `app/actions/admin.ts` (deactivateUser) |
| POST | `/api/v1/users/:id/suspend` | Suspend user. | `app/actions/admin.ts` (suspendUser) |
| GET | `/api/v1/roles` | List roles. | `app/actions/admin.ts` (getRoles) |
| GET | `/api/v1/roles/:id` | Get role by id. | `app/actions/admin.ts` (getRoleById) |
| POST | `/api/v1/roles` | Create role. | `app/actions/admin.ts` (createRole) |
| PUT | `/api/v1/roles/:id` | Update role. | `app/actions/admin.ts` (updateRole) |
| DELETE | `/api/v1/roles/:id` | Delete role. | `app/actions/admin.ts` (deleteRole) |
| GET | `/api/v1/permissions` | List permissions. | `app/actions/admin.ts` (getPermissions) |
| GET | `/api/v1/permissions/:id` | Get permission by id. | `app/actions/admin.ts` (getPermissionById) |
| POST | `/api/v1/permissions` | Create permission. | `app/actions/admin.ts` (createPermission) |
| PUT | `/api/v1/permissions/:id` | Update permission. | `app/actions/admin.ts` (updatePermission) |
| DELETE | `/api/v1/permissions/:id` | Delete permission. | `app/actions/admin.ts` (deletePermission) |

### Customer

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/customer` | List/search customers. Query: `page`, `limit`, `email`, etc. | `app/actions/customer.ts` (getCustomerByEmail), `app/admin/customer/page.tsx`, `app/admin/page.tsx` |
| GET | `/customer/:id` | Get customer by id. | `app/admin/customer/[id]/page.tsx` |
| POST | `/customer` | Create customer. | `app/actions/customer.ts` (createCustomer) |

### Supplier

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/supplier` | List suppliers. Query: `page`, `limit`, etc. | `app/admin/supplier/page.tsx`, `app/admin/account/page.tsx`, `app/supplier/account/page.tsx`, `app/admin/products/new/page.tsx`, `app/admin/page.tsx` |
| GET | `/supplier/:id` | Get supplier by id. | `app/admin/supplier/[id]/page.tsx` |
| GET | `/supplier/bank-account` | List bank accounts. Query: `supplier_id`. | `app/actions/bankAccounts.ts` (listBankAccounts) |
| POST | `/supplier/bank-account` | Create bank account. | `app/actions/bankAccounts.ts` (createBankAccount) |
| PUT | `/supplier/bank-account/:id` | Update bank account. | `app/actions/bankAccounts.ts` (updateBankAccount) |
| DELETE | `/supplier/bank-account/:id` | Delete bank account. | `app/actions/bankAccounts.ts` (deleteBankAccount) |

### Product

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/products` | List products. Query: `supplier_id`, `category_id`, `is_active`, `limit`, `offset`. | `app/actions/product.ts` (ListProducts) |
| GET | `/products/supplier` | List products for authenticated supplier. Query: `category_id`, `is_active`, `limit`, `offset`. | `app/actions/product.ts` (ListSupplierProducts) |
| GET | `/product` | Get product by id. Query: `id`. | `app/actions/product.ts` (GetProduct), `app/admin/orders/[id]/page.tsx`, `app/supplier/orders/[id]/page.tsx` |
| POST | `/product` | Create product. | `app/actions/product.ts` (CreateProduct) |
| PUT | `/product` | Update product. Query: `id`. | `app/actions/product.ts` (UpdateProduct) |
| POST | `/product/grn` | Create/record GRN (goods received note). | `app/admin/products/page.tsx`, `app/supplier/products/page.tsx` |
| PATCH | `/product/price` | Update product price. | `app/admin/products/page.tsx`, `app/supplier/products/page.tsx` |

### Catalogue

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/catalogue` | List catalogue. Query: `supplier_id`, `category_id`, `is_active`, `limit`, `offset`. | `app/actions/catalogue.ts` (GetAllCatalogues) |

### Category

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/category` | List categories. | `app/actions/category.ts` (GetAllCategories) |
| POST | `/category` | Create category. Body: `{ name }`. | `app/actions/category.ts` (Create) |
| PATCH | `/category/:id` | Update category. Body: `{ name }`. | `app/actions/category.ts` (Update) |
| DELETE | `/category` | Delete category. Query: `id`. | `app/actions/category.ts` (Delete) |

### Orders & Invoice

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/orders/customer` | List customer orders. Query: `customer_id`, `status`, `limit`, `offset`. | `app/actions/orders.ts` (ListCustomerOrders, fetchOrders), `app/admin/orders/page.tsx`, `app/admin/page.tsx` |
| GET | `/order/supplier` | List supplier orders. Query: `status`, `limit`, `offset`. | `app/actions/orders.ts` (ListSupplierOrders, fetchOrders), `app/admin/orders/page.tsx`, `app/supplier/orders/page.tsx` |
| GET | `/order` | Get order by id. Query: `id`. | `app/actions/orders.ts` (getOrderById) |
| POST | `/order` | Create order. Body: `{ customer_id, items: [{ product_id, quantity }], ... }`. | `app/actions/orders.ts` (createOrder) |
| PATCH | `/order` | Update order. Query: `id`, `command` or `payment_status`. | `app/actions/orders.ts` (updateOrderStatus, updateOrderPaymentStatus) |
| GET | `/invoice` | Get invoice for order. Query: `order_id`. | `app/actions/orders.ts` (getInvoice) |

### Pricing plans

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/api/v1/plan` | List pricing plans (subscription plans). | `app/actions/pricingPlans.ts` (fetchPricingPlans) |

### Admin dashboard fallback

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| GET | `/api/v1/orders` | Fallback when `/orders/customer` fails (admin dashboard). | `app/admin/page.tsx` |

---

## Next.js API routes (internal)

These are handled by the Next.js app (same host as the web app), not the backend service.

| Method | Path | Description | Used in |
|--------|------|-------------|---------|
| POST | `/api/email` | Send contact form email. Body: `{ email, name, message }`. Uses Gmail (env: `EMAIL`, `PASSWORD`). | Contact form / email flow |

---

## Third-party APIs

| Service | Method | URL / usage | Used in |
|---------|--------|-------------|---------|
| Cloudinary | POST | `https://api.cloudinary.com/v1_1/{cloud_name}/image/upload` | `components/ImageUpload.tsx` – image upload with `file`, `upload_preset`, `api_key` (config: `cloud_name`, `api_key` in component). |

---

## Auth and base URL

- **Client:** `lib/axios.ts` builds the base URL from `NEXT_PUBLIC_BASE_URL` (browser uses e.g. `http://localhost:8090`; Docker hostnames are replaced for browser requests). The axios instance attaches the Bearer token from `localStorage` (key `access_token`) for backend calls; `/api/v1/auth/login` is not authenticated.
- **Server (NextAuth):** `lib/auth.ts` uses `API_BASE_URL` or `NEXT_PUBLIC_BASE_URL` for login, refresh, and SSO calls to the backend.

All backend paths in the tables above are relative to the backend base URL. Next.js routes are relative to the web app origin.
