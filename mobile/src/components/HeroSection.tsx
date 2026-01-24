import { useNavigate, Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronRight, Search, SlidersHorizontal, X, Grid3x3 } from 'lucide-react';
import { useState, useEffect, useRef } from 'react';
import { GetAllSuppliers } from '@/lib/api/supplier';
import { GetAllCategories } from '@/lib/api/category';
import { Supplier, Category } from '@/lib/types';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

// Mock suppliers data for testing
const mockSuppliers: Supplier[] = [
  {
    id: 1,
    name: 'TechSupply Co.',
    tin: 'TIN001',
    latitude: '9.1450',
    longitude: '38.7610',
    general_zone: 'Central',
    region: 'Addis Ababa',
    woreda: 'Bole',
    user: [1],
    is_active: true,
  },
  {
    id: 2,
    name: 'Industrial Solutions',
    tin: 'TIN002',
    latitude: '9.1450',
    longitude: '38.7610',
    general_zone: 'Central',
    region: 'Addis Ababa',
    woreda: 'Kirkos',
    user: [2],
    is_active: true,
  },
  {
    id: 3,
    name: 'Global Equipment',
    tin: 'TIN003',
    latitude: '9.1450',
    longitude: '38.7610',
    general_zone: 'Central',
    region: 'Addis Ababa',
    woreda: 'Arada',
    user: [3],
    is_active: true,
  },
  {
    id: 4,
    name: 'Quality Tools Ltd',
    tin: 'TIN004',
    latitude: '9.1450',
    longitude: '38.7610',
    general_zone: 'Central',
    region: 'Addis Ababa',
    woreda: 'Lideta',
    user: [4],
    is_active: true,
  },
  {
    id: 5,
    name: 'Professional Supplies',
    tin: 'TIN005',
    latitude: '9.1450',
    longitude: '38.7610',
    general_zone: 'Central',
    region: 'Addis Ababa',
    woreda: 'Nifas Silk',
    user: [5],
    is_active: true,
  },
  {
    id: 6,
    name: 'Business Essentials',
    tin: 'TIN006',
    latitude: '9.1450',
    longitude: '38.7610',
    general_zone: 'Central',
    region: 'Addis Ababa',
    woreda: 'Yeka',
    user: [6],
    is_active: true,
  },
];

interface HeroSectionProps {
  onCategoryChange?: (categoryId: number) => void;
}

export const HeroSection = ({ onCategoryChange }: HeroSectionProps) => {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('All');
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filterCategory, setFilterCategory] = useState<number | null>(null);
  const navigate = useNavigate();
  const suppliersScrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        const [supplierList, categoryList] = await Promise.all([
          GetAllSuppliers().catch(() => []),
          GetAllCategories().catch(() => [])
        ]);
        // Only show active suppliers, limit to 6 for display
        const activeSuppliers = supplierList.filter(s => s.is_active).slice(0, 6);
        setSuppliers(activeSuppliers.length > 0 ? activeSuppliers : mockSuppliers);
        setCategories(categoryList);
      } catch (error) {
        console.error('Failed to load data:', error);
        setSuppliers(mockSuppliers);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate('/search');
    }
  };

  const scrollSuppliers = (direction: 'left' | 'right') => {
    if (suppliersScrollRef.current) {
      const scrollAmount = 200;
      suppliersScrollRef.current.scrollBy({
        left: direction === 'right' ? scrollAmount : -scrollAmount,
        behavior: 'smooth'
      });
    }
  };

  const handleApplyFilters = () => {
    setIsFilterOpen(false);
    // Apply category filter locally
    if (filterCategory !== null) {
      setSelectedCategory(filterCategory.toString());
      // Notify parent component of category change
      if (onCategoryChange) {
        onCategoryChange(filterCategory);
      }
    } else {
      setSelectedCategory('All');
      // Notify parent component to show all categories
      if (onCategoryChange) {
        onCategoryChange(0);
      }
    }
  };

  const handleClearFilters = () => {
    setFilterCategory(null);
    setSelectedCategory('All');
    // Notify parent component to show all categories
    if (onCategoryChange) {
      onCategoryChange(0);
    }
  };

  return (
    <section className="relative w-full">
      {/* Fixed Search Bar */}
      <div className="fixed top-0 left-0 right-0 z-40 bg-background/95 backdrop-blur-sm border-b border-border">
        <div className="px-6 sm:px-12 py-3">
          <form onSubmit={handleSearch} className="max-w-xl mx-auto">
            <div className="flex items-center gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-primary z-10" />
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Search products, brands..."
                  className="w-full h-12 pl-11 pr-4 bg-background rounded-md border border-border text-sm placeholder:text-gray-500 text-black focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                />
              </div>
              <button
                type="button"
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  // Sync filterCategory with current selectedCategory when opening
                  if (selectedCategory === 'All') {
                    setFilterCategory(null);
                  } else {
                    const categoryId = parseInt(selectedCategory);
                    setFilterCategory(isNaN(categoryId) ? null : categoryId);
                  }
                  setIsFilterOpen(true);
                }}
                className="h-12 w-12 flex items-center justify-center bg-background rounded-md border border-border text-primary hover:bg-secondary transition-colors btn-press"
                aria-label="Filter"
              >
                <SlidersHorizontal className="w-5 h-5" />
              </button>
            </div>
          </form>
        </div>
      </div>

      {/* Fixed Category Filters */}
      {categories.length > 0 && (
        <div className="fixed top-[72px] left-0 right-0 z-40 bg-card border-b border-border">
          <div className="px-4 py-3">
            <div className="flex gap-2 overflow-x-auto scrollbar-hide -mx-4 px-4">
              <button
                onClick={() => {
                  setSelectedCategory('All');
                  // Notify parent component to show all categories
                  if (onCategoryChange) {
                    onCategoryChange(0);
                  }
                }}
                className={cn(
                  "px-4 py-2 text-sm font-medium rounded-sm whitespace-nowrap btn-press transition-colors",
                  selectedCategory === 'All'
                    ? "bg-primary text-primary-foreground"
                    : "bg-secondary text-muted-foreground hover:text-foreground"
                )}
              >
                All
              </button>
              {categories.map((category) => (
                <button
                  key={category.id}
                  onClick={() => {
                    setSelectedCategory(category.id.toString());
                    // Notify parent component of category change
                    if (onCategoryChange) {
                      onCategoryChange(category.id);
                    }
                  }}
                  className={cn(
                    "px-4 py-2 text-sm font-medium rounded-sm whitespace-nowrap btn-press transition-colors",
                    selectedCategory === category.id.toString()
                      ? "bg-primary text-primary-foreground"
                      : "bg-secondary text-muted-foreground hover:text-foreground"
                  )}
                >
                  {category.name}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Spacer to account for fixed search bar and category filters */}
      <div className={categories.length > 0 ? "pt-[120px]" : "pt-[72px]"}></div>

      {/* Suppliers List */}
      <div>
        <div className="px-4 pt-3">
          <h3 className="text-xs font-semibold text-muted-foreground tracking-wide mb-2">
            Top Suppliers
          </h3>
        </div>
        <div className="border-b border-border">
          <div className="px-4 py-2">
            {loading ? (
              <div className="flex gap-2 overflow-x-auto pb-2">
                {[...Array(6)].map((_, i) => (
                  <div key={i} className="shrink-0 w-20 h-16 bg-muted rounded-md animate-pulse" />
                ))}
              </div>
            ) : suppliers.length === 0 ? (
              <p className="text-xs text-muted-foreground">No suppliers available</p>
            ) : (
              <div className="relative">
                <div 
                  ref={suppliersScrollRef}
                  className="flex gap-2 overflow-x-auto pb-2 scrollbar-hide"
                >
                  {suppliers.map((supplier, index) => (
                    <motion.div
                      key={supplier.id}
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="shrink-0 bg-white rounded-md px-3 py-2 border border-border min-w-[100px]"
                    >
                      <p className="text-xs font-semibold text-black truncate mb-0.5">
                        {supplier.name}
                      </p>
                      <p className="text-[10px] text-gray-600 truncate">
                        {supplier.region}
                      </p>
                    </motion.div>
                  ))}
                </div>
                <button
                  onClick={() => scrollSuppliers('right')}
                  className="absolute right-0 top-1/2 -translate-y-1/2 bg-white/90 backdrop-blur-sm rounded-full p-2 shadow-md hover:bg-white transition-colors btn-press z-10"
                  aria-label="Scroll right"
                >
                  <ChevronRight className="w-4 h-4 text-black" />
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Categories Section */}
      {categories && categories.length > 0 && (
        <div className="border-b border-border">
          <div className="px-4 pt-3">
            <h3 className="text-xs font-semibold text-muted-foreground tracking-wide mb-2">
              Categories
            </h3>
          </div>
          <div className="px-4 py-2">
            <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
              {categories.filter(cat => cat && cat.id && cat.name).map((category, index) => (
                <Link
                  key={`category-${category.id}-${index}`}
                  to={`/search?category=${category.id}`}
                  className="group shrink-0"
                  onClick={() => {
                    if (onCategoryChange && category.id) {
                      onCategoryChange(category.id);
                    }
                  }}
                >
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.05 }}
                    className="bg-card border border-border rounded-md p-3 hover:border-primary hover:shadow-md transition-all btn-press w-[100px] flex flex-col items-center justify-center min-h-[90px]"
                  >
                    <div className="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center mb-2 group-hover:bg-primary/20 transition-colors">
                      <Grid3x3 className="w-4 h-4 text-primary" />
                    </div>
                    <h3 className="text-[10px] font-semibold text-center line-clamp-2 leading-tight">
                      {category.name || 'Category'}
                    </h3>
                  </motion.div>
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Filter Drawer */}
      <AnimatePresence>
        {isFilterOpen && (
          <>
            {/* Overlay */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsFilterOpen(false)}
              className="fixed inset-0 bg-black/40 z-[100]"
            />

            {/* Drawer */}
            <motion.div
              initial={{ x: '100%' }}
              animate={{ x: 0 }}
              exit={{ x: '100%' }}
              transition={{ type: 'spring', damping: 30, stiffness: 300 }}
              className="fixed right-0 top-0 bottom-0 w-full max-w-md bg-background z-[100] flex flex-col shadow-2xl"
            >
              {/* Header */}
              <div className="flex items-center justify-between p-4 border-b border-border">
                <h2 className="font-display text-lg">Filters</h2>
                <button
                  onClick={() => setIsFilterOpen(false)}
                  className="p-2 hover:bg-secondary rounded-sm btn-press transition-colors"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>

              {/* Filter Content */}
              <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {/* Category Filter */}
                <div>
                  <h3 className="text-sm font-semibold mb-3">Categories</h3>
                  {loading ? (
                    <div className="space-y-2">
                      {[...Array(5)].map((_, i) => (
                        <div key={i} className="h-10 bg-muted rounded-md animate-pulse" />
                      ))}
                    </div>
                  ) : categories.length === 0 ? (
                    <p className="text-sm text-muted-foreground">No categories available</p>
                  ) : (
                    <div className="space-y-2">
                      <button
                        onClick={() => setFilterCategory(null)}
                        className={cn(
                          "w-full text-left px-4 py-2 rounded-md text-sm btn-press transition-colors",
                          filterCategory === null
                            ? "bg-primary text-primary-foreground"
                            : "bg-secondary text-foreground hover:bg-secondary/80"
                        )}
                      >
                        All Categories
                      </button>
                      {categories.map((category) => (
                        <button
                          key={category.id}
                          onClick={() => setFilterCategory(category.id)}
                          className={cn(
                            "w-full text-left px-4 py-2 rounded-md text-sm btn-press transition-colors",
                            filterCategory === category.id
                              ? "bg-primary text-primary-foreground"
                              : "bg-secondary text-foreground hover:bg-secondary/80"
                          )}
                        >
                          {category.name}
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              </div>

              {/* Footer */}
              <div className="p-4 border-t border-border space-y-3 safe-bottom">
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    onClick={handleClearFilters}
                    className="flex-1"
                  >
                    Clear
                  </Button>
                  <Button
                    onClick={handleApplyFilters}
                    className="flex-1"
                  >
                    Apply Filters
                  </Button>
                </div>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </section>
  );
};
