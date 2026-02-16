"use client";
import { useEffect, useState, useMemo, useTransition } from "react";
import { Filter, ShoppingBagIcon, X } from "lucide-react";
import Link from "next/link";
import { IoWarning } from "react-icons/io5";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

import { GetAllCategories } from "@/app/actions/category";
import { ListProducts, ProductResponse } from "@/app/actions/product";
import { Button } from "@/components/ui/button";
import { CardSkeleton } from "@/components/CardSkeleton";
import { ProductCard } from "@/components/ProductCard";
import { Catalogue, Image } from "@/lib/types";
import useSearchStore from "@/lib/store/useSearchStore";
import { Input } from "@/components/ui/input";

const ITEMS_PER_PAGE = 12;

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

const Product = () => {
  const [currentPage, setCurrentPage] = useState(0);
  const [error, setError] = useState<string>();
  const [categories, setCategories] = useState([{ id: 0, name: "All" }]);
  const [loading, setLoading] = useState(true);
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [selectedCategory, setSelectedCategory] = useState(0);
  const [isPending, startTransition] = useTransition();
  const [totalProducts, setTotalProducts] = useState(0);
  const searchQuery = useSearchStore((state) => state.searchQuery);
  const clearSearch = useSearchStore((state) => state.clearSearch);
  const [minPrice, setMinPrice] = useState<number | ''>('');
  const [maxPrice, setMaxPrice] = useState<number | ''>('');

  const loadProducts = async () => {
    setLoading(true);
    try {
      console.log('🔄 Starting to load products...');
      
      // Load categories and products in parallel
      const [fetchedCategories, productResponse] = await Promise.all([
        GetAllCategories().catch(() => []),
        ListProducts({ is_active: true, limit: 1000, offset: 0 }).catch(() => ({ 
          products: [], 
          total: 0, 
          limit: 1000, 
          offset: 0 
        }))
      ]);

      // If there are more products than the limit, fetch all pages
      let allProducts = productResponse.products || [];
      const total = productResponse.total || 0;
      const limit = productResponse.limit || 1000;
      
      if (total > allProducts.length) {
        console.log(`📦 Fetching all products: ${total} total, ${allProducts.length} loaded so far`);
        const remainingCount = total - allProducts.length;
        const remainingPages = Math.ceil(remainingCount / limit);
        
        // Fetch remaining pages
        const remainingPromises = [];
        for (let page = 1; page <= remainingPages; page++) {
          const offset = page * limit;
          remainingPromises.push(
            ListProducts({ is_active: true, limit, offset }).catch((error) => {
              console.error(`❌ Failed to fetch products page ${page} (offset ${offset}):`, error);
              return { products: [], total: 0, limit, offset };
            })
          );
        }
        
        const remainingResponses = await Promise.all(remainingPromises);
        remainingResponses.forEach((response) => {
          if (response.products && response.products.length > 0) {
            allProducts = [...allProducts, ...response.products];
          }
        });
        
        console.log(`✅ Loaded ${allProducts.length} out of ${total} total products`);
      }
      
      // Convert ProductResponse[] to Catalogue[] format for compatibility (same as mobile project)
      const catalogueProducts: Catalogue[] = allProducts.map((p: ProductResponse) => {
        // Parse attributes if it's a string
        let parsedAttributes = p.attributes;
        if (typeof p.attributes === 'string') {
          try {
            parsedAttributes = JSON.parse(p.attributes);
          } catch {
            parsedAttributes = {};
          }
        }
        
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
      }).filter((p: Catalogue) => p.id != null && p.id !== undefined);
      
      startTransition(() => {
        setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
        setProducts(catalogueProducts);
        setTotalProducts(total);
        setLoading(false);
      });
    } catch (error: any) {
      console.error("Failed to load products:", error);
      setError(error.message || 'Failed to load products');
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProducts();
  }, []);

  // Calculate price range from all products
  const priceRange = useMemo(() => {
    if (products.length === 0) return { min: 0, max: 0 };
    
    const prices: number[] = [];
    products.forEach((product) => {
      if (product.price) {
        prices.push(product.price);
      }
      if (product.configurables && product.configurables.length > 0) {
        product.configurables.forEach((config) => {
          if (config.price) {
            prices.push(config.price);
          }
        });
      }
    });
    
    if (prices.length === 0) return { min: 0, max: 0 };
    
    return {
      min: Math.min(...prices),
      max: Math.max(...prices),
    };
  }, [products]);

  const filteredProducts = useMemo(() => {
    let filtered = products;
    
    // Filter by category
    if (selectedCategory !== 0) {
      filtered = filtered.filter((product) => {
        const categoryIds = product.configurable_attributes?.category_ids || [];
        return categoryIds.includes(selectedCategory);
      });
    }
    
    // Filter by search query
    if (searchQuery && searchQuery.trim() !== '') {
      const query = searchQuery.toLowerCase().trim();
      filtered = filtered.filter((product) => {
        const nameMatch = product.name?.toLowerCase().includes(query);
        const descMatch = product.desc?.toLowerCase().includes(query);
        return nameMatch || descMatch;
      });
    }
    
    // Filter by price range
    if (minPrice !== '' || maxPrice !== '') {
      filtered = filtered.filter((product) => {
        // Get product price (use min price from configurables if available)
        let productPrice = product.price || 0;
        
        if (product.configurables && product.configurables.length > 0) {
          const prices = product.configurables.map((c) => c.price || 0).filter((p) => p > 0);
          if (prices.length > 0) {
            productPrice = Math.min(...prices); // Use minimum price for filtering
          }
        }
        
        const minCheck = minPrice === '' || productPrice >= minPrice;
        const maxCheck = maxPrice === '' || productPrice <= maxPrice;
        
        return minCheck && maxCheck;
      });
    }
    
    return filtered;
  }, [products, selectedCategory, searchQuery, minPrice, maxPrice]);

  const paginatedProducts = useMemo(() => {
    const startIndex = currentPage * ITEMS_PER_PAGE;
    return filteredProducts.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredProducts, currentPage]);

  const totalPages = Math.ceil(filteredProducts.length / ITEMS_PER_PAGE);

  const handlePageChange = (newPage: number) => {
    if (newPage >= 0 && newPage < totalPages) {
      startTransition(() => {
        setCurrentPage(newPage);
      });
    }
  };

  const handleCategoryChange = (categoryId: number) => {
    startTransition(() => {
      setSelectedCategory(categoryId);
      setCurrentPage(0); // Reset to first page when category changes
    });
  };

  // Reset to first page when search query or price filter changes
  useEffect(() => {
    setCurrentPage(0);
  }, [searchQuery, minPrice, maxPrice]);

  if (error) {
    return (
      <div className="w-full flex flex-col gap-6 items-center justify-center text-center min-h-screen">
        <div className="bg-red-50 rounded-full p-8 mb-4">
          <IoWarning className="w-16 h-16 text-red-900" />
        </div>
        <div>
          <h4 className="text-2xl font-semibold text-gray-900 mb-2">
            Something went wrong while fetching products
          </h4>
        </div>
        <Button size="lg" className="px-8" onClick={loadProducts}>
          Retry
        </Button>
      </div>
    );
  }

  const handlePriceChange = (type: 'min' | 'max', value: string) => {
    const numValue = value === '' ? '' : parseFloat(value);
    if (type === 'min') {
      setMinPrice(numValue);
    } else {
      setMaxPrice(numValue);
    }
  };

  const clearPriceFilter = () => {
    setMinPrice('');
    setMaxPrice('');
  };

  return (
    <div className="w-full max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      {/* Header Section */}
      <div className="flex flex-col items-center mb-12">
        <h2 className="text-3xl font-medium mb-4">Products</h2>
        <div className="h-1 w-20 bg-primary mb-8"></div>
        {/* Search Query Indicator */}
        {searchQuery && (
          <div className="mb-4 flex items-center gap-2 px-4 py-2 bg-primary/10 rounded-full">
            <span className="text-sm text-gray-700">
              Search results for: <strong>"{searchQuery}"</strong>
            </span>
            <button
              onClick={() => {
                clearSearch();
                setCurrentPage(0);
              }}
              className="ml-2 p-1 hover:bg-primary/20 rounded-full transition-colors"
              aria-label="Clear search"
            >
              <X className="h-4 w-4 text-gray-600" />
            </button>
          </div>
        )}
      </div>

      {/* Filters and Content Section */}
      <div className="flex flex-col lg:flex-row gap-8 items-start">
        {/* Price Filter Sidebar - Left */}
        <aside className="w-full lg:w-64 flex-shrink-0">
          <div className="bg-white rounded-lg border border-gray-200 p-6 sticky top-24">
            <h3 className="text-lg font-semibold mb-4 text-gray-900">Price Filter</h3>
            
            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium text-gray-700 mb-2 block">
                  Min Price (ETB)
                </label>
                <Input
                  type="number"
                  placeholder={`${priceRange.min.toLocaleString()}`}
                  value={minPrice}
                  onChange={(e) => handlePriceChange('min', e.target.value)}
                  min={priceRange.min}
                  max={priceRange.max}
                  className="w-full"
                />
              </div>
              
              <div>
                <label className="text-sm font-medium text-gray-700 mb-2 block">
                  Max Price (ETB)
                </label>
                <Input
                  type="number"
                  placeholder={`${priceRange.max.toLocaleString()}`}
                  value={maxPrice}
                  onChange={(e) => handlePriceChange('max', e.target.value)}
                  min={priceRange.min}
                  max={priceRange.max}
                  className="w-full"
                />
              </div>
              
              {(minPrice !== '' || maxPrice !== '') && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={clearPriceFilter}
                  className="w-full"
                >
                  <X className="h-4 w-4 mr-2" />
                  Clear Price Filter
                </Button>
              )}
              
              <div className="pt-2 text-xs text-gray-500">
                Range: {priceRange.min.toLocaleString()} - {priceRange.max.toLocaleString()} ETB
              </div>
            </div>
          </div>
        </aside>

        {/* Main Content */}
        <div className="flex-1">
          {/* Category Filters */}
          <div className="flex items-center justify-between mb-6 w-full">
            <div className="flex gap-4 items-center">
              <Filter className="w-5 h-5 shrink-0" />
              <div className="flex gap-4 overflow-x-auto no-scrollbar w-full">
                {categories.map((category) => (
                  <button
                    key={category.id}
                    onClick={() => handleCategoryChange(category.id)}
                    disabled={isPending}
                    className={`shrink-0 px-4 py-2 rounded-lg border text-sm transition-colors whitespace-nowrap disabled:opacity-50 disabled:cursor-not-allowed ${
                      selectedCategory === category.id
                        ? "bg-primary text-white"
                        : "bg-secondary text-primary hover:bg-opacity-80"
                    }`}
                  >
                    {category.name}
                  </button>
                ))}
              </div>
            </div>
          </div>

        {loading ? (
          <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-6">
              {[...Array(12)].map((_, i) => (
                <CardSkeleton key={i} />
              ))}
            </div>
          </div>
        ) : (
          <>
            {filteredProducts.length === 0 ? (
              <div className="w-full flex flex-col gap-6 items-center justify-center text-center py-16">
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
              <div className="space-y-6">
                <section id="products" className="container mx-auto px-4">
                  <div className={`grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-6 transition-opacity duration-200 ${isPending ? 'opacity-70' : 'opacity-100'}`}>
                    {paginatedProducts.map((product) => (
                      <ProductCard product={product} key={product.id || product.name} />
                    ))}
                  </div>
                </section>

                {totalPages > 1 && (
                  <Pagination>
                    <PaginationContent>
                      <PaginationItem>
                        <PaginationPrevious
                          onClick={() => handlePageChange(currentPage - 1)}
                          className={
                            currentPage === 0
                              ? "pointer-events-none text-primary opacity-50"
                              : "text-primary"
                          }
                        />
                      </PaginationItem>

                      {Array.from({ length: Math.min(5, totalPages) }).map(
                        (_, index) => {
                          let pageNum;
                          if (totalPages <= 5) {
                            pageNum = index;
                          } else if (currentPage <= 2) {
                            pageNum = index;
                          } else if (currentPage >= totalPages - 3) {
                            pageNum = totalPages - 5 + index;
                          } else {
                            pageNum = currentPage - 2 + index;
                          }

                          return (
                            <PaginationItem key={pageNum} className="text-primary">
                              <PaginationLink
                                isActive={currentPage === pageNum}
                                onClick={() => handlePageChange(pageNum)}
                              >
                                {pageNum + 1}
                              </PaginationLink>
                            </PaginationItem>
                          );
                        }
                      )}

                      <PaginationItem>
                        <PaginationNext
                          onClick={() => handlePageChange(currentPage + 1)}
                          className={
                            currentPage === totalPages - 1
                              ? "pointer-events-none text-primary opacity-50"
                              : "text-primary"
                          }
                        />
                      </PaginationItem>
                    </PaginationContent>
                  </Pagination>
                )}
              </div>
            )}
          </>
        )}
        </div>
      </div>
    </div>
  );
};

export default Product;
