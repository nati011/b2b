
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

// Empty middleware - no authentication checks
export function middleware(request: NextRequest) {
  return NextResponse.next();
}

export const config = {
  matcher: [],
}