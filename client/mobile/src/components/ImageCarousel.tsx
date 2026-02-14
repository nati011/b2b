import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface ImageCarouselProps {
  images: string[];
  alt: string;
}

export const ImageCarousel: React.FC<ImageCarouselProps> = ({ images, alt }) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [direction, setDirection] = useState(0);

  // Handle empty images array
  const safeImages = images && images.length > 0 ? images : ['/placeholder.png'];
  const currentImage = safeImages[currentIndex] || '/placeholder.png';

  const slideVariants = {
    enter: (direction: number) => ({
      x: direction > 0 ? '100%' : '-100%',
      opacity: 0,
    }),
    center: {
      x: 0,
      opacity: 1,
    },
    exit: (direction: number) => ({
      x: direction < 0 ? '100%' : '-100%',
      opacity: 0,
    }),
  };

  const navigate = (newDirection: number) => {
    setDirection(newDirection);
    setCurrentIndex((prev) => {
      let next = prev + newDirection;
      if (next < 0) next = safeImages.length - 1;
      if (next >= safeImages.length) next = 0;
      return next;
    });
  };

  const handleDragEnd = (_: any, { offset, velocity }: any) => {
    const swipeThreshold = 50;
    const swipeVelocity = 500;

    if (offset.x < -swipeThreshold || velocity.x < -swipeVelocity) {
      navigate(1);
    } else if (offset.x > swipeThreshold || velocity.x > swipeVelocity) {
      navigate(-1);
    }
  };

  return (
    <div className="relative aspect-[3/4] bg-secondary overflow-hidden">
      <AnimatePresence initial={false} custom={direction} mode="popLayout">
        <motion.img
          key={currentIndex}
          src={currentImage}
          alt={`${alt} - Image ${currentIndex + 1}`}
          custom={direction}
          variants={slideVariants}
          initial="enter"
          animate="center"
          exit="exit"
          transition={{ type: 'spring', stiffness: 300, damping: 30 }}
          drag="x"
          dragConstraints={{ left: 0, right: 0 }}
          dragElastic={0.2}
          onDragEnd={handleDragEnd}
          className="absolute inset-0 w-full h-full object-cover cursor-grab active:cursor-grabbing"
          onError={(e) => {
            // Fallback to placeholder if image fails to load
            e.currentTarget.src = '/placeholder.png';
          }}
        />
      </AnimatePresence>

      {/* Navigation buttons */}
      {safeImages.length > 1 && (
        <>
          <button
            onClick={() => navigate(-1)}
            className="absolute left-3 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/90 backdrop-blur-sm rounded-full flex items-center justify-center btn-press shadow-lg opacity-0 hover:opacity-100 transition-opacity"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <button
            onClick={() => navigate(1)}
            className="absolute right-3 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/90 backdrop-blur-sm rounded-full flex items-center justify-center btn-press shadow-lg opacity-0 hover:opacity-100 transition-opacity"
          >
            <ChevronRight className="w-5 h-5" />
          </button>
        </>
      )}

      {/* Dots */}
      {safeImages.length > 1 && (
        <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
          {safeImages.map((_, index) => (
            <button
              key={index}
              onClick={() => {
                setDirection(index > currentIndex ? 1 : -1);
                setCurrentIndex(index);
              }}
              className={`w-2 h-2 rounded-full transition-all ${
                index === currentIndex
                  ? 'bg-foreground w-6'
                  : 'bg-foreground/30 hover:bg-foreground/50'
              }`}
            />
          ))}
        </div>
      )}
    </div>
  );
};
