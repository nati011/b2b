import { withAuth } from "next-auth/middleware"
import { NextResponse } from "next/server"

export default withAuth(
  function middleware(req) {
    // Add custom middleware logic here if needed
    return NextResponse.next()
  },
  {
    callbacks: {
      authorized: ({ token, req }) => {
        // Allow access to protected routes if user is authenticated
        if (req.nextUrl.pathname.startsWith('/auth')) {
          return true // Allow access to auth pages
        }
        
        // Require authentication for all other routes
        return !!token
      },
    },
    pages: {
      signIn: '/auth/signin',
    },
  }
)

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