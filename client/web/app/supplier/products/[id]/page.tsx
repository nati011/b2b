"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Package, Edit, DollarSign, Box, Calendar, Hash } from "lucide-react";
import { GetProduct, ProductResponse } from "@/app/actions/product";
import Link from "next/link";

export default function ProductDetailPage() {
  const params = useParams();
  const productId = params?.id as string;
  const [product, setProduct] = useState<ProductResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (productId) {
      fetchProduct();
    }
  }, [productId]);

  const fetchProduct = async () => {
    try {
      setLoading(true);
      setError(null);
      const productData = await GetProduct(Number(productId));
      setProduct(productData);
    } catch (err: any) {
      console.error("Error fetching product:", err);
      setError(err.message || "Failed to load product details");
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-ET', {
      style: 'currency',
      currency: 'ETB',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="mb-8">
          <Button
            variant="ghost"
            asChild
            className="mb-4"
          >
            <Link href="/supplier/products">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Products
            </Link>
          </Button>
          <div className="flex items-center gap-3 mb-2">
            <div className="h-10 w-1 bg-primary rounded-full"></div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Product Details</h1>
          </div>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <p className="text-red-600 mb-4">{error}</p>
              <Button onClick={fetchProduct}>Try Again</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="space-y-6">
        <div className="mb-8">
          <Button
            variant="ghost"
            asChild
            className="mb-4"
          >
            <Link href="/supplier/products">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Products
            </Link>
          </Button>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8 text-gray-500">
              Product not found
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="mb-8">
        <Button
          variant="ghost"
          asChild
          className="mb-4"
        >
          <Link href="/supplier/products">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Products
          </Link>
        </Button>
        <div className="flex items-center justify-between">
          <div>
            <div className="flex items-center gap-3 mb-2">
              <div className="h-10 w-1 bg-primary rounded-full"></div>
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Product Details</h1>
            </div>
            <p className="text-gray-600 dark:text-gray-300 mt-2 ml-4">Product #{product.id}</p>
          </div>
          <Button
            variant="default"
            asChild
            className="gap-2"
          >
            <Link href={`/supplier/products/${product.id}/edit`}>
              <Edit className="h-4 w-4" />
              Edit Product
            </Link>
          </Button>
        </div>
      </div>

      {/* Product Overview */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors dark:bg-gray-800 dark:border-gray-700">
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <Package className="h-5 w-5 text-primary" />
            Product Overview
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
            {/* Status */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <Package className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Status</span>
              </div>
              <div className="mt-1">
                <Badge 
                  variant="outline" 
                  className={`text-sm ${
                    product.is_active 
                      ? "bg-primary/10 text-primary-700 border-primary-700" 
                      : "bg-gray-500/10 text-gray-700 border-gray-700"
                  }`}
                >
                  {product.is_active ? "Active" : "Inactive"}
                </Badge>
              </div>
            </div>

            {/* Price */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <DollarSign className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Price</span>
              </div>
              <div className="mt-1">
                <div className="text-xl md:text-2xl font-bold text-primary break-words">
                  {product.price ? formatCurrency(product.price) : 'N/A'}
                </div>
              </div>
            </div>

            {/* Available Stock */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <Box className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Available</span>
              </div>
              <div className="mt-1">
                <div className="text-xl md:text-2xl font-bold text-primary">
                  {product.available_quantity || 0}
                </div>
                {product.unit && (
                  <p className="text-xs text-muted-foreground mt-0.5">{product.unit}</p>
                )}
              </div>
            </div>

            {/* Total Stock */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <Box className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Total</span>
              </div>
              <div className="mt-1">
                <div className="text-xl md:text-2xl font-bold text-primary">
                  {product.total_quantity || 0}
                </div>
                {product.unit && (
                  <p className="text-xs text-muted-foreground mt-0.5">{product.unit}</p>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Product Information */}
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Package className="h-5 w-5 text-primary" />
              <CardTitle>Product Information</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Product ID</p>
              <p className="font-semibold">#{product.id}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Product Name</p>
              <p className="font-semibold text-lg">{product.name}</p>
            </div>
            {product.external_id && (
              <div>
                <p className="text-sm text-muted-foreground">External ID</p>
                <p className="font-semibold">{product.external_id}</p>
              </div>
            )}
            {product.description && (
              <div>
                <p className="text-sm text-muted-foreground mb-2">Description</p>
                <p className="text-sm leading-relaxed">{product.description}</p>
              </div>
            )}
            {product.unit && (
              <div>
                <p className="text-sm text-muted-foreground">Unit</p>
                <p className="font-semibold">{product.unit}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Stock & Supplier Information */}
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Box className="h-5 w-5 text-primary" />
              <CardTitle>Stock & Supplier</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Supplier ID</p>
              <p className="font-semibold">#{product.supplier_id}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Total Quantity</p>
              <p className="font-semibold text-lg">{product.total_quantity || 0}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Reserved Quantity</p>
              <p className="font-semibold">{product.reserved_quantity || 0}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Available Quantity</p>
              <p className="font-semibold text-lg text-primary">
                {product.available_quantity || 0}
              </p>
            </div>
            {product.category_ids && product.category_ids.length > 0 && (
              <div>
                <p className="text-sm text-muted-foreground mb-2">Categories</p>
                <div className="flex flex-wrap gap-2">
                  {product.category_ids.map((catId) => (
                    <Badge key={catId} variant="outline">
                      #{catId}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Timestamps */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors">
        <CardHeader>
          <div className="flex items-center gap-2">
            <Calendar className="h-5 w-5 text-primary" />
            <CardTitle>Timestamps</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <p className="text-sm text-muted-foreground">Created At</p>
              <p className="font-semibold">
                {new Date(product.created_at).toLocaleString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Last Updated</p>
              <p className="font-semibold">
                {new Date(product.updated_at).toLocaleString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Attributes (if available) */}
      {product.attributes && Object.keys(product.attributes).length > 0 && (
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Hash className="h-5 w-5 text-primary" />
              <CardTitle>Product Attributes</CardTitle>
            </div>
            <CardDescription>Additional product information</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {Object.entries(
                typeof product.attributes === 'string' 
                  ? JSON.parse(product.attributes) 
                  : product.attributes
              ).map(([key, value]) => (
                <div key={key} className="flex justify-between items-start py-2 border-b last:border-0">
                  <p className="text-sm font-medium text-muted-foreground capitalize">
                    {key.replace(/_/g, ' ')}
                  </p>
                  <p className="text-sm font-semibold text-right max-w-md">
                    {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                  </p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}



