"use client";
import { useState, useEffect, useMemo } from "react";
import { Filter, ShoppingBagIcon, ArrowRight } from "lucide-react";
import Link from "next/link";
import { GetAllCategories } from "@/app/actions/category";
import { ListProducts, ProductResponse } from "@/app/actions/product";
import { Catalogue, Image } from "@/lib/types";
import { CardSkeleton } from "@/components/CardSkeleton";
import { ProductCard } from "@/components/ProductCard";
import { Button } from "@/components/ui/button";
import useSearchStore from "@/lib/store/useSearchStore";

// Helper function to extract images from attributes
const extractImagesFromAttributes = (attributes: any): Image[] => {
  if (!attributes) return [];
  
  // Parse attributes if it's a string
  let attrs = attributes;
  if (typeof attributes === 'string') {
    try {
      attrs = JSON.parse(attributes);
    } catch {
      return [];
    }
  }
  
  if (attrs.image_url && typeof attrs.image_url === 'string') {
    return [{ ImageUrl: attrs.image_url, BlurHash: '' }];
  }
  if (!attrs.images) return [];

  // Handle different image formats
  if (Array.isArray(attrs.images)) {
    return attrs.images
      .map((img: any) => {
        if (typeof img === 'string') {
          return { ImageUrl: img, BlurHash: '' };
        }
        if (img?.ImageUrl || img?.url || img?.imageUrl) {
          return {
            ImageUrl: img.ImageUrl || img.url || img.imageUrl,
            BlurHash: img.BlurHash || img.blurHash || ''
          };
        }
        return null;
      })
      .filter((img: Image | null): img is Image => img !== null);
  }
  
  // Single image as string
  if (typeof attrs.images === 'string') {
    return [{ ImageUrl: attrs.images, BlurHash: '' }];
  }
  
  // Single image as object
  if (attrs.images?.ImageUrl || attrs.images?.url || attrs.images?.imageUrl) {
    return [{
      ImageUrl: attrs.images.ImageUrl || attrs.images.url || attrs.images.imageUrl,
      BlurHash: attrs.images.BlurHash || attrs.images.blurHash || ''
    }];
  }
  
  return [];
};

type VariantItem = { option: string; value: string; price_modifier?: number };

function buildConfigurablesFromProduct(
  p: ProductResponse,
  parsedAttributes: Record<string, unknown>,
  images: Image[]
): { id: number; name: string; desc: string; price: number; stock: number; external_id: string; attributes: Record<string, string>; images: Image[]; supplier_id: number; categories: number[]; is_active: boolean }[] {
  const basePrice = p.price ?? 0;
  const totalQty = p.available_quantity ?? 0;
  const variants = parsedAttributes?.variants as VariantItem[] | undefined;
  if (Array.isArray(variants) && variants.length > 0 && totalQty > 0) {
    const optionName = variants[0]?.option ?? 'Option';
    const stockPerVariant = Math.max(0, Math.floor(totalQty / variants.length));
    return variants.map((v, idx) => ({
      id: p.id,
      name: p.name,
      desc: p.description || '',
      price: basePrice + (v.price_modifier ?? 0),
      stock: idx < variants.length - 1 ? stockPerVariant : Math.max(0, totalQty - stockPerVariant * (variants.length - 1)),
      external_id: `${p.external_id || ''}-${v.value}`.replace(/\s+/g, '_'),
      attributes: { ...(parsedAttributes as Record<string, string>), [optionName]: v.value },
      images,
      supplier_id: p.supplier_id,
      categories: p.category_ids || [],
      is_active: p.is_active
    }));
  }
  return totalQty > 0 ? [{
    id: p.id,
    name: p.name,
    desc: p.description || '',
    price: basePrice,
    stock: totalQty,
    external_id: p.external_id || '',
    attributes: parsedAttributes as Record<string, string> || {},
    images,
    supplier_id: p.supplier_id,
    categories: p.category_ids || [],
    is_active: p.is_active
  }] : [];
}

export default function ProductGrid() {
  const [categories, setCategories] = useState([{ id: 0, name: "All" }]);
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [selectedCategory, setSelectedCategory] = useState(0);
  const searchQuery = useSearchStore((state) => state.searchQuery);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        // Load both in parallel for better performance
        const [fetchedCategories, productResponse] = await Promise.all([
          GetAllCategories().catch(() => []),
          ListProducts({ is_active: true, limit: 100 }).catch(() => ({ products: [], total: 0, limit: 100, offset: 0 }))
        ]);
        
        setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
        
        // Convert ProductResponse[] to Catalogue[] format for compatibility
        // Store category_ids in configurable_attributes for filtering
        const catalogueProducts: Catalogue[] = (productResponse.products || []).map((p: ProductResponse) => {
          // Parse attributes if it's a string
          let parsedAttributes = p.attributes;
          if (typeof p.attributes === 'string') {
            try {
              parsedAttributes = JSON.parse(p.attributes);
            } catch {
              parsedAttributes = {};
            }
          }
          
          // Extract images from attributes (including image_url)
          const images = extractImagesFromAttributes(parsedAttributes);

          const configurables = buildConfigurablesFromProduct(p, parsedAttributes as Record<string, unknown>, images);

          return {
            id: p.id,
            name: p.name,
            desc: p.description || '',
            price: p.price,
            is_active: p.is_active,
            images: images,
            configurable_attributes: {
              ...(parsedAttributes || {}),
              category_ids: p.category_ids || [], // Store category_ids for filtering
              available_quantity: p.available_quantity, // Store available quantity for stock display
              total_quantity: p.total_quantity,
              reserved_quantity: p.reserved_quantity
            },
            configurables: configurables
          };
        });
        
        setProducts(catalogueProducts);
      } catch (error) {
        console.error("Failed to load data:", error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  const displayedProducts = useMemo(() => {
    if (!products) return [];

    let filtered = products;
    
    // Filter by search query
    if (searchQuery && searchQuery.trim() !== '') {
      const query = searchQuery.toLowerCase().trim();
      filtered = filtered.filter((product) => {
        const nameMatch = product.name?.toLowerCase().includes(query);
        const descMatch = product.desc?.toLowerCase().includes(query);
        // Also check configurables for search matches
        const configurableMatch = product.configurables?.some((config) =>
          config.name?.toLowerCase().includes(query) ||
          config.desc?.toLowerCase().includes(query)
        );
        return nameMatch || descMatch || configurableMatch;
      });
    }
    
    // Filter by category
    if (selectedCategory !== 0) {
      // Filter by category_ids stored in configurable_attributes
      filtered = filtered.filter((product) => {
        const categoryIds = product.configurable_attributes?.category_ids || [];
        return categoryIds.includes(selectedCategory);
      });
    }

    // When searching, show all matching results. Otherwise, limit to 4 for display
    return searchQuery.trim() ? filtered : filtered.slice(0, 4);
  }, [products, selectedCategory, searchQuery]);

  return (
    <section id="products" className="container mx-auto px-4 sm:px-6 lg:px-8 py-20 bg-white">
      <div className="flex flex-col items-center mb-16">
        <p className="text-sm font-semibold text-primary mb-3 uppercase tracking-wide">
          Our Products
        </p>
        <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4 text-center">
          Featured Products
        </h2>
        <div className="h-1 w-24 bg-primary mb-12 rounded-full"></div>
        <div className="flex items-center justify-between mb-8 w-full max-w-6xl">
          <div className="flex gap-3 items-center w-full">
            <Filter className="w-5 h-5 shrink-0 text-gray-600" />
            <div className="flex gap-3 overflow-x-auto no-scrollbar w-full pb-2">
              {categories.map((category) => (
                <button
                  key={category.id}
                  onClick={() => setSelectedCategory(category.id)}
                  className={`shrink-0 px-6 py-2.5 rounded-full text-sm font-medium transition-all duration-300 whitespace-nowrap ${
                    selectedCategory === category.id
                      ? "bg-primary text-white shadow-lg"
                      : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                  }`}
                >
                  {category.name}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 lg:gap-8 max-w-7xl mx-auto">
        {loading ? (
          // Show 4 skeleton cards in grid
          [...Array(4)].map((_, i) => <CardSkeleton key={i} />)
        ) : displayedProducts.length === 0 ? (
          // Empty state
          <div className="col-span-full flex flex-col gap-6 items-center justify-center text-center py-16">
            <div className="bg-gray-50 rounded-full p-8 mb-4">
              <ShoppingBagIcon className="w-16 h-16 text-gray-400" />
            </div>
            <div>
              <h4 className="text-2xl font-semibold text-gray-900 mb-2">
                No products found
              </h4>
              <p className="text-gray-600 mb-6 max-w-md">
                {selectedCategory === 0
                  ? "There are currently no products available"
                  : "No products found in this category"}
              </p>
            </div>
          </div>
        ) : (
          displayedProducts.map((product) => (
            <ProductCard product={product} key={product.id || product.name} />
          ))
        )}
      </div>

      <div className="flex justify-center mt-12">
        <Link href="/product" prefetch={true}>
          <Button 
            size="lg"
            className="bg-primary text-white hover:bg-primary/90 border-2 border-primary transition-all duration-300 px-8 py-6 text-base font-semibold h-auto shadow-none"
          >
            View All Products
            <ArrowRight className="ml-2 h-5 w-5" />
          </Button>
        </Link>
      </div>
    </section>
  );
}