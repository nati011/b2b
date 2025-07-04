import axiosIns from '@/app/libs/axios'
import { UserIdentity } from '@/app/libs/types'

// Cache for user identity to avoid unnecessary API calls
class UserIdentityManager {
  private static instance: UserIdentityManager
  private cache: Map<string, { data: UserIdentity; timestamp: number }> = new Map()
  private readonly CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  private constructor() {}

  static getInstance(): UserIdentityManager {
    if (!UserIdentityManager.instance) {
      UserIdentityManager.instance = new UserIdentityManager()
    }
    return UserIdentityManager.instance
  }

  // Get user identity from cache or fetch from API
  async getUserIdentity(userId: string, forceRefresh = false): Promise<UserIdentity | null> {
    const cached = this.cache.get(userId)
    const now = Date.now()

    // Return cached data if it's still valid and not forcing refresh
    if (!forceRefresh && cached && (now - cached.timestamp) < this.CACHE_DURATION) {
      return cached.data
    }

    try {
      const response = await axiosIns.get<{ user: UserIdentity }>('/identity/user')
      
      if (response.data?.user) {
        // Cache the new data
        this.cache.set(userId, {
          data: response.data.user,
          timestamp: now
        })
        return response.data.user
      }

      return null
    } catch (error) {
      console.error('Error fetching user identity:', error)
      
      // Return cached data if available, even if expired
      if (cached) {
        console.warn('Using expired cached user identity due to API error')
        return cached.data
      }
      
      return null
    }
  }

  // Clear cache for a specific user
  clearCache(userId: string): void {
    this.cache.delete(userId)
  }

  // Clear all cache
  clearAllCache(): void {
    this.cache.clear()
  }

  // Update cache with new data
  updateCache(userId: string, userIdentity: UserIdentity): void {
    this.cache.set(userId, {
      data: userIdentity,
      timestamp: Date.now()
    })
  }

  // Check if cache is valid for a user
  isCacheValid(userId: string): boolean {
    const cached = this.cache.get(userId)
    if (!cached) return false
    
    const now = Date.now()
    return (now - cached.timestamp) < this.CACHE_DURATION
  }

  // Get cache info
  getCacheInfo(userId: string): { isValid: boolean; age: number } | null {
    const cached = this.cache.get(userId)
    if (!cached) return null
    
    const now = Date.now()
    const age = now - cached.timestamp
    
    return {
      isValid: age < this.CACHE_DURATION,
      age
    }
  }
}

// Export singleton instance
export const userIdentityManager = UserIdentityManager.getInstance()

// Utility functions
export async function refreshUserIdentity(userId: string): Promise<UserIdentity | null> {
  return await userIdentityManager.getUserIdentity(userId, true)
}

export function clearUserIdentityCache(userId?: string): void {
  if (userId) {
    userIdentityManager.clearCache(userId)
  } else {
    userIdentityManager.clearAllCache()
  }
}

export function updateUserIdentityCache(userId: string, userIdentity: UserIdentity): void {
  userIdentityManager.updateCache(userId, userIdentity)
}

export function isUserIdentityCacheValid(userId: string): boolean {
  return userIdentityManager.isCacheValid(userId)
} 