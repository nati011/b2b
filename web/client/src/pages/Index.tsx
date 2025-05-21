
import { Hero } from "@/components/Hero";
import { ProductGrid } from "@/components/ProductGrid";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { Testimonials } from "@/components/Testimonials";
import { WhyChooseUs } from "@/components/WhyChooseUs";
import Contact from "@/components/ContactUs";

const Index = () => {
  return (
    <>

      <main>
        <Hero />
        <ProductGrid />
        <WhyChooseUs />
        <Testimonials />
      </main>
    </>
  );
};

export default Index;
