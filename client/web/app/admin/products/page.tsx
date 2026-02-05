"use client";
import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationEllipsis,
} from "@/components/ui/pagination";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Eye, Edit, Plus, MoreVertical, FileText, DollarSign } from "lucide-react";
import { ListProducts, ListSupplierProducts, ProductResponse, ProductListResponse } from "@/app/actions/product";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { toast } from "sonner";
import axios from "@/lib/axios";

const ITEMS_PER_PAGE = 10;

export default function AdminProductsPage() {
  const { user } = useAuth();
  const [products, setProducts] = useState<ProductResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [grnDialogOpen, setGrnDialogOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<ProductResponse | null>(null);
  const [grnFormData, setGrnFormData] = useState({
    quantity: "",
    notes: "",
  });
  const [submittingGrn, setSubmittingGrn] = useState(false);
  const [priceDialogOpen, setPriceDialogOpen] = useState(false);
  const [priceFormData, setPriceFormData] = useState({
    new_price: "",
    reason: "",
  });
  const [submittingPrice, setSubmittingPrice] = useState(false);
  
  // Check if user is a supplier
  const isSupplier = user?.roles?.includes("supplier") || false;

  useEffect(() => {
    fetchProducts(currentPage);
  }, [currentPage]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-ET', {
      style: 'currency',
      currency: 'ETB',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const fetchProducts = async (page: number) => {
    try {
      setLoading(true);
      const offset = (page - 1) * ITEMS_PER_PAGE;
      
      let response: ProductListResponse;
      
      if (isSupplier) {
        // For suppliers, use the dedicated supplier products endpoint
        // This endpoint requires authentication and automatically filters by supplier
        response = await ListSupplierProducts({ 
          limit: ITEMS_PER_PAGE, 
          offset: offset 
        });
      } else {
        // For admins, use the regular products endpoint (returns all products)
        response = await ListProducts({ 
          limit: ITEMS_PER_PAGE, 
          offset: offset 
        });
      }
      
      setProducts(response.products || []);
      setTotal(response.total || 0);
      setTotalPages(Math.ceil((response.total || 0) / ITEMS_PER_PAGE));
    } catch (error: any) {
      console.error("Error fetching products:", error);
      setProducts([]);
      setTotal(0);
      setTotalPages(0);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="mb-4 sm:mb-8">
        <div className="flex items-center gap-2 sm:gap-3 mb-2">
          <div className="h-8 sm:h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white">Products</h1>
        </div>
        <p className="text-sm sm:text-base text-gray-600 dark:text-white mt-2 ml-3 sm:ml-4">
          {isSupplier 
            ? "View and manage your products" 
            : "View and manage all products"}
        </p>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <CardTitle>{isSupplier ? "Your Products" : "All Products"}</CardTitle>
              <CardDescription>
                {isSupplier 
                  ? "Products from your supplier account" 
                  : "A list of all products in the system"}
              </CardDescription>
            </div>
            <Button asChild className="w-full sm:w-auto">
              <Link href="/admin/products/new">
                <Plus className="h-4 w-4 mr-2" />
                New Product
              </Link>
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {/* Desktop Table View */}
          <div className="hidden md:block relative w-full overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Price</TableHead>
                  <TableHead>Stock</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-gray-500">
                      <div className="flex items-center justify-center py-4">
                        <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
                        <span className="ml-2">Loading products...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : products.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-gray-500 py-8">
                      No products found
                    </TableCell>
                  </TableRow>
                ) : (
                  products.map((product) => (
                    <TableRow key={product.id} className="hover:bg-muted/50 transition-colors">
                      <TableCell className="font-medium">{product.name}</TableCell>
                      <TableCell className="max-w-md truncate">
                        {product.description || "-"}
                      </TableCell>
                      <TableCell>{formatCurrency(product.price || 0)}</TableCell>
                      <TableCell>{product.available_quantity || product.total_quantity || 0}</TableCell>
                      <TableCell>
                        <Badge
                          variant="outline"
                          className={product.is_active 
                            ? "bg-primary/10 text-primary-700 border-primary-700" 
                            : "bg-gray-500/10 text-gray-700 border-gray-700"
                          }
                        >
                          {product.is_active ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-8 w-8 p-0"
                              >
                                <MoreVertical className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem asChild>
                                <Link href={`/admin/products/${product.id}`} className="flex items-center">
                                  <Eye className="h-4 w-4 mr-2" />
                                  View
                                </Link>
                              </DropdownMenuItem>
                              <DropdownMenuItem asChild>
                                <Link href={`/admin/products/${product.id}/edit`} className="flex items-center">
                                  <Edit className="h-4 w-4 mr-2" />
                                  Edit
                                </Link>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => {
                                  setSelectedProduct(product);
                                  setPriceFormData({ new_price: (product.price || 0).toString(), reason: "" });
                                  setPriceDialogOpen(true);
                                }}
                              >
                                <DollarSign className="h-4 w-4 mr-2" />
                                Update Price
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => {
                                  setSelectedProduct(product);
                                  setGrnFormData({ quantity: "", notes: "" });
                                  setGrnDialogOpen(true);
                                }}
                              >
                                <FileText className="h-4 w-4 mr-2" />
                                GRN
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Mobile Card View */}
          <div className="md:hidden space-y-4">
            {loading ? (
              <div className="flex items-center justify-center py-8">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
                <span className="ml-2 text-gray-500">Loading products...</span>
              </div>
            ) : products.length === 0 ? (
              <div className="text-center text-gray-500 py-8">
                No products found
              </div>
            ) : (
              products.map((product) => (
                <Card key={product.id} className="border border-gray-200">
                  <CardContent className="p-4">
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex-1">
                        <h3 className="font-semibold text-lg mb-1">{product.name}</h3>
                        <p className="text-sm text-gray-600 line-clamp-2">
                          {product.description || "No description"}
                        </p>
                      </div>
                      <Badge
                        variant="outline"
                        className={`ml-2 shrink-0 ${
                          product.is_active 
                            ? "bg-primary/10 text-primary-700 border-primary-700" 
                            : "bg-gray-500/10 text-gray-700 border-gray-700"
                        }`}
                      >
                        {product.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </div>
                    <div className="space-y-2 mb-4">
                      <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Price:</span>
                        <span className="font-semibold">{formatCurrency(product.price || 0)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-sm text-gray-600">Stock:</span>
                        <span className="text-sm font-medium">
                          {product.available_quantity || product.total_quantity || 0} units
                        </span>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button
                            variant="outline"
                            size="sm"
                            className="flex-1"
                          >
                            <MoreVertical className="h-4 w-4 mr-2" />
                            Actions
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem asChild>
                            <Link href={`/admin/products/${product.id}`} className="flex items-center">
                              <Eye className="h-4 w-4 mr-2" />
                              View
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem asChild>
                            <Link href={`/admin/products/${product.id}/edit`} className="flex items-center">
                              <Edit className="h-4 w-4 mr-2" />
                              Edit
                            </Link>
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => {
                              setSelectedProduct(product);
                              setPriceFormData({ new_price: (product.price || 0).toString(), reason: "" });
                              setPriceDialogOpen(true);
                            }}
                          >
                            <DollarSign className="h-4 w-4 mr-2" />
                            Update Price
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => {
                              setSelectedProduct(product);
                              setGrnFormData({ quantity: "", notes: "" });
                              setGrnDialogOpen(true);
                            }}
                          >
                            <FileText className="h-4 w-4 mr-2" />
                            GRN
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="mt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
              <div className="text-xs sm:text-sm text-muted-foreground text-center sm:text-left">
                Showing {((currentPage - 1) * ITEMS_PER_PAGE) + 1} to {Math.min(currentPage * ITEMS_PER_PAGE, total)} of {total} products
              </div>
              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious 
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        if (currentPage > 1) {
                          setCurrentPage(currentPage - 1);
                        }
                      }}
                      className={currentPage === 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
                    />
                  </PaginationItem>
                  
                  {/* Page numbers */}
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => {
                    // Show first page, last page, current page, and pages around current
                    if (
                      page === 1 ||
                      page === totalPages ||
                      (page >= currentPage - 1 && page <= currentPage + 1)
                    ) {
                      return (
                        <PaginationItem key={page}>
                          <PaginationLink
                            href="#"
                            onClick={(e) => {
                              e.preventDefault();
                              setCurrentPage(page);
                            }}
                            isActive={currentPage === page}
                            className="cursor-pointer"
                          >
                            {page}
                          </PaginationLink>
                        </PaginationItem>
                      );
                    } else if (
                      page === currentPage - 2 ||
                      page === currentPage + 2
                    ) {
                      return (
                        <PaginationItem key={page}>
                          <PaginationEllipsis />
                        </PaginationItem>
                      );
                    }
                    return null;
                  })}
                  
                  <PaginationItem>
                    <PaginationNext 
                      href="#"
                      onClick={(e) => {
                        e.preventDefault();
                        if (currentPage < totalPages) {
                          setCurrentPage(currentPage + 1);
                        }
                      }}
                      className={currentPage === totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </CardContent>
      </Card>

      {/* GRN Dialog */}
      <Dialog open={grnDialogOpen} onOpenChange={setGrnDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Goods Receiving Note (GRN)</DialogTitle>
            <DialogDescription>
              Record goods received for {selectedProduct?.name || "product"}
            </DialogDescription>
          </DialogHeader>
          <form
            onSubmit={async (e) => {
              e.preventDefault();
              if (!selectedProduct) return;

              try {
                setSubmittingGrn(true);
                const response = await axios.post("/product/grn", {
                  product_id: selectedProduct.id,
                  quantity: parseInt(grnFormData.quantity),
                  notes: grnFormData.notes || undefined,
                });
                toast.success("Goods receiving note created successfully");
                setGrnDialogOpen(false);
                setGrnFormData({ quantity: "", notes: "" });
                setSelectedProduct(null);
                // Refresh products list
                fetchProducts(currentPage);
              } catch (error: any) {
                console.error("Error creating GRN:", error);
                toast.error(
                  error.response?.data?.message ||
                    error.message ||
                    "Failed to create goods receiving note"
                );
              } finally {
                setSubmittingGrn(false);
              }
            }}
          >
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="grn-quantity">Quantity Received *</Label>
                <Input
                  id="grn-quantity"
                  type="number"
                  min="1"
                  required
                  value={grnFormData.quantity}
                  onChange={(e) =>
                    setGrnFormData({ ...grnFormData, quantity: e.target.value })
                  }
                  placeholder="Enter quantity"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="grn-notes">Notes (Optional)</Label>
                <textarea
                  id="grn-notes"
                  className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  value={grnFormData.notes}
                  onChange={(e) =>
                    setGrnFormData({ ...grnFormData, notes: e.target.value })
                  }
                  placeholder="Additional notes about the goods received"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setGrnDialogOpen(false);
                  setGrnFormData({ quantity: "", notes: "" });
                  setSelectedProduct(null);
                }}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submittingGrn}>
                {submittingGrn ? "Creating..." : "Create GRN"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Update Price Dialog */}
      <Dialog open={priceDialogOpen} onOpenChange={setPriceDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Update Product Price</DialogTitle>
            <DialogDescription>
              Update the price for {selectedProduct?.name || "product"}
            </DialogDescription>
          </DialogHeader>
          <form
            onSubmit={async (e) => {
              e.preventDefault();
              if (!selectedProduct) return;

              try {
                setSubmittingPrice(true);
                const response = await axios.patch("/product/price", {
                  product_id: selectedProduct.id,
                  new_price: parseFloat(priceFormData.new_price),
                  reason: priceFormData.reason || undefined,
                });
                toast.success("Product price updated successfully");
                setPriceDialogOpen(false);
                setPriceFormData({ new_price: "", reason: "" });
                setSelectedProduct(null);
                // Refresh products list
                fetchProducts(currentPage);
              } catch (error: any) {
                console.error("Error updating price:", error);
                toast.error(
                  error.response?.data?.message ||
                    error.message ||
                    "Failed to update product price"
                );
              } finally {
                setSubmittingPrice(false);
              }
            }}
          >
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="current-price">Current Price (ETB)</Label>
                <Input
                  id="current-price"
                  type="number"
                  step="0.01"
                  value={selectedProduct?.price || 0}
                  disabled
                  className="bg-muted cursor-not-allowed"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="new-price">New Price (ETB) *</Label>
                <Input
                  id="new-price"
                  type="number"
                  step="0.01"
                  min="0"
                  required
                  value={priceFormData.new_price}
                  onChange={(e) =>
                    setPriceFormData({ ...priceFormData, new_price: e.target.value })
                  }
                  placeholder="Enter new price"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="price-reason">Reason (Optional)</Label>
                <textarea
                  id="price-reason"
                  className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  value={priceFormData.reason}
                  onChange={(e) =>
                    setPriceFormData({ ...priceFormData, reason: e.target.value })
                  }
                  placeholder="Reason for price change (e.g., market adjustment, cost change)"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setPriceDialogOpen(false);
                  setPriceFormData({ new_price: "", reason: "" });
                  setSelectedProduct(null);
                }}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={submittingPrice}>
                {submittingPrice ? "Updating..." : "Update Price"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}

