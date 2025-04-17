import { useState, useEffect } from "react";
import { Filter, ShoppingCart, Search, User } from "lucide-react";
import { Link } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { fetchCategories } from "@/api/CategoryApi";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useCart } from "@/contexts/CartContext";
import { fetchProducts } from "@/api/ProductApi";

const Products = () => {
  const [categories, setCategories] = useState([{ id: 0, name: "All" }]);
  const [products, setProducts] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(0);
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [showSearch, setShowSearch] = useState(false);
  const { getTotalItems } = useCart();
  const productsPerPage = 9;

  useEffect(() => {
    const loadCategories = async () => {
      try {
        const fetchedCategories = await fetchCategories();
        setCategories([{ id: 0, name: "All" }, ...fetchedCategories]);
      } catch (error) {
        console.error("Error fetching categories:", error);
      }
    };

    loadCategories();
  }, []);

  useEffect(() => {
    const loadProducts = async () => {
      try {
        const productList = await fetchProducts();
        // console.log(productList);
        setProducts(productList);
      } catch (error) {
        console.error("Failed to load products:", error);
      }
    };

    loadProducts();
  }, []);

  const getPriceRange = (product) => {
    if (!product.configurables || product.configurables.length === 0) {
      return `$${product.price}`;
    }

    const prices = product.configurables.map((config) => config.price);
    const minPrice = Math.min(...prices);
    const maxPrice = Math.max(...prices);

    if (minPrice === maxPrice) {
      return `$${minPrice}`;
    } else {
      return `$${minPrice} - $${maxPrice}`;
    }
  };

  const startIndex = (currentPage - 1) * productsPerPage;

  // Filter products based on selected category and search query
  const filteredProducts = products.filter((product) => {
    const matchesCategory =
      selectedCategory === 0 || // "All" category
      product.configurables.some((configurable) =>
        configurable.categories.includes(selectedCategory)
      );

    const matchesSearchQuery = product.name
      .toLowerCase()
      .includes(searchQuery.toLowerCase());

    return matchesCategory && matchesSearchQuery;
  });

  const totalPages = Math.ceil(filteredProducts.length / productsPerPage);
  const displayedProducts = filteredProducts.slice(
    startIndex,
    startIndex + productsPerPage
  );

  return (
    <div className='min-h-screen bg-white flex flex-col'>
      <header className='fixed top-0 left-0 right-0 bg-white z-50 border-b border-gray-100'>
        <nav className='container mx-auto px-4'>
          <div className='flex justify-between items-center h-16'>
            <Link
              to='/'
              className='text-xl font-semibold text-primary shrink-0'
            >
              Store
            </Link>

            {/* Desktop Navigation */}
            <div className='hidden md:flex items-center gap-4'>
              <div className='relative w-64'>
                <Search className='absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4' />
                <Input
                  type='text'
                  placeholder='Search products...'
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className='pl-10'
                />
              </div>
              <Link
                to='/login'
                className='p-2 hover:bg-secondary rounded-full transition-colors'
              >
                <User className='w-5 h-5' />
              </Link>
              <Link
                to='/cart'
                className='relative p-2 hover:bg-secondary rounded-full transition-colors'
              >
                <ShoppingCart className='w-5 h-5' />
                {getTotalItems() > 0 && (
                  <span className='absolute -top-1 -right-1 bg-primary text-white text-xs w-5 h-5 rounded-full flex items-center justify-center'>
                    {getTotalItems()}
                  </span>
                )}
              </Link>
            </div>

            {/* Mobile Navigation */}
            <div className='flex md:hidden items-center gap-2'>
              <Button
                variant='ghost'
                size='icon'
                className='md:hidden'
                onClick={() => setShowSearch(!showSearch)}
              >
                <Search className='w-5 h-5' />
              </Button>
              <Link
                to='/cart'
                className='relative p-2 hover:bg-secondary rounded-full transition-colors'
              >
                <ShoppingCart className='w-5 h-5' />
                {getTotalItems() > 0 && (
                  <span className='absolute -top-1 -right-1 bg-primary text-white text-xs w-5 h-5 rounded-full flex items-center justify-center'>
                    {getTotalItems()}
                  </span>
                )}
              </Link>
              <Link
                to='/login'
                className='p-2 hover:bg-secondary rounded-full transition-colors'
              >
                <User className='w-5 h-5' />
              </Link>
            </div>
          </div>

          {/* Mobile Search Bar */}
          {showSearch && (
            <div className='md:hidden px-4 py-3 border-t border-gray-100'>
              <div className='relative'>
                <Search className='absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4' />
                <Input
                  type='text'
                  placeholder='Search products...'
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className='pl-10 w-full'
                />
              </div>
            </div>
          )}
        </nav>
      </header>

      <main className='container mx-auto px-4 pt-24 flex-1'>
        <div className='flex items-center gap-4 mb-8 w-full'>
          <Filter className='w-5 h-5 shrink-0' />
          <div className='flex gap-4 overflow-x-auto no-scrollbar w-full'>
            {categories.map((category) => (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={`shrink-0 px-4 py-2 rounded-full text-sm transition-colors whitespace-nowrap ${
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

        <div className='grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8 mb-8'>
          {displayedProducts.map((product) => {
            return (
              <Link
                key={product.id}
                to={`/product/${product.id}`}
                state={{ product }}
                className='group animate-fade-in'
              >
                <div className='aspect-square overflow-hidden rounded-lg bg-secondary mb-4'>
                  <img
                    src={product.images[0]}
                    alt={product.name}
                    className='w-full h-full object-cover transform transition-transform group-hover:scale-105'
                  />
                </div>
                <h3 className='text-lg font-medium text-primary mb-2'>
                  {product.name}
                </h3>
                <p className='text-sm text-primary'>{getPriceRange(product)}</p>
              </Link>
            );
          })}
        </div>

        <Pagination className='my-8'>
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
                className={
                  currentPage === 1 ? "pointer-events-none opacity-50" : ""
                }
              />
            </PaginationItem>
            {[...Array(totalPages)].map((_, i) => (
              <PaginationItem key={i + 1}>
                <PaginationLink
                  onClick={() => setCurrentPage(i + 1)}
                  isActive={currentPage === i + 1}
                >
                  {i + 1}
                </PaginationLink>
              </PaginationItem>
            ))}
            <PaginationItem>
              <PaginationNext
                onClick={() =>
                  setCurrentPage((prev) => Math.min(totalPages, prev + 1))
                }
                className={
                  currentPage === totalPages
                    ? "pointer-events-none opacity-50"
                    : ""
                }
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </main>

      <footer className='mt-auto border-t border-gray-100'>
        <div className='container mx-auto px-4 py-4'>
          <p className='text-sm text-gray-500 text-center'>
            © {new Date().getFullYear()} Store. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Products;
