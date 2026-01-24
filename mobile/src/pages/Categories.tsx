import { motion } from 'framer-motion';
import { useState, useEffect } from 'react';
import { Grid3x3, ArrowRight } from 'lucide-react';
import { GetAllCategories } from '@/lib/api/category';
import { Category } from '@/lib/types';
import { Link } from 'react-router-dom';

const Categories = () => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadCategories = async () => {
      setLoading(true);
      try {
        const categoryList = await GetAllCategories().catch(() => []);
        setCategories(categoryList);
      } catch (error) {
        console.error('Failed to load categories:', error);
      } finally {
        setLoading(false);
      }
    };

    loadCategories();
  }, []);

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Header */}
      <div className="p-4 border-b border-border">
        <h1 className="font-display text-xl font-semibold">Categories</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Browse products by category
        </p>
      </div>

      {/* Categories List */}
      {loading ? (
        <div className="p-4">
          <div className="grid grid-cols-2 gap-3">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="bg-card border border-border rounded-md aspect-[4/3] animate-pulse" />
            ))}
          </div>
        </div>
      ) : categories.length === 0 ? (
        <div className="flex flex-col items-center justify-center text-center py-16 px-4">
          <div className="bg-muted rounded-full p-8 mb-4">
            <Grid3x3 className="w-12 h-12 text-muted-foreground" />
          </div>
          <h3 className="text-lg font-semibold mb-2">No categories available</h3>
          <p className="text-sm text-muted-foreground mb-6 max-w-sm">
            Categories will appear here when available.
          </p>
          <Link
            to="/"
            className="px-6 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium btn-press"
          >
            Go Home
          </Link>
        </div>
      ) : (
        <div className="p-4">
          <div className="grid grid-cols-2 gap-3">
            {categories.map((category, index) => (
              <Link
                key={category.id}
                to={`/search?category=${category.id}`}
                className="group"
              >
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.05 }}
                  className="bg-card border border-border rounded-md p-5 hover:border-primary hover:shadow-md transition-all btn-press h-full flex flex-col items-center justify-center min-h-[120px]"
                >
                  <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center mb-3 group-hover:bg-primary/20 transition-colors">
                    <Grid3x3 className="w-6 h-6 text-primary" />
                  </div>
                  <h3 className="text-sm font-semibold text-center mb-1">{category.name}</h3>
                  <ArrowRight className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors mt-1" />
                </motion.div>
              </Link>
            ))}
          </div>
        </div>
      )}
    </motion.div>
  );
};

export default Categories;

