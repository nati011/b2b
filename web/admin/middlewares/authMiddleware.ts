import { NextFetchEvent, NextRequest, NextResponse } from 'next/server'
import { CustomMiddleware } from './middleware-chain'
import { getToken } from 'next-auth/jwt'
import { ACL } from '@/lib/constants'
import { signOut } from 'next-auth/react'

const publicPaths = [
  '/auth/signin',
  '/auth/signup',
  '/auth/error',
  '/auth/verify',
  '/forbidden',
  '/404',
  '/500'
]

function isPublicPath(pathname: string): boolean {
  return publicPaths.some(path => 
    pathname.startsWith(path) || pathname === path
  )
}

function isProtectedPath(pathname: string): boolean {
  return !isPublicPath(pathname)
}

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
    console.log(token)
    if ((token?.error != null && token?.error != "ForbiddenError") || !token ) {
      return NextResponse.redirect(new URL('/auth/signin', request.url))
    }

    const requiredPermission = ACL[pathname]
    if (requiredPermission) {
      let userPermissions: any[] = []
      
      const tokenWithUser = token as any
      if (tokenWithUser.user?.permissions?.List && Array.isArray(tokenWithUser.user.permissions.List)) {
        userPermissions = tokenWithUser.user.permissions.List
      }

      const hasPermission = userPermissions.some(
        (perm: any) => requiredPermission.includes(perm.Name)
      )

      if (!hasPermission) {
        return NextResponse.redirect(new URL('/forbidden', request.url))
      }
    }

    return middleware(request, event, response)
  }
}