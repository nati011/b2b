
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

// Middleware to protect admin routes
export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Protect admin routes - redirect to login if not authenticated
  // Note: Actual role checking happens in the admin layout component
  if (pathname.startsWith('/admin')) {
    // The admin layout will handle authentication and authorization
    // This middleware just ensures the route exists
    return NextResponse.next();
  }

  return NextResponse.next();
}

export const config = {
  matcher: '/admin/:path*',
}