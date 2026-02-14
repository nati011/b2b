import { useRef } from 'react';
import { motion } from 'framer-motion';
import { brands } from '@/data/products';
import { ChevronLeft, ChevronRight } from 'lucide-react';

export const BrandCarousel = () => {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const scrollAmount = 200;
      scrollRef.current.scrollBy({
        left: direction === 'left' ? -scrollAmount : scrollAmount,
        behavior: 'smooth',
      });
    }
  };

  return (
    <section className="py-8 px-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="font-display text-xl">Trending Brands</h2>
        <div className="flex gap-2">
          <button
            onClick={() => scroll('left')}
            className="p-2 border border-border rounded-sm btn-press hover:bg-secondary transition-colors"
          >
            <ChevronLeft className="w-4 h-4" />
          </button>
          <button
            onClick={() => scroll('right')}
            className="p-2 border border-border rounded-sm btn-press hover:bg-secondary transition-colors"
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      </div>

      <div
        ref={scrollRef}
        className="flex gap-4 overflow-x-auto scrollbar-hide snap-x snap-mandatory -mx-4 px-4"
      >
        {brands.map((brand, index) => (
          <motion.div
            key={brand.id}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
            className="flex-shrink-0 snap-start"
          >
            <div className="w-32 h-32 sm:w-36 sm:h-36 img-soft rounded-sm flex items-center justify-center overflow-hidden btn-press cursor-pointer">
              <img
                src={brand.logo}
                alt={brand.name}
                className="w-full h-full object-cover opacity-90 hover:opacity-100 transition-opacity"
              />
            </div>
            <p className="text-xs font-medium tracking-wide text-center mt-2 text-muted-foreground">
              {brand.name}
            </p>
          </motion.div>
        ))}
      </div>
    </section>
  );
};
