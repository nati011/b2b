"use client";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { signIn } from "next-auth/react";
import { PiSpinner } from "react-icons/pi";
import { toast } from "sonner";
import Image from "next/image";

import { InitResetPassword, Login, getSavedCredentials } from "@/app/actions/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Eye, EyeOff } from "lucide-react";

interface Errors {
  email?: string;
  password?: string;
  general?: string;
}

export default function SupplierLoginPage() {
  const searchParams = useSearchParams();
  const [errors, setErrors] = useState<Errors>({});
  const [loading, setLoading] = useState(false);
  const callbackUrl = searchParams.get("callbackUrl") || "/supplier";
  const error = searchParams.get("error");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const router = useRouter();

  // Prefetch supplier routes on mount
  useEffect(() => {
    router.prefetch("/supplier");
    router.prefetch("/supplier/products");
    router.prefetch("/supplier/orders");
  }, [router]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setErrors({});

    if (!email) {
      setErrors((prev) => ({ ...prev, email: "Email is required" }));
      return;
    }
    if (!password) {
      setErrors((prev) => ({ ...prev, password: "Password is required" }));
      return;
    }

    try {
      setLoading(true);
      
      // Try direct login with basic auth first
      try {
        await Login({ email, password });
        // Login successful - token is stored in localStorage
        // Store credentials in localStorage for future use
        if (typeof window !== 'undefined') {
          localStorage.setItem('saved_email', email);
          localStorage.setItem('saved_password', password);
          localStorage.setItem('user_email', email);
          // Update login status
          setIsLoggedIn(true);
          
          // Check user roles and redirect based on role
          const userRolesStr = localStorage.getItem('user_roles');
          if (userRolesStr) {
            try {
              const userRoles: string[] = JSON.parse(userRolesStr);
              
              // Redirect based on role
              if (userRoles.includes('supplier')) {
                toast.success("Supplier login successful!");
                router.push('/supplier');
                return;
              } else if (userRoles.includes('admin') || userRoles.includes('superadmin')) {
                // Admin should go to admin portal
                toast.success("Login successful!");
                router.push('/admin');
                return;
              } else {
                // User doesn't have supplier role
                setErrors((prev) => ({ ...prev, general: "Access denied. Supplier privileges required." }));
                setTimeout(() => {
                  setErrors({});
                }, 5000);
                return;
              }
            } catch (e) {
              console.error('Error parsing user roles:', e);
            }
          } else {
            // No roles found, deny access
            setErrors((prev) => ({ ...prev, general: "Access denied. Admin privileges required." }));
            setTimeout(() => {
              setErrors({});
            }, 5000);
            return;
          }
        }
        toast.success("Admin login successful!");
        router.push(callbackUrl);
      } catch (loginError: any) {
        // If direct login fails, fall back to NextAuth (for backward compatibility)
        console.log("Direct login failed, trying NextAuth:", loginError);
        const result = await signIn("credentials", {
          email,
          password,
          redirect: false,
          callbackUrl,
        });
        
        if (result?.error) {
          setErrors((prev) => ({ ...prev, general: loginError.message || "Invalid Credentials" }));
          setTimeout(() => {
            setErrors({});
          }, 5000);
        } else if (result?.ok) {
          // Store credentials for session use when using NextAuth
          if (typeof window !== 'undefined') {
            localStorage.setItem('saved_email', email);
            localStorage.setItem('saved_password', password);
            localStorage.setItem('user_email', email);
            // Update login status
            setIsLoggedIn(true);
            // Dispatch custom event to notify AuthContext of login
            window.dispatchEvent(new Event('auth-state-changed'));
            
            // For NextAuth, decode the token from the response to get roles
            // The token should be in the session, but we can also check localStorage
            // Wait a moment for the session to be established, then check roles
            setTimeout(async () => {
              // Try to get roles from localStorage first (if Login function stored them)
              let userRoles: string[] = [];
              const userRolesStr = localStorage.getItem('user_roles');
              if (userRolesStr) {
                try {
                  userRoles = JSON.parse(userRolesStr);
                } catch (e) {
                  console.error('Error parsing user roles from localStorage:', e);
                }
              }
              
              // If no roles in localStorage, try to get from session
              if (userRoles.length === 0) {
                try {
                  const { getSession } = await import('next-auth/react');
                  const session = await getSession();
                  if (session?.user && (session.user as any).roles) {
                    userRoles = (session.user as any).roles;
                    // Store in localStorage for consistency
                    localStorage.setItem('user_roles', JSON.stringify(userRoles));
                  }
                } catch (e) {
                  console.error('Error getting session:', e);
                }
              }
              
              // Redirect based on role
              if (userRoles.includes('admin') || userRoles.includes('superadmin')) {
                router.push('/admin');
                return;
              } else if (userRoles.includes('supplier')) {
                router.push('/supplier');
                return;
              } else {
                // User doesn't have admin role
                setErrors((prev) => ({ ...prev, general: "Access denied. Admin privileges required." }));
                setTimeout(() => {
                  setErrors({});
                }, 5000);
                return;
              }
            }, 200);
          } else {
            router.push(callbackUrl);
          }
        }
      }
    } catch (error: any) {
      setErrors((prev) => ({ ...prev, general: error.message || "An error occurred during login" }));
      setTimeout(() => {
        setErrors({});
      }, 5000);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (error) {
      setErrors((prev) => ({ ...prev, general: "Invalid Credentials" }));
    }
  }, [error]);

  // Check if user is logged in based on user_email and saved_password in localStorage
  const checkLoginStatus = () => {
    if (typeof window !== 'undefined') {
      const userEmail = localStorage.getItem('user_email');
      const savedPassword = localStorage.getItem('saved_password');
      const isLoggedInStatus = !!(userEmail && savedPassword);
      setIsLoggedIn(isLoggedInStatus);
      return isLoggedInStatus;
    }
    return false;
  };

  // Check login status on mount and when localStorage changes
  useEffect(() => {
    checkLoginStatus();

    // Listen for storage changes (e.g., from other tabs or after login)
    const handleStorageChange = () => {
      checkLoginStatus();
    };

    // Listen for custom auth state change event
    const handleAuthStateChange = () => {
      checkLoginStatus();
    };

    window.addEventListener('storage', handleStorageChange);
    window.addEventListener('auth-state-changed', handleAuthStateChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('auth-state-changed', handleAuthStateChange);
    };
  }, []);

  // Redirect away from login page if user is already logged in and has admin role
  useEffect(() => {
    // Check for access_token to ensure user is actually authenticated
    const accessToken = localStorage.getItem('access_token');
    if (isLoggedIn && accessToken) {
      const userRolesStr = localStorage.getItem('user_roles');
      if (userRolesStr) {
        try {
          const userRoles: string[] = JSON.parse(userRolesStr);
          
          if (userRoles.includes('supplier')) {
            router.push('/supplier');
            return;
          } else if (userRoles.includes('admin') || userRoles.includes('superadmin')) {
            router.push('/admin');
            return;
          }
        } catch (e) {
          console.error('Error parsing user roles:', e);
        }
      }
    }
  }, [isLoggedIn, router]);

  // Load saved credentials on mount
  useEffect(() => {
    const savedCredentials = getSavedCredentials();
    if (savedCredentials) {
      setEmail(savedCredentials.email);
      setPassword(savedCredentials.password);
    }
  }, []);

  const handleInitResetPassword = async () => {
    try {
      const success = await InitResetPassword(email);
      toast.success(success);
      router.push("/check-your-email");
    } catch (error: any) {
      toast.error(error.message);
      router.push("/check-your-email");
    }
  };

  return (
    <div className="flex flex-col gap-4 p-6 md:p-10 my-32">
      <div className="flex flex-1 items-center justify-center">
        <div className="w-full max-w-lg">
          <form className="flex flex-col gap-6" onSubmit={handleSubmit}>
            <div className="flex flex-col items-center gap-2 text-left">
              <div className="flex items-center gap-3 mb-2">
                <Image 
                  src="/logo.svg" 
                  alt="Efoyeta Logo" 
                  width={32} 
                  height={32} 
                  className="h-8 w-8"
                />
                <h1 className="text-2xl font-bold">Supplier Login</h1>
              </div>
              <p className="text-balance text-sm text-muted-foreground">
                Enter your supplier credentials to access the supplier portal
              </p>
            </div>
            <div className="grid gap-6">
              {errors.general && (
                <div className="bg-red-100/[0.2] rounded-md p-2 border border-red-900">
                  <p className="text-red-900 text-center font-medium">
                    {errors.general}
                  </p>
                </div>
              )}
              <div className="grid gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="admin@example.com"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value.toLowerCase())}
                  onBlur={(e) => {
                    if (!e.target.value) {
                      setErrors((prev) => ({
                        ...prev,
                        email: "Email is required",
                      }));
                    } else {
                      setErrors((prev) => ({ ...prev, email: "" }));
                    }
                  }}
                />
              </div>
              <div className="grid gap-2">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <div
                    onClick={() => {
                      handleInitResetPassword();
                    }}
                    className="ml-auto text-sm underline-offset-4 hover:underline cursor-pointer"
                  >
                    <p>Forgot your password?</p>
                  </div>
                </div>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="********"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    onBlur={(e) => {
                      if (!e.target.value) {
                        setErrors((prev) => ({
                          ...prev,
                          password: "Password is required",
                        }));
                      } else {
                        setErrors((prev) => ({ ...prev, password: "" }));
                      }
                    }}
                    className="pr-10"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                    aria-label={showPassword ? "Hide password" : "Show password"}
                  >
                    {showPassword ? (
                      <EyeOff className="h-5 w-5" />
                    ) : (
                      <Eye className="h-5 w-5" />
                    )}
                  </button>
                </div>
              </div>
              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? (
                  <div className="flex gap-2">
                    <PiSpinner className="animate-spin" />
                    <span>Loading</span>
                  </div>
                ) : (
                  <span>Login to Supplier Portal</span>
                )}
              </Button>
            </div>
            <div className="text-center text-sm">
              <p className="text-muted-foreground">
                Admin user?{" "}
                <a href="/admin/login" className="underline underline-offset-4 text-primary">
                  Login here
                </a>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

