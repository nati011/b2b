import { Star } from "lucide-react";

interface Testimonial {
  id: string;
  name: string;
  position: string;
  company: string;
  content: string;
  rating: number;
  image: string;
}

export default function Testimonials() {
  const testimonials: Testimonial[] = [
    {
      id: "1",
      name: "Sarah Johnson",
      position: "Retail Operations Manager",
      company: "Urban Boutique Group",
      content:
        "Efoyeta's wholesale platform has transformed our inventory management. Their curated collections consistently outperform other suppliers, and the analytics dashboard gives us valuable insights for better purchasing decisions.",
      rating: 5,
      image:
        "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=400&q=80",
    },
    {
      id: "2",
      name: "Michael Chen",
      position: "Owner",
      company: "Modern Home Store",
      content:
        "Since partnering with Efoyeta two years ago, our average basket size has increased by 28%. The quality of their products and reliability of their supply chain has made them our go-to wholesale partner.",
      rating: 5,
      image:
        "https://images.unsplash.com/photo-1560250097-0b93528c311a?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=400&q=80",
    },
    {
      id: "3",
      name: "Emma Rodriguez",
      position: "Purchasing Director",
      company: "HomeStyle Inc.",
      content:
        "The flexible minimum order quantities and seasonal discounts have allowed us to experiment with new product categories with minimal risk. Their account management team provides exceptional support.",
      rating: 4,
      image:
        "https://images.unsplash.com/photo-1580489944761-15a19d654956?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=400&q=80",
    },
  ];

  return (
    <section className="py-16 container mx-auto px-4">
      <div className="flex flex-col items-center mb-12">
        <h2 className="text-3xl font-medium mb-4">What Our Retailers Say</h2>
        <div className="h-1 w-20 bg-primary mb-6"></div>
        <p className="text-lg text-center text-muted-foreground mb-4 max-w-2xl">
          Trusted by hundreds of retailers across the country
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {testimonials.map((testimonial) => (
          <div
            key={testimonial.id}
            className="bg-secondary/20 p-6 rounded-lg border"
          >
            <div className="flex items-center mb-4">
              {[...Array(5)].map((_, i) => (
                <Star
                  key={i}
                  className={`h-5 w-5 ${
                    i < testimonial.rating
                      ? "fill-primary text-primary"
                      : "fill-muted text-muted"
                  }`}
                />
              ))}
            </div>

            <p className="text-muted-foreground mb-6 italic">
              "{testimonial.content}"
            </p>

            <div className="flex items-center">
              <div className="w-12 h-12 rounded-full overflow-hidden mr-4">
                <img
                  src={testimonial.image}
                  alt={testimonial.name}
                  className="w-full h-full object-cover"
                />
              </div>
              <div>
                <h4 className="font-medium">{testimonial.name}</h4>
                <p className="text-sm text-muted-foreground">
                  {testimonial.position}, {testimonial.company}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
