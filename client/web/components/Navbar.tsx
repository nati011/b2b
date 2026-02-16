'use client'
import { ShoppingBag, Menu, X, Search, LogIn, CircleUser, Shield } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState, useEffect, useCallback } from "react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import useCartStore from "@/lib/store/useCartStore";
import useSearchStore from "@/lib/store/useSearchStore";
import Image from "next/image";
import { useAuth } from "@/context/AuthContext";

export const Navbar = () => {
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [searchOpen, setSearchOpen] = useState(false);
    const [searchInput, setSearchInput] = useState('');
    const { isAuthenticated, user } = useAuth();
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [isMounted, setIsMounted] = useState(false);
    const router = useRouter();
    const pathname = usePathname();
    const setSearchQuery = useSearchStore((state) => state.setSearchQuery);
    const searchQuery = useSearchStore((state) => state.searchQuery);

    const toggleSearch = () => {
        setSearchOpen(!searchOpen);
        if (!searchOpen) {
            // When opening, populate with current search query
            setSearchInput(searchQuery);
        } else {
            // When closing, only clear the input field, keep the query for filtering
            setSearchInput('');
        }
    };

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        const trimmedQuery = searchInput.trim();
        setSearchQuery(trimmedQuery);
        setSearchOpen(false);
        // Navigate to products page if not already there
        if (pathname !== '/product') {
            router.push('/product');
        }
    };

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setSearchInput(value);
        // Update search query in real-time for immediate filtering
        setSearchQuery(value.trim());
        // Navigate to products page if not already there and user is typing
        if (value.trim() && pathname !== '/product') {
            router.push('/product');
        }
    };

    // Only subscribe to totalItems to prevent unnecessary re-renders
    const totalItems = useCartStore((state) => state.totalItems)

    const navLinks = [
        { name: "Home", href: "/" },
        { name: "Products", href: "/product" },
        { name: "Orders", href: "/orders" },
        { name: "About Us", href: "/about" },
        { name: "Contact", href: "/contact" },
        { name: "Return & Refund", href: "/return-refund-policy" },
    ];

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
        // Mark as mounted to prevent hydration mismatch
        setIsMounted(true);
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

    // Prefetch all navigation links on mount and hover
    useEffect(() => {
        navLinks.forEach((link) => {
            router.prefetch(link.href);
        });
        // Prefetch common authenticated routes
        if (isAuthenticated || isLoggedIn) {
            router.prefetch("/profile");
            router.prefetch("/orders");
        }
    }, [router, isAuthenticated, isLoggedIn]);

    // Handle link hover for prefetching
    const handleLinkHover = useCallback((href: string) => {
        router.prefetch(href);
    }, [router]);

    // Check if user is admin
    const isAdmin = useCallback(() => {
        return user?.roles && (user.roles.includes("admin") || user.roles.includes("superadmin"));
    }, [user]);

    // Check if user is supplier
    const isSupplier = useCallback(() => {
        return user?.roles && user.roles.includes("supplier");
    }, [user]);

    return (
        <header className="sticky top-0 z-50 w-full bg-white/95 backdrop-blur-sm border-b border-gray-100 print:hidden">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-20">
                    {/* Logo */}
                    <Link 
                        href="/" 
                        className="flex items-center flex-shrink-0"
                        prefetch={true}
                    >
                        <Image 
                            src='/logo.png' 
                            width={180} 
                            height={90} 
                            alt="Efoyeta Store Logo" 
                            className="h-16 md:h-20 w-auto"
                            priority
                            loading="eager"
                        />
                        <span className="sr-only">EFOYETA STORE</span>
                    </Link>

                    {/* Navigation - Centered */}
                    <nav className="flex items-center justify-center flex-1 px-4 sm:px-6 lg:px-8">
                        <div className="hidden lg:flex items-center space-x-1">
                            {navLinks.map((link) => (
                                <Link
                                    key={link.name}
                                    href={link.href}
                                    prefetch={true}
                                    onMouseEnter={() => handleLinkHover(link.href)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 transition-colors duration-200 rounded-md hover:bg-gray-50"
                                >
                                    {link.name}
                                </Link>
                            ))}
                        </div>
                    </nav>

                    {/* Right Side Actions */}
                    <div className="flex items-center gap-1 md:gap-2 flex-shrink-0">
                        {/* Search */}
                        {searchOpen ? (
                            <div className="absolute inset-0 left-0 right-0 px-4 flex items-center justify-center bg-white z-50 border-b border-gray-200">
                                <div className="w-full max-w-2xl">
                                    <form onSubmit={handleSearch} className="flex items-center gap-2">
                                        <div className="relative flex-grow">
                                            <Input
                                                placeholder="Search products..."
                                                className="w-full pr-10"
                                                autoFocus
                                                value={searchInput}
                                                onChange={handleSearchChange}
                                            />
                                            {searchInput && (
                                                <Button
                                                    type="button"
                                                    variant="ghost"
                                                    size="icon"
                                                    className="absolute right-0 top-0 h-full px-2 hover:bg-transparent"
                                                    onClick={() => {
                                                        setSearchInput('');
                                                        setSearchQuery('');
                                                    }}
                                                >
                                                    <X className="h-4 w-4 text-gray-500" />
                                                </Button>
                                            )}
                                        </div>
                                        <Button
                                            type="submit"
                                            variant="ghost"
                                            size="icon"
                                        >
                                            <Search className="h-5 w-5" />
                                        </Button>
                                        <Button
                                            type="button"
                                            variant="ghost"
                                            size="icon"
                                            onClick={toggleSearch}
                                        >
                                            <X className="h-5 w-5" />
                                        </Button>
                                    </form>
                                </div>
                            </div>
                        ) : (
                            <>
                                <Button
                                    onClick={toggleSearch}
                                    variant="ghost"
                                    size="icon"
                                    className="hidden md:flex text-gray-600 hover:text-gray-900"
                                >
                                    <Search className="h-5 w-5" />
                                </Button>

                                {/* Cart */}
                                <Link 
                                    href="/cart" 
                                    className="relative" 
                                    prefetch={true}
                                    onMouseEnter={() => handleLinkHover("/cart")}
                                >
                                    <Button 
                                        variant="ghost" 
                                        size="icon" 
                                        className="text-gray-600 hover:text-gray-900"
                                    >
                                        <ShoppingBag className="h-5 w-5" />
                                        {totalItems > 0 && (
                                            <span className="absolute -top-1 -right-1 bg-primary text-white rounded-full w-5 h-5 flex items-center justify-center text-xs font-semibold animate-in fade-in zoom-in duration-200">
                                                {totalItems}
                                            </span>
                                        )}
                                    </Button>
                                </Link>

                                {/* Admin Panel Link */}
                                {isMounted && (isAuthenticated || isLoggedIn) && isAdmin() && (
                                    <Link 
                                        href="/admin" 
                                        className="relative" 
                                        prefetch={true}
                                        onMouseEnter={() => handleLinkHover("/admin")}
                                    >
                                        <Button 
                                            variant="ghost" 
                                            size="icon"
                                            className="text-gray-600 hover:text-gray-900"
                                            title="Admin Panel"
                                        >
                                            <Shield className="h-5 w-5" />
                                        </Button>
                                    </Link>
                                )}
                                {/* Supplier Portal Link */}
                                {isMounted && (isAuthenticated || isLoggedIn) && isSupplier() && (
                                    <Link 
                                        href="/supplier" 
                                        className="relative" 
                                        prefetch={true}
                                        onMouseEnter={() => handleLinkHover("/supplier")}
                                    >
                                        <Button 
                                            variant="ghost" 
                                            size="icon"
                                            className="text-gray-600 hover:text-gray-900"
                                            title="Supplier Portal"
                                        >
                                            <Shield className="h-5 w-5" />
                                        </Button>
                                    </Link>
                                )}

                                {/* User Menu */}
                                {/* Only check auth state after mount to prevent hydration mismatch */}
                                {isMounted && (isAuthenticated || isLoggedIn) ? (
                                    <Button 
                                        variant="ghost" 
                                        size="icon"
                                        className="text-gray-600 hover:text-gray-900"
                                        title="Profile"
                                        onClick={() => router.push("/profile")}
                                    >
                                        <CircleUser className="h-5 w-5" />
                                    </Button>
                                ) : (
                                    <Link 
                                        href="/login"
                                        prefetch={true}
                                        onMouseEnter={() => handleLinkHover("/login")}
                                    >
                                        <Button 
                                            variant="ghost" 
                                            size="sm"
                                            className="text-gray-700 hover:text-gray-900 hover:bg-gray-50 hidden md:flex font-medium"
                                        >
                                            Sign In
                                        </Button>
                                        <Button 
                                            variant="ghost" 
                                            size="icon"
                                            className="md:hidden text-gray-600 hover:text-gray-900"
                                        >
                                            <LogIn className="h-5 w-5" />
                                        </Button>
                                    </Link>
                                )}

                                {/* Mobile Menu */}
                                <Sheet open={mobileMenuOpen} onOpenChange={setMobileMenuOpen}>
                                    <SheetTrigger asChild>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="lg:hidden text-gray-600 hover:text-gray-900"
                                        >
                                            <Menu className="h-5 w-5" />
                                            <span className="sr-only">Open menu</span>
                                        </Button>
                                    </SheetTrigger>
                                    <SheetContent side="right" className="w-80">
                                        <div className="flex flex-col space-y-4 mt-8">
                                        {navLinks.map((link) => (
                                            <Link
                                                key={link.name}
                                                href={link.href}
                                                prefetch={true}
                                                onMouseEnter={() => handleLinkHover(link.href)}
                                                className="text-lg font-medium text-gray-700 hover:text-primary transition-colors py-2"
                                                onClick={() => setMobileMenuOpen(false)}
                                            >
                                                {link.name}
                                            </Link>
                                        ))}
                                        {isMounted && (isAuthenticated || isLoggedIn) && isAdmin() && (
                                            <Link 
                                                href="/admin" 
                                                prefetch={true}
                                                onMouseEnter={() => handleLinkHover("/admin")}
                                                onClick={() => setMobileMenuOpen(false)}
                                                className="text-lg font-medium text-gray-700 hover:text-primary transition-colors py-2 flex items-center gap-2"
                                            >
                                                <Shield className="h-5 w-5" />
                                                Admin Panel
                                            </Link>
                                        )}
                                        {isMounted && (isAuthenticated || isLoggedIn) && isSupplier() && (
                                            <Link 
                                                href="/supplier" 
                                                prefetch={true}
                                                onMouseEnter={() => handleLinkHover("/supplier")}
                                                onClick={() => setMobileMenuOpen(false)}
                                                className="text-lg font-medium text-gray-700 hover:text-primary transition-colors py-2 flex items-center gap-2"
                                            >
                                                <Shield className="h-5 w-5" />
                                                Supplier Portal
                                            </Link>
                                        )}
                                        {isMounted && !(isAuthenticated || isLoggedIn) && (
                                            <Link 
                                                href="/login" 
                                                prefetch={true}
                                                onMouseEnter={() => handleLinkHover("/login")}
                                                onClick={() => setMobileMenuOpen(false)}
                                            >
                                                <Button className="w-full justify-start" variant="outline">
                                                    Sign In
                                                </Button>
                                            </Link>
                                        )}
                                            <div className="pt-4 border-t border-gray-200">
                                                <Button
                                                    onClick={toggleSearch}
                                                    variant="outline"
                                                    className="flex items-center justify-start w-full"
                                                >
                                                    <Search className="h-5 w-5 mr-2" />
                                                    Search Products
                                                </Button>
                                            </div>
                                        </div>
                                    </SheetContent>
                                </Sheet>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </header>
    );
};