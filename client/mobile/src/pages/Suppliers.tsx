import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Link, useNavigate } from 'react-router-dom';
import { ArrowLeft, Building2, Mail, Phone, Search } from 'lucide-react';
import { GetAllSuppliers } from '@/lib/api/supplier';
import type { Supplier } from '@/lib/types';

const Suppliers = () => {
  const navigate = useNavigate();
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    const fetchSuppliers = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // Fetch all suppliers with a high limit to get all of them
        const response = await GetAllSuppliers({ limit: 1000 });
        const supplierItems = response.items || [];
        
        const allSuppliers = supplierItems
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
        
        setSuppliers(allSuppliers);
      } catch (err: any) {
        console.error('Failed to fetch suppliers:', err);
        setError(err.message || 'Failed to load suppliers');
      } finally {
        setLoading(false);
      }
    };

    fetchSuppliers();
  }, []);

  const filteredSuppliers = suppliers.filter((supplier) =>
    supplier.business_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    supplier.support_email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    supplier.support_phone?.includes(searchQuery)
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-background">
        <div className="text-center space-y-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
          <p className="text-sm text-muted-foreground">Loading suppliers...</p>
        </div>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="page-transition pb-20"
    >
      {/* Header */}
      <div className="sticky top-0 bg-background z-10 border-b border-border">
        <div className="flex items-center p-4">
          <button
            onClick={() => navigate(-1)}
            className="p-2 -ml-2 hover:bg-secondary rounded-sm btn-press"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <h1 className="font-display text-lg ml-2">All Suppliers</h1>
        </div>

        {/* Search Bar */}
        <div className="px-4 pb-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search suppliers..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full h-10 pl-10 pr-4 bg-secondary rounded-lg text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/20"
            />
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="px-4 py-6">
        {error ? (
          <div className="text-center py-12 text-muted-foreground">
            <p className="text-sm">{error}</p>
          </div>
        ) : filteredSuppliers.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground">
            <p className="text-sm">
              {searchQuery ? 'No suppliers found matching your search' : 'No suppliers available'}
            </p>
          </div>
        ) : (
          <>
            <div className="mb-4">
              <p className="text-sm text-muted-foreground">
                {filteredSuppliers.length} {filteredSuppliers.length === 1 ? 'supplier' : 'suppliers'} found
              </p>
            </div>
            <div className="grid grid-cols-1 gap-4">
              {filteredSuppliers.map((supplier, index) => (
                <Link
                  key={supplier.id}
                  to={`/supplier/${supplier.id}`}
                  className="block"
                >
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="bg-card border border-border rounded-lg p-4 hover:border-primary/50 hover:shadow-md transition-all duration-300 cursor-pointer"
                  >
                    <div className="flex items-start gap-4">
                      <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center shrink-0">
                        <Building2 className="w-6 h-6 text-primary" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <h3 className="font-semibold text-base text-foreground mb-2">
                          {supplier.business_name}
                        </h3>
                        <div className="space-y-1.5">
                          {supplier.support_phone && (
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                              <Phone className="w-4 h-4 shrink-0" />
                              <span className="truncate">{supplier.support_phone}</span>
                            </div>
                          )}
                          {supplier.support_email && (
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                              <Mail className="w-4 h-4 shrink-0" />
                              <span className="truncate">{supplier.support_email}</span>
                            </div>
                          )}
                        </div>
                        {supplier.status && (
                          <div className="mt-2">
                            <span className={`inline-block text-xs px-2 py-1 rounded-full ${
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
                  </motion.div>
                </Link>
              ))}
            </div>
          </>
        )}
      </div>
    </motion.div>
  );
};

export default Suppliers;

