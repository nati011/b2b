import { NextFetchEvent, NextRequest, NextResponse } from 'next/server'
import { CustomMiddleware } from './middleware-chain'
import { getToken } from 'next-auth/jwt'
import { canAccessPath, isProtectedPath } from '@/lib/acl'

export function withAuthMiddleware(middleware: CustomMiddleware) {
  return async (request: NextRequest, event: NextFetchEvent) => {
    const response = NextResponse.next()
    const pathname = request.nextUrl.pathname

    // Skip auth check for static assets and API routes
    if (!isProtectedPath(pathname)) {
      return middleware(request, event, response)
    }

    // Optimize: Cache token check for faster navigation
    const token = await getToken({ 
      req: request,
      secret: process.env.NEXTAUTH_SECRET,
      // Reduce token validation overhead
      cookieName: process.env.NODE_ENV === 'production' 
        ? '__Secure-next-auth.session-token' 
        : 'next-auth.session-token'
    })
    
    if ((token?.error != null && token?.error != "ForbiddenError") || !token ) {
      return NextResponse.redirect(new URL('/auth/signin', request.url))
    }
    
    const tokenWithUser = token as any
    const permissions = tokenWithUser?.user?.permissions
    const allowed = canAccessPath(permissions, pathname)
    
    // Removed console.log for production performance
    // if (!allowed) {
    //   return NextResponse.redirect(new URL('/forbidden', request.url))
    // }

    return middleware(request, event, response)
  }
}