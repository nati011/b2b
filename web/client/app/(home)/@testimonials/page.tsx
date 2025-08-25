import { Check } from "lucide-react";
import { Button } from "@/components/ui/button";

interface PricingTier {
  id: string;
  name: string;
  price: string;
  period: string;
  description: string;
  features: string[];
  popular?: boolean;
}

export default function PricingSection(){
  const pricingTiers: PricingTier[] = [
    {
      id: "monthly",
      name: "Monthly",
      price: "$2,500",
      period: "per month",
      description: "Perfect for getting started with minimal commitment",
      features: [
        "Full product catalog access",
        "Basic marketing materials",
        "Email support",
        "Monthly reporting"
      ]
    },
    {
      id: "yearly",
      name: "Annual",
      price: "$11,500",
      period: "per year",
      description: "Most popular choice for established distributors",
      features: [
        "Everything in Monthly",
        "Priority customer support",
        "Quarterly business reviews",
        "Advanced marketing materials",
        "Volume discounts"
      ],
      popular: true
    },
    {
      id: "two-year",
      name: "Two Year",
      price: "$22,000",
      period: "for 24 months",
      description: "Best value for long-term partnerships",
      features: [
        "Everything in Annual",
        "Dedicated account manager",
        "Custom marketing support",
        "Exclusive product previews",
        "Maximum volume discounts",
        "Territory protection"
      ]
    }
  ];

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

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {pricingTiers.map(tier => (
            <div 
              key={tier.id} 
              className={`relative bg-background p-8 rounded-lg border ${
                tier.popular ? "border-primary shadow-lg scale-105" : "border-border"
              }`}
            >
              {tier.popular && (
                <div className="absolute -top-4 left-1/2 transform -translate-x-1/2 bg-primary text-primary-foreground px-4 py-1 rounded-full text-sm font-medium">
                  Recommended
                </div>
              )}
              
              <div className="text-center mb-6">
                <h3 className="text-xl font-medium mb-2">{tier.name}</h3>
                <div className="text-3xl font-bold text-primary mb-1">{tier.price}</div>
                <div className="text-sm text-muted-foreground">{tier.period}</div>
                <p className="text-sm text-muted-foreground mt-4">{tier.description}</p>
              </div>

              <ul className="space-y-3 mb-8">
                {tier.features.map((feature, index) => (
                  <li key={index} className="flex items-center">
                    <Check className="h-4 w-4 text-primary mr-3 flex-shrink-0" />
                    <span className="text-sm">{feature}</span>
                  </li>
                ))}
              </ul>

              <Button 
                variant={tier.popular ? "default" : "secondary"}
                className={tier.popular ? "default" : "secondary"}
              >
                Get Started
              </Button>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};