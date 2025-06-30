import { useSession, signIn, signOut } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { useCallback, useEffect } from 'react'
import { toast } from 'sonner'
import { RefreshAttemptManager } from '@/app/utils/token'
import { UserIdentity } from '@/app/libs/types'

// Extended session type
interface ExtendedSession {
  user?: {
    id: string
    name: string
    email: string
    username: string
    roles: string[]
  } | null
  accessToken?: string
  userIdentity?: UserIdentity
  error?: string
  expires: string
}

export function useAuth() {
  const { data: session, status, update } = useSession()
  const router = useRouter()
  const refreshManager = RefreshAttemptManager.getInstance()
  
  const extendedSession = session as ExtendedSession

  // Check if user is authenticated
  const isAuthenticated = status === 'authenticated' && !!extendedSession?.user

  // Check if user is loading
  const isLoading = status === 'loading'

  // Check if user has a specific role
  const hasRole = useCallback((role: string) => {
    if (!extendedSession?.user?.roles) return false
    return extendedSession.user.roles.includes(role)
  }, [extendedSession])

  // Check if user has any of the specified roles
  const hasAnyRole = useCallback((roles: string[]) => {
    if (!extendedSession?.user?.roles) return false
    return roles.some(role => extendedSession.user!.roles.includes(role))
  }, [extendedSession])

  // Check if user has all of the specified roles
  const hasAllRoles = useCallback((roles: string[]) => {
    if (!extendedSession?.user?.roles) return false
    return roles.every(role => extendedSession.user!.roles.includes(role))
  }, [extendedSession])

  // Check if user has a specific permission
  const hasPermission = useCallback((permissionName: string) => {
    if (!extendedSession?.userIdentity?.permissions) return false
    return extendedSession.userIdentity.permissions.some(
      permission => permission.Name === permissionName
    )
  }, [extendedSession])

  // Check if user has any of the specified permissions
  const hasAnyPermission = useCallback((permissionNames: string[]) => {
    if (!extendedSession?.userIdentity?.permissions) return false
    return permissionNames.some(permissionName =>
      extendedSession.userIdentity!.permissions.some(
        permission => permission.Name === permissionName
      )
    )
  }, [extendedSession])

  // Check if user has all of the specified permissions
  const hasAllPermissions = useCallback((permissionNames: string[]) => {
    if (!extendedSession?.userIdentity?.permissions) return false
    return permissionNames.every(permissionName =>
      extendedSession.userIdentity!.permissions.some(
        permission => permission.Name === permissionName
      )
    )
  }, [extendedSession])

  // Sign in function
  const login = useCallback(async (email: string, password: string, callbackUrl?: string) => {
    try {
      const result = await signIn('credentials', {
        email,
        password,
        redirect: false,
      })

      if (result?.error) {
        toast.error(result.error)
        return { success: false, error: result.error }
      }

      if (result?.ok) {
        // Reset refresh attempts on successful login
        if (extendedSession?.user?.id) {
          refreshManager.resetAttempts(extendedSession.user.id)
        }
        toast.success('Successfully signed in')
        router.push(callbackUrl || '/')
        return { success: true }
      }

      return { success: false, error: 'Sign in failed' }
    } catch (error) {
      console.error('Sign in error:', error)
      toast.error('An unexpected error occurred')
      return { success: false, error: 'An unexpected error occurred' }
    }
  }, [router, extendedSession?.user?.id, refreshManager])

  // Sign out function
  const logout = useCallback(async (callbackUrl?: string) => {
    try {
      // Reset refresh attempts on logout
      if (extendedSession?.user?.id) {
        refreshManager.resetAttempts(extendedSession.user.id)
      }
      
      await signOut({
        callbackUrl: callbackUrl || '/auth/signin',
        redirect: true,
      })
      toast.success('Successfully signed out')
    } catch (error) {
      console.error('Sign out error:', error)
      toast.error('Error signing out')
    }
  }, [extendedSession?.user?.id, refreshManager])

  // Refresh session with attempt limiting
  const refreshSession = useCallback(async () => {
    if (!extendedSession?.user?.id) {
      console.error('No user ID available for refresh attempt tracking')
      return
    }

    const userId = extendedSession.user.id

    if (!refreshManager.canAttemptRefresh(userId)) {
      console.error('Maximum refresh attempts reached')
      toast.error('Too many refresh attempts. Please sign in again.')
      await logout()
      return
    }

    try {
      refreshManager.recordAttempt(userId)
      await update()
      
      // Reset attempts on successful refresh
      if (status === 'authenticated') {
        refreshManager.resetAttempts(userId)
      }
    } catch (error) {
      console.error('Session refresh error:', error)
      
      // Check if max attempts reached after failed refresh
      if (refreshManager.isMaxAttemptsReached(userId)) {
        toast.error('Session expired due to too many refresh attempts. Please sign in again.')
        await logout()
      }
    }
  }, [update, extendedSession?.user?.id, status, logout, refreshManager])

  // Handle authentication errors
  useEffect(() => {
    if (extendedSession?.error) {
      console.error('Authentication error:', extendedSession.error)
      
      if (extendedSession.error === 'RefreshAccessTokenError') {
        const userId = extendedSession.user?.id
        
        if (userId && refreshManager.isMaxAttemptsReached(userId)) {
          toast.error('Session expired due to too many refresh attempts. Please sign in again.')
        } else {
          toast.error('Your session has expired. Please sign in again.')
        }
        logout()
      }
    }
  }, [extendedSession?.error, extendedSession?.user?.id, logout, refreshManager])

  // Auto-redirect if not authenticated
  useEffect(() => {
    if (status === 'unauthenticated') {
      // Check if we're on the client side before accessing window
      if (typeof window !== 'undefined') {
        const currentPath = window.location.pathname
        if (!currentPath.startsWith('/auth')) {
          router.push('/auth/signin')
        }
      } else {
        // On server side, just redirect to signin
        router.push('/auth/signin')
      }
    }
  }, [status, router])

  return {
    // State
    session: extendedSession,
    user: extendedSession?.user,
    userIdentity: extendedSession?.userIdentity,
    isAuthenticated,
    isLoading,
    status,
    
    // Actions
    login,
    logout,
    refreshSession,
    
    // Role checks
    hasRole,
    hasAnyRole,
    hasAllRoles,
    
    // Permission checks
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    
    // Utilities
    accessToken: extendedSession?.accessToken,
    
    // Refresh attempt info
    refreshAttempts: extendedSession?.user?.id ? refreshManager.getAttempts(extendedSession.user.id) : 0,
    canRefresh: extendedSession?.user?.id ? refreshManager.canAttemptRefresh(extendedSession.user.id) : true,
  }
} 