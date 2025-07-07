"use client";
import { useEffect, useState, useMemo } from "react";
import { Filter, ShoppingBagIcon } from "lucide-react";
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
import { GetAllCatalogues } from "@/app/actions/catalogue";
import { Button } from "@/components/ui/button";
import { CardSkeleton } from "@/components/CardSkeleton";
import { ProductCard } from "@/components/ProductCard";
import { Catalogue } from "@/lib/types";

const ITEMS_PER_PAGE = 12; // Number of products per page

const Product = () => {
  const [currentPage, setCurrentPage] = useState(0);
  const [error, setError] = useState<string>();
  const [categories, setCategories] = useState([{ id: 0, name: "All" }]);
  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [selectedCategory, setSelectedCategory] = useState(0);

  const loadCategories = async () => {
    try {
      const fetchedCategories = await GetAllCategories();
      setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
    } catch (error) {
      console.error("Error fetching categories:", error);
    }
  };

  const loadProducts = async () => {
    setLoading(true);
    try {
      await loadCategories();
      const productList = await GetAllCatalogues();
      setProducts(productList);
    } catch (error: any) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProducts();
  }, []);
  const filteredProducts = useMemo(() => {
    if (selectedCategory === 0) return products;
    return products.filter((product) => {
      return product.configurables.some((configurable) =>
        configurable.categories?.includes(selectedCategory)
      );
    });
  }, [products, selectedCategory]);

  const paginatedProducts = useMemo(() => {
    const startIndex = currentPage * ITEMS_PER_PAGE;
    return filteredProducts.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredProducts, currentPage]);

  const totalPages = Math.ceil(filteredProducts.length / ITEMS_PER_PAGE);

  const handlePageChange = (newPage: number) => {
    if (newPage >= 0 && newPage < totalPages) {
      setCurrentPage(newPage);
    }
  };

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

  return (
    <div className="w-full max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div>
        <div className="flex flex-col items-center mb-12">
          <h2 className="text-3xl font-medium mb-4">Products</h2>
          <div className="h-1 w-20 bg-primary mb-8"></div>
          <div className="flex items-center justify-between mb-6 mx-auto">
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

        {loading ? (
          <div className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
              {[...Array(10)].map((_, i) => (
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
                  <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-4 gap-6">
                    {paginatedProducts.map((product, index) => (
                      <ProductCard product={product} key={index} />
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
  );
};

export default Product;
