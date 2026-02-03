"use client";
import { useEffect, useState } from "react";
import { useSession, signOut } from "next-auth/react";
import { useRouter } from "next/navigation";
import { User, Mail, Phone, MapPin, Settings, LogOut, Loader2, ChevronRight } from "lucide-react";
import { getCustomerByEmail } from "@/app/actions/customer";
import { getAccessToken, Logout } from "@/app/actions/auth";
import { useAuth } from "@/context/AuthContext";

const menuItems = [
  { icon: Settings, label: 'Settings', description: 'Account preferences' },
];

interface CustomerData {
  id?: number;
  full_name: string;
  email: string;
  phone_number?: string;
  city?: string;
  region?: string;
  woreda?: string;
  status?: string;
  is_active?: boolean;
}

export default function Profile() {
  const { data: session, status } = useSession();
  const { isAuthenticated: authContextIsAuthenticated, user: authUser, isLoading: authLoading } = useAuth();
  const router = useRouter();
  const [customerData, setCustomerData] = useState<CustomerData>({
    full_name: "",
    email: "",
    phone_number: "",
    city: "",
    region: "",
    woreda: "",
    status: "",
  });
  const [isLoading, setIsLoading] = useState(true);

  // Get user from session, authContext, or localStorage
  const getUser = () => {
    if (session?.user) return session.user;
    if (authUser) return authUser;
    
    // Fallback to localStorage if available
    if (typeof window !== 'undefined') {
      const storedEmail = localStorage.getItem('saved_email') || localStorage.getItem('user_email');
      const storedName = localStorage.getItem('user_name');
      if (storedEmail) {
        return {
          name: storedName || '',
          email: storedEmail,
        };
      }
    }
    return null;
  };

  const user = getUser();

  // Redirect to login if not authenticated (only after all loading is complete)
  useEffect(() => {
    // Wait for all auth checks to complete
    if (authLoading || status === 'loading') {
      return;
    }

    // Check authentication - must have either session, token, authContext, or localStorage credentials
    const hasSession = !!session?.user;
    const hasToken = !!getAccessToken();
    
    // Also check localStorage for user_email and saved_password
    let hasLocalStorageAuth = false;
    if (typeof window !== 'undefined') {
      const userEmail = localStorage.getItem('user_email');
      const savedPassword = localStorage.getItem('saved_password');
      hasLocalStorageAuth = !!(userEmail && savedPassword);
    }
    
    const isAuth = authContextIsAuthenticated || hasSession || hasToken || hasLocalStorageAuth;

    if (!isAuth) {
      router.push("/login");
      return;
    }
  }, [authContextIsAuthenticated, authLoading, status, session, router]);

  // Load customer data
  useEffect(() => {
    const loadCustomerData = async () => {
      // Wait for all loading to complete
      if (status === "loading" || authLoading) {
        return;
      }

      // Check if user is authenticated - be more lenient
      const hasSession = !!session?.user;
      const hasToken = !!getAccessToken();
      
      // Also check localStorage for user_email and saved_password
      let hasLocalStorageAuth = false;
      if (typeof window !== 'undefined') {
        const userEmail = localStorage.getItem('user_email');
        const savedPassword = localStorage.getItem('saved_password');
        hasLocalStorageAuth = !!(userEmail && savedPassword);
      }
      
      const isAuth = authContextIsAuthenticated || hasSession || hasToken || hasLocalStorageAuth;

      if (!isAuth || !user) {
        // Don't redirect here, let the redirect effect handle it
        setIsLoading(false);
        return;
      }

      try {
        const storedName = user.name || "";
        const storedEmail = user.email || "";

        // Initialize with data from session
        const initialData: CustomerData = {
          full_name: storedName || "User",
          email: storedEmail || "",
          phone_number: "",
          city: "",
          region: "",
          woreda: "",
          status: "",
        };

        // Try to fetch customer details from API
        if (storedEmail) {
          try {
            const customer = await getCustomerByEmail(storedEmail);

            if (customer) {
              const customerDataToStore = {
                id: customer.id,
                full_name: customer.full_name || storedName,
                email: customer.email || storedEmail,
                phone_number: customer.phone_number || "",
                city: customer.city || "",
                region: customer.region || "",
                woreda: customer.woreda || "",
                status: customer.status || "",
                is_active: customer.is_active,
              };
              
              setCustomerData(customerDataToStore);
              
              // Store in localStorage for use in checkout
              if (typeof window !== 'undefined') {
                localStorage.setItem('customer_id', customer.id.toString());
                localStorage.setItem('customer_data', JSON.stringify(customerDataToStore));
              }
            } else {
              setCustomerData(initialData);
            }
          } catch (error) {
            console.error("Error fetching customer data:", error);
            setCustomerData(initialData);
          }
        } else {
          setCustomerData(initialData);
        }
      } catch (error) {
        console.error("Error loading customer data:", error);
        setCustomerData({
          full_name: user?.name || "User",
          email: user?.email || "",
          phone_number: "",
          city: "",
          region: "",
          woreda: "",
          status: "",
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadCustomerData();
  }, [user, authContextIsAuthenticated, status, authLoading, session, router]);

  const handleSignOut = () => {
    // Clear basic auth token
    Logout(true);
    // Sign out from NextAuth
    signOut({ callbackUrl: "/login" });
  };

  if (isLoading || status === "loading" || authLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading profile...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background pb-20">
      {/* Profile Header */}
      <div className="bg-card border-b border-border">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex items-center gap-4 mb-4">
            <div className="w-20 h-20 bg-primary rounded-full flex items-center justify-center">
              <User className="w-10 h-10 text-primary-foreground" />
            </div>
            <div className="flex-1">
              <h1 className="text-xl font-semibold mb-1">
                {customerData.full_name || "User"}
              </h1>
              <p className="text-sm text-muted-foreground mb-1">
                {customerData.email || "No email"}
              </p>
              {customerData.status && (
                <p className="text-xs text-muted-foreground capitalize">
                  {customerData.status}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Personal Information */}
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 mt-6">
        <h2 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide mb-3 px-2">
          Personal Information
        </h2>
        <div className="bg-card border border-border rounded-md overflow-hidden">
          <div className="flex items-center p-4 border-b border-border">
            <div className="flex items-center gap-3">
              <Mail className="w-5 h-5 text-muted-foreground" />
              <div>
                <p className="text-xs text-muted-foreground mb-0.5">Email</p>
                <p className="text-sm font-medium">
                  {customerData.email || "Not provided"}
                </p>
              </div>
            </div>
          </div>

          {/* Customer Details Section */}
          <div className="p-4 space-y-3">
            {/* Full Name */}
            {customerData.full_name && (
              <div className="flex items-start gap-3">
                <User className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">
                    Full Name
                  </p>
                  <p className="text-sm font-medium">{customerData.full_name}</p>
                </div>
              </div>
            )}

            {/* Email */}
            {customerData.email && (
              <div className="flex items-start gap-3">
                <Mail className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Email</p>
                  <p className="text-sm font-medium">{customerData.email}</p>
                </div>
              </div>
            )}

            {/* Phone Number */}
            {customerData.phone_number && (
              <div className="flex items-start gap-3">
                <Phone className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">
                    Phone Number
                  </p>
                  <p className="text-sm font-medium">
                    {customerData.phone_number}
                  </p>
                </div>
              </div>
            )}

            {/* Woreda */}
            {customerData.woreda && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Woreda</p>
                  <p className="text-sm font-medium">{customerData.woreda}</p>
                </div>
              </div>
            )}

            {/* City */}
            {customerData.city && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">City</p>
                  <p className="text-sm font-medium">{customerData.city}</p>
                </div>
              </div>
            )}

            {/* Region */}
            {customerData.region && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 shrink-0" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-muted-foreground mb-0.5">Region</p>
                  <p className="text-sm font-medium">{customerData.region}</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Menu Items */}
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 mt-6">
        <h2 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide mb-3 px-2">
          Account
        </h2>
        <div className="bg-card border border-border rounded-md overflow-hidden">
          {menuItems.map((item, index) => {
            const Icon = item.icon;
            return (
              <button
                key={item.label}
                className={`w-full flex items-center justify-between p-4 transition-colors hover:bg-secondary ${
                  index !== menuItems.length - 1 ? 'border-b border-border' : ''
                }`}
              >
                <div className="flex items-center gap-3">
                  <Icon className="w-5 h-5 text-muted-foreground" />
                  <div className="text-left">
                    <p className="text-sm font-medium">
                      {item.label}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {item.description}
                    </p>
                  </div>
                </div>
                <ChevronRight className="w-4 h-4 text-muted-foreground" />
              </button>
            );
          })}
        </div>
      </div>

      {/* Sign Out */}
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 mt-6">
        <button 
          onClick={handleSignOut}
          className="w-full flex items-center justify-center gap-2 p-4 text-sm text-destructive hover:bg-destructive/10 border border-destructive/20 rounded-md transition-colors"
        >
          <LogOut className="w-4 h-4" />
          Sign Out
        </button>
      </div>

      {/* App Info */}
      <div className="text-center mt-8 mb-4 text-xs text-muted-foreground">
        <p>Version 1.0.0</p>
      </div>
    </div>
  );
}
