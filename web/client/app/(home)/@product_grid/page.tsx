"use client";
import { useState, useEffect, useMemo } from "react";
import { Filter, ShoppingBagIcon, ArrowRight } from "lucide-react";
import Link from "next/link";
import { GetAllCategories } from "@/app/actions/category";
import { GetAllCatalogues } from "@/app/actions/catalogue";
import { Catalogue } from "@/lib/types";
import { CardSkeleton } from "@/components/CardSkeleton";
import { ProductCard } from "@/components/ProductCard";
import { Button } from "@/components/ui/button";

export default function ProductGrid() {
  const [categories, setCategories] = useState([{ id: 0, name: "All" }]);
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [selectedCategory, setSelectedCategory] = useState(0);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        // Load both in parallel for better performance
        const [fetchedCategories, productList] = await Promise.all([
          GetAllCategories().catch(() => []),
          GetAllCatalogues().catch(() => [])
        ]);
        
        setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
        setProducts(productList);
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
    if (selectedCategory !== 0) {
      filtered = products.filter((product) =>
        product.configurables.some((configurable) =>
          configurable.categories?.includes(selectedCategory)
        )
      );
    }

    return filtered.slice(0, 4);
  }, [products, selectedCategory]);

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
            <ProductCard product={product} key={product.name} />
          ))
        )}
      </div>

      <div className="flex justify-center mt-12">
        <Link href="/product" prefetch={true}>
          <Button 
            variant="outline" 
            size="lg"
            className="border-2 hover:bg-primary hover:text-white hover:border-primary transition-all duration-300"
          >
            View All Products
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </Link>
      </div>
    </section>
  );
}