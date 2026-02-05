"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft, User, Mail, Phone, MapPin, Calendar, CheckCircle2, XCircle, Clock } from "lucide-react";
import axios from "@/lib/axios";
import Link from "next/link";

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

export default function CustomerDetailPage() {
  const params = useParams();
  const customerId = params?.id as string;
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (customerId) {
      fetchCustomer();
    }
  }, [customerId]);

  const fetchCustomer = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await axios.get(`/customer/${customerId}`);
      
      // Handle different response formats
      const customerData = response.data?.customer || response.data;
      
      if (!customerData || !customerData.id) {
        throw new Error("Invalid customer data received");
      }
      
      setCustomer({
        id: customerData.id,
        full_name: customerData.full_name || customerData.fullName,
        email: customerData.email,
        phone_number: customerData.phone_number || customerData.phoneNumber,
        city: customerData.city,
        region: customerData.region,
        woreda: customerData.woreda,
        status: customerData.status,
        is_active: customerData.is_active !== undefined ? customerData.is_active : customerData.isActive,
        created_at: customerData.created_at || customerData.createdAt,
        updated_at: customerData.updated_at || customerData.updatedAt,
      });
    } catch (err: any) {
      console.error("Error fetching customer:", err);
      setError(err.response?.data?.message || err.message || "Failed to load customer details");
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string, isActive: boolean) => {
    if (!isActive) {
      return (
        <Badge variant="secondary" className="flex items-center gap-1.5">
          <XCircle className="h-3.5 w-3.5" />
          Inactive
        </Badge>
      );
    }
    const statusLower = status?.toLowerCase() || '';
    const statusColors: Record<string, { bg: string; text: string; border: string; icon: React.ReactNode }> = {
      active: {
        bg: "bg-emerald-500/10",
        text: "text-emerald-600",
        border: "border-emerald-500/20",
        icon: <CheckCircle2 className="h-3.5 w-3.5" />
      },
      inactive: {
        bg: "bg-gray-500/10",
        text: "text-gray-600",
        border: "border-gray-500/20",
        icon: <XCircle className="h-3.5 w-3.5" />
      },
      suspended: {
        bg: "bg-red-500/10",
        text: "text-red-600",
        border: "border-red-500/20",
        icon: <XCircle className="h-3.5 w-3.5" />
      },
      pending: {
        bg: "bg-yellow-500/10",
        text: "text-yellow-600",
        border: "border-yellow-500/20",
        icon: <Clock className="h-3.5 w-3.5" />
      },
    };
    const colorClass = statusColors[statusLower] || statusColors.inactive;
    return (
      <Badge variant="outline" className={`${colorClass.bg} ${colorClass.text} ${colorClass.border} flex items-center gap-1.5`}>
        {colorClass.icon}
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
            asChild
            className="mb-4"
          >
            <Link href="/admin/customer">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Customers
            </Link>
          </Button>
          <div className="flex items-center gap-3 mb-2">
            <div className="h-10 w-1 bg-primary rounded-full"></div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Customer Details</h1>
          </div>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <p className="text-red-600 mb-4">{error}</p>
              <Button onClick={fetchCustomer}>Try Again</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!customer) {
    return (
      <div className="space-y-6">
        <div className="mb-8">
          <Button
            variant="ghost"
            asChild
            className="mb-4"
          >
            <Link href="/admin/customer">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Customers
            </Link>
          </Button>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8 text-gray-500">
              Customer not found
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
          <Link href="/admin/customer">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Customers
          </Link>
        </Button>
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Customer Details</h1>
        </div>
        <p className="text-gray-600 dark:text-gray-400 mt-2 ml-4">Customer #{customer.id}</p>
      </div>

      {/* Customer Overview */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors dark:bg-gray-800 dark:border-gray-700">
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <User className="h-5 w-5 text-primary" />
            Customer Overview
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
            {/* Status */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <CheckCircle2 className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Status</span>
              </div>
              <div className="mt-1">
                {getStatusBadge(customer.status, customer.is_active)}
              </div>
            </div>

            {/* Active Status */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <User className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Active</span>
              </div>
              <div className="mt-1">
                <Badge variant={customer.is_active ? "default" : "secondary"} className="text-sm">
                  {customer.is_active ? "Yes" : "No"}
                </Badge>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

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
              <p className="font-semibold">#{customer.id}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Full Name</p>
              <p className="font-semibold text-lg">{customer.full_name || 'N/A'}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Status</p>
              <div className="mt-1">
                {getStatusBadge(customer.status, customer.is_active)}
              </div>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Active</p>
              <Badge variant={customer.is_active ? "default" : "secondary"} className="mt-1">
                {customer.is_active ? "Active" : "Inactive"}
              </Badge>
            </div>
          </CardContent>
        </Card>

        {/* Contact Information */}
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Mail className="h-5 w-5 text-primary" />
              <CardTitle>Contact Information</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground flex items-center gap-2">
                <Mail className="h-4 w-4" />
                Email
              </p>
              <p className="font-semibold mt-1">
                {customer.email ? (
                  <a href={`mailto:${customer.email}`} className="text-primary hover:underline">
                    {customer.email}
                  </a>
                ) : (
                  <span className="text-muted-foreground">N/A</span>
                )}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground flex items-center gap-2">
                <Phone className="h-4 w-4" />
                Phone Number
              </p>
              <p className="font-semibold mt-1">
                {customer.phone_number ? (
                  <a href={`tel:${customer.phone_number}`} className="text-primary hover:underline">
                    {customer.phone_number}
                  </a>
                ) : (
                  <span className="text-muted-foreground">N/A</span>
                )}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Location Information */}
      {(customer.city || customer.region || customer.woreda) && (
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <MapPin className="h-5 w-5 text-primary" />
              <CardTitle>Location Information</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-3">
              {customer.city && (
                <div>
                  <p className="text-sm text-muted-foreground">City</p>
                  <p className="font-semibold">{customer.city}</p>
                </div>
              )}
              {customer.region && (
                <div>
                  <p className="text-sm text-muted-foreground">Region</p>
                  <p className="font-semibold">{customer.region}</p>
                </div>
              )}
              {customer.woreda && (
                <div>
                  <p className="text-sm text-muted-foreground">Woreda</p>
                  <p className="font-semibold">{customer.woreda}</p>
                </div>
              )}
            </div>
            {customer.city && customer.region && (
              <div className="mt-4 pt-4 border-t">
                <p className="text-sm text-muted-foreground">Full Address</p>
                <p className="font-semibold">
                  {[customer.city, customer.region, customer.woreda].filter(Boolean).join(", ")}
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      )}

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
                {new Date(customer.created_at).toLocaleString('en-US', {
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
                {new Date(customer.updated_at).toLocaleString('en-US', {
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
    </div>
  );
}

