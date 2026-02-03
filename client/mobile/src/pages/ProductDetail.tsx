import { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  ArrowLeft, 
  Share2, 
  ShoppingCart, 
  Check, 
  Package,
  Info,
  Minus,
  Plus,
  AlertCircle
} from 'lucide-react';
import { GetProduct } from '@/lib/api/product';
import { useCart } from '@/context/CartContext';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { ProductResponse } from '@/lib/api/product';
import type { Product } from '@/lib/types';

const ProductDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { addItem, openCart, items } = useCart();
  
  const [product, setProduct] = useState<ProductResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSize, setSelectedSize] = useState<string | null>(null);
  const [selectedColor, setSelectedColor] = useState<string | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [activeTab, setActiveTab] = useState<'details' | 'specs'>('details');
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [imageLoading, setImageLoading] = useState(false);
  const [loadedImages, setLoadedImages] = useState<Set<string>>(new Set());
  const preloadingRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    const fetchProduct = async () => {
      if (!id) {
        setError('Product ID is required');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);
        
        const productId = Number.parseInt(id, 10);
        if (Number.isNaN(productId) || productId <= 0) {
          throw new Error(`Invalid product ID: ${id}`);
        }
        
        const data = await GetProduct(productId);
        
        if (!data?.id) {
          throw new Error('Invalid product data received');
        }
        
        // Parse attributes if it's a string
        if (data.attributes && typeof data.attributes === 'string') {
          try {
            data.attributes = JSON.parse(data.attributes);
          } catch {
            data.attributes = {};
          }
        }
        
        setProduct(data);
        
        // Auto-select first options if available
        const attrs = data.attributes || {};
        
        const getInitialColors = () => {
          if (!attrs.color) return [];
          if (Array.isArray(attrs.color)) return attrs.color;
          return [attrs.color];
        };
        
        const getInitialSizes = () => {
          if (!attrs.size) return [];
          if (Array.isArray(attrs.size)) return attrs.size;
          return [attrs.size];
        };
        
        const colors = getInitialColors();
        const sizes = getInitialSizes();
        
        if (colors.length === 1) setSelectedColor(colors[0]);
        if (sizes.length === 1) setSelectedSize(sizes[0]);
      } catch (err: any) {
        console.error('Failed to fetch product:', err);
        setError(err.message || 'Failed to load product');
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [id]);

  // Extract product attributes (safe to call even if product is null)
  const attributes = product?.attributes || {};
  
  const getColorArray = () => {
    if (!attributes.color) return [];
    return Array.isArray(attributes.color) ? attributes.color : [attributes.color];
  };
  
  const getSizeArray = () => {
    if (!attributes.size) return [];
    return Array.isArray(attributes.size) ? attributes.size : [attributes.size];
  };
  
  const getImageArray = () => {
    if (!attributes.images) return [];
    
    if (Array.isArray(attributes.images)) {
      return attributes.images
        .map((img: any) => {
          if (typeof img === 'string') return img;
          return img?.ImageUrl || img?.url || img?.imageUrl || '';
        })
        .filter(Boolean);
    }
    
    if (typeof attributes.images === 'string') {
      return [attributes.images];
    }
    
    const singleImage = attributes.images?.ImageUrl || attributes.images?.url || attributes.images?.imageUrl || '';
    return singleImage ? [singleImage] : [];
  };
  
  const colors = getColorArray();
  const sizes = getSizeArray();
  const images = getImageArray();
  const currentImage = images[currentImageIndex] || '/placeholder.png';

  // Preload images when they change
  useEffect(() => {
    if (images.length === 0 || !product) return;
    
    images.forEach((imageUrl) => {
      if (!imageUrl || imageUrl === '/placeholder.png') return;
      
      // Skip if already loaded or currently preloading
      if (loadedImages.has(imageUrl) || preloadingRef.current.has(imageUrl)) {
        return;
      }
      
      // Mark as preloading
      preloadingRef.current.add(imageUrl);
      
      const img = new Image();
      img.onload = () => {
        preloadingRef.current.delete(imageUrl);
        setLoadedImages(current => new Set(current).add(imageUrl));
      };
      img.onerror = () => {
        preloadingRef.current.delete(imageUrl);
      };
      img.src = imageUrl;
    });
  }, [images, product, loadedImages]);

  // Reset loading state when image changes
  useEffect(() => {
    if (!product) return;
    
    // Check if image is already preloaded
    if (loadedImages.has(currentImage)) {
      setImageLoading(false);
      return;
    }
    
    // For placeholder, don't show loading
    if (currentImage === '/placeholder.png') {
      setImageLoading(false);
      return;
    }
    
    // Set loading state for new images
    setImageLoading(true);
    
    // Check if image is already in browser cache (quick check)
    const img = new Image();
    let cancelled = false;
    
    img.onload = () => {
      if (!cancelled) {
        setImageLoading(false);
        setLoadedImages(prev => new Set(prev).add(currentImage));
      }
    };
    
    img.onerror = () => {
      if (!cancelled) {
        setImageLoading(false);
      }
    };
    
    img.src = currentImage;
    
    // Fallback timeout in case image never loads
    const timeout = setTimeout(() => {
      if (!cancelled) {
        setImageLoading(false);
      }
    }, 3000);
    
    return () => {
      cancelled = true;
      clearTimeout(timeout);
    };
  }, [currentImage, loadedImages, product]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-background">
        <div className="text-center space-y-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
          <p className="text-sm text-muted-foreground">Loading product...</p>
        </div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-4 bg-background">
        <div className="text-center space-y-4 max-w-sm">
          <AlertCircle className="w-16 h-16 text-destructive mx-auto" />
          <h2 className="text-xl font-semibold text-foreground">Product Not Found</h2>
          <p className="text-sm text-muted-foreground">
            {error || 'The product you\'re looking for doesn\'t exist or has been removed.'}
          </p>
        </div>
        <Button onClick={() => navigate(-1)} variant="outline" className="mt-6">
          Go Back
        </Button>
      </div>
    );
  }

  const handleAddToCart = () => {
    const finalColor = selectedColor || (colors.length === 1 ? colors[0] : 'Default');
    const finalSize = selectedSize || (sizes.length === 1 ? sizes[0] : 'Standard');
    
    const cartProduct: Product = {
      id: product.id,
      name: product.name,
      description: product.description,
      external_id: product.external_id,
      attributes: product.attributes,
      unit: product.unit,
      is_active: product.is_active,
      supplier_id: product.supplier_id,
      price: product.price,
      total_quantity: product.total_quantity,
      reserved_quantity: product.reserved_quantity,
      available_quantity: product.available_quantity,
      category_ids: product.category_ids,
      created_at: product.created_at,
      updated_at: product.updated_at,
    };
    
    // Add multiple quantities
    for (let i = 0; i < quantity; i++) {
      addItem(cartProduct, finalSize, finalColor);
    }
  };

  const canAddToCart = product.is_active && 
    (colors.length === 0 || selectedColor || colors.length === 1) &&
    (sizes.length === 0 || selectedSize || sizes.length === 1);

  const maxQuantity = product.available_quantity || 99;

  return (
    <div className="min-h-screen bg-background pb-24">
      {/* Header Bar */}
      <div className="sticky top-0 z-30 bg-background/95 backdrop-blur-md border-b border-border">
        <div className="flex items-center justify-between px-4 py-3">
          <button
            onClick={() => navigate(-1)}
            className="p-2 -ml-2 rounded-lg hover:bg-secondary transition-colors"
            aria-label="Go back"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <div className="flex items-center gap-2">
            <button
              onClick={openCart}
              className="relative p-2 rounded-lg hover:bg-secondary transition-colors"
              aria-label="Open cart"
            >
              <ShoppingCart className="w-5 h-5" />
              {items.length > 0 && (
                <span className="absolute -top-1 -right-1 w-5 h-5 bg-primary text-primary-foreground text-xs font-semibold rounded-full flex items-center justify-center">
                  {items.reduce((sum, item) => sum + item.quantity, 0)}
                </span>
              )}
            </button>
            <button
              className="p-2 rounded-lg hover:bg-secondary transition-colors"
              aria-label="Share"
            >
              <Share2 className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      {/* Product Image */}
      <div className="relative w-full aspect-square bg-secondary overflow-hidden">
        <AnimatePresence mode="wait" initial={false}>
          <motion.img
            key={`${currentImageIndex}-${currentImage}`}
            src={currentImage}
            alt={product.name}
            initial={{ opacity: 0.3 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2, ease: 'easeInOut' }}
            className="w-full h-full object-cover"
            style={{ minHeight: '100%', minWidth: '100%' }}
            onLoad={() => {
              setImageLoading(false);
              setLoadedImages(prev => new Set(prev).add(currentImage));
            }}
            onError={(e) => {
              const target = e.currentTarget;
              setImageLoading(false);
              // Prevent infinite loop by checking if already on placeholder
              if (target.src && !target.src.includes('placeholder.png')) {
                target.src = '/placeholder.png';
              }
            }}
            loading="eager"
          />
        </AnimatePresence>
        
        {/* Loading overlay - only show if actually loading */}
        {imageLoading && !loadedImages.has(currentImage) && currentImage !== '/placeholder.png' && (
          <div className="absolute inset-0 flex items-center justify-center bg-secondary/80 backdrop-blur-sm z-10">
            <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
          </div>
        )}

        {/* Image Indicators */}
        {images.length > 1 && (
          <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                  {images.map((img, index) => (
                    <button
                      key={`${img}-${index}`}
                      onClick={() => setCurrentImageIndex(index)}
                      className={cn(
                        "h-1.5 rounded-full transition-all",
                        index === currentImageIndex
                          ? "w-6 bg-foreground"
                          : "w-1.5 bg-foreground/30"
                      )}
                      aria-label={`Go to image ${index + 1}`}
                    />
                  ))}
          </div>
        )}

        {/* Stock Badge */}
        {!product.is_active && (
          <div className="absolute top-4 right-4 bg-destructive text-destructive-foreground px-3 py-1.5 rounded-full text-xs font-medium">
            Out of Stock
          </div>
        )}
      </div>

      {/* Product Info */}
      <div className="px-4 pt-6 space-y-6">
        {/* Title & Price Section */}
        <div className="space-y-2">
          {attributes.brand && (
            <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
              {attributes.brand}
            </p>
          )}
          <h1 className="text-2xl font-bold text-foreground leading-tight">
            {product.name}
          </h1>
          <div className="flex items-baseline gap-3">
            {product.price ? (
              <>
                <span className="text-3xl font-bold text-foreground">
                  {product.price.toLocaleString()} ETB
                </span>
                {product.unit && (
                  <span className="text-sm text-muted-foreground">/ {product.unit}</span>
                )}
              </>
            ) : (
              <span className="text-lg text-muted-foreground">Price on request</span>
            )}
          </div>
        </div>

        {/* Availability */}
        {product.available_quantity > 0 && (
          <div className="flex items-center gap-1.5 text-sm text-green-600">
            <Check className="w-4 h-4" />
            <span>In Stock ({product.available_quantity} available)</span>
          </div>
        )}

        {/* Color Selection */}
        {colors.length > 0 && (
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <label className="text-sm font-semibold text-foreground">
                Color {selectedColor && <span className="text-muted-foreground font-normal">({selectedColor})</span>}
              </label>
            </div>
            <div className="flex gap-2 flex-wrap">
              {colors.map((color: string) => (
                <button
                  key={color}
                  onClick={() => setSelectedColor(color)}
                  className={cn(
                    "px-4 py-2.5 text-sm font-medium rounded-lg border-2 transition-all",
                    selectedColor === color
                      ? "border-primary bg-primary text-primary-foreground shadow-md"
                      : "border-border bg-card hover:border-primary/50"
                  )}
                >
                  {color}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Size Selection */}
        {sizes.length > 0 && (
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <label className="text-sm font-semibold text-foreground">
                Size {selectedSize && <span className="text-muted-foreground font-normal">({selectedSize})</span>}
              </label>
            </div>
            <div className="flex gap-2 flex-wrap">
              {sizes.map((size: string) => (
                <button
                  key={size}
                  onClick={() => setSelectedSize(size)}
                  className={cn(
                    "min-w-[56px] h-12 px-4 text-sm font-medium rounded-lg border-2 transition-all",
                    selectedSize === size
                      ? "border-primary bg-primary text-primary-foreground shadow-md"
                      : "border-border bg-card hover:border-primary/50"
                  )}
                >
                  {size}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Quantity Selector */}
        <div className="space-y-3">
          <div className="text-sm font-semibold text-foreground">Quantity</div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 border border-border rounded-lg">
              <button
                onClick={() => setQuantity(Math.max(1, quantity - 1))}
                disabled={quantity <= 1}
                className="p-2 hover:bg-secondary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                aria-label="Decrease quantity"
              >
                <Minus className="w-4 h-4" />
              </button>
              <span className="w-12 text-center font-semibold">{quantity}</span>
              <button
                onClick={() => setQuantity(Math.min(maxQuantity, quantity + 1))}
                disabled={quantity >= maxQuantity}
                className="p-2 hover:bg-secondary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                aria-label="Increase quantity"
              >
                <Plus className="w-4 h-4" />
              </button>
            </div>
            <span className="text-sm text-muted-foreground">
              {maxQuantity} available
            </span>
          </div>
        </div>

        {/* Tabs for Details & Specifications */}
        <Card className="border-border">
          <div className="flex border-b border-border">
            <button
              onClick={() => setActiveTab('details')}
              className={cn(
                "flex-1 px-4 py-3 text-sm font-medium transition-colors border-b-2",
                activeTab === 'details'
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              )}
            >
              <div className="flex items-center justify-center gap-2">
                <Info className="w-4 h-4" />
                Details
              </div>
            </button>
            <button
              onClick={() => setActiveTab('specs')}
              className={cn(
                "flex-1 px-4 py-3 text-sm font-medium transition-colors border-b-2",
                activeTab === 'specs'
                  ? "border-primary text-primary"
                  : "border-transparent text-muted-foreground hover:text-foreground"
              )}
            >
              <div className="flex items-center justify-center gap-2">
                <Package className="w-4 h-4" />
                Specifications
              </div>
            </button>
          </div>
          <CardContent className="p-4">
            <AnimatePresence mode="wait">
              {activeTab === 'details' ? (
                <motion.div
                  key="details"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.2 }}
                  className="space-y-3"
                >
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    {product.description || 'No description available for this product.'}
                  </p>
                  {attributes.details && (
                    <ul className="space-y-2">
                      {(Array.isArray(attributes.details) ? attributes.details : [attributes.details]).map(
                        (detail: string) => (
                          <li key={detail} className="flex items-start gap-2 text-sm text-muted-foreground">
                            <Check className="w-4 h-4 mt-0.5 text-primary shrink-0" />
                            <span>{detail}</span>
                          </li>
                        )
                      )}
                    </ul>
                  )}
                </motion.div>
              ) : (
                <motion.div
                  key="specs"
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.2 }}
                  className="space-y-3"
                >
                  <div className="grid gap-3">
                    {product.unit && (
                      <div className="flex justify-between py-2 border-b border-border">
                        <span className="text-sm font-medium text-muted-foreground">Unit</span>
                        <span className="text-sm text-foreground">{product.unit}</span>
                      </div>
                    )}
                    {attributes.material && (
                      <div className="flex justify-between py-2 border-b border-border">
                        <span className="text-sm font-medium text-muted-foreground">Material</span>
                        <span className="text-sm text-foreground">{attributes.material}</span>
                      </div>
                    )}
                    {attributes.weight && (
                      <div className="flex justify-between py-2 border-b border-border">
                        <span className="text-sm font-medium text-muted-foreground">Weight</span>
                        <span className="text-sm text-foreground">{attributes.weight}</span>
                      </div>
                    )}
                    <div className="flex justify-between py-2 border-b border-border">
                      <span className="text-sm font-medium text-muted-foreground">Available Quantity</span>
                      <span className="text-sm text-foreground">{product.available_quantity}</span>
                    </div>
                    {product.external_id && (
                      <div className="flex justify-between py-2">
                        <span className="text-sm font-medium text-muted-foreground">SKU</span>
                        <span className="text-sm text-foreground font-mono">{product.external_id}</span>
                      </div>
                    )}
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </CardContent>
        </Card>
      </div>

      {/* Fixed Bottom Action Bar */}
      <div className="fixed bottom-0 left-0 right-0 z-[60] bg-background border-t border-border safe-bottom">
        <div className="px-4 py-3 flex items-center gap-3">
          <Button
            onClick={handleAddToCart}
            disabled={!canAddToCart}
            className="flex-1 h-12 text-base font-semibold shadow-lg"
            size="lg"
          >
            <ShoppingCart className="w-5 h-5 mr-2" />
            {(() => {
              if (canAddToCart) {
                const quantityText = quantity > 1 ? ` (${quantity})` : '';
                return `Add to Cart${quantityText}`;
              }
              if (colors.length > 1 || sizes.length > 1) {
                return 'Select Options';
              }
              return 'Unavailable';
            })()}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ProductDetail;
