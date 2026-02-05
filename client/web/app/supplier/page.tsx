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
import { Badge } from "@/components/ui/badge";
import { 
  Package, 
  ShoppingCart, 
  DollarSign, 
  TrendingUp,
  Activity,
  ArrowUpRight,
  Eye
} from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { OrderResponse, ListSupplierOrders } from "@/app/actions/orders";
import { ListSupplierProducts } from "@/app/actions/product";

interface DashboardStats {
  totalRevenue: number;
  todayRevenue: number;
  totalOrders: number;
  pendingOrders: number;
  completedOrders: number;
  cancelledOrders: number;
  totalProducts: number;
  activeProducts: number;
  averageOrderValue: number;
}

export default function SupplierDashboard() {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats>({
    totalRevenue: 0,
    todayRevenue: 0,
    totalOrders: 0,
    pendingOrders: 0,
    completedOrders: 0,
    cancelledOrders: 0,
    totalProducts: 0,
    activeProducts: 0,
    averageOrderValue: 0,
  });
  const [recentOrders, setRecentOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        
        // Fetch supplier-specific data
        const [productsRes, ordersRes] = await Promise.allSettled([
          ListSupplierProducts({ limit: 1000, offset: 0 }),
          ListSupplierOrders({ limit: 5, offset: 0 }),
        ]);

        // Process products
        const products = productsRes.status === "fulfilled" 
          ? (productsRes.value.products || [])
          : [];
        
        const activeProducts = products.filter((p: any) => p.is_active !== false).length;

        // Process orders
        let orders: OrderResponse[] = [];
        if (ordersRes.status === "fulfilled") {
          const orderData = ordersRes.value;
          if (orderData.orders) {
            orders = orderData.orders;
          }
        }

        // Calculate revenue and order statistics
        const totalRevenue = orders.reduce((sum, order) => sum + (order.total || 0), 0);
        
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const todayRevenue = orders
          .filter(order => {
            const orderDate = new Date(order.created_at);
            orderDate.setHours(0, 0, 0, 0);
            return orderDate.getTime() === today.getTime();
          })
          .reduce((sum, order) => sum + (order.total || 0), 0);

        const pendingOrders = orders.filter(o => 
          o.status?.toLowerCase() === 'pending' || o.status?.toLowerCase() === 'processing'
        ).length;
        
        const completedOrders = orders.filter(o => 
          o.status?.toLowerCase() === 'completed' || o.status?.toLowerCase() === 'delivered'
        ).length;
        
        const cancelledOrders = orders.filter(o => 
          o.status?.toLowerCase() === 'cancelled'
        ).length;

        const averageOrderValue = orders.length > 0 ? totalRevenue / orders.length : 0;

        // Get recent orders (limit to 5)
        const sortedOrders = [...orders].sort((a, b) => 
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
        const recent = sortedOrders.slice(0, 5);

        setStats({
          totalRevenue,
          todayRevenue,
          totalOrders: orders.length,
          pendingOrders,
          completedOrders,
          cancelledOrders,
          totalProducts: products.length,
          activeProducts,
          averageOrderValue,
        });

        setRecentOrders(recent);
      } catch (error) {
        console.error("Error fetching dashboard stats:", error);
      } finally {
        setLoading(false);
      }
    };

    if (user) {
      fetchDashboardData();
    }
  }, [user]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-ET', {
      style: 'currency',
      currency: 'ETB',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const getStatusBadge = (status: string) => {
    const statusLower = status?.toLowerCase() || '';
    
    const statusColors: Record<string, string> = {
      pending: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
      cancelled: "bg-red-500/10 text-red-600 border-red-500/20",
      canceled: "bg-red-500/10 text-red-600 border-red-500/20",
      delivered: "bg-emerald-500/10 text-emerald-600 border-emerald-500/20",
      completed: "bg-emerald-500/10 text-emerald-600 border-emerald-500/20",
      processing: "bg-blue-500/10 text-blue-600 border-blue-500/20",
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

  return (
    <div className="space-y-4 sm:space-y-6">
      {/* Header */}
      <div className="mb-4 sm:mb-8">
        <div className="flex items-center gap-2 sm:gap-3 mb-2">
          <div className="h-8 sm:h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white">Supplier Dashboard</h1>
        </div>
        <p className="text-sm sm:text-base text-gray-600 dark:text-white mt-2 ml-3 sm:ml-4">Welcome back, {user?.name || user?.email}</p>
      </div>

      {/* Revenue & Key Metrics */}
      <div className="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary break-words">{formatCurrency(stats.totalRevenue)}</div>
            <p className="text-xs text-muted-foreground">All-time sales</p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Today's Revenue</CardTitle>
            <TrendingUp className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary break-words">{formatCurrency(stats.todayRevenue)}</div>
            <p className="text-xs text-muted-foreground">Sales today</p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Orders</CardTitle>
            <ShoppingCart className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.totalOrders}</div>
            <p className="text-xs text-muted-foreground">All orders</p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Order Value</CardTitle>
            <Activity className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary break-words">{formatCurrency(stats.averageOrderValue)}</div>
            <p className="text-xs text-muted-foreground">Per order</p>
          </CardContent>
        </Card>
      </div>

      {/* Order Status & Product Metrics */}
      <div className="grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending Orders</CardTitle>
            <ArrowUpRight className="h-4 w-4 text-yellow-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.pendingOrders}</div>
            <p className="text-xs text-muted-foreground">Awaiting processing</p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Completed Orders</CardTitle>
            <ArrowUpRight className="h-4 w-4 text-emerald-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.completedOrders}</div>
            <p className="text-xs text-muted-foreground">Successfully delivered</p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Products</CardTitle>
            <Package className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.totalProducts}</div>
            <p className="text-xs text-muted-foreground">
              {stats.activeProducts} active
            </p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cancelled Orders</CardTitle>
            <ArrowUpRight className="h-4 w-4 text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.cancelledOrders}</div>
            <p className="text-xs text-muted-foreground">Cancelled</p>
          </CardContent>
        </Card>
      </div>

      {/* Recent Orders */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="h-8 w-1 bg-primary rounded-full"></div>
            <div>
              <CardTitle>Recent Orders</CardTitle>
              <CardDescription>
                Latest orders for your products
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {recentOrders.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              No orders found
            </div>
          ) : (
            <>
              {/* Desktop Table View */}
              <div className="hidden md:block overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Order ID</TableHead>
                      <TableHead>Customer ID</TableHead>
                      <TableHead>Total</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Payment</TableHead>
                      <TableHead>Date</TableHead>
                      <TableHead className="text-right">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {recentOrders.map((order) => (
                      <TableRow key={order.id}>
                        <TableCell className="font-medium">#{order.id}</TableCell>
                        <TableCell>{order.customer_id}</TableCell>
                        <TableCell className="font-semibold">
                          {formatCurrency(order.total || 0)}
                        </TableCell>
                        <TableCell>{getStatusBadge(order.status || '')}</TableCell>
                        <TableCell>
                          <Badge variant="outline">
                            {order.payment_status || 'N/A'}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          {new Date(order.created_at).toLocaleDateString('en-US', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit',
                          })}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="sm"
                            asChild
                            className="hover:bg-primary/10 hover:text-primary"
                          >
                            <Link href={`/supplier/orders/${order.id}`}>
                              <Eye className="h-4 w-4 mr-2" />
                              View
                            </Link>
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              {/* Mobile Card View */}
              <div className="md:hidden space-y-4">
                {recentOrders.map((order) => (
                  <Card key={order.id} className="border border-gray-200">
                    <CardContent className="p-4">
                      <div className="flex items-start justify-between mb-3">
                        <div>
                          <p className="font-semibold text-lg">Order #{order.id}</p>
                          <p className="text-sm text-gray-500">Customer: {order.customer_id}</p>
                        </div>
                        {getStatusBadge(order.status || '')}
                      </div>
                      <div className="space-y-2 mb-4">
                        <div className="flex justify-between">
                          <span className="text-sm text-gray-600">Total:</span>
                          <span className="font-semibold">{formatCurrency(order.total || 0)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-sm text-gray-600">Payment:</span>
                          <Badge variant="outline" className="text-xs">
                            {order.payment_status || 'N/A'}
                          </Badge>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-sm text-gray-600">Date:</span>
                          <span className="text-sm text-gray-900">
                            {new Date(order.created_at).toLocaleDateString('en-US', {
                              year: 'numeric',
                              month: 'short',
                              day: 'numeric',
                              hour: '2-digit',
                              minute: '2-digit',
                            })}
                          </span>
                        </div>
                      </div>
                      <Button
                        variant="outline"
                        size="sm"
                        className="w-full"
                        asChild
                      >
                        <Link href={`/supplier/orders/${order.id}`}>
                          <Eye className="h-4 w-4 mr-2" />
                          View Details
                        </Link>
                      </Button>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

