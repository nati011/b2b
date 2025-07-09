import {chain} from '@/middlewares/middleware-chain'
import { withAuthMiddleware } from './middlewares/authMiddleware'


export default chain([withAuthMiddleware])

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - images (public images)
     * - auth (authentication pages - handled by the callback)
     */
    '/((?!api|_next/static|_next/image|favicon.ico|images).*)',
  ],
}