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
      title: "Wide Variety of Products",
      description: "Shop a broad range of carefully selected items for every need.",
      icon: <Shield className="h-8 w-8 text-primary" />,
    },
    {
      id: "2",
      title: "Honest Delivery Times",
      description: "Clear, reliable delivery timelines you can plan around.",
      icon: <Award className="h-8 w-8 text-primary" />,
    },
    {
      id: "3",
      title: "Trusted Partner Suppliers",
      description: "We work closely with reliable suppliers to maintain quality.",
      icon: <Handshake className="h-8 w-8 text-primary" />,
    },
    {
      id: "4",
      title: "Support via WhatsApp and Telegram",
      description: "Get quick help and updates from our friendly support team.",
      icon: <BookOpen className="h-8 w-8 text-primary" />,
    },
  ];

  return (
    <section className="py-20 bg-gradient-to-b from-white to-gray-50">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <p className="text-sm font-semibold text-primary uppercase tracking-wider mb-3">
            Why shop with us
          </p>
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Why Shop With Efoyeta Store
          </h2>
          <div className="h-1 w-24 bg-primary mx-auto mb-6 rounded-full"></div>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Wide selection, honest delivery times, trusted suppliers, and responsive support.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {features.map((feature) => (
            <div
              key={feature.id}
              className="bg-white p-8 rounded-2xl border border-gray-100 shadow-sm hover:shadow-xl transition-all duration-300 hover:-translate-y-2 text-center group"
            >
              <div className="flex justify-center mb-6">
                <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                  {feature.icon}
                </div>
              </div>
              <h3 className="text-xl font-semibold mb-3 text-gray-900">{feature.title}</h3>
              <p className="text-gray-600 leading-relaxed">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
