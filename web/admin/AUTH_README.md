# Authentication System Documentation

This document outlines the improved authentication system for the admin dashboard, including token rotation, NextAuth configuration, and security best practices.

## Overview

The authentication system has been completely refactored to provide:
- **Secure token rotation** with automatic refresh
- **Improved error handling** and user feedback
- **Type-safe implementation** with proper TypeScript types
- **Better security** with secure cookies and proper session management
- **Enhanced user experience** with loading states and error messages

## Architecture

### 1. NextAuth Configuration (`/app/api/auth/[...nextauth]/route.ts`)

The main authentication configuration includes:

#### Providers
- **Credentials Provider**: Email/password authentication
- **Google OAuth**: Social login (optional, configured via environment variables)

#### Token Management
- **Automatic token refresh**: Tokens are refreshed before expiration
- **Error handling**: Proper handling of refresh failures
- **Session management**: Secure session storage with JWT strategy

#### Security Features
- **Secure cookies**: HTTP-only cookies in production
- **CSRF protection**: Built-in NextAuth CSRF protection
- **Session timeout**: 72-hour session duration
- **Environment-based configuration**: Different settings for dev/prod

### 2. Token Utilities (`/app/utils/token.ts`)

Utility functions for token management:

```typescript
// Token validation
isTokenExpired(token: string): boolean
getTokenExpirationTime(token: string): number

// User extraction
extractUserFromToken(token: string): UserInfo

// Role checking
hasRole(token: string, role: string): boolean
hasAnyRole(token: string, roles: string[]): boolean
hasAllRoles(token: string, roles: string[]): boolean

// Token refresh
refreshToken(refreshToken: string, apiUrl: string): Promise<...>
```

### 3. Authentication Hook (`/hooks/useAuth.ts`)

Custom React hook providing authentication state and actions:

```typescript
const {
  // State
  user,
  isAuthenticated,
  isLoading,
  session,
  
  // Actions
  login,
  logout,
  refreshSession,
  
  // Role checks
  hasRole,
  hasAnyRole,
  hasAllRoles,
  
  // Utilities
  accessToken,
} = useAuth()
```

### 4. Axios Configuration (`/app/libs/axios.ts`)

Enhanced HTTP client with automatic token management:

- **Request interceptor**: Automatically adds authentication headers
- **Response interceptor**: Handles 401 errors with automatic token refresh
- **Error handling**: Comprehensive error handling with user feedback
- **Timeout configuration**: 30-second request timeout

### 5. Middleware (`/middleware.ts`)

Route protection with NextAuth middleware:

- **Protected routes**: All routes except auth pages require authentication
- **Automatic redirects**: Unauthenticated users redirected to signin
- **Custom authorization**: Extensible authorization logic

## Environment Variables

Required environment variables:

```env
# NextAuth
NEXTAUTH_SECRET=your-secret-key
NEXTAUTH_URL=http://localhost:3000

# API Configuration
NEXT_BASE_URL=https://your-api-url.com

# Google OAuth (optional)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

## Usage Examples

### Basic Authentication Check

```typescript
import { useAuth } from '@/hooks/useAuth'

function MyComponent() {
  const { isAuthenticated, user, isLoading } = useAuth()
  
  if (isLoading) return <div>Loading...</div>
  if (!isAuthenticated) return <div>Please sign in</div>
  
  return <div>Welcome, {user?.name}!</div>
}
```

### Role-Based Access Control

```typescript
import { useAuth } from '@/hooks/useAuth'

function AdminPanel() {
  const { hasRole, hasAnyRole } = useAuth()
  
  // Check for specific role
  if (!hasRole('admin')) {
    return <div>Access denied</div>
  }
  
  // Check for multiple roles
  if (!hasAnyRole(['admin', 'moderator'])) {
    return <div>Insufficient permissions</div>
  }
  
  return <div>Admin panel content</div>
}
```

### API Calls with Authentication

```typescript
import axiosInstance from '@/app/libs/axios'

// Token is automatically added to requests
const response = await axiosInstance.get('/api/users')
```

### Manual Token Refresh

```typescript
import { useAuth } from '@/hooks/useAuth'

function TokenManager() {
  const { refreshSession } = useAuth()
  
  const handleRefresh = async () => {
    await refreshSession()
  }
  
  return <button onClick={handleRefresh}>Refresh Session</button>
}
```

## Refresh Attempt Limiting

The authentication system implements refresh attempt limiting to prevent abuse and improve security:

### Configuration
- **Maximum attempts**: 3 refresh attempts per user
- **Reset interval**: Attempts reset every 5 minutes
- **Tracking**: Per-user attempt tracking with user ID

### Implementation
- **NextAuth level**: Token refresh attempts tracked in JWT
- **Axios level**: Request retry attempts limited to 2 per request
- **Global manager**: Singleton class manages attempts across the application

### Behavior
1. **Successful refresh**: Attempt counter resets to 0
2. **Failed refresh**: Attempt counter increments
3. **Max attempts reached**: User automatically logged out
4. **Periodic reset**: All attempts reset every 5 minutes

### Usage Example

```typescript
import { useAuth } from '@/hooks/useAuth'

function TokenStatus() {
  const { refreshAttempts, canRefresh, refreshSession } = useAuth()
  
  return (
    <div>
      <p>Refresh attempts: {refreshAttempts}/3</p>
      <p>Can refresh: {canRefresh ? 'Yes' : 'No'}</p>
      <button 
        onClick={refreshSession}
        disabled={!canRefresh}
      >
        Refresh Session
      </button>
    </div>
  )
}
```

## Security Best Practices

### 1. Token Security
- Tokens are stored in HTTP-only cookies
- Automatic token rotation prevents token reuse
- Secure cookie settings in production
- **Refresh attempt limiting**: Maximum 3 refresh attempts per user to prevent abuse

### 2. Session Management
- 72-hour session timeout
- Automatic session refresh with attempt limiting
- Proper session cleanup on logout
- **Rate limiting**: Refresh attempts are tracked and limited per user

### 3. Error Handling
- Comprehensive error logging
- User-friendly error messages
- Automatic logout on authentication failures
- **Attempt tracking**: Failed refresh attempts are tracked and limited

### 4. CSRF Protection
- NextAuth built-in CSRF protection
- Secure cookie configuration
- Proper request validation

## Troubleshooting

### Common Issues

1. **Token Refresh Failures**
   - Check API endpoint availability
   - Verify refresh token validity
   - Check network connectivity

2. **Session Expiration**
   - Tokens automatically refresh before expiration
   - Manual refresh available via `refreshSession()`
   - Automatic logout on refresh failure

3. **Role-Based Access Issues**
   - Verify user roles in token payload
   - Check role checking functions
   - Ensure proper role assignment

### Debug Mode

Enable debug mode in development:

```typescript
// In NextAuth config
debug: process.env.NODE_ENV === 'development'
```

## Migration Guide

### From Old Authentication System

1. **Update imports**: Replace direct NextAuth usage with `useAuth` hook
2. **Update session access**: Use `user` object from hook instead of `session.user`
3. **Update API calls**: Use `axiosInstance` instead of custom axios setup
4. **Update role checks**: Use hook methods instead of manual token decoding

### Breaking Changes

- Session structure changed to include `accessToken` and `error` fields
- Token refresh logic moved to NextAuth callbacks
- API client now handles authentication automatically

## Performance Considerations

- **Token caching**: Tokens are cached in memory for performance
- **Lazy loading**: Authentication state loaded on demand
- **Optimistic updates**: UI updates immediately, validated in background
- **Minimal re-renders**: Hook optimized to prevent unnecessary re-renders

## Future Enhancements

- [ ] Multi-factor authentication support
- [ ] Session analytics and monitoring
- [ ] Advanced role-based permissions
- [ ] Audit logging for authentication events
- [ ] Integration with external identity providers 