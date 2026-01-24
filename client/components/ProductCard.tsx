"use client";
import React, { useMemo } from "react";
import useCatalogueStore from "@/lib/store/useCatalogueStore";
import { Catalogue } from "@/lib/types";
import Image from "next/image";
import Link from "next/link";

interface Props {
  product: Catalogue;
}

export const ProductCard: React.FC<Props> = React.memo(({ product }) => {
  const setProduct = useCatalogueStore((state) => state.setProduct);

  const priceRange = useMemo(() => {
    if (!product.configurables || product.configurables.length === 0) {
      return `$${product.price}`;
    }

    const prices = product.configurables.map((config) => config.price);
    const minPrice = Math.min(...prices);
    const maxPrice = Math.max(...prices);

    return minPrice === maxPrice
      ? `${minPrice.toLocaleString()} ETB`
      : `${minPrice.toLocaleString()} ETB - ${maxPrice.toLocaleString()} ETB`;
  }, [product.configurables, product.price]);

  const isOutOfStock = useMemo(() => {
    return product.configurables.every(
      (configurable) => configurable.stock === 0
    );
  }, [product.configurables]);

  return (
    <Link
      href={`/product/${product.name}`}
      className="group animate-fade-in"
      onClick={() => setProduct(product)}
    >
      <div className="bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100">
        <div className="aspect-square overflow-hidden bg-gray-50 relative">
          <Image
            width={400}
            height={400}
            src={product?.images?.[0]?.ImageUrl ?? ""}
            alt={product.name}
            className="w-full h-full object-cover transform transition-transform duration-300 group-hover:scale-110"
            loading="lazy"
            sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 25vw"
            quality={85}
          />
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
          <p className="text-lg font-bold text-primary">
            {priceRange}
          </p>
        </div>
      </div>
    </Link>
  );
});

ProductCard.displayName = 'ProductCard';
