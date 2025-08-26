"use client";
import { useState, useEffect } from "react";
import { Filter } from "lucide-react";
import { MdOutlineKeyboardDoubleArrowRight } from "react-icons/md";
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
  const [products, setProducts] = useState<Catalogue[]>();
  const [selectedCategory, setSelectedCategory] = useState(0);

  useEffect(() => {
    const loadCategories = async () => {
      try {
        const fetchedCategories = await GetAllCategories();
        setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
      } catch (error) {
        console.error("Error fetching categories:", error);
      }
    };

    loadCategories();
  }, []);

  useEffect(() => {
    const loadProducts = async () => {
      setLoading(true);
      try {
        const productList = await GetAllCatalogues();
        setProducts(productList);
      } catch (error) {
        console.error("Failed to load products:", error);
      } finally {
        setLoading(false);
      }
    };

    loadProducts();
  }, []);

  const displayedProducts = products ? products.slice(0, 4) : [];
  return (
    <section id="products" className="container mx-auto px-4 py-16">
      <div className="flex flex-col items-center mb-12">
        <h2 className="text-3xl font-medium mb-4">Featured Products</h2>
        <div className="h-1 w-20 bg-primary mb-8"></div>
        <div className="flex items-center justify-between mb-6 w-full">
          <div className="flex gap-4 items-center">
            <Filter className="w-5 h-5 shrink-0" />
            <div className="flex gap-4 overflow-x-auto no-scrollbar w-full">
              {categories.map((category) => (
                <button
                  key={category.id}
                  onClick={() => setSelectedCategory(category.id)}
                  className={`shrink-0 px-4 py-2 rounded-lg border text-sm transition-colors whitespace-nowrap ${
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
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-6">
        {loading ? (
          <>
            <CardSkeleton />
            <CardSkeleton />
            <CardSkeleton />
            <CardSkeleton />
          </>
        ) : (
          <>
            {/* @ts-ignore */}
            {displayedProducts.map((product: Catalogue, index: number) => (
              <ProductCard product={product} key={index} />
            ))}
          </>
        )}
      </div>
      <div className="w-full flex items-center">
      <Link href="/product">
          <Button variant="outline">
            Load All
          </Button>
        </Link>
      </div>

    </section>
  );
}