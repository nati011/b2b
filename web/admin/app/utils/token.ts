import { jwtDecode } from "jwt-decode"

// Token types
export interface TokenPayload {
  exp: number
  iat: number
  sub: string
  email: string
  name: string
  preferred_username: string
  given_name: string
  family_name: string
  realm_access: {
    roles: string[]
  }
}

export interface UserInfo {
  id: string
  name: string
  email: string
  username: string
  roles: string[]
}

// Token utility functions
export function decodeToken(token: string): TokenPayload {
  try {
    return jwtDecode<TokenPayload>(token)
  } catch (error) {
    throw new Error("Invalid token format")
  }
}

export function isTokenExpired(token: string): boolean {
  try {
    const decoded = decodeToken(token)
    const currentTime = Math.floor(Date.now() / 1000)
    return decoded.exp < currentTime
  } catch (error) {
    return true // Consider invalid tokens as expired
  }
}

export function getTokenExpirationTime(token: string): number {
  try {
    const decoded = decodeToken(token)
    return decoded.exp * 1000 // Convert to milliseconds
  } catch (error) {
    return 0
  }
}

export function extractUserFromToken(token: string): UserInfo {
  const decoded = decodeToken(token)
  
  return {
    id: decoded.sub,
    name: decoded.name || `${decoded.given_name || ''} ${decoded.family_name || ''}`.trim(),
    email: decoded.email,
    username: decoded.preferred_username,
    roles: decoded.realm_access?.roles || []
  }
}

export function hasRole(token: string, role: string): boolean {
  try {
    const decoded = decodeToken(token)
    return decoded.realm_access?.roles?.includes(role) || false
  } catch (error) {
    return false
  }
}

export function hasAnyRole(token: string, roles: string[]): boolean {
  try {
    const decoded = decodeToken(token)
    const userRoles = decoded.realm_access?.roles || []
    return roles.some(role => userRoles.includes(role))
  } catch (error) {
    return false
  }
}

export function hasAllRoles(token: string, roles: string[]): boolean {
  try {
    const decoded = decodeToken(token)
    const userRoles = decoded.realm_access?.roles || []
    return roles.every(role => userRoles.includes(role))
  } catch (error) {
    return false
  }
}

// Token refresh utilities
export async function refreshToken(refreshToken: string, apiUrl: string): Promise<{
  accessToken: string
  refreshToken: string
  user: UserInfo
}> {
  const response = await fetch(`${apiUrl}/api/v1/auth/refresh`, {
    method: 'POST',
    body: JSON.stringify({ refresh_token: refreshToken }),
  })

  if (!response.ok) {
    throw new Error(`Token refresh failed: ${response.status}`)
  }

  const data = await response.json()
  
  if (!data.body?.access_token) {
    throw new Error('Invalid refresh response')
  }

  const user = extractUserFromToken(data.body.access_token)
  
  return {
    accessToken: data.body.access_token,
    refreshToken: data.body.refresh_token || refreshToken,
    user
  }
}

// Refresh attempt management
export class RefreshAttemptManager {
  private static instance: RefreshAttemptManager
  private attempts: Map<string, number> = new Map()
  private readonly MAX_ATTEMPTS = 3
  private readonly RESET_INTERVAL = 5 * 60 * 1000 // 5 minutes

  private constructor() {
    // Reset attempts periodically
    setInterval(() => {
      this.attempts.clear()
    }, this.RESET_INTERVAL)
  }

  static getInstance(): RefreshAttemptManager {
    if (!RefreshAttemptManager.instance) {
      RefreshAttemptManager.instance = new RefreshAttemptManager()
    }
    return RefreshAttemptManager.instance
  }

  canAttemptRefresh(userId: string): boolean {
    const currentAttempts = this.attempts.get(userId) || 0
    return currentAttempts < this.MAX_ATTEMPTS
  }

  recordAttempt(userId: string): void {
    const currentAttempts = this.attempts.get(userId) || 0
    this.attempts.set(userId, currentAttempts + 1)
  }

  resetAttempts(userId: string): void {
    this.attempts.delete(userId)
  }

  getAttempts(userId: string): number {
    return this.attempts.get(userId) || 0
  }

  isMaxAttemptsReached(userId: string): boolean {
    return this.getAttempts(userId) >= this.MAX_ATTEMPTS
  }
} 