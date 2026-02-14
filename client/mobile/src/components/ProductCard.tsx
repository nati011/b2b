import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Product } from '@/types/product';
import { Badge } from '@/components/ui/badge';

interface ProductCardProps {
  product: Product;
  index?: number;
}

export const ProductCard: React.FC<ProductCardProps> = ({ product, index = 0 }) => {
  const hasDiscount = product.originalPrice && product.originalPrice > product.price;
  const discountPercent = hasDiscount 
    ? Math.round((1 - product.price / product.originalPrice!) * 100) 
    : 0;

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.4 }}
    >
      <Link to={`/product/${product.id}`} className="block group">
        <div className="bg-card border border-border rounded-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
          {/* Image Container */}
          <div className="aspect-square bg-secondary overflow-hidden relative">
            <img
              src={product.images[0]}
              alt={product.name}
              className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              loading="lazy"
            />
            {/* Badges */}
            <div className="absolute top-2 left-2 flex flex-col gap-1">
              {product.isNew && (
                <Badge variant="default" className="text-[10px] px-2 py-0.5">
                  NEW
                </Badge>
              )}
              {hasDiscount && (
                <Badge variant="destructive" className="text-[10px] px-2 py-0.5">
                  -{discountPercent}%
                </Badge>
              )}
            </div>
          </div>

          {/* Content */}
          <div className="p-3 space-y-1.5">
            <p className="text-[10px] text-muted-foreground font-medium tracking-wide uppercase">
              {product.brand}
            </p>
            <h3 className="text-sm font-semibold leading-snug line-clamp-2 text-foreground">
              {product.name}
            </h3>
            <p className="text-xs text-muted-foreground line-clamp-1">
              {product.category}
            </p>
            <div className="flex items-center gap-2 pt-1">
              <span className="text-base font-bold text-primary">
                {product.price.toLocaleString()} ETB
              </span>
              {hasDiscount && (
                <span className="text-xs text-muted-foreground line-through">
                  {product.originalPrice?.toLocaleString()} ETB
                </span>
              )}
            </div>
          </div>
        </div>
      </Link>
    </motion.div>
  );
};
