"use client";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useEffect } from "react";
import usePlanstore from "@/lib/store/usePricingPlan";
import Link from "next/link";

export default function PricingSection(){
  const { plans, loading, error, fetchPricingPlan } = usePlanstore();

  useEffect(() => {
    void fetchPricingPlan();
  }, [fetchPricingPlan]);

  const formatPrice = (price: number) => `${price.toLocaleString()} ETB`;
  const termToPeriod = (termInMonth: number) => {
    if (termInMonth === 1) return "per month";
    if (termInMonth === 12) return "per year";
    return `for ${termInMonth} months`;
  };

  const isRecommended = (name: string) => /annual|annunal/i.test(name);

  return (
    <section className="py-16 bg-secondary/30">
      <div className="container mx-auto px-4">
        <div className="flex flex-col items-center mb-12">
          <h2 className="text-3xl font-medium mb-4">Distributor Pricing</h2>
          <div className="h-1 w-20 bg-primary mb-6"></div>
          <p className="text-lg text-center text-muted-foreground mb-4 max-w-2xl">
            Choose the package that fits your business needs and growth goals
          </p>
        </div>

        {error && (
          <div className="text-center text-destructive mb-8">{error}</div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {loading && plans.length === 0 && (
            Array.from({ length: 3 }).map((_, idx) => (
              <div key={idx} className="relative bg-background p-8 rounded-lg border border-border">
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
              key={plan.name}
              className={`relative bg-background p-8 rounded-lg border ${isRecommended(plan.name) ? "border-primary shadow-lg scale-105" : "border-border"}`}
            >
              {isRecommended(plan.name) && (
                <div className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-primary text-primary-foreground px-4 py-1 rounded-full text-sm font-medium">
                  Recommended
                </div>
              )}
              <div className="text-center mb-6">
                <h3 className="text-xl font-medium mb-2">{plan.name}</h3>
                <div className="text-3xl font-bold text-primary mb-1">{formatPrice(plan.price)}</div>
                <div className="text-sm text-muted-foreground">{termToPeriod(plan.term_in_month)}</div>
                <p className="text-sm text-muted-foreground mt-4">{String((plan as any).desc ?? "")}</p>
              </div>

              <Link href={`/signup/distributor/?plan_id=${plan.id}`}>
              <Button 
                variant="default"
                className="w-full"
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