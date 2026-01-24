import { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import { Search as SearchIcon, SlidersHorizontal, X } from 'lucide-react';
import { ProductCard } from '@/components/ProductCard';
import { products } from '@/data/products';
import { ProductCategory } from '@/types/product';
import { cn } from '@/lib/utils';

const categories: Array<ProductCategory | 'All'> = ['All', 'Equipment', 'Tools', 'Accessories', 'Footwear'];

const Search = () => {
  const [query, setQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<ProductCategory | 'All'>('All');

  const filteredProducts = useMemo(() => {
    return products.filter((product) => {
      const matchesQuery = query === '' || 
        product.name.toLowerCase().includes(query.toLowerCase()) ||
        product.brand.toLowerCase().includes(query.toLowerCase());
      
      const matchesCategory = selectedCategory === 'All' || product.category === selectedCategory;

      return matchesQuery && matchesCategory;
    });
  }, [query, selectedCategory]);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20 pt-2"
    >
      {/* Search Header */}
      <div className="sticky top-0 bg-background z-10 px-4 pb-4 space-y-3">
        {/* Search Input */}
        <div className="relative">
          <SearchIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search products, brands..."
            className="w-full h-11 pl-10 pr-10 bg-secondary rounded-sm text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-foreground/10"
          />
          {query && (
            <button
              onClick={() => setQuery('')}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>

        {/* Category Filters */}
        <div className="flex gap-2 overflow-x-auto scrollbar-hide -mx-4 px-4">
          {categories.map((category) => (
            <button
              key={category}
              onClick={() => setSelectedCategory(category)}
              className={cn(
                "px-4 py-2 text-sm font-medium rounded-sm whitespace-nowrap btn-press transition-colors",
                selectedCategory === category
                  ? "bg-primary text-primary-foreground"
                  : "bg-secondary text-muted-foreground hover:text-foreground"
              )}
            >
              {category}
            </button>
          ))}
        </div>
      </div>

      {/* Results */}
      <div className="px-4">
        <div className="flex items-center justify-between mb-4">
          <p className="text-sm text-muted-foreground">
            {filteredProducts.length} result{filteredProducts.length !== 1 ? 's' : ''}
          </p>
          <button className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground btn-press">
            <SlidersHorizontal className="w-4 h-4" />
            Sort
          </button>
        </div>

        {filteredProducts.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-muted-foreground">No products found</p>
            <button
              onClick={() => {
                setQuery('');
                setSelectedCategory('All');
              }}
              className="mt-2 text-sm underline hover:no-underline"
            >
              Clear filters
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-4">
            {filteredProducts.map((product, index) => (
              <ProductCard key={product.id} product={product} index={index} />
            ))}
          </div>
        )}
      </div>
    </motion.div>
  );
};

export default Search;
