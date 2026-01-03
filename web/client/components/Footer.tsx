"use client";

import { Instagram } from "lucide-react";
import { PiTelegramLogo } from "react-icons/pi";
import { IoLogoTiktok } from "react-icons/io5";
import Link from "next/link";
import Image from 'next/image'

export const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <footer className="bg-gray-900 text-gray-300 border-t border-gray-800">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-12">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 md:gap-12 items-start">
                    {/* Left: Logo and Tagline */}
                    <div className="flex flex-col items-center md:items-start">
                        <Link href="/" className="inline-block mb-1">
                            <Image 
                                src='/logo.png' 
                                width={100} 
                                height={50} 
                                alt="Efoyeta Store Logo" 
                                className="brightness-0 invert"
                            />
                        </Link>
                        <p className="text-gray-400 text-sm text-center md:text-left max-w-xs leading-relaxed -mt-1">
                            Turning Small Capital into Big Opportunities.
                        </p>
                    </div>

                    {/* Center: Quick Links */}
                    <div className="flex flex-col items-center md:items-start">
                        <h3 className="text-white font-semibold text-sm mb-4 hidden md:block">
                            Quick Links
                        </h3>
                        <div className="flex flex-col items-center md:items-start gap-3">
                            <Link href="/" className="text-gray-400 hover:text-primary transition-colors text-sm font-medium">
                                Home
                            </Link>
                            <Link href="/product" prefetch={true} className="text-gray-400 hover:text-primary transition-colors text-sm font-medium">
                                Shop
                            </Link>
                            <Link href="/signup/distributor" className="text-gray-400 hover:text-primary transition-colors text-sm font-medium">
                                Partner Portal
                            </Link>
                            <Link href="/contact" className="text-gray-400 hover:text-primary transition-colors text-sm font-medium">
                                Contact
                            </Link>
                        </div>
                    </div>

                    {/* Third: Mobile App */}
                    <div className="flex flex-col items-center md:items-start">
                        <div className="text-center md:text-left mb-4">
                            <h3 className="text-white font-semibold text-sm mb-1">
                                Get the App
                            </h3>
                            <p className="text-gray-400 text-xs">
                                Download our mobile app
                            </p>
                        </div>
                        <div className="flex flex-col gap-3">
                            {/* Google Play - Coming Soon */}
                            <a 
                                href="#" 
                                onClick={(e) => e.preventDefault()}
                                className="inline-flex items-center justify-center px-4 py-2.5 bg-gray-800 text-white rounded-lg opacity-60 cursor-not-allowed transition-all duration-300 relative"
                                title="Coming Soon"
                            >
                                <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.53,12.9 20.18,13.18L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
                                </svg>
                                <div className="text-left">
                                    <div className="text-xs opacity-80">Get it on</div>
                                    <div className="text-xs font-semibold">Google Play</div>
                                </div>
                                <span className="absolute -top-1 -right-1 bg-primary text-white text-xs px-1.5 py-0.5 rounded-full">Soon</span>
                            </a>
                            {/* App Store - Coming Soon */}
                            <a 
                                href="#" 
                                onClick={(e) => e.preventDefault()}
                                className="inline-flex items-center justify-center px-4 py-2.5 bg-gray-800 text-white rounded-lg opacity-60 cursor-not-allowed transition-all duration-300 relative"
                                title="Coming Soon"
                            >
                                <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.24 0-1.44.62-2.2.44-3.06-.4C1.79 15.25 2.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>
                                </svg>
                                <div className="text-left">
                                    <div className="text-xs opacity-80">Download on the</div>
                                    <div className="text-xs font-semibold">App Store</div>
                                </div>
                                <span className="absolute -top-1 -right-1 bg-primary text-white text-xs px-1.5 py-0.5 rounded-full">Soon</span>
                            </a>
                        </div>
                    </div>

                    {/* Right: Social Links */}
                    <div className="flex flex-col items-center md:items-start">
                        <div className="text-center md:text-left mb-4">
                            <h3 className="text-white font-semibold text-sm mb-1">
                                Follow Us
                            </h3>
                            <p className="text-gray-400 text-xs">
                                Contact us on our social media
                            </p>
                        </div>
                        <div className="flex items-center gap-3">
                            <a 
                                href="https://www.tiktok.com/@efoyeta" 
                                className="p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                aria-label="TikTok"
                                title="TikTok"
                            >
                                <IoLogoTiktok className="w-5 h-5" />
                            </a>
                            <a 
                                href="https://t.me/efoyetastore1" 
                                className="p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                aria-label="Telegram"
                                title="Telegram"
                            >
                                <PiTelegramLogo className="w-5 h-5" />
                            </a>
                            <a 
                                href="https://www.instagram.com/efoyetastore?igsh=MTJodTNtZXRwaGgwOQ==" 
                                className="p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110 hover:shadow-lg hover:shadow-primary/20"
                                aria-label="Instagram"
                                title="Instagram"
                            >
                                <Instagram className="w-5 h-5" />
                            </a>
                        </div>
                    </div>
                </div>

                {/* Copyright */}
                <div className="mt-8 pt-8 border-t border-gray-800">
                    <div className="text-center text-sm text-gray-500">
                        &copy; {currentYear}{" "}
                        <a 
                            href="https://www.efoyetastore.com/" 
                            className="hover:text-primary transition-colors"
                        >
                            Efoyeta Store PLC.
                        </a>{" "}
                        All rights reserved.
                    </div>
                </div>
            </div>
        </footer>
    );
};