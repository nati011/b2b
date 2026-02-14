"use client";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import Link from "next/link";

export default function Hero() {
  return (
    <section className="w-full relative bg-gradient-to-b from-gray-50 via-white to-white overflow-hidden">
      {/* Background decorative elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        {/* Subtle grid pattern */}
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#e5e7eb_1px,transparent_1px),linear-gradient(to_bottom,#e5e7eb_1px,transparent_1px)] bg-[size:4rem_4rem] opacity-20"></div>
        
        {/* Geometric shapes */}
        <div className="absolute top-20 -left-20 w-96 h-96 bg-primary/5 rounded-full blur-3xl"></div>
        <div className="absolute bottom-20 -right-20 w-96 h-96 bg-primary/5 rounded-full blur-3xl"></div>
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-primary/3 rounded-full blur-3xl"></div>
        
        {/* Subtle diagonal lines */}
        <div className="absolute top-0 left-0 w-full h-full">
          <div className="absolute top-1/4 left-0 w-px h-64 bg-gradient-to-b from-transparent via-primary/10 to-transparent"></div>
          <div className="absolute top-1/3 right-0 w-px h-64 bg-gradient-to-b from-transparent via-primary/10 to-transparent"></div>
        </div>
      </div>
      
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32 lg:py-40 relative z-10">
        <div className="max-w-4xl mx-auto text-center">
          {/* Badge */}
          <div className="inline-flex items-center justify-center mb-6">
            <span className="px-4 py-2 bg-primary/10 text-primary text-sm font-semibold tracking-wide rounded-full">
              <span className="typing-text">Turning Small Capital into Big Opportunities</span>
            </span>
          </div>

          {/* Main Heading */}
          <h1 className="text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-bold text-gray-900 mb-6 leading-tight">
            Welcome to{" "}
            <span className="text-primary">Efoyetastore</span>
          </h1>

          {/* Description */}
          <p className="text-lg md:text-xl text-gray-600 mb-10 max-w-2xl mx-auto leading-relaxed">
            Your one-stop online store for high-quality products delivered directly to your doorstep. We
            bring a wide selection of items through our trusted partner suppliers, combining convenience,
            affordability, and reliability.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-16">
            <Button
              size="lg"
              className="bg-primary hover:bg-primary/90 text-white shadow-lg hover:shadow-xl transition-all duration-300 px-8 py-6 text-base font-semibold"
              asChild
            >
              <Link href="/product" prefetch={true}>
                Explore Products <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
            </Button>
            <Button 
              variant="outline" 
              size="lg"
              className="border-2 border-gray-300 hover:bg-gray-50 hover:border-primary transition-all duration-300 px-8 py-6 text-base font-semibold"
              asChild
            >
              <Link href="#pricing">Become a Partner</Link>
            </Button>
          </div>

        </div>
      </div>

      {/* Mobile App Banner */}
      <div className="w-full bg-gradient-to-r from-primary/10 via-primary/5 to-primary/10 border-t border-b border-primary/20 mt-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 max-w-6xl mx-auto">
            <div className="flex items-center gap-4 flex-1">
              <div className="flex items-center justify-center w-14 h-14 bg-primary rounded-xl shadow-lg flex-shrink-0">
                <svg className="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-bold text-gray-900 mb-1">
                  Get the Mobile App
                </h3>
                <p className="text-sm text-gray-600">
                  Shop wholesale products on the go with our mobile app
                </p>
              </div>
            </div>
            <div className="flex flex-col sm:flex-row gap-3 flex-shrink-0">
              {/* Google Play - Coming Soon */}
              <a 
                href="#" 
                onClick={(e) => e.preventDefault()}
                className="inline-flex items-center justify-center px-6 py-3 bg-white text-gray-900 rounded-lg opacity-60 cursor-not-allowed transition-all duration-300 shadow-lg relative"
                title="Coming Soon"
              >
                <svg className="w-6 h-6 mr-2" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.53,12.9 20.18,13.18L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
                </svg>
                <div className="text-left">
                  <div className="text-xs opacity-80">Get it on</div>
                  <div className="text-sm font-semibold">Google Play</div>
                </div>
                <span className="absolute -top-2 -right-2 bg-primary text-white text-xs px-2 py-0.5 rounded-full">Soon</span>
              </a>
              {/* App Store - Coming Soon */}
              <a 
                href="#" 
                onClick={(e) => e.preventDefault()}
                className="inline-flex items-center justify-center px-6 py-3 bg-white text-gray-900 rounded-lg opacity-60 cursor-not-allowed transition-all duration-300 shadow-lg relative"
                title="Coming Soon"
              >
                <svg className="w-6 h-6 mr-2" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.24 0-1.44.62-2.2.44-3.06-.4C1.79 15.25 2.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>
                </svg>
                <div className="text-left">
                  <div className="text-xs opacity-80">Download on the</div>
                  <div className="text-sm font-semibold">App Store</div>
                </div>
                <span className="absolute -top-2 -right-2 bg-primary text-white text-xs px-2 py-0.5 rounded-full">Soon</span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
