"use client";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { 
  LayoutDashboard, 
  Package, 
  ShoppingCart,
  LogOut,
  User,
  ChevronLeft,
  ChevronRight,
  Menu,
  X,
  CreditCard,
  Moon,
  Sun
} from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";
import { Logout } from "@/app/actions/auth";
import { Sheet, SheetContent } from "@/components/ui/sheet";

const supplierRoutes = [
  { href: "/supplier", label: "Dashboard", icon: LayoutDashboard },
  { href: "/supplier/products", label: "Products", icon: Package },
  { href: "/supplier/orders", label: "Orders", icon: ShoppingCart },
  { href: "/supplier/account", label: "Account Management", icon: CreditCard },
];

export default function SupplierLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  // Handle mounting for theme
  useEffect(() => {
    setMounted(true);
  }, []);

  // Load collapsed state from localStorage
  useEffect(() => {
    const savedState = localStorage.getItem('supplier-sidebar-collapsed');
    if (savedState !== null) {
      setIsCollapsed(savedState === 'true');
    }
  }, []);

  // Save collapsed state to localStorage
  useEffect(() => {
    localStorage.setItem('supplier-sidebar-collapsed', isCollapsed.toString());
  }, [isCollapsed]);

  // Close mobile menu when route changes
  useEffect(() => {
    setIsMobileMenuOpen(false);
  }, [pathname]);

  const toggleSidebar = () => {
    setIsCollapsed(!isCollapsed);
  };

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  // Skip authentication check for login page
  const isLoginPage = pathname === "/supplier/login";

  useEffect(() => {
    // Don't redirect if we're on the login page
    if (isLoginPage) {
      return;
    }

    if (!isLoading && (!isAuthenticated || !user)) {
      router.push("/supplier/login?callbackUrl=/supplier");
      return;
    }

    // Check if user has supplier role
    if (!isLoading && user && !isSupplierRole(user.roles || [])) {
      // If user is admin, redirect to admin portal
      if (isAdminRole(user.roles || [])) {
        router.push("/admin");
        return;
      }
      router.push("/supplier/login?callbackUrl=/supplier");
      return;
    }
  }, [user, isAuthenticated, isLoading, router, isLoginPage]);

  // If it's the login page, render without the supplier layout wrapper
  if (isLoginPage) {
    return <>{children}</>;
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated || !user || !isSupplierRole(user.roles || [])) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Supplier Header */}
      <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 shadow-sm sticky top-0 z-40">
        <div className="max-w-7xl mx-auto">
          <div className="flex justify-between items-center h-16 px-4 sm:px-6 lg:px-8">
            {/* Mobile Menu Button */}
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleMobileMenu}
              className={cn(
                "lg:hidden hover:bg-primary/10 hover:text-primary min-w-[44px] min-h-[44px]",
                isMobileMenuOpen && "bg-primary/10 text-primary"
              )}
              aria-label={isMobileMenuOpen ? "Close menu" : "Open menu"}
              aria-expanded={isMobileMenuOpen}
            >
              {isMobileMenuOpen ? (
                <X className="h-5 w-5" />
              ) : (
                <Menu className="h-5 w-5" />
              )}
            </Button>

            <div className="flex items-center gap-3 ml-auto">
              {/* Dark Mode Toggle */}
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
                className="h-9 w-9 rounded-full hover:bg-primary/10 hover:text-primary focus:ring-2 focus:ring-primary/20 transition-all duration-200"
                aria-label="Toggle theme"
              >
                {mounted ? (
                  theme === "dark" ? (
                    <Sun className="h-4 w-4" />
                  ) : (
                    <Moon className="h-4 w-4" />
                  )
                ) : (
                  <Moon className="h-4 w-4" />
                )}
              </Button>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-9 rounded-full px-3 hover:bg-primary/10 dark:hover:bg-white/10 hover:text-primary dark:hover:text-white focus:ring-2 focus:ring-primary/20 transition-all duration-200 cursor-pointer active:scale-95 flex items-center gap-2"
                  >
                    <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center transition-all duration-200 hover:bg-primary/20 dark:hover:bg-white/20 hover:scale-105 cursor-pointer ring-2 ring-transparent hover:ring-primary/20 dark:hover:ring-white/20">
                      <User className="h-4 w-4 text-primary dark:text-white transition-transform duration-200" />
                    </div>
                    <span className="text-sm text-gray-600 dark:text-gray-300 hidden sm:block">{user.email}</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-64 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 dark:bg-gray-800">
                  <DropdownMenuItem
                    onClick={() => router.push("/supplier/profile")}
                    className="px-4 py-3 cursor-pointer focus:bg-accent"
                  >
                    <div className="flex items-center gap-3 w-full">
                      <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center flex-shrink-0">
                        <User className="h-5 w-5 text-primary" />
                      </div>
                      <div className="flex flex-col min-w-0 flex-1">
                        <p className="text-sm font-semibold truncate">{user.name || "User"}</p>
                        <p className="text-xs text-muted-foreground truncate">{user.email}</p>
                      </div>
                    </div>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => {
                      Logout(false);
                      window.location.href = "/supplier/login";
                    }}
                    className="text-destructive focus:text-destructive cursor-pointer mx-1 my-1 rounded-md"
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    Logout
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      </header>

      <div className="flex relative">
        {/* Desktop Sidebar */}
        <aside className={cn(
          "bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 min-h-[calc(100vh-4rem)] shadow-sm transition-all duration-300 relative hidden lg:block fixed left-0 top-16 bottom-0 overflow-y-auto overflow-x-visible",
          isCollapsed ? "w-16" : "w-64"
        )}>
          {/* Desktop Toggle Button */}
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleSidebar}
            className={cn(
              "absolute -right-3 top-4 z-[100] h-7 w-7 rounded-full border border-gray-200 dark:border-gray-700 dark:border-white/20 bg-white dark:bg-white dark:text-gray-900 shadow-sm hover:bg-primary/10 dark:hover:bg-white/90 hover:border-primary hover:text-primary dark:hover:text-gray-900 transition-all",
              isCollapsed && "rotate-180"
            )}
            aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
          >
            {isCollapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
          </Button>

          <nav className={cn(
            "p-4 space-y-1 transition-all duration-300",
            isCollapsed && "px-2"
          )}>
            {supplierRoutes.map((route) => {
              const Icon = route.icon;
              const isActive = pathname === route.href || 
                (route.href !== "/supplier" && pathname.startsWith(route.href));
              
              return (
                <Link
                  key={route.href}
                  href={route.href}
                  className={cn(
                    "flex items-center gap-3 rounded-lg text-sm font-medium transition-all duration-200 group relative",
                    "min-h-[44px] touch-manipulation",
                    isCollapsed ? "px-2 py-2 justify-center" : "px-4 py-3",
                    isActive
                      ? "bg-primary/10 dark:bg-primary/20 text-primary border-l-4 border-primary shadow-sm"
                      : "text-gray-600 dark:text-gray-300 hover:bg-primary/5 dark:hover:bg-primary/10 hover:text-primary active:bg-primary/10"
                  )}
                  title={isCollapsed ? route.label : undefined}
                >
                  <Icon className={cn(
                    "h-5 w-5 shrink-0",
                    isActive && "text-primary"
                  )} />
                  {!isCollapsed && (
                    <span className="whitespace-nowrap">{route.label}</span>
                  )}
                  {isCollapsed && (
                    <div className="absolute left-full ml-2 px-2 py-1 bg-gray-900 dark:bg-gray-700 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity whitespace-nowrap z-50 shadow-lg">
                      {route.label}
                    </div>
                  )}
                </Link>
              );
            })}
          </nav>
        </aside>

        {/* Mobile Sidebar Sheet */}
        <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
          <SheetContent side="left" className="w-72 p-0 sm:w-80 [&>button]:hidden dark:bg-gray-800 dark:border-gray-700">
            <div className="flex flex-col h-full">
              <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
                <div className="flex-1"></div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setIsMobileMenuOpen(false)}
                  className="h-8 w-8 ml-auto shrink-0"
                >
                  <X className="h-5 w-5" />
                </Button>
              </div>

              <nav className="flex-1 p-4 space-y-1 overflow-y-auto bg-white dark:bg-gray-800">
                {supplierRoutes.map((route) => {
                  const Icon = route.icon;
                  const isActive = pathname === route.href || 
                    (route.href !== "/supplier" && pathname.startsWith(route.href));
                  
                  return (
                    <Link
                      key={route.href}
                      href={route.href}
                      onClick={() => setIsMobileMenuOpen(false)}
                      className={cn(
                        "flex items-center gap-3 rounded-lg text-base font-medium transition-all duration-200 px-4 py-3 min-h-[48px] touch-manipulation",
                        "active:scale-[0.98]",
                        isActive
                          ? "bg-primary/10 dark:bg-primary/20 text-primary border-l-4 border-primary shadow-sm"
                          : "text-gray-700 dark:text-gray-300 hover:bg-primary/5 dark:hover:bg-primary/10 hover:text-primary active:bg-primary/10"
                      )}
                    >
                      <Icon className={cn(
                        "h-6 w-6 shrink-0",
                        isActive && "text-primary"
                      )} />
                      <span className="font-medium">{route.label}</span>
                    </Link>
                  );
                })}
              </nav>

              <div className="p-4 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
                <Button
                  variant="outline"
                  className="w-full justify-start h-12 text-base"
                  onClick={() => {
                    setIsMobileMenuOpen(false);
                    Logout(false);
                    window.location.href = "/supplier/login";
                  }}
                >
                  <LogOut className="h-5 w-5 mr-2" />
                  Logout
                </Button>
              </div>
            </div>
          </SheetContent>
        </Sheet>

        {/* Main Content */}
        <main className={cn(
          "flex-1 p-4 sm:p-6 lg:p-8 transition-all duration-300 w-full",
          isCollapsed ? "lg:ml-16" : "lg:ml-64"
        )}>
          {children}
        </main>
      </div>
    </div>
  );
}

function isSupplierRole(roles: string[]): boolean {
  return roles.includes("supplier");
}

function isAdminRole(roles: string[]): boolean {
  return roles.includes("admin") || roles.includes("superadmin");
}

