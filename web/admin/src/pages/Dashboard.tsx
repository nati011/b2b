
import { Card } from "@/components/ui/card";
import {
  Users,
  ShoppingCart,
  DollarSign,
  TrendingUp,
  ChevronRight,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { mockOrders, mockProducts, mockCustomers } from "@/lib/mock-data";

const statsData = [
  {
    label: "Total Customers",
    value: "1,234",
    trend: "+12.5%",
    icon: Users,
    color: "text-blue-500",
  },
  {
    label: "Total Orders",
    value: "842",
    trend: "+23.1%",
    icon: ShoppingCart,
    color: "text-green-500",
  },
  {
    label: "Total Revenue",
    value: "$45,678",
    trend: "+18.7%",
    icon: DollarSign,
    color: "text-purple-500",
  },
  {
    label: "Growth Rate",
    value: "24.5%",
    trend: "+5.2%",
    icon: TrendingUp,
    color: "text-orange-500",
  },
];

const popularProducts = [
  { id: 7, name: "Coffee Maker", sales: 230, price: "$89.99", category: "Appliances" },
  { id: 12, name: "LED Monitor", sales: 215, price: "$299.99", category: "Electronics" },
  { id: 3, name: "Super Tool", sales: 198, price: "$79.99", category: "Tools" },
];

const loyalCustomers = [
  { id: 8, name: "Maria Garcia", totalSpent: "$8,900", orders: 89, lastPurchase: "2024-02-13" },
  { id: 11, name: "Daniel Martinez", totalSpent: "$9,100", orders: 91, lastPurchase: "2024-02-10" },
  { id: 15, name: "Richard Clark", totalSpent: "$8,700", orders: 87, lastPurchase: "2024-02-06" },
];

const Dashboard = () => {
  return (
    <div className="container animate-fadeIn">
      <h1 className="text-3xl font-semibold mb-8">Dashboard</h1>
      
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {statsData.map((stat, index) => (
          <Card key={index} className="p-6 glass">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-muted-foreground">{stat.label}</p>
                <h3 className="text-2xl font-semibold mt-2">{stat.value}</h3>
                <p className="text-sm text-green-500 mt-1">{stat.trend}</p>
              </div>
              <div className={cn("p-3 rounded-full glass", stat.color)}>
                <stat.icon className="h-6 w-6" />
              </div>
            </div>
          </Card>
        ))}
      </div>

      {/* Tables Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        {/* Latest Orders */}
        <Card className="p-6 glass">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Latest Orders</h2>
            <button className="text-sm text-primary flex items-center hover:underline">
              View All <ChevronRight className="h-4 w-4 ml-1" />
            </button>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left border-b">
                  <th className="pb-3">Order ID</th>
                  <th className="pb-3">Customer</th>
                  <th className="pb-3">Status</th>
                </tr>
              </thead>
              <tbody>
                {mockOrders.slice(0, 3).map((order) => (
                  <tr key={order.id} className="border-b last:border-0">
                    <td className="py-3">#{order.id}</td>
                    <td className="py-3">{order.customer}</td>
                    <td className="py-3">
                      <span className={cn(
                        "px-2 py-1 rounded-full text-xs",
                        {
                          "bg-green-100 text-green-800": order.status === "completed",
                          "bg-yellow-100 text-yellow-800": order.status === "pending",
                          "bg-blue-100 text-blue-800": order.status === "processing",
                        }
                      )}>
                        {order.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>

        {/* Popular Products */}
        <Card className="p-6 glass">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Popular Products</h2>
            <button className="text-sm text-primary flex items-center hover:underline">
              View All <ChevronRight className="h-4 w-4 ml-1" />
            </button>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left border-b">
                  <th className="pb-3">Product</th>
                  <th className="pb-3">Sales</th>
                  <th className="pb-3">Price</th>
                </tr>
              </thead>
              <tbody>
                {popularProducts.map((product) => (
                  <tr key={product.id} className="border-b last:border-0">
                    <td className="py-3">{product.name}</td>
                    <td className="py-3">{product.sales}</td>
                    <td className="py-3">{product.price}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>

        {/* Loyal Customers */}
        <Card className="p-6 glass">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Loyal Customers</h2>
            <button className="text-sm text-primary flex items-center hover:underline">
              View All <ChevronRight className="h-4 w-4 ml-1" />
            </button>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left border-b">
                  <th className="pb-3">Customer</th>
                  <th className="pb-3">Orders</th>
                  <th className="pb-3">Total Spent</th>
                </tr>
              </thead>
              <tbody>
                {loyalCustomers.map((customer) => (
                  <tr key={customer.id} className="border-b last:border-0">
                    <td className="py-3">{customer.name}</td>
                    <td className="py-3">{customer.orders}</td>
                    <td className="py-3">{customer.totalSpent}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Dashboard;
