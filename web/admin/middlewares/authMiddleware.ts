    import { NextFetchEvent, NextRequest, NextResponse } from 'next/server'
    import { CustomMiddleware } from './middleware-chain'

    import { getToken } from 'next-auth/jwt'
import { ACL } from '@/lib/constants'

    const protectedPaths = ['/']

    export function withAuthMiddleware(middleware: CustomMiddleware) {
        return async (request: NextRequest, event: NextFetchEvent) => {
            const response = NextResponse.next()

            const token = await getToken({ req: request })
            const pathname = request.nextUrl.pathname
            // @ts-ignore
            if(token?.error && protectedPaths.includes(pathname)){
                const signInUrl = new URL('/auth/signin', request.url)
                signInUrl.searchParams.set('callbackUrl', request.nextUrl.pathname)
                return NextResponse.redirect(signInUrl)
            }

            if ((!token) && protectedPaths.includes(pathname)) {
                const signInUrl = new URL('/auth/signin', request.url)
                signInUrl.searchParams.set('callbackUrl', pathname)
                return NextResponse.redirect(signInUrl)
            }
            const requiredPermission = ACL[pathname]
            if (requiredPermission) {
                let userPermissions: any[] = []
                if (token && Array.isArray((token as any).user?.permissions?.List)) {
                    userPermissions = (token as any).user?.permissions?.List
                }


                const hasPermission = userPermissions.some(
                    (perm: any) =>  requiredPermission.includes(perm.Name)
                )

                if (!hasPermission) {
                    return NextResponse.redirect(new URL('/forbidden', request.url))
                }
            }

            return middleware(request, event, response)
        }
    }
