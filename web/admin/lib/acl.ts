// Centralized ACL utilities for admin app
import { ACL } from "./constants";

export type PermissionLike = { Name?: string; name?: string; Resource?: string; resource?: string } | string;

// Public routes that never require auth
const PUBLIC_PATHS = [
  "/auth/signin",
  "/auth/signup",
  "/auth/error",
  "/auth/verify",
  "/forbidden",
  "/404",
  "/500",
];

export function isPublicPath(pathname: string): boolean {
  return PUBLIC_PATHS.some((p) => pathname === p || pathname.startsWith(p));
}

export function isProtectedPath(pathname: string): boolean {
  if (isPublicPath(pathname)) return false;
  // If an explicit ACL entry exists or the path is not public, treat as protected
  return Boolean(matchAclKey(pathname));
}

export function normalizePath(pathname: string): string {
  if (!pathname) return "/";
  // replace numeric segments and common UUID-like segments with {param}
  const segments = pathname.split("/").filter(Boolean);
  const normalized = segments.map((seg) => {
    if (/^\d+$/.test(seg)) return "{param}";
    if (/^[0-9a-fA-F-]{8,}$/.test(seg)) return "{param}";
    return seg;
  });
  return "/" + normalized.join("/");
}

function matchAclKey(pathname: string): string | null {
  const normalized = normalizePath(pathname);
  if (ACL[normalized]) return normalized;
  // direct match fallback
  if (ACL[pathname]) return pathname;
  // try prefix matches for nested routes
  const keys = Object.keys(ACL).sort((a, b) => b.length - a.length);
  for (const key of keys) {
    if (normalized.startsWith(key)) return key;
  }
  return null;
}

export function extractPermissionNames(permissions: any): string[] {
  if (!permissions) return [];
  // Support both flat arrays and { List: [] }
  const list = Array.isArray(permissions)
    ? permissions
    : Array.isArray(permissions.List)
    ? permissions.List
    : [];
  return list
    .map((p: PermissionLike) => {
      if (typeof p === 'string') return p;
      return p?.Name || p?.name;
    })
    .filter((n: string | undefined): n is string => Boolean(n));
}

export function hasAnyPermission(permissions: any, required: string[]): boolean {
  const names = extractPermissionNames(permissions);
  if (required.length === 0) return true;
  return required.some((r) => names.includes(r));
}

export function hasAllPermissions(permissions: any, required: string[]): boolean {
  const names = extractPermissionNames(permissions);
  if (required.length === 0) return true;
  return required.every((r) => names.includes(r));
}

export function requiredPermissionsForPath(pathname: string): string[] {
  const key = matchAclKey(pathname);
  if (!key) return [];
  return ACL[key] || [];
}

export function canAccessPath(permissions: any, pathname: string): boolean {
  console.log(permissions)
  console.log(pathname)
  if (isPublicPath(pathname)) return true;
  const required = requiredPermissionsForPath(pathname);
  if (required.length === 0) return true;
  console.log(required)
  console.log(hasAnyPermission(permissions, required))
  return hasAnyPermission(permissions, required);
}


