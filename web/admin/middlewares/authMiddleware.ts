import { NextFetchEvent, NextRequest, NextResponse } from 'next/server'
import { CustomMiddleware } from './middleware-chain'
import { getToken } from 'next-auth/jwt'
import { canAccessPath, isProtectedPath } from '@/lib/acl'

export function withAuthMiddleware(middleware: CustomMiddleware) {
  return async (request: NextRequest, event: NextFetchEvent) => {
    const response = NextResponse.next()
    const pathname = request.nextUrl.pathname

    if (!isProtectedPath(pathname)) {
      return middleware(request, event, response)
    }

    const token = await getToken({ 
      req: request,
      secret: process.env.NEXTAUTH_SECRET 
    })
    if ((token?.error != null && token?.error != "ForbiddenError") || !token ) {
      return NextResponse.redirect(new URL('/auth/signin', request.url))
    }
    
    const tokenWithUser = token as any
    const permissions = tokenWithUser?.user?.permissions
    const allowed = canAccessPath(permissions, pathname)
    console.log(allowed, "Allowed")
    // if (!allowed) {
    //   return NextResponse.redirect(new URL('/forbidden', request.url))
    // }

    return middleware(request, event, response)
  }
}