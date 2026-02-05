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
import { Input } from "@/components/ui/input";
import { Eye, Search } from "lucide-react";
import axios from "@/lib/axios";
import Link from "next/link";

const ITEMS_PER_PAGE = 5;

interface Customer {
  id: string;
  full_name: string;
  email: string;
  phone_number: string;
  city?: string;
  region?: string;
  woreda?: string;
  status: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export default function AdminCustomerPage() {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    fetchCustomers(currentPage, searchQuery);
  }, [currentPage, searchQuery]);

  const fetchCustomers = async (page: number, search: string = "") => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: ITEMS_PER_PAGE.toString(),
      });
      if (search.trim()) {
        params.append("search", search.trim());
      }
      const response = await axios.get(`/customer?${params.toString()}`);
      const data = response.data;
      
      // Handle different response formats
      let customersList: Customer[] = [];
      let totalCount = 0;
      
      if (data.items && Array.isArray(data.items)) {
        customersList = data.items;
        totalCount = data.total || customersList.length;
        setTotalPages(data.total_pages || Math.ceil(totalCount / ITEMS_PER_PAGE));
      } else if (data.body?.items && Array.isArray(data.body.items)) {
        customersList = data.body.items;
        totalCount = data.body.total || customersList.length;
        setTotalPages(data.body.total_pages || Math.ceil(totalCount / ITEMS_PER_PAGE));
      } else if (Array.isArray(data)) {
        customersList = data;
        totalCount = data.length;
        setTotalPages(Math.ceil(totalCount / ITEMS_PER_PAGE));
      }
      
      setCustomers(customersList);
      setTotal(totalCount);
    } catch (error: any) {
      console.error("Error fetching customers:", error);
      setCustomers([]);
      setTotal(0);
      setTotalPages(0);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string, isActive: boolean) => {
    if (!isActive) {
      return <Badge variant="secondary">Inactive</Badge>;
    }
    const statusLower = status?.toLowerCase() || '';
    const statusColors: Record<string, string> = {
      active: "bg-emerald-500/10 text-emerald-600 border-emerald-500/20",
      inactive: "bg-gray-500/10 text-gray-600 border-gray-500/20",
      suspended: "bg-red-500/10 text-red-600 border-red-500/20",
    };
    const colorClass = statusColors[statusLower] || "bg-gray-500/10 text-gray-600 border-gray-500/20";
    return <Badge variant="outline" className={colorClass}>{status?.charAt(0).toUpperCase() + status?.slice(1) || 'Unknown'}</Badge>;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-4 sm:mb-8">
        <div className="flex items-center gap-2 sm:gap-3 mb-2">
          <div className="h-8 sm:h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white">Customers</h1>
        </div>
        <p className="text-sm sm:text-base text-gray-600 dark:text-white mt-2 ml-3 sm:ml-4">
          View and manage all customers
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Customers</CardTitle>
          <CardDescription>
            A list of all customers in the system
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Search Input */}
          <div className="mb-6">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
              <Input
                type="text"
                placeholder="Search by name, email, or phone number..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setCurrentPage(1); // Reset to first page when searching
                }}
                className="pl-10"
              />
            </div>
          </div>
          {customers.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              No customers found
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>ID</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Email</TableHead>
                    <TableHead>Phone</TableHead>
                    <TableHead>Location</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {customers.map((customer) => (
                    <TableRow key={customer.id}>
                      <TableCell className="font-medium">#{customer.id}</TableCell>
                      <TableCell>{customer.full_name}</TableCell>
                      <TableCell>{customer.email}</TableCell>
                      <TableCell>{customer.phone_number || 'N/A'}</TableCell>
                      <TableCell>
                        {customer.city && customer.region 
                          ? `${customer.city}, ${customer.region}`
                          : customer.city || customer.region || 'N/A'}
                      </TableCell>
                      <TableCell>{getStatusBadge(customer.status, customer.is_active)}</TableCell>
                      <TableCell>
                        {new Date(customer.created_at).toLocaleDateString('en-US', {
                          year: 'numeric',
                          month: 'short',
                          day: 'numeric',
                        })}
                      </TableCell>
                      <TableCell>
                        <Button
                          variant="ghost"
                          size="sm"
                          asChild
                        >
                          <Link href={`/admin/customer/${customer.id}`}>
                            <Eye className="h-4 w-4" />
                          </Link>
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="mt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
              <div className="text-xs sm:text-sm text-muted-foreground text-center sm:text-left">
                Showing {((currentPage - 1) * ITEMS_PER_PAGE) + 1} to {Math.min(currentPage * ITEMS_PER_PAGE, total)} of {total} customers
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

