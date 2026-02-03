import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Catalogue } from '@/lib/types';
import { Badge } from '@/components/ui/badge';

interface CatalogueCardProps {
  catalogue: Catalogue;
  index?: number;
}

export const CatalogueCard: React.FC<CatalogueCardProps> = ({ catalogue, index = 0 }) => {
  // Get the first configurable product for display
  const firstConfigurable = catalogue.configurables?.[0];
  const imageUrl = firstConfigurable?.images?.[0]?.ImageUrl || catalogue.images?.[0]?.ImageUrl || '';
  const price = firstConfigurable?.price || catalogue.price || 0;
  const productName = firstConfigurable?.name || catalogue.name;
  
  // Ensure we have a valid numeric ID for the product detail page
  const productId = catalogue.id || firstConfigurable?.id;
  
  // Only render link if we have a valid numeric ID
  if (!productId || (typeof productId !== 'number' && isNaN(Number(productId)))) {
    console.warn('CatalogueCard: Invalid product ID', { catalogue, productId });
    return (
      <motion.div
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: index * 0.05, duration: 0.4 }}
      >
        <div className="bg-card border border-border rounded-md overflow-hidden opacity-50">
          <div className="aspect-square bg-secondary overflow-hidden relative">
            <img
              src={imageUrl}
              alt={productName}
              className="w-full h-full object-cover"
              loading="lazy"
            />
          </div>
          <div className="p-3 space-y-1.5">
            <h3 className="text-sm font-semibold leading-snug line-clamp-2 text-foreground">
              {productName}
            </h3>
            <p className="text-xs text-muted-foreground line-clamp-2">
              {catalogue.desc || firstConfigurable?.desc || ''}
            </p>
            <div className="flex items-center gap-2 pt-1">
              <span className="text-base font-bold text-primary">
                {price.toLocaleString()} ETB
              </span>
            </div>
          </div>
        </div>
      </motion.div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05, duration: 0.4 }}
    >
      <Link to={`/product/${productId}`} className="block group">
        <div className="bg-card border border-border rounded-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
          {/* Image Container */}
          <div className="aspect-square bg-secondary overflow-hidden relative">
            <img
              src={imageUrl}
              alt={productName}
              className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              loading="lazy"
            />
            {/* Badges */}
            <div className="absolute top-2 left-2 flex flex-col gap-1">
              {catalogue.is_active && (
                <Badge variant="default" className="text-[10px] px-2 py-0.5">
                  Active
                </Badge>
              )}
            </div>
          </div>

          {/* Content */}
          <div className="p-3 space-y-1.5">
            <h3 className="text-sm font-semibold leading-snug line-clamp-2 text-foreground">
              {productName}
            </h3>
            <p className="text-xs text-muted-foreground line-clamp-2">
              {catalogue.desc || firstConfigurable?.desc || ''}
            </p>
            <div className="flex items-center gap-2 pt-1">
              <span className="text-base font-bold text-primary">
                {price.toLocaleString()} ETB
              </span>
            </div>
          </div>
        </div>
      </Link>
    </motion.div>
  );
};

