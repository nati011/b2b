"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { User, Mail, Loader2, Shield } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { Logout } from "@/app/actions/auth";

export default function AdminProfile() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && (!isAuthenticated || !user)) {
      router.push("/admin/login?callbackUrl=/admin/profile");
      return;
    }

    // Check if user has admin role
    if (!isLoading && user && !isAdmin(user.roles || [])) {
      router.push("/admin/login?callbackUrl=/admin/profile");
      return;
    }
  }, [user, isAuthenticated, isLoading, router]);

  function isAdmin(roles: string[]): boolean {
    // Only allow admin role (not supplier)
    // Check both lowercase and original case for compatibility
    const normalizedRoles = roles.map(r => r.toLowerCase());
    return normalizedRoles.includes("admin") || normalizedRoles.includes("superadmin");
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading profile...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Profile</h1>
        <p className="text-sm text-muted-foreground mt-1">View and manage your account information</p>
      </div>

      {/* Profile Header Card */}
      <div className="bg-card border border-border rounded-lg shadow-sm mb-6">
        <div className="p-6">
          <div className="flex items-center gap-4">
            <div className="h-20 w-20 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
              <User className="h-10 w-10 text-primary" />
            </div>
            <div className="flex-1 min-w-0">
              <h2 className="text-xl font-semibold mb-1 truncate">
                {user.name || "User"}
              </h2>
              <p className="text-sm text-muted-foreground mb-2 truncate">
                {user.email || "No email"}
              </p>
              {user.roles && user.roles.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {user.roles.map((role) => (
                    <span
                      key={role}
                      className="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-md bg-primary/10 text-primary"
                    >
                      <Shield className="h-3 w-3" />
                      {role}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Personal Information */}
      <div className="bg-card border border-border rounded-lg shadow-sm">
        <div className="p-6">
          <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide mb-4">
            Personal Information
          </h3>
          <div className="space-y-4">
            {/* Name */}
            <div className="flex items-start gap-3">
              <User className="w-5 h-5 text-muted-foreground mt-0.5 shrink-0" />
              <div className="flex-1 min-w-0">
                <p className="text-xs text-muted-foreground mb-1">Full Name</p>
                <p className="text-sm font-medium">{user.name || "Not provided"}</p>
              </div>
            </div>

            {/* Email */}
            <div className="flex items-start gap-3">
              <Mail className="w-5 h-5 text-muted-foreground mt-0.5 shrink-0" />
              <div className="flex-1 min-w-0">
                <p className="text-xs text-muted-foreground mb-1">Email</p>
                <p className="text-sm font-medium">{user.email || "Not provided"}</p>
              </div>
            </div>

            {/* User ID */}
            {user.id && (
              <div className="flex items-start gap-3">
                <Shield className="w-5 h-5 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-1">User ID</p>
                  <p className="text-sm font-medium font-mono">{user.id}</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

