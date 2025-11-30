"use client";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import Link from "next/link";

export default function Hero() {
  return (
    <section className="w-full bg-gradient-to-b from-gray-50 via-white to-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32 lg:py-40">
        <div className="max-w-4xl mx-auto text-center">
          {/* Badge */}
          <div className="inline-flex items-center justify-center mb-6">
            <span className="px-4 py-2 bg-primary/10 text-primary text-sm font-semibold rounded-full uppercase tracking-wide">
              Welcome to Efoyeta Store
            </span>
          </div>

          {/* Main Heading */}
          <h1 className="text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-bold text-gray-900 mb-6 leading-tight">
            Turning Small Capital into{" "}
            <span className="text-primary">Big Opportunities</span>
          </h1>

          {/* Description */}
          <p className="text-lg md:text-xl text-gray-600 mb-10 max-w-2xl mx-auto leading-relaxed">
            Elevate your retail business with our exclusive wholesale collections. 
            Access premium products, competitive pricing, and dedicated support to maximize your profitability.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mb-16">
            <Button
              size="lg"
              className="bg-primary hover:bg-primary/90 text-white shadow-lg hover:shadow-xl transition-all duration-300 px-8 py-6 text-base font-semibold"
              asChild
            >
              <Link href="#products">
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

          {/* Stats */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-8 md:gap-12 pt-12 border-t border-gray-200">
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold text-primary mb-2">500+</div>
              <div className="text-sm md:text-base text-gray-600 font-medium">Products</div>
            </div>
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold text-primary mb-2">1000+</div>
              <div className="text-sm md:text-base text-gray-600 font-medium">Retailers</div>
            </div>
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold text-primary mb-2">50%</div>
              <div className="text-sm md:text-base text-gray-600 font-medium">Savings</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
