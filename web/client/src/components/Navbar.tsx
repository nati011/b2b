
import { Link } from "react-router-dom";
import { ShoppingBag, Menu, X, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useCart } from "@/contexts/CartContext";
import { useState } from "react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";

export const Navbar = () => {
    const { getTotalItems } = useCart();
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [searchOpen, setSearchOpen] = useState(false);

    const toggleSearch = () => {
        setSearchOpen(!searchOpen);
    };

    const navLinks = [
        { name: "Home", href: "/" },
        { name: "Shop", href: "/#products" },
        { name: "About", href: "/about" },
        { name: "Contact", href: "/contact" },
    ];

    return (
        <header className="sticky top-0 z-50 w-full bg-white bg-opacity-95 backdrop-blur-sm border-b">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    {/* Logo */}
                    <div className="flex-shrink-0">
                        <Link to="/" className="text-lg md:text-xl font-bold">
                            Efoyeta Store
                        </Link>
                    </div>

                    {/* Desktop Navigation */}
                    <nav className="hidden md:flex space-x-8">
                        {navLinks.map((link) => (
                            <Link
                                key={link.name}
                                to={link.href}
                                className="text-sm font-medium hover:text-primary transition-colors"
                            >
                                {link.name}
                            </Link>
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

                                <Link to="/cart" className="ml-4 relative">
                                    <Button variant="ghost" size="icon">
                                        <ShoppingBag className="h-5 w-5" />
                                        {getTotalItems() > 0 && (
                                            <span className="absolute -top-2 -right-2 bg-primary text-white rounded-full w-5 h-5 flex items-center justify-center text-xs">
                                                {getTotalItems()}
                                            </span>
                                        )}
                                    </Button>
                                </Link>

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
                                                <Link
                                                    key={link.name}
                                                    to={link.href}
                                                    className="text-lg font-medium"
                                                    onClick={() => setMobileMenuOpen(false)}
                                                >
                                                    {link.name}
                                                </Link>
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