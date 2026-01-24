import { useState, useEffect, useMemo, useRef } from 'react';
import { motion } from 'framer-motion';
import { HeroSection } from '@/components/HeroSection';
import { CatalogueCard } from '@/components/CatalogueCard';
import { GetAllCatalogues } from '@/lib/api/catalogue';
import { GetAllCategories } from '@/lib/api/category';
import { GetAllSuppliers } from '@/lib/api/supplier';
import { Catalogue, Category, Supplier } from '@/lib/types';
import { Filter, Grid3x3, ChevronRight, Building2, MapPin } from 'lucide-react';
import { Link } from 'react-router-dom';

// Sample categories for fallback
const sampleCategories: Category[] = [
  { id: 1, name: 'Electronics' },
  { id: 2, name: 'Office Supplies' },
  { id: 3, name: 'Tools & Equipment' },
  { id: 4, name: 'Furniture' },
  { id: 5, name: 'Safety & Security' },
  { id: 6, name: 'Industrial' },
];

// Sample suppliers for fallback
const sampleSuppliers: Supplier[] = [
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

const Home = () => {
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number>(0);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const categoriesScrollRef = useRef<HTMLDivElement>(null);
  const suppliersScrollRef = useRef<HTMLDivElement>(null);
  const [showScrollButton, setShowScrollButton] = useState(false);
  const [showSuppliersScrollButtons, setShowSuppliersScrollButtons] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const [productList, categoryList, supplierList] = await Promise.all([
          GetAllCatalogues().catch(() => []),
          GetAllCategories().catch(() => []),
          GetAllSuppliers().catch(() => [])
        ]);
        setProducts(productList);
        setCategories(categoryList);
        const activeSuppliers = supplierList.filter(s => s.is_active).slice(0, 6);
        setSuppliers(activeSuppliers);
      } catch (error) {
        console.error('Failed to load data:', error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  useEffect(() => {
    const checkScroll = () => {
      if (categoriesScrollRef.current) {
        const { scrollWidth, clientWidth } = categoriesScrollRef.current;
        setShowScrollButton(scrollWidth > clientWidth);
      }
    };

    checkScroll();
    window.addEventListener('resize', checkScroll);
    return () => window.removeEventListener('resize', checkScroll);
  }, [categories, loading]);

  useEffect(() => {
    const checkSuppliersScroll = () => {
      if (suppliersScrollRef.current) {
        const { scrollWidth, clientWidth } = suppliersScrollRef.current;
        setShowSuppliersScrollButtons(scrollWidth > clientWidth);
      }
    };

    checkSuppliersScroll();
    window.addEventListener('resize', checkSuppliersScroll);
    return () => window.removeEventListener('resize', checkSuppliersScroll);
  }, [suppliers, loading]);

  const scrollCategories = (direction: 'left' | 'right') => {
    if (categoriesScrollRef.current) {
      const scrollAmount = 200;
      categoriesScrollRef.current.scrollBy({
        left: direction === 'right' ? scrollAmount : -scrollAmount,
        behavior: 'smooth'
      });
    }
  };

  const scrollSuppliers = (direction: 'left' | 'right') => {
    if (suppliersScrollRef.current) {
      const scrollAmount = 300;
      suppliersScrollRef.current.scrollBy({
        left: direction === 'right' ? scrollAmount : -scrollAmount,
        behavior: 'smooth'
      });
    }
  };

  const displayedProducts = useMemo(() => {
    let filtered = products;

    // Filter by search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter((product) => {
        const nameMatch = product.name?.toLowerCase().includes(query);
        const descMatch = product.desc?.toLowerCase().includes(query);
        const configurableMatch = product.configurables?.some((config) =>
          config.name?.toLowerCase().includes(query)
        );
        return nameMatch || descMatch || configurableMatch;
      });
    }

    // Filter by category
    if (selectedCategory !== 0) {
      filtered = filtered.filter((product) =>
        product.configurables?.some((configurable) =>
          configurable.categories?.includes(selectedCategory)
        )
      );
    }

    // When searching, show all results. Otherwise, limit to 4 for display
    return searchQuery.trim() ? filtered : filtered.slice(0, 4);
  }, [products, searchQuery, selectedCategory]);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      <HeroSection 
        onCategoryChange={setSelectedCategory}
        onSearchChange={setSearchQuery}
      />

      {/* Categories Section */}
      {!searchQuery.trim() && categories.length > 0 && (
        <section className="px-4 py-6">
          <div className="flex items-center gap-2 mb-4">
            <h2 className="font-display text-lg font-semibold">Categories</h2>
          </div>
          <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
            {categories.map((category, index) => (
              <Link
                key={category.id}
                to={`/search?category=${category.id}`}
                className="group shrink-0"
                onClick={() => setSelectedCategory(category.id)}
              >
                <motion.div
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: index * 0.05 }}
                  className="bg-card border border-border rounded-md p-4 hover:border-primary hover:shadow-md transition-all btn-press w-[120px] flex flex-col items-center justify-center min-h-[100px]"
                >
                  <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center mb-2 group-hover:bg-primary/20 transition-colors">
                    <Grid3x3 className="w-5 h-5 text-primary" />
                  </div>
                  <h3 className="text-xs font-semibold text-center line-clamp-2">{category.name}</h3>
                </motion.div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Categories Section */}
      <div className="mt-2 mb-0 pb-4 relative">
        {(categories.length > 0 || (!loading && categories.length === 0)) && (
          <div className="flex items-center justify-between px-4 mb-4">
            <h2 className="text-base font-semibold text-foreground tracking-tight">Shop by Category</h2>
            <Link
              to="/categories"
              className="text-xs text-muted-foreground hover:text-primary font-medium transition-colors flex items-center gap-1"
            >
              View All
              <span className="text-[10px] leading-none">→</span>
            </Link>
          </div>
        )}
        <div className="relative">
          {loading ? (
            <div ref={categoriesScrollRef} className="flex gap-2 overflow-x-auto pb-1 scrollbar-hide px-4">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="bg-muted/30 border border-border/50 rounded-xl w-[90px] h-[52px] shrink-0 animate-pulse" />
              ))}
            </div>
          ) : (
            <div ref={categoriesScrollRef} className="flex gap-2 overflow-x-auto pb-1 scrollbar-hide px-4">
              {(categories.length > 0 ? categories : sampleCategories).map((category, index) => (
                <Link
                  key={category.id}
                  to={`/search?category=${category.id}`}
                  className="group shrink-0"
                  onClick={() => setSelectedCategory(category.id)}
                >
                  <motion.div
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.02, duration: 0.2, ease: "easeOut" }}
                    className="bg-gradient-to-br from-card to-card/80 border border-border/60 rounded-xl p-2.5 hover:border-primary/50 hover:from-primary/5 hover:to-primary/10 hover:shadow-md transition-all duration-200 btn-press w-[90px] h-[52px] flex flex-col items-center justify-center shadow-sm"
                  >
                    <h3 className="text-[11px] font-semibold text-center line-clamp-2 leading-tight text-foreground group-hover:text-primary transition-colors duration-200">{category.name}</h3>
                  </motion.div>
                </Link>
              ))}
            </div>
          )}
          {showScrollButton && (
            <button
              onClick={() => scrollCategories('right')}
              className="absolute right-2 top-1/2 -translate-y-1/2 bg-background/80 backdrop-blur-sm border border-border rounded-full p-2 shadow-lg hover:bg-background hover:shadow-xl transition-all z-10"
              aria-label="Scroll right"
            >
              <ChevronRight className="w-4 h-4 text-foreground" />
            </button>
          )}
        </div>
      </div>

      {/* Banner Section */}
      {!searchQuery.trim() && (
        <section className="px-4 pb-6 pt-0">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="relative bg-gradient-to-br from-primary via-primary/95 to-primary/90 rounded-xl overflow-hidden shadow-lg border border-primary/20"
          >
            {/* Decorative pattern overlay */}
            <div className="absolute inset-0 opacity-10">
              <div className="absolute top-0 right-0 w-32 h-32 bg-white rounded-full -mr-16 -mt-16"></div>
              <div className="absolute bottom-0 left-0 w-24 h-24 bg-white rounded-full -ml-12 -mb-12"></div>
            </div>
            
            <div className="relative p-4 md:p-5">
              <div className="mb-3">
                <div className="text-[10px] font-medium text-primary-foreground/80 uppercase tracking-wider mb-1">
                  Welcome to
                </div>
                <h2 className="text-2xl font-bold text-primary-foreground mb-2 leading-tight">
                  Efoyeta
                </h2>
                <div className="w-10 h-0.5 bg-primary-foreground/30 rounded-full mb-2"></div>
              </div>
              
              <p className="text-primary-foreground/90 mb-4 text-xs leading-relaxed">
                We bring a wide selection of items through our trusted partner suppliers, combining convenience, affordability, and reliability.
              </p>
              
              <button
                onClick={() => {
                  const productsSection = document.querySelector('[data-section="products"]');
                  if (productsSection) {
                    productsSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
                  } else {
                    window.scrollBy({ top: 400, behavior: 'smooth' });
                  }
                }}
                className="inline-flex items-center gap-1.5 px-5 py-2 bg-primary-foreground text-primary rounded-lg text-xs font-semibold shadow-md hover:bg-primary-foreground/95 hover:shadow-lg transition-all duration-200 btn-press"
              >
                Shop Now
                <ChevronRight className="w-3.5 h-3.5" />
              </button>
            </div>
          </motion.div>
        </section>
      )}

      {/* Search Results or Products */}
      {searchQuery.trim() ? (
        <section className="px-4 py-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="font-display text-lg font-semibold">
              Search Results
            </h2>
            <p className="text-sm text-muted-foreground">
              {displayedProducts.length} result{displayedProducts.length !== 1 ? 's' : ''}
            </p>
          </div>

          {loading ? (
            <div className="grid grid-cols-2 gap-3">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="bg-card border border-border rounded-md aspect-square animate-pulse" />
              ))}
            </div>
          ) : displayedProducts.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              <p>No products found</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-3">
              {displayedProducts.map((product, index) => (
                <CatalogueCard key={product.name} catalogue={product} index={index} />
              ))}
            </div>
          )}
        </section>
      ) : (
        <>
          {/* Top Products */}
          <section className="px-4 py-6" data-section="products">
            <div className="flex items-center gap-2 mb-4">
              <h2 className="font-display text-lg font-semibold">Top Products</h2>
            </div>

            {/* Category Filter */}
            {categories.length > 0 && (
              <div className="flex items-center gap-2 mb-4 overflow-x-auto pb-2">
                <Filter className="w-4 h-4 shrink-0 text-muted-foreground" />
                <div className="flex gap-2">
                  <button
                    onClick={() => setSelectedCategory(0)}
                    className={`shrink-0 px-4 py-1.5 rounded-full text-xs font-medium transition-all ${
                      selectedCategory === 0
                        ? 'bg-primary text-primary-foreground'
                        : 'bg-muted text-muted-foreground'
                    }`}
                  >
                    All
                  </button>
                  {categories.map((category) => (
                    <button
                      key={category.id}
                      onClick={() => setSelectedCategory(category.id)}
                      className={`shrink-0 px-4 py-1.5 rounded-full text-xs font-medium transition-all whitespace-nowrap ${
                        selectedCategory === category.id
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      {category.name}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {loading ? (
              <div className="grid grid-cols-2 gap-3">
                {[...Array(4)].map((_, i) => (
                  <div key={i} className="bg-card border border-border rounded-md aspect-square animate-pulse" />
                ))}
              </div>
            ) : displayedProducts.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No products found</p>
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-3">
                {displayedProducts.map((product, index) => (
                  <CatalogueCard key={product.name} catalogue={product} index={index} />
                ))}
              </div>
            )}
          </section>

          {/* Top Suppliers */}
          <section className="px-4 py-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-display text-lg font-semibold">Top Suppliers</h2>
              {(suppliers.length > 0 || (!loading && suppliers.length === 0)) && (
                <Link
                  to="/suppliers"
                  className="text-xs text-muted-foreground hover:text-primary font-medium transition-colors"
                >
                  View All
                </Link>
              )}
            </div>

            <div className="relative">
              {loading ? (
                <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide px-4 -mx-4">
                  {[...Array(4)].map((_, i) => (
                    <div key={i} className="bg-card border border-border rounded-lg p-4 w-[280px] h-32 shrink-0 animate-pulse" />
                  ))}
                </div>
              ) : (suppliers.length > 0 ? suppliers : sampleSuppliers).length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <p className="text-sm">No suppliers available</p>
                </div>
              ) : (
                <>
                  <div ref={suppliersScrollRef} className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide px-4 -mx-4">
                    {(suppliers.length > 0 ? suppliers : sampleSuppliers).map((supplier, index) => (
                      <motion.div
                        key={supplier.id}
                        initial={{ opacity: 0, x: 20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.05, duration: 0.3 }}
                        className="bg-card border border-border rounded-lg p-4 hover:border-primary/50 hover:shadow-md transition-all duration-300 w-[280px] shrink-0"
                      >
                        <div className="flex flex-col gap-3 h-full">
                          <div className="flex items-start gap-3">
                            <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center shrink-0">
                              <Building2 className="w-5 h-5 text-primary" />
                            </div>
                            <div className="flex-1 min-w-0">
                              <h3 className="font-semibold text-sm text-foreground mb-1 line-clamp-2">
                                {supplier.name}
                              </h3>
                              <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                <MapPin className="w-3.5 h-3.5 shrink-0" />
                                <span className="truncate">{supplier.region}, {supplier.woreda}</span>
                              </div>
                            </div>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                  {showSuppliersScrollButtons && (
                    <button
                      onClick={() => scrollSuppliers('right')}
                      className="absolute right-2 top-1/2 -translate-y-1/2 bg-background/80 backdrop-blur-sm border border-border rounded-full p-2 shadow-lg hover:bg-background hover:shadow-xl transition-all z-10"
                      aria-label="Scroll right"
                    >
                      <ChevronRight className="w-4 h-4 text-foreground" />
                    </button>
                  )}
                </>
              )}
            </div>
          </section>

          {/* Featured Products */}
          <section className="px-4 py-6">
            <div className="flex items-center gap-2 mb-4">
              <h2 className="font-display text-lg font-semibold">Featured Products</h2>
            </div>

            {/* Category Filter */}
            {categories.length > 0 && (
              <div className="flex items-center gap-2 mb-4 overflow-x-auto pb-2">
                <Filter className="w-4 h-4 shrink-0 text-muted-foreground" />
                <div className="flex gap-2">
                  <button
                    onClick={() => setSelectedCategory(0)}
                    className={`shrink-0 px-4 py-1.5 rounded-full text-xs font-medium transition-all ${
                      selectedCategory === 0
                        ? 'bg-primary text-primary-foreground'
                        : 'bg-muted text-muted-foreground'
                    }`}
                  >
                    All
                  </button>
                  {categories.map((category) => (
                    <button
                      key={category.id}
                      onClick={() => setSelectedCategory(category.id)}
                      className={`shrink-0 px-4 py-1.5 rounded-full text-xs font-medium transition-all whitespace-nowrap ${
                        selectedCategory === category.id
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      {category.name}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {loading ? (
              <div className="grid grid-cols-2 gap-3">
                {[...Array(4)].map((_, i) => (
                  <div key={i} className="bg-card border border-border rounded-md aspect-square animate-pulse" />
                ))}
              </div>
            ) : displayedProducts.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No products found</p>
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-3">
                {displayedProducts.map((product, index) => (
                  <CatalogueCard key={product.name} catalogue={product} index={index} />
                ))}
              </div>
            )}
          </section>
        </>
      )}
    </motion.div>
  );
};

export default Home;
