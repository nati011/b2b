
export { default } from "next-auth/middleware"

export const config = {
    matcher: ['/cart', '/profile', '/checkout/successful']
}