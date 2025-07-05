import { Shield, Award, Handshake, BookOpen } from "lucide-react";

interface Feature {
  id: string;
  title: string;
  description: string;
  icon: React.ReactNode;
}

export default function WhyChooseUs() {
  const features: Feature[] = [
    {
      id: "1",
      title: "Wholesale Pricing",
      description:
        "Competitive pricing with volume discounts starting at just $500 minimum order value",
      icon: <Shield className="h-8 w-8 text-primary" />,
    },
    {
      id: "2",
      title: "Quality Guaranteed",
      description:
        "Every product backed by our satisfaction guarantee with hassle-free returns",
      icon: <Award className="h-8 w-8 text-primary" />,
    },
    {
      id: "3",
      title: "Retailer Support",
      description:
        "Dedicated account managers and marketing materials to boost your sales",
      icon: <Handshake className="h-8 w-8 text-primary" />,
    },
    {
      id: "4",
      title: "Trend Forecasting",
      description:
        "Quarterly trend reports and data-driven recommendations for your store",
      icon: <BookOpen className="h-8 w-8 text-primary" />,
    },
  ];

  return (
    <section className="py-16 bg-secondary/30">
      <div className="container mx-auto px-4">
        <div className="flex flex-col items-center mb-12">
          <h2 className="text-3xl font-medium mb-4">Why Partner With Us</h2>
          <div className="h-1 w-20 bg-primary mb-6"></div>
          <p className="text-lg text-center text-muted-foreground mb-4 max-w-2xl">
            Discover the advantages of making Efoyeta your wholesale partner
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {features.map((feature) => (
            <div
              key={feature.id}
              className="bg-background p-6 rounded-lg border text-center shadow-sm"
            >
              <div className="flex justify-center mb-4">{feature.icon}</div>
              <h3 className="text-xl font-medium mb-3">{feature.title}</h3>
              <p className="text-muted-foreground">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
