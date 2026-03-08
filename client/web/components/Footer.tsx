"use client";

import { Instagram, Youtube } from "lucide-react";
import { PiTelegramLogo } from "react-icons/pi";
import { IoLogoTiktok } from "react-icons/io5";
import Link from "next/link";
import { useRouter } from "next/navigation";
import Image from 'next/image'

export const Footer = () => {
    const currentYear = new Date().getFullYear();
    const router = useRouter();

    // Prefetch footer links on hover
    const handleLinkHover = (href: string) => {
        router.prefetch(href);
    };

    return (
        <footer className="bg-gray-900 text-gray-300 border-t border-gray-800">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 sm:gap-10 lg:gap-12">
                    {/* Left: Logo and Tagline */}
                    <div className="flex flex-col items-center sm:items-start space-y-3 sm:space-y-4">
                        <Link 
                            href="/" 
                            className="inline-block"
                            prefetch={true}
                            onMouseEnter={() => handleLinkHover("/")}
                        >
                            <Image 
                                src='/logo.png' 
                                width={100} 
                                height={50} 
                                alt="Efoyeta Store Logo" 
                                className="brightness-0 invert w-20 sm:w-24 lg:w-28 h-12 sm:h-14 lg:h-16 object-contain"
                                loading="lazy"
                            />
                        </Link>
                        <p className="text-gray-400 text-xs sm:text-sm text-center sm:text-left max-w-xs leading-relaxed">
                            Your one-stop online store for high-quality products delivered directly to your doorstep.
                        </p>
                    </div>

                    {/* Center: Quick Links */}
                    <div className="flex flex-col items-center sm:items-start space-y-3 sm:space-y-4">
                        <h3 className="text-white font-semibold text-sm sm:text-base">
                            Quick Links
                        </h3>
                        <nav className="flex flex-row flex-wrap items-center justify-center sm:flex-col sm:items-start gap-2 sm:gap-3 w-full">
                            <Link 
                                href="/" 
                                prefetch={true}
                                onMouseEnter={() => handleLinkHover("/")}
                                className="text-gray-400 hover:text-primary transition-colors text-sm font-medium text-center sm:text-left py-1 whitespace-nowrap"
                            >
                                Home
                            </Link>
                            <Link 
                                href="/product" 
                                prefetch={true}
                                onMouseEnter={() => handleLinkHover("/product")}
                                className="text-gray-400 hover:text-primary transition-colors text-sm font-medium text-center sm:text-left py-1 whitespace-nowrap"
                            >
                                Products
                            </Link>
                            <Link 
                                href="/contact" 
                                prefetch={true}
                                onMouseEnter={() => handleLinkHover("/contact")}
                                className="text-gray-400 hover:text-primary transition-colors text-sm font-medium text-center sm:text-left py-1 whitespace-nowrap"
                            >
                                Contact
                            </Link>
                            <Link 
                                href="/return-refund-policy" 
                                prefetch={true}
                                onMouseEnter={() => handleLinkHover("/return-refund-policy")}
                                className="text-gray-400 hover:text-primary transition-colors text-sm font-medium text-center sm:text-left py-1 whitespace-nowrap"
                            >
                                Return & Refund Policy
                            </Link>
                        </nav>
                    </div>

                    {/* Mobile: Get the App and Follow Us side by side */}
                    <div className="grid grid-cols-2 gap-4 sm:contents">
                        {/* Third: Mobile App */}
                        <div className="flex flex-col items-center sm:items-start space-y-3 sm:space-y-4">
                            <div className="text-center sm:text-left w-full">
                                <h3 className="text-white font-semibold text-sm sm:text-base mb-1 sm:mb-2">
                                    Get the App
                                </h3>
                                <p className="text-gray-400 text-xs sm:text-sm">
                                    Download our mobile app
                                </p>
                            </div>
                            <div className="flex flex-col gap-2 sm:gap-3 w-full sm:w-auto">
                                {/* Google Play - Coming Soon */}
                                <a 
                                    href="#" 
                                    onClick={(e) => e.preventDefault()}
                                    className="inline-flex items-center justify-center px-2 sm:px-4 py-2 sm:py-2.5 sm:py-3 bg-gray-800 text-white rounded-lg opacity-60 cursor-not-allowed transition-all duration-300 relative w-full sm:w-auto"
                                    title="Coming Soon"
                                >
                                    <svg className="w-3 h-3 sm:w-4 sm:h-4 lg:w-5 lg:h-5 mr-1 sm:mr-2 flex-shrink-0" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.53,12.9 20.18,13.18L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
                                    </svg>
                                    <div className="text-left flex-1 hidden sm:block">
                                        <div className="text-xs opacity-80">Get it on</div>
                                        <div className="text-xs font-semibold">Google Play</div>
                                    </div>
                                    <span className="absolute -top-1 -right-1 bg-primary text-white text-xs px-1 sm:px-1.5 py-0.5 rounded-full text-[10px] sm:text-xs">Soon</span>
                                </a>
                                {/* App Store - Coming Soon */}
                                <a 
                                    href="#" 
                                    onClick={(e) => e.preventDefault()}
                                    className="inline-flex items-center justify-center px-2 sm:px-4 py-2 sm:py-2.5 sm:py-3 bg-gray-800 text-white rounded-lg opacity-60 cursor-not-allowed transition-all duration-300 relative w-full sm:w-auto"
                                    title="Coming Soon"
                                >
                                    <svg className="w-3 h-3 sm:w-4 sm:h-4 lg:w-5 lg:h-5 mr-1 sm:mr-2 flex-shrink-0" viewBox="0 0 24 24" fill="currentColor">
                                        <path d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.24 0-1.44.62-2.2.44-3.06-.4C1.79 15.25 2.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>
                                    </svg>
                                    <div className="text-left flex-1 hidden sm:block">
                                        <div className="text-xs opacity-80">Download on the</div>
                                        <div className="text-xs font-semibold">App Store</div>
                                    </div>
                                    <span className="absolute -top-1 -right-1 bg-primary text-white text-xs px-1 sm:px-1.5 py-0.5 rounded-full text-[10px] sm:text-xs">Soon</span>
                                </a>
                            </div>
                        </div>

                        {/* Right: Social Links */}
                        <div className="flex flex-col items-center sm:items-start space-y-3 sm:space-y-4">
                            <div className="text-center sm:text-left w-full">
                                <h3 className="text-white font-semibold text-sm sm:text-base mb-1 sm:mb-2">
                                    Follow Us
                                </h3>
                                <p className="text-gray-400 text-xs sm:text-sm hidden sm:block">
                                    Contact us on our social media
                                </p>
                            </div>
                            <div className="flex items-center gap-2 sm:gap-3 justify-center sm:justify-start">
                                <a 
                                    href="https://www.tiktok.com/@efoyeta" 
                                    className="p-2 sm:p-2.5 lg:p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                    aria-label="TikTok"
                                    title="TikTok"
                                >
                                    <IoLogoTiktok className="w-4 h-4 sm:w-5 sm:h-5" />
                                </a>
                                <a 
                                    href="https://t.me/efoyetastore1" 
                                    className="p-2 sm:p-2.5 lg:p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                    aria-label="Telegram"
                                    title="Telegram"
                                >
                                    <PiTelegramLogo className="w-4 h-4 sm:w-5 sm:h-5" />
                                </a>
                                <a 
                                    href="https://www.instagram.com/efoyetastore?igsh=MTJodTNtZXRwaGgwOQ==" 
                                    className="p-2 sm:p-2.5 lg:p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                    aria-label="Instagram"
                                    title="Instagram"
                                >
                                    <Instagram className="w-4 h-4 sm:w-5 sm:h-5" />
                                </a>
                                <a 
                                    href="https://www.youtube.com/@efoyetastore" 
                                    className="p-2 sm:p-2.5 lg:p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                    aria-label="YouTube"
                                    title="YouTube"
                                >
                                    <Youtube className="w-4 h-4 sm:w-5 sm:h-5" />
                                </a>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Payment Partners */}
                <div className="mt-8 sm:mt-10 pt-6 sm:pt-8 border-t border-gray-800">
                    <h3 className="text-white font-semibold text-sm sm:text-base text-center mb-4 sm:mb-5">
                        Payment Partners
                    </h3>
                    <div className="flex flex-wrap items-center justify-center gap-6 sm:gap-8">
                        <Image
                            src="/telebirr-logo.png"
                            width={120}
                            height={48}
                            alt="Telebirr"
                            className="h-10 sm:h-12 object-contain"
                            loading="lazy"
                        />
                        <Image
                            src="/cbe-logo.png"
                            width={80}
                            height={80}
                            alt="Commercial Bank of Ethiopia"
                            className="h-10 sm:h-12 w-auto object-contain"
                            loading="lazy"
                        />
                    </div>
                </div>

                {/* Copyright */}
                <div className="mt-6 sm:mt-8 pt-6 sm:pt-8 border-t border-gray-800">
                    <div className="text-center text-xs sm:text-sm text-gray-500 px-4">
                        &copy; {currentYear}{" "}
                        <a 
                            href="https://www.efoyetastore.com/" 
                            className="hover:text-primary transition-colors"
                        >
                            Efoyeta Store General Trading PLC
                        </a>{" "}
                        All rights reserved.
                    </div>
                </div>
            </div>
        </footer>
    );
};