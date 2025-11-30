import React from "react";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import Link from "next/link";

export default function OverViewLayout({
  hero,
  product_grid,
  why_choose_us,
  testimonials,
}: {
  hero: React.ReactNode;
  product_grid: React.ReactNode;
  why_choose_us: React.ReactNode;
  testimonials: React.ReactNode;
}) {
  return (
    <>
      {hero}
      {product_grid}
      {why_choose_us}
      {testimonials}
      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-primary to-primary/90 text-white">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-4xl md:text-5xl font-bold mb-6">
              Ready to Grow Your Business?
            </h2>
            <p className="text-xl text-white/90 mb-8 leading-relaxed">
              Join thousands of retailers who trust Efoyeta Store for their wholesale needs. 
              Start your journey today and unlock exclusive benefits.
            </p>
            <div className="flex justify-center">
              <Button
                size="lg"
                variant="secondary"
                className="bg-white text-primary hover:bg-gray-100 shadow-lg hover:shadow-xl transition-all duration-300 px-8 py-6 text-lg font-semibold"
                asChild
              >
                <Link href="/product" prefetch={true}>
                  Browse Products <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>
    </>
  );
}
