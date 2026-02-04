"use client";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ArrowLeft, Package, User, Calendar, DollarSign, Truck, CreditCard, Eye } from "lucide-react";
import { getOrderById, updateOrderStatus, updateOrderPaymentStatus, OrderResponse } from "@/app/actions/orders";
import { ProductResponse } from "@/app/actions/product";
import Link from "next/link";
import { toast } from "sonner";
import axios from "@/lib/axios";

export default function OrderDetailPage() {
  const params = useParams();
  const router = useRouter();
  const orderId = params?.id as string;
  const [order, setOrder] = useState<OrderResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [products, setProducts] = useState<Record<number, ProductResponse>>({});

  useEffect(() => {
    if (orderId) {
      fetchOrder();
    }
  }, [orderId]);

  const fetchOrder = async () => {
    try {
      setLoading(true);
      setError(null);
      const orderData = await getOrderById(Number(orderId));
      setOrder(orderData);
      
      // Fetch product details for all items in the order
      await fetchProductDetails(orderData);
    } catch (err: any) {
      console.error("Error fetching order:", err);
      setError(err.message || "Failed to load order details");
    } finally {
      setLoading(false);
    }
  };

  const fetchProductDetails = async (orderData: OrderResponse) => {
    const orderItems = getOrderItemsFromOrder(orderData);
    const productIds = orderItems
      .map(item => item.product_id)
      .filter((id): id is number => id !== undefined && id !== 0);
    
    if (productIds.length === 0) return;

    // Fetch all products in parallel
    const productPromises = productIds.map(async (productId) => {
      try {
        const response = await axios.get(`/product?id=${productId}`);
        return { id: productId, product: response.data };
      } catch (err) {
        console.error(`Error fetching product ${productId}:`, err);
        return null;
      }
    });

    const productResults = await Promise.all(productPromises);
    const productsMap: Record<number, ProductResponse> = {};
    
    productResults.forEach((result) => {
      if (result && result.product) {
        productsMap[result.id] = result.product;
      }
    });

    setProducts(productsMap);
  };

  const getOrderItemsFromOrder = (orderData: OrderResponse) => {
    // First try to use items from API response
    if (orderData.items && orderData.items.length > 0) {
      return orderData.items.map(item => ({
        product_id: item.product_id,
        quantity: item.quantity,
        price: item.price,
      }));
    }
    
    // Fallback to cart_snapshot
    try {
      const cartSnapshot = typeof orderData.cart_snapshot === 'string' 
        ? JSON.parse(orderData.cart_snapshot) 
        : orderData.cart_snapshot;
      
      if (Array.isArray(cartSnapshot)) {
        return cartSnapshot.map((item: any) => ({
          product_id: item.product_id || item.ProductId || 0,
          quantity: item.quantity || item.Quantity || 0,
          price: item.price || item.ProductPrice || 0,
        }));
      }
    } catch (e) {
      console.warn('Failed to parse cart_snapshot:', e);
    }
    
    return [];
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-ET', {
      style: 'currency',
      currency: 'ETB',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const getStatusBadge = (status: string) => {
    const statusLower = status?.toLowerCase() || '';
    
    // Color codes for different statuses
    const statusColors: Record<string, string> = {
      pending: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
      cancelled: "bg-red-500/10 text-red-600 border-red-500/20",
      canceled: "bg-red-500/10 text-red-600 border-red-500/20",
      delivered: "bg-green-500/10 text-green-600 border-green-500/20",
      completed: "bg-green-500/10 text-green-600 border-green-500/20",
      processing: "bg-blue-500/10 text-blue-600 border-blue-500/20",
    };
    
    const colorClass = statusColors[statusLower] || "bg-gray-500/10 text-gray-600 border-gray-500/20";
    
    return (
      <Badge variant="outline" className={colorClass}>
        {status?.charAt(0).toUpperCase() + status?.slice(1) || 'Unknown'}
      </Badge>
    );
  };

  const parseCartSnapshot = () => {
    if (!order?.cart_snapshot) return [];
    
    try {
      const cartSnapshot = typeof order.cart_snapshot === 'string' 
        ? JSON.parse(order.cart_snapshot) 
        : order.cart_snapshot;
      
      if (Array.isArray(cartSnapshot)) {
        return cartSnapshot;
      }
      return [];
    } catch (e) {
      console.warn('Failed to parse cart_snapshot:', e);
      return [];
    }
  };

  const parseCustomerSnapshot = () => {
    if (!order?.customer_snapshot) return null;
    
    try {
      const customerSnapshot = typeof order.customer_snapshot === 'string'
        ? JSON.parse(order.customer_snapshot)
        : order.customer_snapshot;
      return customerSnapshot;
    } catch (e) {
      console.warn('Failed to parse customer_snapshot:', e);
      return null;
    }
  };

  const getOrderItems = () => {
    // First try to use items from API response
    if (order?.items && order.items.length > 0) {
      return order.items;
    }
    
    // Fallback to cart_snapshot
    const cartItems = parseCartSnapshot();
    return cartItems.map((item: any) => ({
      product_id: item.product_id || item.ProductId || 0,
      quantity: item.quantity || item.Quantity || 0,
      price: item.price || item.ProductPrice || 0,
      product_name: item.product_name || item.ProductName || 'Unknown Product',
    }));
  };

  const handleMarkAsDelivered = async () => {
    if (!orderId || !order) return;
    
    try {
      setUpdating(true);
      const updatedOrder = await updateOrderStatus(orderId, "delivered");
      setOrder(updatedOrder);
      toast.success("Order marked as delivered");
    } catch (err: any) {
      console.error("Error updating order status:", err);
      toast.error(err.message || "Failed to mark order as delivered");
    } finally {
      setUpdating(false);
    }
  };

  const handleMarkAsPaid = async () => {
    if (!orderId || !order) return;
    
    try {
      setUpdating(true);
      const updatedOrder = await updateOrderPaymentStatus(orderId, "paid");
      setOrder(updatedOrder);
      toast.success("Order marked as paid");
    } catch (err: any) {
      console.error("Error updating payment status:", err);
      toast.error(err.message || "Failed to mark order as paid");
    } finally {
      setUpdating(false);
    }
  };

  const getPaymentStatusBadge = (status: string) => {
    const statusLower = status?.toLowerCase() || '';
    
    const statusColors: Record<string, string> = {
      paid: "bg-green-500/10 text-green-600 border-green-500/20",
      unpaid: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
      pending: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
      failed: "bg-red-500/10 text-red-600 border-red-500/20",
      refunded: "bg-gray-500/10 text-gray-600 border-gray-500/20",
    };
    
    const colorClass = statusColors[statusLower] || "bg-gray-500/10 text-gray-600 border-gray-500/20";
    
    return (
      <Badge variant="outline" className={colorClass}>
        {status?.charAt(0).toUpperCase() + status?.slice(1) || 'Unknown'}
      </Badge>
    );
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
            onClick={() => router.back()}
            className="mb-4"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Orders
          </Button>
          <div className="flex items-center gap-3 mb-2">
            <div className="h-10 w-1 bg-primary rounded-full"></div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Order Details</h1>
          </div>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <p className="text-red-600 mb-4">{error}</p>
              <Button onClick={fetchOrder}>Try Again</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="space-y-6">
        <div className="mb-8">
          <Button
            variant="ghost"
            onClick={() => router.back()}
            className="mb-4"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Orders
          </Button>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8 text-gray-500">
              Order not found
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const customerInfo = parseCustomerSnapshot();
  const orderItems = getOrderItems();
  const subtotal = orderItems.reduce((sum, item) => sum + ((item.price || 0) * (item.quantity || 0)), 0);
  const total = order.total || subtotal;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="mb-8">
        <Button
          variant="ghost"
          asChild
          className="mb-4"
        >
          <Link href="/supplier/orders">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Orders
          </Link>
        </Button>
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Order Details</h1>
        </div>
        <p className="text-gray-600 dark:text-gray-300 mt-2 ml-4">Order #{order.id}</p>
      </div>

      {/* Action Buttons */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-semibold mb-1">Order Actions</h3>
              <p className="text-sm text-muted-foreground">
                Update order status and payment information
              </p>
            </div>
            <div className="flex gap-3">
              <Button
                onClick={handleMarkAsDelivered}
                disabled={updating || order.status?.toLowerCase() === 'delivered'}
                variant={order.status?.toLowerCase() === 'delivered' ? "outline" : "default"}
                className="gap-2"
              >
                <Truck className="h-4 w-4" />
                {order.status?.toLowerCase() === 'delivered' ? 'Already Delivered' : 'Mark as Delivered'}
              </Button>
              <Button
                onClick={handleMarkAsPaid}
                disabled={updating || order.payment_status?.toLowerCase() === 'paid'}
                variant={order.payment_status?.toLowerCase() === 'paid' ? "outline" : "default"}
                className="gap-2"
              >
                <CreditCard className="h-4 w-4" />
                {order.payment_status?.toLowerCase() === 'paid' ? 'Already Paid' : 'Mark as Paid'}
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Order Overview */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Order Status</CardTitle>
            <Package className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-lg font-semibold">
              {getStatusBadge(order.status)}
            </div>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Amount</CardTitle>
            <DollarSign className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary break-words">{formatCurrency(total)}</div>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Payment Status</CardTitle>
            <CreditCard className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-lg font-semibold">
              {getPaymentStatusBadge(order.payment_status || 'unpaid')}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Customer Information */}
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <User className="h-5 w-5 text-primary" />
              <CardTitle>Customer Information</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Customer ID</p>
              <p className="font-semibold">#{order.customer_id}</p>
            </div>
            {customerInfo && (
              <>
                {customerInfo.full_name && (
                  <div>
                    <p className="text-sm text-muted-foreground">Full Name</p>
                    <p className="font-semibold">{customerInfo.full_name}</p>
                  </div>
                )}
                {customerInfo.email && (
                  <div>
                    <p className="text-sm text-muted-foreground">Email</p>
                    <p className="font-semibold break-words">{customerInfo.email}</p>
                  </div>
                )}
                {(customerInfo.phone || customerInfo.phone_number) && (
                  <div>
                    <p className="text-sm text-muted-foreground">Phone Number</p>
                    <p className="font-semibold">{customerInfo.phone || customerInfo.phone_number}</p>
                  </div>
                )}
                {(customerInfo.city || customerInfo.region || customerInfo.woreda) && (
                  <div className="pt-2 border-t">
                    <p className="text-sm text-muted-foreground mb-2">Address</p>
                    <div className="space-y-1">
                      {customerInfo.city && (
                        <p className="font-semibold text-sm">City: <span className="font-normal">{customerInfo.city}</span></p>
                      )}
                      {customerInfo.region && (
                        <p className="font-semibold text-sm">Region: <span className="font-normal">{customerInfo.region}</span></p>
                      )}
                      {customerInfo.woreda && (
                        <p className="font-semibold text-sm">Woreda: <span className="font-normal">{customerInfo.woreda}</span></p>
                      )}
                    </div>
                  </div>
                )}
                {(customerInfo.status || customerInfo.is_active !== undefined) && (
                  <div className="pt-2 border-t">
                    <p className="text-sm text-muted-foreground mb-2">Status</p>
                    <div className="flex items-center gap-2">
                      {customerInfo.status && (
                        <Badge variant="outline" className={
                          customerInfo.status?.toLowerCase() === 'active' 
                            ? "bg-green-500/10 text-green-600 border-green-500/20"
                            : customerInfo.status?.toLowerCase() === 'inactive'
                            ? "bg-gray-500/10 text-gray-600 border-gray-500/20"
                            : "bg-yellow-500/10 text-yellow-600 border-yellow-500/20"
                        }>
                          {customerInfo.status.charAt(0).toUpperCase() + customerInfo.status.slice(1)}
                        </Badge>
                      )}
                      {customerInfo.is_active !== undefined && (
                        <Badge variant="outline" className={
                          customerInfo.is_active
                            ? "bg-green-500/10 text-green-600 border-green-500/20"
                            : "bg-gray-500/10 text-gray-600 border-gray-500/20"
                        }>
                          {customerInfo.is_active ? 'Active' : 'Inactive'}
                        </Badge>
                      )}
                    </div>
                  </div>
                )}
              </>
            )}
            {!customerInfo && (
              <div className="text-sm text-muted-foreground italic">
                Customer information not available in snapshot
              </div>
            )}
          </CardContent>
        </Card>

        {/* Order Information */}
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Calendar className="h-5 w-5 text-primary" />
              <CardTitle>Order Information</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Order ID</p>
              <p className="font-semibold">#{order.id}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Created At</p>
              <p className="font-semibold">
                {new Date(order.created_at).toLocaleString('en-US', {
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
                {new Date(order.updated_at).toLocaleString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </p>
            </div>
            {order.confirmation_status && (
              <div>
                <p className="text-sm text-muted-foreground">Confirmation Status</p>
                <p className="font-semibold">
                  <Badge variant="outline">
                    {order.confirmation_status}
                  </Badge>
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Order Items */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors">
        <CardHeader>
          <div className="flex items-center gap-2">
            <Package className="h-5 w-5 text-primary" />
            <CardTitle>Order Items</CardTitle>
          </div>
          <CardDescription>Items included in this order</CardDescription>
        </CardHeader>
        <CardContent>
          {orderItems.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              No items found in this order
            </div>
          ) : (
            <div className="relative w-full overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Product ID</TableHead>
                    <TableHead>Product Name</TableHead>
                    <TableHead className="text-right">Quantity</TableHead>
                    <TableHead className="text-right">Unit Price</TableHead>
                    <TableHead className="text-right">Total</TableHead>
                    <TableHead className="text-center">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {orderItems.map((item, index) => {
                    const itemTotal = (item.price || 0) * (item.quantity || 0);
                    const productId = item.product_id;
                    return (
                      <TableRow key={index} className="hover:bg-muted/50 transition-colors">
                        <TableCell className="font-medium">
                          {productId ? (
                            <Link 
                              href={`/supplier/products/${productId}`}
                              className="text-primary hover:underline transition-colors"
                            >
                              #{productId}
                            </Link>
                          ) : (
                            'N/A'
                          )}
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            <div className="font-medium">
                              {products[productId]?.name || (item as any).product_name || 'Unknown Product'}
                            </div>
                            {products[productId]?.description && (
                              <div className="text-sm text-muted-foreground line-clamp-2">
                                {products[productId].description}
                              </div>
                            )}
                            {products[productId]?.unit && (
                              <div className="text-xs text-muted-foreground">
                                Unit: {products[productId].unit}
                              </div>
                            )}
                          </div>
                        </TableCell>
                        <TableCell className="text-right">{item.quantity || 0}</TableCell>
                        <TableCell className="text-right">
                          {formatCurrency(item.price || 0)}
                        </TableCell>
                        <TableCell className="text-right font-semibold">
                          {formatCurrency(itemTotal)}
                        </TableCell>
                        <TableCell className="text-center">
                          {productId ? (
                            <Button
                              variant="ghost"
                              size="sm"
                              asChild
                              className="hover:bg-primary/10 hover:text-primary"
                            >
                              <Link href={`/supplier/products/${productId}`}>
                                <Eye className="h-4 w-4 mr-2" />
                                View
                              </Link>
                            </Button>
                          ) : (
                            <span className="text-muted-foreground text-sm">N/A</span>
                          )}
                        </TableCell>
                      </TableRow>
                    );
                  })}
                  <TableRow className="border-t-2 font-semibold">
                    <TableCell colSpan={5} className="text-right">
                      Subtotal:
                    </TableCell>
                    <TableCell className="text-right">
                      {formatCurrency(subtotal)}
                    </TableCell>
                  </TableRow>
                  {order.total && order.total !== subtotal && (
                    <TableRow className="font-semibold">
                      <TableCell colSpan={5} className="text-right">
                        Total:
                      </TableCell>
                      <TableCell className="text-right text-primary">
                        {formatCurrency(total)}
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}


