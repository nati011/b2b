import {chain} from '@/middlewares/middleware-chain'
import { withAuthMiddleware } from './middlewares/authMiddleware'


export default chain([withAuthMiddleware])

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|images).*)',
  ],
}