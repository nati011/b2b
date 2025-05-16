import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselPrevious,
    CarouselNext
} from "@/components/ui/carousel";

export const Hero = () => {
    const slides = [
        {
            id: 1,
            title: "Premium home goods for modern retailers",
            subtitle: "Curated collections of high-margin products",
            description: "Elevate your retail space with our exclusive wholesale collections designed to maximize profitability and customer engagement.",
            cta: "View Wholesale Catalog",
            secondaryCta: "Become a Partner",
            image: "https://images.unsplash.com/photo-1532372320572-cda25653a26d?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
            bgColor: "bg-white"
        },
        {
            id: 2,
            title: "Exclusive wholesale pricing",
            subtitle: "Up to 50% off retail prices",
            description: "Partner with us to access exclusive wholesale pricing and boost your profit margins with our premium product lines.",
            cta: "See Pricing",
            secondaryCta: "Contact Sales",
            image: "https://images.unsplash.com/photo-1556228453-efd6c1ff04f6?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
            bgColor: "bg-white"
        },
        {
            id: 3,
            title: "Full retailer support",
            subtitle: "Marketing materials included",
            description: "Get access to ready-to-use marketing materials, product training, and dedicated account management.",
            cta: "Learn More",
            secondaryCta: "Schedule Demo",
            image: "https://images.unsplash.com/photo-1586023492125-27b2c045efd7?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
            bgColor: "bg-white"
        }
    ];

    const [api, setApi] = useState<any>(null);
    const [current, setCurrent] = useState(0);

    useEffect(() => {
        if (!api) return;

        const interval = setInterval(() => {
            api.scrollNext();
        }, 5000); // Auto-rotate every 5 seconds

        api.on("select", () => {
            setCurrent(api.selectedScrollSnap());
        });

        return () => {
            clearInterval(interval);
            api.off("select");
        };
    }, [api]);

    return (
        <section className="w-full overflow-hidden">
            {/* Main Hero Carousel */}
            <Carousel
                className="w-full"
                opts={{
                    loop: true,
                    duration: 50
                }}
                setApi={setApi}
            >
                <CarouselContent>
                    {slides.map((slide) => (
                        <CarouselItem key={slide.id}>
                            <div className={`w-full ${slide.bgColor} transition-all duration-500`}>
                                <div className="container mx-auto grid grid-cols-1 md:grid-cols-2 min-h-[600px]">
                                    <div className="flex flex-col justify-center px-8 py-16 order-2 md:order-1">
                                        <span className="bg-primary/10 text-primary px-3 py-1 rounded-md text-sm font-medium w-fit mb-4">
                                            Wholesale Only
                                        </span>

                                        <h2 className="text-lg md:text-xl text-primary/80 font-normal mb-2 animate-slideUp">
                                            {slide.subtitle}
                                        </h2>

                                        <h1 className="text-3xl md:text-4xl lg:text-5xl font-medium tracking-tight mb-6 animate-fadeIn">
                                            {slide.title}
                                        </h1>

                                        <p className="text-lg text-muted-foreground mb-8 max-w-md animate-slideUp">
                                            {slide.description}
                                        </p>

                                        <div className="flex flex-col sm:flex-row gap-4 animate-slideUp" style={{ animationDelay: "0.2s" }}>
                                            <Button size="lg" className="bg-primary hover:bg-primary/90" asChild>
                                                <a href="#products">
                                                    {slide.cta} <ArrowRight className="ml-2 h-4 w-4" />
                                                </a>
                                            </Button>
                                            <Button variant="outline" size="lg">
                                                {slide.secondaryCta}
                                            </Button>
                                        </div>
                                    </div>

                                    <div className="relative order-1 md:order-2">
                                        <div className="h-[300px] md:h-[600px] w-full">
                                            <div
                                                className="h-full w-full bg-cover bg-center animate-fadeIn transition-all duration-500"
                                                style={{ backgroundImage: `url(${slide.image})` }}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </CarouselItem>
                    ))}
                </CarouselContent>

                <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 z-10 flex gap-2 justify-center">
                    {slides.map((_, index) => (
                        <button
                            key={index}
                            onClick={() => api?.scrollTo(index)}
                            className={`w-3 h-3 rounded-full transition-all ${current === index ? "bg-primary w-6" : "bg-primary/30"
                                }`}
                            aria-label={`Go to slide ${index + 1}`}
                        />
                    ))}
                </div>

                <div className="hidden md:flex absolute bottom-8 right-8 z-10 gap-2">
                    <CarouselPrevious className="static translate-y-0 h-10 w-10" />
                    <CarouselNext className="static translate-y-0 h-10 w-10" />
                </div>
            </Carousel>

            {/* Highlights Section */}
            {/* <div className="container mx-auto px-4 py-12 grid grid-cols-1 md:grid-cols-3 gap-8 border-b">
                <div className="flex items-center animate-slideUp" style={{ animationDelay: "0.3s" }}>
                    <div className="rounded-full bg-primary/10 p-3 mr-4">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary"><path d="M21.2 8.4c.5.38.8.96.8 1.6 0 1.1-.9 2-2 2H3c-1.1 0-2-.9-2-2 0-.64.3-1.22.8-1.6" /><path d="m5.5 8.4 1.1-3.36a1 1 0 0 1 .95-.64h8.9c.45 0 .85.29.95.64L18.5 8.4" /><path d="M4 14h.01" /><path d="M8 14h.01" /><path d="M12 14h.01" /><path d="M16 14h.01" /><path d="M20 14h.01" /><path d="M4 19h.01" /><path d="M8 19h.01" /><path d="M12 19h.01" /><path d="M16 19h.01" /><path d="M20 19h.01" /></svg>
                    </div>
                    <div>
                        <h3 className="font-medium">Wholesale Pricing</h3>
                        <p className="text-sm text-muted-foreground">Up to 50% off retail</p>
                    </div>
                </div>
                <div className="flex items-center animate-slideUp" style={{ animationDelay: "0.4s" }}>
                    <div className="rounded-full bg-primary/10 p-3 mr-4">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary"><path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a2 2 0 0 1-2 2Zm0 0a2 2 0 0 1-2-2v-9c0-1.1.9-2 2-2h2" /><path d="M18 14h-8" /><path d="M15 18h-5" /><path d="M10 6h8v4h-8V6Z" /></svg>
                    </div>
                    <div>
                        <h3 className="font-medium">Low MOQ</h3>
                        <p className="text-sm text-muted-foreground">Starting at $500</p>
                    </div>
                </div>
                <div className="flex items-center animate-slideUp" style={{ animationDelay: "0.5s" }}>
                    <div className="rounded-full bg-primary/10 p-3 mr-4">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary"><path d="M7 10v12" /><path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H4a2 2 0 0 1-2-2v-8a2 2 0 0 1 2-2h2.76a2 2 0 0 0 1.79-1.11L12 2h0a3.13 3.13 0 0 1 3 3.88Z" /></svg>
                    </div>
                    <div>
                        <h3 className="font-medium">Retailer Support</h3>
                        <p className="text-sm text-muted-foreground">Marketing materials included</p>
                    </div>
                </div>
            </div> */}
        </section>
    );
};
