"use client";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useEffect, useState } from "react";
import usePlanstore from "@/lib/store/usePricingPlan";
import Link from "next/link";

export default function PricingSection(){
  const { plans, loading, error, fetchPricingPlan } = usePlanstore();
  const [mounted, setMounted] = useState(false);

  // Ensure this only runs on client-side
  useEffect(() => {
    setMounted(true);
    console.log('PricingSection mounted on client, fetching plans...');
    void fetchPricingPlan();
  }, []); // Empty dependency array - only run on mount

  const formatPrice = (price: number) => `${price.toLocaleString()} ETB`;
  const termToPeriod = (termInMonth: number) => {
    if (termInMonth === 1) return "per month";
    if (termInMonth === 12) return "per year";
    return `for ${termInMonth} months`;
  };

  const isRecommended = (name: string) => /year|annual|annunal/i.test(name) && !/two/i.test(name);

  return (
    <section id="pricing" className="py-20 bg-white scroll-mt-20">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <p className="text-sm font-semibold text-primary uppercase tracking-wider mb-3">
            Pricing Plans
          </p>
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Distributor Pricing
          </h2>
          <div className="h-1 w-24 bg-primary mx-auto mb-6 rounded-full"></div>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Choose the package that fits your business needs and growth goals
          </p>
        </div>

        {error && (
          <div className="text-center mb-8">
            <div className="text-destructive mb-4">{error}</div>
            <Button 
              onClick={() => void fetchPricingPlan()} 
              variant="outline"
              className="mt-2"
            >
              Retry
            </Button>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          {loading && plans.length === 0 && (
            Array.from({ length: 3 }).map((_, idx) => (
              <div key={idx} className="relative bg-white p-8 rounded-2xl border border-gray-200 shadow-sm">
                <div className="text-center mb-6">
                  <Skeleton className="h-6 w-32 mx-auto mb-2" />
                  <Skeleton className="h-8 w-24 mx-auto mb-1" />
                  <Skeleton className="h-4 w-28 mx-auto" />
                  <div className="mt-4 space-y-2">
                    <Skeleton className="h-3 w-56 mx-auto" />
                    <Skeleton className="h-3 w-44 mx-auto" />
                  </div>
                </div>
                <div className="space-y-3 mb-8">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-5/6" />
                  <Skeleton className="h-4 w-2/3" />
                </div>
                <Skeleton className="h-9 w-full" />
              </div>
            ))
          )}

          {!loading && plans.length === 0 && !error && (
            <div className="md:col-span-3 text-center text-muted-foreground">No pricing plans available.</div>
          )}

          {plans.map(plan => (
            <div 
              key={plan.id}
              className={`relative bg-white p-8 rounded-2xl border transition-all duration-300 hover:shadow-xl ${
                isRecommended(plan.name) 
                  ? "border-primary shadow-lg scale-105 ring-2 ring-primary/20" 
                  : "border-gray-200 shadow-sm hover:border-primary/50"
              }`}
            >
              {isRecommended(plan.name) && (
                <div className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-primary text-white px-6 py-2 rounded-full text-sm font-semibold shadow-lg">
                  Recommended
                </div>
              )}
              <div className="text-center mb-8">
                <h3 className="text-2xl font-bold text-gray-900 mb-4">{plan.name}</h3>
                <div className="text-4xl font-bold text-primary mb-2">{formatPrice(plan.price)}</div>
                <div className="text-sm text-gray-600 mb-4">{termToPeriod(plan.term_in_month)}</div>
                <p className="text-sm text-gray-600 mt-4">{plan.desc || ""}</p>
              </div>

              <Link href={`/signup/distributor/?plan_id=${plan.id}`}>
                <Button 
                  variant={isRecommended(plan.name) ? "default" : "outline"}
                  className={`w-full py-6 text-lg font-semibold ${
                    isRecommended(plan.name)
                      ? "bg-primary hover:bg-primary/90 text-white shadow-lg"
                      : "border-2 hover:bg-primary hover:text-white hover:border-primary"
                  }`}
                >
                  Get Started
                </Button>
              </Link>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};