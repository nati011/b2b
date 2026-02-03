import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { ArrowLeft, Building2, Phone, AlertCircle } from 'lucide-react';
import { GetSupplier } from '@/lib/api/supplier';
import { ListProducts } from '@/lib/api/product';
import { CatalogueCard } from '@/components/CatalogueCard';
import type { Supplier } from '@/lib/types';
import type { ProductResponse } from '@/lib/api/product';
import type { Catalogue } from '@/lib/types';

const SupplierDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const [supplier, setSupplier] = useState<Supplier | null>(null);
  const [products, setProducts] = useState<Catalogue[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      if (!id) {
        setError('Supplier ID is required');
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        setError(null);
        
        const supplierId = Number.parseInt(id, 10);
        if (Number.isNaN(supplierId) || supplierId <= 0) {
          throw new Error(`Invalid supplier ID: ${id}`);
        }
        
        console.log('Fetching supplier:', supplierId);
        
        // Fetch supplier and products in parallel
        const [supplierData, productsData] = await Promise.all([
          GetSupplier(supplierId).catch((err) => {
            console.error('GetSupplier error:', err);
            throw new Error(typeof err === 'string' ? err : err?.message || 'Failed to fetch supplier');
          }),
          ListProducts({ supplier_id: supplierId, is_active: true }).catch((err) => {
            console.error('ListProducts error:', err);
            // Don't fail the whole page if products fail to load
            return { products: [] };
          })
        ]);
        
        console.log('Supplier data received:', supplierData);
        
        if (!supplierData || !supplierData.id) {
          throw new Error('Invalid supplier data received');
        }
        
        // Map supplier data to Supplier type
        const supplier: Supplier = {
          id: supplierData.id,
          business_name: supplierData.business_name,
          status: supplierData.status,
          support_email: supplierData.support_email,
          support_phone: supplierData.support_phone,
          is_active: supplierData.is_active,
          created_at: supplierData.created_at,
          updated_at: supplierData.updated_at,
        };
        
        setSupplier(supplier);
        
        // Convert ProductResponse[] to Catalogue[] format
        const catalogueProducts: Catalogue[] = productsData.products?.map((p: ProductResponse) => ({
          id: p.id,
          name: p.name,
          desc: p.description || '',
          price: p.price || 0,
          is_active: p.is_active,
          images: [],
          configurable_attributes: p.attributes || {},
          configurables: []
        })).filter((p: Catalogue) => p.id != null && p.id !== undefined) || [];
        
        setProducts(catalogueProducts);
      } catch (err: any) {
        console.error('Failed to fetch supplier data:', err);
        const errorMessage = typeof err === 'string' ? err : err?.message || 'Failed to load supplier';
        setError(errorMessage);
        setSupplier(null);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-background">
        <div className="text-center space-y-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
          <p className="text-sm text-muted-foreground">Loading supplier...</p>
        </div>
      </div>
    );
  }

  if (error || !supplier) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-4 bg-background">
        <div className="text-center space-y-4 max-w-sm">
          <AlertCircle className="w-16 h-16 text-destructive mx-auto" />
          <h2 className="text-xl font-semibold text-foreground">Supplier Not Found</h2>
          <p className="text-sm text-muted-foreground">
            {error || 'The supplier you\'re looking for doesn\'t exist or has been removed.'}
          </p>
        </div>
        <button
          onClick={() => navigate(-1)}
          className="mt-6 px-4 py-2 bg-primary text-primary-foreground rounded-lg text-sm font-medium"
        >
          Go Back
        </button>
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
          <h1 className="font-display text-lg ml-2">Supplier Details</h1>
        </div>
      </div>

      {/* Supplier Info */}
      <div className="px-4 pt-6 pb-4">
        <div className="bg-card border border-border rounded-lg p-6 space-y-4">
          <div className="flex items-start gap-4">
            <div className="w-16 h-16 bg-primary/10 rounded-lg flex items-center justify-center shrink-0">
              <Building2 className="w-8 h-8 text-primary" />
            </div>
            <div className="flex-1">
              <h2 className="font-display text-xl font-semibold mb-2">{supplier.business_name}</h2>
              {supplier.status && (
                <span className={`inline-block text-xs px-3 py-1 rounded-full mb-3 ${
                  supplier.status === 'active' 
                    ? 'bg-green-500/10 text-green-600' 
                    : 'bg-gray-500/10 text-gray-600'
                }`}>
                  {supplier.status}
                </span>
              )}
            </div>
          </div>

          {/* Contact Information */}
          <div className="space-y-3 pt-2 border-t border-border">
            {supplier.support_phone && (
              <div className="flex items-center gap-3">
                <Phone className="w-4 h-4 text-muted-foreground" />
                <span className="text-sm">{supplier.support_phone}</span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Products Section */}
      <section className="px-4 py-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-display text-lg font-semibold">Products</h2>
          <span className="text-xs text-muted-foreground">
            {products.length} {products.length === 1 ? 'product' : 'products'}
          </span>
        </div>

        {products.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground">
            <p className="text-sm">No products available from this supplier</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-3">
            {products.map((product, index) => (
              <CatalogueCard key={product.id || product.name} catalogue={product} index={index} />
            ))}
          </div>
        )}
      </section>
    </motion.div>
  );
};

export default SupplierDetail;

