
import { useState } from "react";
import { Filter, ShoppingCart, Search, User } from "lucide-react";
import { Link } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useCart } from "@/contexts/CartContext";

const categories = ["All", "Electronics", "Clothing", "Accessories", "Footwear"];

const mockProducts = [
  {
    id: 1,
    name: "Wireless Headphones",
    price: 199.99,
    description: "Premium wireless headphones with active noise cancellation and up to 30 hours of battery life.",
    images: ["https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500&q=80"],
    type: "electronics",
    colors: ["Black", "White", "Blue"],
  },
  {
    id: 2,
    name: "Smart Watch",
    price: 299.99,
    description: "Advanced smartwatch with health tracking features and AMOLED display.",
    images: ["https://images.unsplash.com/photo-1546868871-7041f2a55e12?w=500&q=80"],
    type: "electronics",
    colors: ["Black", "Silver", "Gold"],
  },
  {
    id: 3,
    name: "Cotton T-Shirt",
    price: 29.99,
    description: "Comfortable 100% cotton t-shirt for everyday wear.",
    images: ["https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500&q=80"],
    category: "Clothing",
    sizes: ["XS", "S", "M", "L", "XL"],
    colors: ["White", "Black", "Gray", "Navy"],
  },
  {
    id: 4,
    name: "Leather Wallet",
    price: 49.99,
    description: "Genuine leather wallet with multiple card slots and coin pocket.",
    images: ["https://images.unsplash.com/photo-1627123424574-724758594e93?w=500&q=80"],
    category: "Accessories",
    colors: ["Brown", "Black"],
  },
  {
    id: 5,
    name: "Running Shoes",
    price: 89.99,
    description: "Lightweight running shoes with responsive cushioning.",
    images: ["https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=500&q=80"],
    category: "Footwear",
    sizes: ["7", "8", "9", "10", "11", "12"],
    colors: ["Black/Red", "Blue/White", "Gray/Yellow"],
  },
  {
    id: 6,
    name: "Wireless Earbuds",
    price: 159.99,
    description: "True wireless earbuds with premium sound quality.",
    images: ["https://images.unsplash.com/photo-1572569511254-d8f925fe2cbb?w=500&q=80"],
    category: "Electronics",
    colors: ["White", "Black"],
  },
  {
    id: 7,
    name: "Denim Jeans",
    price: 79.99,
    description: "Classic fit denim jeans with stretch comfort.",
    images: ["https://images.unsplash.com/photo-1542272604-787c3835535d?w=500&q=80"],
    category: "Clothing",
    sizes: ["30x30", "32x32", "34x32", "36x32"],
    colors: ["Blue", "Black", "Gray"],
  },
  {
    id: 8,
    name: "Backpack",
    price: 69.99,
    description: "Durable backpack with laptop compartment and multiple pockets.",
    images: ["https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=500&q=80"],
    category: "Accessories",
    colors: ["Black", "Navy", "Gray"],
  },
  {
    id: 9,
    name: "Smart Speaker",
    price: 129.99,
    description: "Voice-controlled smart speaker with premium sound.",
    images: ["https://images.unsplash.com/photo-1543512214-318c7553f230?w=500&q=80"],
    category: "Electronics",
    colors: ["Black", "White"],
  },
  {
    id: 10,
    name: "Summer Dress",
    price: 59.99,
    description: "Lightweight summer dress with floral pattern.",
    images: ["https://images.unsplash.com/photo-1585487000160-6ebcfceb0d03?w=500&q=80"],
    category: "Clothing",
    sizes: ["XS", "S", "M", "L"],
    colors: ["Blue Floral", "Pink Floral", "White"],
  },
  {
    id: 11,
    name: "Sunglasses",
    price: 149.99,
    description: "Polarized sunglasses with UV protection.",
    images: ["https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=500&q=80"],
    category: "Accessories",
    colors: ["Black/Gold", "Tortoise/Brown"],
  },
  {
    id: 12,
    name: "Fitness Tracker",
    price: 89.99,
    description: "Water-resistant fitness tracker with heart rate monitoring.",
    images: ["https://images.unsplash.com/photo-1575311373937-040b8e1fd5b6?w=500&q=80"],
    category: "Electronics",
    colors: ["Black", "Blue", "Pink"],
  },
];

const Products = () => {
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [showSearch, setShowSearch] = useState(false);
  const { getTotalItems } = useCart();
  const productsPerPage = 9;

  const filteredProducts = mockProducts
    .filter(product => 
      (selectedCategory === "All" || product.category === selectedCategory) &&
      product.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

  const totalPages = Math.ceil(filteredProducts.length / productsPerPage);
  const startIndex = (currentPage - 1) * productsPerPage;
  const displayedProducts = filteredProducts.slice(startIndex, startIndex + productsPerPage);

  return (
    <div className="min-h-screen bg-white flex flex-col">
      <header className="fixed top-0 left-0 right-0 bg-white z-50 border-b border-gray-100">
        <nav className="container mx-auto px-4">
          <div className="flex justify-between items-center h-16">
            <Link to="/" className="text-xl font-semibold text-primary shrink-0">
              Store
            </Link>

            {/* Desktop Navigation */}
            <div className="hidden md:flex items-center gap-4">
              <div className="relative w-64">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  type="text"
                  placeholder="Search products..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Link 
                to="/login" 
                className="p-2 hover:bg-secondary rounded-full transition-colors"
              >
                <User className="w-5 h-5" />
              </Link>
              <Link
                to="/cart"
                className="relative p-2 hover:bg-secondary rounded-full transition-colors"
              >
                <ShoppingCart className="w-5 h-5" />
                {getTotalItems() > 0 && (
                  <span className="absolute -top-1 -right-1 bg-primary text-white text-xs w-5 h-5 rounded-full flex items-center justify-center">
                    {getTotalItems()}
                  </span>
                )}
              </Link>
            </div>

            {/* Mobile Navigation */}
            <div className="flex md:hidden items-center gap-2">
              <Button
                variant="ghost"
                size="icon"
                className="md:hidden"
                onClick={() => setShowSearch(!showSearch)}
              >
                <Search className="w-5 h-5" />
              </Button>
              <Link
                to="/cart"
                className="relative p-2 hover:bg-secondary rounded-full transition-colors"
              >
                <ShoppingCart className="w-5 h-5" />
                {getTotalItems() > 0 && (
                  <span className="absolute -top-1 -right-1 bg-primary text-white text-xs w-5 h-5 rounded-full flex items-center justify-center">
                    {getTotalItems()}
                  </span>
                )}
              </Link>
              <Link 
                to="/login" 
                className="p-2 hover:bg-secondary rounded-full transition-colors"
              >
                <User className="w-5 h-5" />
              </Link>
            </div>
          </div>

          {/* Mobile Search Bar */}
          {showSearch && (
            <div className="md:hidden px-4 py-3 border-t border-gray-100">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  type="text"
                  placeholder="Search products..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 w-full"
                />
              </div>
            </div>
          )}
        </nav>
      </header>

      <main className="container mx-auto px-4 pt-24 flex-1">
        <div className="flex items-center gap-4 mb-8 w-full">
          <Filter className="w-5 h-5 shrink-0" />
          <div className="flex gap-4 overflow-x-auto no-scrollbar w-full">
            {categories.map((category) => (
              <button
                key={category}
                onClick={() => setSelectedCategory(category)}
                className={`shrink-0 px-4 py-2 rounded-full text-sm transition-colors whitespace-nowrap ${
                  selectedCategory === category
                    ? "bg-primary text-white"
                    : "bg-secondary text-primary hover:bg-opacity-80"
                }`}
              >
                {category}
              </button>
            ))}
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8 mb-8">
          {displayedProducts.map((product) => (
            <Link
              key={product.id}
              to={`/product/${product.id}`}
              className="group animate-fade-in"
            >
              <div className="aspect-square overflow-hidden rounded-lg bg-secondary mb-4">
                <img
                  src={product.images[0]}
                  alt={product.name}
                  className="w-full h-full object-cover transform transition-transform group-hover:scale-105"
                />
              </div>
              <h3 className="text-lg font-medium text-primary mb-2">
                {product.name}
              </h3>
              <p className="text-sm text-primary">${product.price}</p>
            </Link>
          ))}
        </div>

        <Pagination className="my-8">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious 
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                className={currentPage === 1 ? "pointer-events-none opacity-50" : ""}
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
                onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                className={currentPage === totalPages ? "pointer-events-none opacity-50" : ""}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </main>

      <footer className="mt-auto border-t border-gray-100">
        <div className="container mx-auto px-4 py-4">
          <p className="text-sm text-gray-500 text-center">
            © {new Date().getFullYear()} Store. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Products;
