import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { HeroSection } from '@/components/HeroSection';
import { CatalogueCard } from '@/components/CatalogueCard';
import { GetAllCatalogues } from '@/lib/api/catalogue';
import { GetAllCategories } from '@/lib/api/category';
import { Catalogue, Category } from '@/lib/types';
import { Filter, Grid3x3 } from 'lucide-react';
import { Link } from 'react-router-dom';

const Home = () => {
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number>(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const [productList, categoryList] = await Promise.all([
          GetAllCatalogues().catch(() => []),
          GetAllCategories().catch(() => [])
        ]);
        setProducts(productList);
        setCategories(categoryList);
      } catch (error) {
        console.error('Failed to load data:', error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  const displayedProducts = selectedCategory === 0
    ? products.slice(0, 4)
    : products.filter((product) =>
        product.configurables?.some((configurable) =>
          configurable.categories?.includes(selectedCategory)
        )
      ).slice(0, 4);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      <HeroSection onCategoryChange={setSelectedCategory} />

      {/* Categories Section */}
      {categories.length > 0 && (
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

      {/* Top Products */}
      <section className="px-4 py-6">
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
    </motion.div>
  );
};

export default Home;
