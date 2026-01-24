'use client'
import { ShoppingBag, Menu, X, Search, LogIn, CircleUser, ShoppingBagIcon, LogOut, Settings } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import Link from "next/link";
import useCartStore from "@/lib/store/useCartStore";
import { signOut, useSession } from "next-auth/react";
import Image from "next/image";

export const Navbar = () => {
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [searchOpen, setSearchOpen] = useState(false);
    const { data: session, status } = useSession();

    const user = session?.user;
    const user_initials = user?.name
        ?.split(" ")
        .map((n) => n[0])
        .join("");

    const toggleSearch = () => {
        setSearchOpen(!searchOpen);
    };

    // Only subscribe to totalItems to prevent unnecessary re-renders
    const totalItems = useCartStore((state) => state.totalItems)

    const navLinks = [
        { name: "Home", href: "/" },
        { name: "Shop", href: "/product" },
        { name: "Partner Portal", href: "/signup/supplier" },
        { name: "Contact", href: "/contact" },
    ];

    return (
        <header className="sticky top-0 z-50 w-full bg-white/95 backdrop-blur-sm border-b border-gray-100 print:hidden">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-20">
                    {/* Logo */}
                    <Link href="/" className="flex items-center flex-shrink-0">
                        <Image 
                            src='/logo.png' 
                            width={180} 
                            height={90} 
                            alt="Efoyeta Store Logo" 
                            className="h-16 md:h-20 w-auto"
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
                                    <div className="flex items-center gap-2">
                                        <Input
                                            placeholder="Search products..."
                                            className="flex-grow"
                                            autoFocus
                                        />
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={toggleSearch}
                                        >
                                            <X className="h-5 w-5" />
                                        </Button>
                                    </div>
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
                                <Link href="/cart" className="relative" prefetch={false}>
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

                                {/* User Menu */}
                                {user ? (
                                    <DropdownMenu>
                                        <DropdownMenuTrigger asChild>
                                            <Button 
                                                variant="ghost" 
                                                size="icon"
                                                className="text-gray-600 hover:text-gray-900"
                                            >
                                                <CircleUser className="h-5 w-5" />
                                            </Button>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent align="end" className="w-56">
                                            <DropdownMenuLabel className="flex gap-3 items-center">
                                                <Avatar>
                                                    <AvatarFallback>
                                                        {user_initials}
                                                    </AvatarFallback>
                                                </Avatar>
                                                <div className="flex flex-col">
                                                    <span className="text-sm font-medium">{user?.name}</span>
                                                </div>
                                            </DropdownMenuLabel>
                                            <DropdownMenuSeparator />
                                            <DropdownMenuGroup>
                                                <Link href="/profile" prefetch={true}>
                                                    <DropdownMenuItem className="gap-2 cursor-pointer">
                                                        <Settings className="h-4 w-4" />
                                                        <span>Profile</span>
                                                    </DropdownMenuItem>
                                                </Link>
                                                <Link href="/orders" prefetch={true}>
                                                    <DropdownMenuItem className="gap-2 cursor-pointer">
                                                        <ShoppingBagIcon className="h-4 w-4" />
                                                        <span>Orders</span>
                                                    </DropdownMenuItem>
                                                </Link>
                                                <DropdownMenuSeparator />
                                                <DropdownMenuItem 
                                                    onClick={() => signOut()}
                                                    className="gap-2 cursor-pointer text-red-600 focus:text-red-600"
                                                >
                                                    <LogOut className="h-4 w-4" />
                                                    <span>Logout</span>
                                                </DropdownMenuItem>
                                            </DropdownMenuGroup>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                ) : (
                                    <Link href="/login">
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
                                                className="text-lg font-medium text-gray-700 hover:text-primary transition-colors py-2"
                                                onClick={() => setMobileMenuOpen(false)}
                                            >
                                                {link.name}
                                            </Link>
                                        ))}
                                        {!user && (
                                            <Link href="/login" onClick={() => setMobileMenuOpen(false)}>
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