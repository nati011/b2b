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
import { Eye } from "lucide-react";
import axios from "@/lib/axios";
import Link from "next/link";
import { OrderResponse } from "@/app/actions/orders";
import { useAuth } from "@/context/AuthContext";

const ITEMS_PER_PAGE = 5;

export default function AdminOrdersPage() {
  const { user } = useAuth();
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  
  // Supplier page - always supplier
  const isSupplier = true;

  useEffect(() => {
    fetchOrders(currentPage);
  }, [currentPage, isSupplier]);

  const fetchOrders = async (page: number) => {
    try {
      setLoading(true);
      const offset = (page - 1) * ITEMS_PER_PAGE;
      
      // Use supplier-specific endpoint if user is a supplier, otherwise use customer endpoint
      const endpoint = isSupplier 
        ? `/order/supplier?limit=${ITEMS_PER_PAGE}&offset=${offset}`
        : `/orders/customer?limit=${ITEMS_PER_PAGE}&offset=${offset}`;
      
      const response = await axios.get(endpoint);
      const orderData = response.data;
      
      // Handle response format
      let ordersList: OrderResponse[] = [];
      if (orderData.orders) {
        ordersList = orderData.orders;
        setTotal(orderData.total || ordersList.length);
      } else if (orderData.body?.orders) {
        ordersList = orderData.body.orders;
        setTotal(orderData.body.total || ordersList.length);
      } else if (Array.isArray(orderData)) {
        ordersList = orderData;
        setTotal(orderData.length);
      } else {
        ordersList = [];
        setTotal(0);
      }
      
      setOrders(ordersList);
      setTotalPages(Math.ceil((orderData.total || orderData.body?.total || ordersList.length) / ITEMS_PER_PAGE));
    } catch (error: any) {
      console.error("Error fetching orders:", error);
      setOrders([]);
      setTotal(0);
      setTotalPages(0);
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

  return (
    <div>
      <div className="mb-4 sm:mb-8">
        <div className="flex items-center gap-2 sm:gap-3 mb-2">
          <div className="h-8 sm:h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white">Orders</h1>
        </div>
        <p className="text-sm sm:text-base text-gray-600 dark:text-white mt-2 ml-3 sm:ml-4">
          {isSupplier 
            ? "View and manage orders for your products" 
            : "View and manage all orders"}
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{isSupplier ? "Your Orders" : "All Orders"}</CardTitle>
          <CardDescription>
            {isSupplier 
              ? "Orders containing products from your supplier account" 
              : "A list of all orders in the system"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Desktop Table View */}
          <div className="hidden md:block relative w-full overflow-x-auto">
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
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-gray-500">
                      <div className="flex items-center justify-center py-4">
                        <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
                        <span className="ml-2">Loading orders...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : orders.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-gray-500 py-8">
                      No orders found
                    </TableCell>
                  </TableRow>
                ) : (
                  orders.map((order) => (
                    <TableRow key={order.id} className="hover:bg-muted/50 transition-colors">
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
                <span className="ml-2 text-gray-500">Loading orders...</span>
              </div>
            ) : orders.length === 0 ? (
              <div className="text-center text-gray-500 py-8">
                No orders found
              </div>
            ) : (
              orders.map((order) => (
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
                      <Link href={`/admin/orders/${order.id}`}>
                        <Eye className="h-4 w-4 mr-2" />
                        View Details
                      </Link>
                    </Button>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="mt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
              <div className="text-xs sm:text-sm text-muted-foreground text-center sm:text-left">
                Showing {((currentPage - 1) * ITEMS_PER_PAGE) + 1} to {Math.min(currentPage * ITEMS_PER_PAGE, total)} of {total} orders
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
    </div>
  );
}

