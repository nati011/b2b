import { useState, useEffect, useMemo, useRef } from 'react';
import { motion } from 'framer-motion';
import { HeroSection } from '@/components/HeroSection';
import { CatalogueCard } from '@/components/CatalogueCard';
import { GetAllCatalogues } from '@/lib/api/catalogue';
import { GetAllCategories } from '@/lib/api/category';
import { GetAllSuppliers } from '@/lib/api/supplier';
import { Catalogue, Category } from '@/lib/types';
import { Filter, Grid3x3, ChevronRight, Building2, MapPin } from 'lucide-react';
import { Link } from 'react-router-dom';


const Home = () => {
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [suppliers, setSuppliers] = useState<any[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number>(0);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const suppliersScrollRef = useRef<HTMLDivElement>(null);
  const [showSuppliersScrollButtons, setShowSuppliersScrollButtons] = useState(false);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        console.log('🔄 Starting to load data...');
        
        // First, fetch with a reasonable limit to get total count
        const initialResponse = await GetAllCatalogues({ limit: 1000, offset: 0 }).catch((error) => {
          console.error('❌ Failed to fetch catalogue:', error);
          return { products: [], total: 0, limit: 1000, offset: 0 };
        });

        // If there are more products than the limit, fetch all pages
        let allProducts = initialResponse.products || [];
        const totalProducts = initialResponse.total || 0;
        const limit = initialResponse.limit || 1000;
        
        if (totalProducts > allProducts.length) {
          console.log(`📦 Fetching all products: ${totalProducts} total, ${allProducts.length} loaded so far`);
          const remainingCount = totalProducts - allProducts.length;
          const remainingPages = Math.ceil(remainingCount / limit);
          
          // Fetch remaining pages
          const remainingPromises = [];
          for (let page = 1; page <= remainingPages; page++) {
            const offset = page * limit;
            remainingPromises.push(
              GetAllCatalogues({ limit, offset }).catch((error) => {
                console.error(`❌ Failed to fetch catalogue page ${page} (offset ${offset}):`, error);
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
          
          console.log(`✅ Loaded ${allProducts.length} out of ${totalProducts} total products`);
        }
        
        const [categoryList, supplierResponse] = await Promise.all([
          GetAllCategories().catch((error) => {
            console.error('❌ Failed to fetch categories:', error);
            return [];
          }),
          GetAllSuppliers({ limit: 100, page: 1 }).catch((error) => {
            console.error('❌ Failed to fetch suppliers:', error);
            console.error('Supplier fetch error details:', {
              message: error.message,
              response: error.response?.data,
              status: error.response?.status,
            });
            return { items: [], total: 0, page: 1, limit: 100, total_pages: 0, has_next: false, has_prev: false };
          })
        ]);
        
        console.log('✅ Received responses:', {
          catalogue: { products: allProducts.length, total: totalProducts },
          categories: categoryList,
          suppliers: supplierResponse
        });
        
        // Convert ProductResponse[] to Catalogue[] format for compatibility
        const catalogueProducts: Catalogue[] = allProducts.map((p: any) => ({
          id: p.id, // Include product ID for navigation (required for product detail page)
          name: p.name,
          desc: p.description || '',
          price: p.price,
          is_active: p.is_active,
          images: [], // Backend doesn't return images yet
          configurable_attributes: p.attributes || {},
          configurables: [] // Backend doesn't return configurables yet
        })).filter((p: Catalogue) => p.id != null && p.id !== undefined) || []; // Filter out products without IDs
        
        console.log('📦 Processed products:', catalogueProducts.length);
        setProducts(catalogueProducts);
        setCategories(categoryList);
        
        // Convert SupplierResponse[] to Supplier[] format for compatibility
        const supplierItems = supplierResponse.items || [];
        console.log('🏢 Raw supplier items from API:', supplierItems.length, supplierResponse);
        
        const activeSuppliers = supplierItems
          .filter((s: any) => s.is_active)
          .map((s: any) => ({
            id: s.id,
            business_name: s.business_name,
            status: s.status,
            support_email: s.support_email,
            support_phone: s.support_phone,
            is_active: s.is_active,
            created_at: s.created_at,
            updated_at: s.updated_at,
          }));
        console.log('🏢 Processed active suppliers:', activeSuppliers.length);
        
        // Only set suppliers if we got data from API, don't fallback to sample
        if (supplierResponse.items && supplierResponse.items.length > 0) {
          setSuppliers(activeSuppliers);
        } else {
          console.warn('⚠️ No suppliers returned from API. Check backend connection and database.');
          setSuppliers([]);
        }
      } catch (error) {
        console.error('❌ Failed to load data:', error);
        // Show user-friendly error message
        if (error instanceof Error) {
          console.error('Error details:', error.message, error.stack);
        }
        setSuppliers([]);
      } finally {
        setLoading(false);
        console.log('✅ Finished loading data');
      }
    };

    loadData();
  }, []);

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
                  Efoyeta Store
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
          
          {/* Suppliers */}
          <section className="px-4 py-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-display text-lg font-semibold">Suppliers</h2>
            </div>

            <div className="relative">
              {loading ? (
                <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide px-4 -mx-4">
                  {[...Array(4)].map((_, i) => (
                    <div key={i} className="bg-card border border-border rounded-lg p-4 w-[280px] h-32 shrink-0 animate-pulse" />
                  ))}
                </div>
              ) : suppliers.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <p className="text-sm">No suppliers available</p>
                </div>
              ) : (
                <>
                  <div ref={suppliersScrollRef} className="flex gap-3 overflow-x-auto pb-2 scrollbar-hide px-4 -mx-4">
                    {suppliers.map((supplier, index) => (
                      <Link
                        key={supplier.id}
                        to={`/supplier/${supplier.id}`}
                        className="block shrink-0"
                      >
                        <motion.div
                          initial={{ opacity: 0, x: 20 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: index * 0.05, duration: 0.3 }}
                          className="bg-card border border-border rounded-lg p-4 hover:border-primary/50 hover:shadow-md transition-all duration-300 w-[280px] cursor-pointer"
                        >
                        <div className="flex flex-col gap-3 h-full">
                          <div className="flex items-start gap-3">
                            <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center shrink-0">
                              <Building2 className="w-5 h-5 text-primary" />
                            </div>
                            <div className="flex-1 min-w-0">
                              <h3 className="font-semibold text-sm text-foreground mb-1 line-clamp-2">
                                {supplier.business_name}
                              </h3>
                              {supplier.support_phone && (
                                <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                  <MapPin className="w-3.5 h-3.5 shrink-0" />
                                  <span className="truncate">{supplier.support_phone}</span>
                                </div>
                              )}
                              {supplier.support_email && !supplier.support_phone && (
                                <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                                  <MapPin className="w-3.5 h-3.5 shrink-0" />
                                  <span className="truncate">{supplier.support_email}</span>
                                </div>
                              )}
                              {supplier.status && (
                                <div className="mt-1">
                                  <span className={`text-[10px] px-2 py-0.5 rounded-full ${
                                    supplier.status === 'active' 
                                      ? 'bg-green-500/10 text-green-600' 
                                      : 'bg-gray-500/10 text-gray-600'
                                  }`}>
                                    {supplier.status}
                                  </span>
                                </div>
                              )}
                            </div>
                          </div>
                        </div>
                        </motion.div>
                      </Link>
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

          {/* Products */}
          <section className="px-4 py-6">
            <div className="flex items-center gap-2 mb-4">
              <h2 className="font-display text-lg font-semibold">Products</h2>
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
