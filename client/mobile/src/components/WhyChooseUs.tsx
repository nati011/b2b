import { Shield, Award, Handshake, BookOpen } from 'lucide-react';
import { motion } from 'framer-motion';

interface Feature {
  id: string;
  title: string;
  description: string;
  icon: React.ReactNode;
}

export const WhyChooseUs = () => {
  const features: Feature[] = [
    {
      id: '1',
      title: 'Wide Variety of Products',
      description: 'Shop a broad range of carefully selected items for every need.',
      icon: <Shield className="h-8 w-8 text-primary" />,
    },
    {
      id: '2',
      title: 'Honest Delivery Times',
      description: 'Clear, reliable delivery timelines you can plan around.',
      icon: <Award className="h-8 w-8 text-primary" />,
    },
    {
      id: '3',
      title: 'Trusted Partner Suppliers',
      description: 'We work closely with reliable suppliers to maintain quality.',
      icon: <Handshake className="h-8 w-8 text-primary" />,
    },
    {
      id: '4',
      title: 'Support via WhatsApp & Email',
      description: 'Get quick help and updates from our friendly support team.',
      icon: <BookOpen className="h-8 w-8 text-primary" />,
    },
  ];

  return (
    <section className="py-12 bg-gradient-to-b from-background to-muted/30">
      <div className="container mx-auto px-4">
        <div className="text-center mb-8">
          <p className="text-xs font-semibold text-primary uppercase tracking-wider mb-2">
            Why shop with us
          </p>
          <h2 className="text-2xl sm:text-3xl font-bold text-foreground mb-3">
            Why Shop With Efoyetastore
          </h2>
          <div className="h-1 w-16 bg-primary mx-auto mb-4 rounded-full"></div>
          <p className="text-sm text-muted-foreground max-w-md mx-auto">
            Wide selection, honest delivery times, trusted suppliers, and responsive support.
          </p>
        </div>

        <div className="grid grid-cols-2 gap-4">
          {features.map((feature, index) => (
            <motion.div
              key={feature.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="bg-card p-4 rounded-xl border border-border shadow-sm hover:shadow-md transition-all duration-300 text-center"
            >
              <div className="flex justify-center mb-3">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                  {feature.icon}
                </div>
              </div>
              <h3 className="text-sm font-semibold mb-2 text-foreground">{feature.title}</h3>
              <p className="text-xs text-muted-foreground leading-relaxed">{feature.description}</p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
};

