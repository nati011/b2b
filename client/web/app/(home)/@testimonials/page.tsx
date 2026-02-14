"use client";
import { Button } from "@/components/ui/button";
import Link from "next/link";

interface PricingPlan {
  id: string;
  name: string;
  price: number;
  period: string;
  isRecommended: boolean;
}

export default function PricingSection(){
  // Hardcoded pricing plans
  const plans: PricingPlan[] = [
    {
      id: "monthly",
      name: "Monthly",
      price: 500,
      period: "per month",
      isRecommended: false,
    },
    {
      id: "yearly",
      name: "Year",
      price: 4500,
      period: "per year",
      isRecommended: true,
    },
  ];

  const formatPrice = (price: number) => `${price.toLocaleString()} ETB`;

  return (
    <section id="pricing" className="py-20 bg-white scroll-mt-20">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <p className="text-sm font-semibold text-primary uppercase tracking-wider mb-3">
            Pricing Plans
          </p>
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Supplier Pricing
          </h2>
          <div className="h-1 w-24 bg-primary mx-auto mb-6 rounded-full"></div>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Choose the package that fits your business needs and growth goals
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
          {plans.map(plan => (
            <div 
              key={plan.id}
              className={`relative bg-white p-8 rounded-2xl border transition-all duration-300 hover:shadow-xl ${
                plan.isRecommended
                  ? "border-primary shadow-lg scale-105 ring-2 ring-primary/20" 
                  : "border-gray-200 shadow-sm hover:border-primary/50"
              }`}
            >
              {plan.isRecommended && (
                <div className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-primary text-white px-6 py-2 rounded-full text-sm font-semibold shadow-lg">
                  Recommended
                </div>
              )}
              <div className="text-center mb-8">
                <h3 className="text-2xl font-bold text-gray-900 mb-4">{plan.name}</h3>
                <div className="text-4xl font-bold text-primary mb-2">{formatPrice(plan.price)}</div>
                <div className="text-sm text-gray-600 mb-4">{plan.period}</div>
              </div>

              <Link href="/contact">
                <Button 
                  variant={plan.isRecommended ? "default" : "outline"}
                  className={`w-full py-6 text-lg font-semibold ${
                    plan.isRecommended
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