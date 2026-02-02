"use client";
import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { 
  Package, 
  ShoppingCart, 
  DollarSign, 
  TrendingUp,
  Users,
  Activity,
  ArrowUpRight,
  Building2
} from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import axios from "@/lib/axios";
import { OrderResponse } from "@/app/actions/orders";
import { ListProducts } from "@/app/actions/product";

interface DashboardStats {
  totalRevenue: number;
  todayRevenue: number;
  totalOrders: number;
  pendingOrders: number;
  completedOrders: number;
  cancelledOrders: number;
  totalProducts: number;
  activeProducts: number;
  totalCustomers: number;
  totalSuppliers: number;
  activeSuppliers: number;
  averageOrderValue: number;
}

export default function AdminDashboard() {
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
    totalCustomers: 0,
    totalSuppliers: 0,
    activeSuppliers: 0,
    averageOrderValue: 0,
  });
  const [loading, setLoading] = useState(true);
  
  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        
        // Admin dashboard - fetch all orders
        const ordersPromise = axios.get("/orders/customer?limit=1000&offset=0").catch(() => 
          axios.get("/api/v1/orders").catch(() => ({ data: { body: [], orders: [] } }))
        );
        
        // Fetch all data in parallel
        const [productsRes, ordersRes, customersRes, suppliersRes] = await Promise.allSettled([
          ListProducts({ limit: 1000, offset: 0 }),
          ordersPromise,
          axios.get("/customer?page=1&limit=1000").catch(() => ({ data: { body: [], items: [] } })),
          axios.get("/supplier?page=1&limit=1000").catch(() => ({ data: { body: [], items: [] } })),
        ]);

        // Process products
        const products = productsRes.status === "fulfilled" 
          ? (productsRes.value.products || [])
          : [];
        
        const activeProducts = products.filter((p: any) => p.is_active !== false).length;

        // Process orders
        let orders: OrderResponse[] = [];
        if (ordersRes.status === "fulfilled") {
          const orderData = ordersRes.value.data;
          if (orderData.orders) {
            orders = orderData.orders;
          } else if (orderData.body) {
            orders = Array.isArray(orderData.body) ? orderData.body : [];
          } else if (Array.isArray(orderData)) {
            orders = orderData;
          }
        }

        // Process customers
        const customerData = customersRes.status === "fulfilled" 
          ? (customersRes.value.data.body || customersRes.value.data)
          : {};
        const customers = customerData.items || customerData.customers || [];
        const totalCustomers = customerData.total || customers.length;

        // Process suppliers
        const supplierData = suppliersRes.status === "fulfilled" 
          ? (suppliersRes.value.data.body || suppliersRes.value.data)
          : {};
        const suppliers = supplierData.items || supplierData.suppliers || [];
        const totalSuppliers = supplierData.total || suppliers.length;
        const activeSuppliers = suppliers.filter((s: any) => s.is_active !== false).length;

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

        setStats({
          totalRevenue,
          todayRevenue,
          totalOrders: orders.length,
          pendingOrders,
          completedOrders,
          cancelledOrders,
          totalProducts: products.length,
          activeProducts,
          totalCustomers: totalCustomers,
          totalSuppliers: totalSuppliers,
          activeSuppliers: activeSuppliers,
          averageOrderValue,
        });
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
          <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white">Admin Dashboard</h1>
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
            <ArrowUpRight className="h-4 w-4 text-green-600" />
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
            <CardTitle className="text-sm font-medium">Total Customers</CardTitle>
            <Users className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.totalCustomers}</div>
            <p className="text-xs text-muted-foreground">All registered customers</p>
          </CardContent>
        </Card>

        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Suppliers</CardTitle>
            <Building2 className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">{stats.totalSuppliers}</div>
            <p className="text-xs text-muted-foreground">
              {stats.activeSuppliers} active
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

