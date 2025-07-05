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

    const {
        cartItems,
        totalItems
    } = useCartStore()

    const navLinks = [
        { name: "Home", href: "/" },
        { name: "Shop", href: "/product" },
        { name: "Contact", href: "/contact" },
    ];

    return (
        <header className="sticky top-0 z-50 w-full bg-white bg-opacity-95 backdrop-blur-sm border-b print:hidden">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    {/* Logo */}
                    
                    <Link href="/" className="flex font-semibold">
                        <Image src='/logo.png' width={150} height={100} alt="logo" />
                        <span className="sr-only">EFOYETA STORE</span>

                    </Link>

                    {/* Desktop Navigation */}
                    <nav className="hidden md:flex space-x-8">
                        {navLinks.map((link) => (
                            <a
                                key={link.name}
                                href={link.href}
                                className="text-sm font-medium hover:text-primary transition-colors"
                            >
                                {link.name}
                            </a>
                        ))}
                    </nav>

                    {/* Search and Cart */}
                    <div className="flex items-center">
                        {searchOpen ? (
                            <div className="absolute inset-0 px-4 flex items-center justify-center bg-white bg-opacity-95 z-10">
                                <div className="w-full max-w-md">
                                    <div className="flex items-center">
                                        <Input
                                            placeholder="Search products..."
                                            className="flex-grow"
                                            autoFocus
                                        />
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={toggleSearch}
                                            className="ml-2"
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
                                    className="hidden md:flex"
                                >
                                    <Search className="h-5 w-5" />
                                </Button>

                                <a href="/cart" className="ml-4 relative">
                                    <Button variant="ghost" size="icon">
                                        <ShoppingBag className="h-5 w-5" />
                                        {totalItems > 0 && (
                                            <span className="absolute -top-2 -right-2 bg-primary text-white rounded-full w-5 h-5 flex items-center justify-center text-xs">
                                                {totalItems}
                                            </span>
                                        )}
                                    </Button>
                                </a>
                                {
                                    user ? (

                                        <DropdownMenu>
                                            <DropdownMenuTrigger asChild>
                                                <Button variant="ghost" size="icon">
                                                    <CircleUser className="h-5 w-5" />
                                                </Button>
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent align="end">
                                                <DropdownMenuLabel className="flex gap-2 items-center">
                                                    <Avatar>

                                                        <AvatarFallback>
                                                            {user_initials}
                                                        </AvatarFallback>
                                                    </Avatar>
                                                    <div className="flex flex-col gap-1">
                                                        {user?.name}
                                                    </div>
                                                </DropdownMenuLabel>
                                                <DropdownMenuSeparator />
                                                <DropdownMenuGroup>
                                                    <Link href="/profile">
                                                        <DropdownMenuItem
                                                            className="gap-2"
                                                        >
                                                            <Settings className="h-5 w-5" />
                                                            <span>Profile</span>
                                                        </DropdownMenuItem>
                                                    </Link>

                                                    <Link href="/orders">
                                                        <DropdownMenuItem
                                                            className="gap-2"
                                                        >
                                                            <ShoppingBagIcon className="h-5 w-5" />
                                                            <span>Orders</span>
                                                        </DropdownMenuItem>
                                                    </Link>

                                                    <DropdownMenuItem onClick={() => signOut()}
                                                    >
                                                        <LogOut className="mr-2 text-red-500" />
                                                        Logout
                                                    </DropdownMenuItem>
                                                </DropdownMenuGroup>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    ) : (
                                        <a href="/login" className="ml-4 relative">
                                            <Button variant="ghost" size="icon">
                                                <LogIn className="h-5 w-5" />
                                            </Button>
                                        </a>
                                    )
                                }



                                {/* Mobile menu button */}
                                <Sheet>
                                    <SheetTrigger asChild>
                                        <Button variant="ghost" size="icon" className="ml-2 md:hidden">
                                            <Menu className="h-5 w-5" />
                                        </Button>
                                    </SheetTrigger>
                                    <SheetContent side="right">
                                        <div className="flex flex-col space-y-6 mt-10">
                                            {navLinks.map((link) => (
                                                <a
                                                    key={link.name}
                                                    href={link.href}
                                                    className="text-lg font-medium"
                                                    onClick={() => setMobileMenuOpen(false)}
                                                >
                                                    {link.name}
                                                </a>
                                            ))}
                                            <Button
                                                onClick={toggleSearch}
                                                variant="outline"
                                                className="flex items-center justify-start w-full"
                                            >
                                                <Search className="h-5 w-5 mr-2" />
                                                Search
                                            </Button>
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