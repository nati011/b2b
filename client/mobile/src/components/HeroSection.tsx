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
  onSearchChange?: (query: string) => void;
}

export const HeroSection = ({ onCategoryChange, onSearchChange }: HeroSectionProps) => {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('All');
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filterCategory, setFilterCategory] = useState<number | null>(null);
  const [isScrolled, setIsScrolled] = useState(false);
  const navigate = useNavigate();
  const suppliersScrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        const [supplierResponse, categoryList] = await Promise.all([
          GetAllSuppliers({ limit: 100, page: 1 }).catch((error) => {
            console.error('❌ Failed to fetch suppliers in HeroSection:', error);
            return { items: [], total: 0, page: 1, limit: 100, total_pages: 0, has_next: false, has_prev: false };
          }),
          GetAllCategories().catch(() => [])
        ]);
        
        // Convert SupplierResponse[] to Supplier[] format for compatibility
        const supplierItems = supplierResponse.items || [];
        const activeSuppliers: Supplier[] = supplierItems
          .filter((s: any) => s.is_active)
          .slice(0, 6)
          .map((s: any) => ({
            id: s.id,
            name: s.business_name,
            tin: '',
            latitude: '',
            longitude: '',
            general_zone: '',
            region: '',
            woreda: '',
            user: [],
            is_active: s.is_active
          }));
        
        // Only use API data, fallback to empty array if no suppliers
        setSuppliers(activeSuppliers.length > 0 ? activeSuppliers : []);
        setCategories(categoryList);
      } catch (error) {
        console.error('Failed to load data:', error);
        setSuppliers([]);
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

  useEffect(() => {
    const handleScroll = () => {
      // Change to brand color when scrolled down more than 50px
      const scrollY = window.scrollY || window.pageYOffset;
      setIsScrolled(scrollY > 50);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    handleScroll(); // Check initial state
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <section className="relative w-full">
      {/* Fixed Search Bar */}
      <div className={`fixed top-0 left-0 right-0 z-40 backdrop-blur-md shadow-sm transition-colors duration-300 ${
        isScrolled ? 'bg-primary/95' : 'bg-background/95'
      }`}>
        <div className="px-4 sm:px-6 py-3">
          <form onSubmit={handleSearch} className="max-w-xl mx-auto">
            <div className="flex items-center gap-2.5">
              <div className="relative flex-1">
                <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4.5 h-4.5 text-muted-foreground z-10" />
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => {
                    const query = e.target.value;
                    setSearchQuery(query);
                    if (onSearchChange) {
                      onSearchChange(query);
                    }
                  }}
                  placeholder="Search products, brands..."
                  className="w-full h-11 pl-11 pr-10 bg-card rounded-lg border border-border/80 text-sm placeholder:text-muted-foreground text-foreground focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all shadow-sm"
                />
                {searchQuery && (
                  <button
                    type="button"
                    onClick={() => {
                      setSearchQuery('');
                      if (onSearchChange) {
                        onSearchChange('');
                      }
                    }}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors z-10 p-1 rounded-sm hover:bg-muted"
                    aria-label="Clear search"
                  >
                    <X className="w-4 h-4" />
                  </button>
                )}
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
                className="h-11 w-11 flex items-center justify-center bg-card rounded-lg border border-border/80 text-primary hover:bg-primary/5 hover:border-primary/40 transition-all btn-press shadow-sm"
                aria-label="Filter"
              >
                <SlidersHorizontal className="w-4.5 h-4.5" />
              </button>
            </div>
          </form>
        </div>
      </div>

      {/* Fixed Category Filters */}
      {categories.length > 0 && (
        <div className={`fixed top-[60px] left-0 right-0 z-40 backdrop-blur-md shadow-sm transition-colors duration-300 ${
          isScrolled ? 'bg-primary/95' : 'bg-background/95'
        }`}>
          <div className="px-4 py-2.5">
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
                  "px-4 py-1.5 text-xs font-medium rounded-lg whitespace-nowrap btn-press transition-all duration-200",
                  selectedCategory === 'All'
                    ? "bg-primary text-primary-foreground shadow-sm"
                    : "bg-card border border-border/80 text-muted-foreground hover:text-foreground hover:border-primary/40 hover:bg-primary/5"
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
                    "px-4 py-1.5 text-xs font-medium rounded-lg whitespace-nowrap btn-press transition-all duration-200",
                    selectedCategory === category.id.toString()
                      ? "bg-primary text-primary-foreground shadow-sm"
                      : "bg-card border border-border/80 text-muted-foreground hover:text-foreground hover:border-primary/40 hover:bg-primary/5"
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
      <div className={categories.length > 0 ? "pt-[108px]" : "pt-[68px]"}></div>

      {/* Filter Dropdown */}
      <AnimatePresence>
        {isFilterOpen && (
          <>
            {/* Overlay */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsFilterOpen(false)}
              className="fixed inset-0 bg-black/20 z-[50]"
            />

            {/* Dropdown */}
            <motion.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2 }}
              className="fixed right-4 top-[88px] w-[calc(100%-2rem)] max-w-sm bg-background border border-border rounded-md shadow-lg z-[100] max-h-[70vh] overflow-hidden"
            >
              {/* Header */}
              <div className="flex items-center justify-between p-4 border-b border-border">
                <h2 className="font-display text-lg font-semibold">Categories</h2>
                <button
                  onClick={() => setIsFilterOpen(false)}
                  className="p-2 hover:bg-secondary rounded-sm btn-press transition-colors"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>

              {/* Filter Content */}
              <div className="overflow-y-auto max-h-[calc(70vh-80px)] p-2">
                {loading ? (
                  <div className="space-y-2 p-2">
                    {[...Array(5)].map((_, i) => (
                      <div key={i} className="h-10 bg-muted rounded-md animate-pulse" />
                    ))}
                  </div>
                ) : categories.length === 0 ? (
                  <p className="text-sm text-muted-foreground p-4 text-center">No categories available</p>
                ) : (
                  <div className="space-y-1">
                    <button
                      onClick={() => {
                        setFilterCategory(null);
                        handleApplyFilters();
                      }}
                      className={cn(
                        "w-full text-left px-4 py-3 rounded-md text-sm btn-press transition-colors",
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
                        onClick={() => {
                          setFilterCategory(category.id);
                          handleApplyFilters();
                        }}
                        className={cn(
                          "w-full text-left px-4 py-3 rounded-md text-sm btn-press transition-colors",
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
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </section>
  );
};
