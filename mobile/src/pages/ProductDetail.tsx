import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { ArrowLeft, ChevronDown, ChevronUp, Heart } from 'lucide-react';
import { getProductById } from '@/data/products';
import { ImageCarousel } from '@/components/ImageCarousel';
import { useCart } from '@/context/CartContext';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

const ProductDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { addItem } = useCart();
  
  const product = getProductById(id || '');
  
  const [selectedSize, setSelectedSize] = useState<string | null>(null);
  const [selectedColor, setSelectedColor] = useState<string | null>(null);
  const [detailsOpen, setDetailsOpen] = useState(true);
  const [isWishlisted, setIsWishlisted] = useState(false);

  if (!product) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-muted-foreground">Product not found</p>
      </div>
    );
  }

  const handleAddToCart = () => {
    if (!selectedSize || !selectedColor) return;
    addItem(product, selectedSize, selectedColor);
  };

  const canAddToCart = selectedSize && selectedColor;

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-28"
    >
      {/* Header */}
      <div className="sticky top-0 z-20 bg-background/80 backdrop-blur-sm">
        <div className="flex items-center justify-between p-4">
          <button
            onClick={() => navigate(-1)}
            className="p-2 -ml-2 hover:bg-secondary rounded-sm btn-press"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <button
            onClick={() => setIsWishlisted(!isWishlisted)}
            className="p-2 -mr-2 hover:bg-secondary rounded-sm btn-press"
          >
            <Heart className={cn("w-5 h-5", isWishlisted && "fill-foreground")} />
          </button>
        </div>
      </div>

      {/* Image Carousel */}
      <ImageCarousel images={product.images} alt={product.name} />

      {/* Product Info */}
      <div className="p-4 space-y-6">
        {/* Title & Price */}
        <div>
          <p className="text-sm text-muted-foreground tracking-wide uppercase mb-1">
            {product.brand}
          </p>
          <h1 className="font-display text-2xl mb-2">{product.name}</h1>
          <div className="flex items-center gap-3">
            <span className="text-xl font-medium">${product.price.toLocaleString()}</span>
            {product.originalPrice && (
              <span className="text-muted-foreground line-through">
                ${product.originalPrice.toLocaleString()}
              </span>
            )}
          </div>
        </div>

        {/* Description */}
        <p className="text-sm text-muted-foreground leading-relaxed">
          {product.description}
        </p>

        {/* Color Selection */}
        <div>
          <p className="text-sm font-medium mb-3">
            Color: <span className="text-muted-foreground font-normal">{selectedColor || 'Select'}</span>
          </p>
          <div className="flex gap-2 flex-wrap">
            {product.colors.map((color) => (
              <button
                key={color}
                onClick={() => setSelectedColor(color)}
                className={cn(
                  "px-4 py-2 text-sm border rounded-sm btn-press transition-colors",
                  selectedColor === color
                    ? "border-primary bg-primary text-primary-foreground"
                    : "border-border hover:border-primary"
                )}
              >
                {color}
              </button>
            ))}
          </div>
        </div>

        {/* Size Selection */}
        <div>
          <p className="text-sm font-medium mb-3">
            Size: <span className="text-muted-foreground font-normal">{selectedSize || 'Select'}</span>
          </p>
          <div className="flex gap-2 flex-wrap">
            {product.sizes.map((size) => (
              <button
                key={size}
                onClick={() => setSelectedSize(size)}
                className={cn(
                  "min-w-[48px] h-12 px-3 text-sm border rounded-sm btn-press transition-colors",
                  selectedSize === size
                    ? "border-primary bg-primary text-primary-foreground"
                    : "border-border hover:border-primary"
                )}
              >
                {size}
              </button>
            ))}
          </div>
        </div>

        {/* Details & Fit */}
        <div className="border-t border-border pt-4">
          <button
            onClick={() => setDetailsOpen(!detailsOpen)}
            className="flex items-center justify-between w-full py-2 btn-press"
          >
            <span className="text-sm font-medium">Details & Fit</span>
            {detailsOpen ? (
              <ChevronUp className="w-4 h-4" />
            ) : (
              <ChevronDown className="w-4 h-4" />
            )}
          </button>
          
          <motion.div
            initial={false}
            animate={{ height: detailsOpen ? 'auto' : 0, opacity: detailsOpen ? 1 : 0 }}
            className="overflow-hidden"
          >
            <div className="py-3 space-y-4">
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-2">Details</p>
                <ul className="space-y-1">
                  {product.details.map((detail, index) => (
                    <li key={index} className="text-sm text-muted-foreground">
                      • {detail}
                    </li>
                  ))}
                </ul>
              </div>
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-2">Fit</p>
                <p className="text-sm text-muted-foreground">{product.fit}</p>
              </div>
            </div>
          </motion.div>
        </div>
      </div>

      {/* Fixed Add to Cart Button */}
      <div className="fixed bottom-14 left-0 right-0 p-4 bg-background border-t border-border safe-bottom">
        <Button
          onClick={handleAddToCart}
          disabled={!canAddToCart}
          className="w-full h-12 text-sm tracking-wide btn-press rounded-sm"
        >
          {canAddToCart ? 'Add to Cart' : 'Select Size & Color'}
        </Button>
      </div>
    </motion.div>
  );
};

export default ProductDetail;
