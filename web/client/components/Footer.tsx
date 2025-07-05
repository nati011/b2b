import { Facebook, Instagram, Twitter } from "lucide-react";
import Link from "next/link";
import Image from 'next/image'

export const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <div className="bg-secondary mt-16 w-full print:hidden">
            <div className="container mx-auto px-4 py-12">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                    <div className="md:col-span-1">
                        <Link href="/" className="text-xl font-bold">
                            <Image src='/logo.png' width={150} height={100} alt="logo" />
                        </Link>
                        <p className="mt-4 text-sm text-muted-foreground">
                            Lorem ipsum dolor sit amet consectetur, adipisicing elit.
                        </p>
                        <div className="flex mt-6 space-x-4">
                            <a href="#" className="text-muted-foreground hover:text-primary transition-colors">
                                <Facebook className="h-5 w-5" />
                            </a>
                            <a href="#" className="text-muted-foreground hover:text-primary transition-colors">
                                <Instagram className="h-5 w-5" />
                            </a>
                            <a href="#" className="text-muted-foreground hover:text-primary transition-colors">
                                <Twitter className="h-5 w-5" />
                            </a>
                        </div>
                    </div>

                    <div>
                        <h3 className="font-medium mb-4">Shop</h3>
                        <ul className="space-y-2 text-sm">
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">All Products</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">New Arrivals</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Featured</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Sale</a></li>
                        </ul>
                    </div>

                    <div>
                        <h3 className="font-medium mb-4">Help</h3>
                        <ul className="space-y-2 text-sm">
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Shipping</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Returns</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">FAQ</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Contact</a></li>
                        </ul>
                    </div>

                    <div>
                        <h3 className="font-medium mb-4">About</h3>
                        <ul className="space-y-2 text-sm">
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Our Story</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Sustainability</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Press</a></li>
                            <li><a href="#" className="text-muted-foreground hover:text-primary transition-colors">Careers</a></li>
                        </ul>
                    </div>
                </div>

                <div className="border-t border-primary mt-12 pt-6">
                    <div className="flex flex-col md:flex-row justify-between items-center">
                        <p className="text-sm text-muted-foreground mb-4 md:mb-0">
                            &copy; {currentYear} Efoyeta. All rights reserved.
                        </p>
                        <div className="flex space-x-6">
                            <a href="#" className="text-xs text-muted-foreground hover:text-primary transition-colors">
                                Privacy Policy
                            </a>
                            <a href="#" className="text-xs text-muted-foreground hover:text-primary transition-colors">
                                Terms of Service
                            </a>
                            <a href="#" className="text-xs text-muted-foreground hover:text-primary transition-colors">
                                Cookies
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};