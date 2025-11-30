import { Facebook, Instagram, Linkedin, Mail, MapPin, Phone, Smartphone, Twitter } from "lucide-react";
import { PiTelegramLogo } from "react-icons/pi";
import { IoLogoTiktok } from "react-icons/io5";
import Link from "next/link";
import Image from 'next/image'

export const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <footer className="bg-gray-900 text-gray-300">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-12 mb-12">
                    {/* Company Info */}
                    <div>
                        <Link href="/" className="inline-block mb-2">
                            <Image 
                                src='/logo.png' 
                                width={120} 
                                height={120} 
                                alt="Efoyeta Store Logo" 
                                className="brightness-0 invert"
                            />
                        </Link>
                        <p className="text-gray-400 leading-relaxed max-w-sm mt-1">
                            Turning Small Capital into Big Opportunities.
                        </p>
                        <div className="flex gap-3">
                            <a 
                                href="https://t.me/efoyetastore1" 
                                className="p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110"
                                aria-label="Telegram"
                            >
                                <PiTelegramLogo className="w-5 h-5" />
                            </a>
                            <a 
                                href="https://www.tiktok.com/@efoyeta" 
                                className="p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110"
                                aria-label="TikTok"
                            >
                                <IoLogoTiktok className="w-5 h-5" />
                            </a>
                            <a 
                                href="https://www.instagram.com/efoyetastore?igsh=MTJodTNtZXRwaGgwOQ==" 
                                className="p-3 bg-gray-800 rounded-lg hover:bg-primary transition-all duration-300 hover:scale-110"
                                aria-label="Instagram"
                            >
                                <Instagram className="w-5 h-5" />
                            </a>
                        </div>
                    </div>

                    {/* Quick Links */}
                    <div>
                        <h3 className="text-white font-semibold text-lg mb-6">Quick Links</h3>
                        <ul className="space-y-3">
                            <li>
                                <Link href="/" className="text-gray-400 hover:text-primary transition-colors">
                                    Home
                                </Link>
                            </li>
                            <li>
                                <Link href="/product" prefetch={true} className="text-gray-400 hover:text-primary transition-colors">
                                    Shop
                                </Link>
                            </li>
                            <li>
                                <Link href="/contact" className="text-gray-400 hover:text-primary transition-colors">
                                    Contact
                                </Link>
                            </li>
                        </ul>
                    </div>

                    {/* Mobile App */}
                    <div>
                        <div className="flex items-center gap-2 mb-4">
                            <Smartphone className="w-5 h-5 text-primary" />
                            <h3 className="text-white font-semibold text-lg">Get the Mobile App</h3>
                        </div>
                        <p className="text-gray-400 mb-6 text-sm leading-relaxed">
                            Shop wholesale products on the go with our mobile app
                        </p>
                        <a 
                            href="https://play.google.com/store" 
                            target="_blank" 
                            rel="noopener noreferrer"
                            className="inline-flex items-center justify-center px-6 py-3 bg-white text-gray-900 rounded-lg hover:bg-gray-100 transition-all duration-300 hover:scale-105 shadow-lg"
                        >
                            <svg className="w-6 h-6 mr-2" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.53,12.9 20.18,13.18L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
                            </svg>
                            <div className="text-left">
                                <div className="text-xs opacity-80">Get it on</div>
                                <div className="text-sm font-semibold">Google Play</div>
                            </div>
                        </a>
                    </div>
                </div>

                {/* Copyright */}
                <div className="pt-8 border-t border-gray-800">
                    <div className="text-center text-sm text-gray-400">
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