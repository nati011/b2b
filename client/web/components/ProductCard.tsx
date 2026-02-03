"use client";
import React, { useMemo } from "react";
import useCatalogueStore from "@/lib/store/useCatalogueStore";
import { Catalogue } from "@/lib/types";
import Image from "next/image";
import Link from "next/link";
import { Check, AlertCircle, ImageIcon } from "lucide-react";
import { Badge } from "@/components/ui/badge";

interface Props {
  product: Catalogue;
}

export const ProductCard: React.FC<Props> = React.memo(({ product }) => {
  const setProduct = useCatalogueStore((state) => state.setProduct);

  const priceRange = useMemo(() => {
    if (!product.configurables || product.configurables.length === 0) {
      return product.price ? `${product.price.toLocaleString()} ETB` : 'Price not available';
    }

    const prices = product.configurables.map((config) => config.price);
    const minPrice = Math.min(...prices);
    const maxPrice = Math.max(...prices);

    return minPrice === maxPrice
      ? `${minPrice.toLocaleString()} ETB`
      : `${minPrice.toLocaleString()} ETB - ${maxPrice.toLocaleString()} ETB`;
  }, [product.configurables, product.price]);

  const stockInfo = useMemo(() => {
    // If no configurables, check available_quantity from configurable_attributes
    if (!product.configurables || product.configurables.length === 0) {
      const availableQuantity = product.configurable_attributes?.available_quantity;
      const stock = availableQuantity !== undefined && availableQuantity !== null
        ? availableQuantity
        : (product.is_active ? 999 : 0);
      
      return {
        stock: stock,
        isAvailable: stock > 0 && product.is_active,
        isLowStock: stock > 0 && stock < 10
      };
    }
    
    // Calculate total stock from all configurables
    const totalStock = product.configurables.reduce((sum, config) => sum + (config.stock || 0), 0);
    const hasAnyStock = product.configurables.some((config) => (config.stock || 0) > 0);
    
    return {
      stock: totalStock,
      isAvailable: hasAnyStock,
      isLowStock: totalStock > 0 && totalStock < 10
    };
  }, [product.configurables, product.configurable_attributes, product.is_active]);

  const isOutOfStock = useMemo(() => {
    return !stockInfo.isAvailable;
  }, [stockInfo.isAvailable]);

  // Get the first valid image URL (non-empty string)
  const imageUrl = useMemo(() => {
    const url = product?.images?.[0]?.ImageUrl;
    return url && typeof url === 'string' && url.trim() !== '' ? url : null;
  }, [product?.images]);

  return (
    <Link
      href={`/product/${product.name}`}
      className="group animate-fade-in"
      prefetch={true}
      onClick={() => setProduct(product)}
    >
      <div className="bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100">
        <div className="aspect-square overflow-hidden bg-gray-50 relative">
          {imageUrl ? (
            <Image
              width={400}
              height={400}
              src={imageUrl}
              alt={product.name}
              className="w-full h-full object-cover transform transition-transform duration-300 group-hover:scale-110"
              loading="lazy"
              sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 25vw"
              quality={85}
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center bg-gray-100">
              <ImageIcon className="w-16 h-16 text-gray-400" />
            </div>
          )}
          {isOutOfStock && (
            <div className="absolute top-4 right-4 bg-red-500 text-white px-3 py-1 rounded-full text-xs font-semibold">
              Out of Stock
            </div>
          )}
        </div>
        <div className="p-5">
          <h3 className="text-lg font-semibold mb-2 text-gray-900 group-hover:text-primary transition-colors line-clamp-2">
            {product.name}
          </h3>
          <p className="text-lg font-bold text-primary mb-2">
            {priceRange}
          </p>
          {/* Stock Availability */}
          <div className="flex items-center gap-2">
            {stockInfo.isAvailable ? (
              <>
                <Check className="w-4 h-4 text-green-600" />
                <span className="text-sm text-gray-600">
                  {stockInfo.stock > 0 
                    ? `${stockInfo.stock} ${stockInfo.stock === 1 ? 'unit' : 'units'} available`
                    : 'In Stock'}
                </span>
                {stockInfo.isLowStock && (
                  <Badge variant="outline" className="ml-auto text-xs bg-yellow-50 border-yellow-300 text-yellow-800">
                    Low Stock
                  </Badge>
                )}
              </>
            ) : (
              <>
                <AlertCircle className="w-4 h-4 text-red-600" />
                <span className="text-sm text-red-600">Out of Stock</span>
              </>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
});

ProductCard.displayName = 'ProductCard';
