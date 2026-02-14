"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Building2, Mail, Phone, Calendar, CheckCircle2, XCircle, Clock } from "lucide-react";
import axios from "@/lib/axios";
import Link from "next/link";

interface Supplier {
  id: number;
  business_name: string;
  status: string;
  support_email?: string;
  support_phone?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export default function SupplierDetailPage() {
  const params = useParams();
  const supplierId = params?.id as string;
  const [supplier, setSupplier] = useState<Supplier | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (supplierId) {
      fetchSupplier();
    }
  }, [supplierId]);

  const fetchSupplier = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await axios.get(`/supplier/${supplierId}`);
      
      // Handle different response formats
      const supplierData = response.data?.supplier || response.data;
      
      if (!supplierData || !supplierData.id) {
        throw new Error("Invalid supplier data received");
      }
      
      setSupplier({
        id: supplierData.id,
        business_name: supplierData.business_name || supplierData.businessName,
        status: supplierData.status,
        support_email: supplierData.support_email || supplierData.supportEmail,
        support_phone: supplierData.support_phone || supplierData.supportPhone,
        is_active: supplierData.is_active !== undefined ? supplierData.is_active : supplierData.isActive,
        created_at: supplierData.created_at || supplierData.createdAt,
        updated_at: supplierData.updated_at || supplierData.updatedAt,
      });
    } catch (err: any) {
      console.error("Error fetching supplier:", err);
      setError(err.response?.data?.message || err.message || "Failed to load supplier details");
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
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
            <Link href="/admin/supplier">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Suppliers
            </Link>
          </Button>
          <div className="flex items-center gap-3 mb-2">
            <div className="h-10 w-1 bg-primary rounded-full"></div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Supplier Details</h1>
          </div>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <p className="text-red-600 mb-4">{error}</p>
              <Button onClick={fetchSupplier}>Try Again</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!supplier) {
    return (
      <div className="space-y-6">
        <div className="mb-8">
          <Button
            variant="ghost"
            asChild
            className="mb-4"
          >
            <Link href="/admin/supplier">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Suppliers
            </Link>
          </Button>
        </div>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8 text-gray-500">
              Supplier not found
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
          <Link href="/admin/supplier">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Suppliers
          </Link>
        </Button>
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Supplier Details</h1>
        </div>
        <p className="text-gray-600 dark:text-gray-400 mt-2 ml-4">Supplier #{supplier.id}</p>
      </div>

      {/* Supplier Overview */}
      <Card className="border-primary/20 hover:border-primary/40 transition-colors dark:bg-gray-800 dark:border-gray-700">
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center gap-2 text-lg">
            <Building2 className="h-5 w-5 text-primary" />
            Supplier Overview
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
                {getStatusBadge(supplier.status)}
              </div>
            </div>

            {/* Active Status */}
            <div className="flex flex-col gap-2 p-3 rounded-lg bg-muted/50 dark:bg-gray-700/50">
              <div className="flex items-center gap-2">
                <div className="p-1.5 rounded-md bg-primary/10 dark:bg-primary/20">
                  <Building2 className="h-3.5 w-3.5 text-primary" />
                </div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Active</span>
              </div>
              <div className="mt-1">
                <Badge variant={supplier.is_active ? "default" : "secondary"} className="text-sm">
                  {supplier.is_active ? "Yes" : "No"}
                </Badge>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Supplier Information */}
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <div className="flex items-center gap-2">
              <Building2 className="h-5 w-5 text-primary" />
              <CardTitle>Business Information</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Supplier ID</p>
              <p className="font-semibold">#{supplier.id}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Business Name</p>
              <p className="font-semibold text-lg">{supplier.business_name || 'N/A'}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Status</p>
              <div className="mt-1">
                {getStatusBadge(supplier.status)}
              </div>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Active</p>
              <Badge variant={supplier.is_active ? "default" : "secondary"} className="mt-1">
                {supplier.is_active ? "Active" : "Inactive"}
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
                Support Email
              </p>
              <p className="font-semibold mt-1">
                {supplier.support_email ? (
                  <a href={`mailto:${supplier.support_email}`} className="text-primary hover:underline">
                    {supplier.support_email}
                  </a>
                ) : (
                  <span className="text-muted-foreground">N/A</span>
                )}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground flex items-center gap-2">
                <Phone className="h-4 w-4" />
                Support Phone
              </p>
              <p className="font-semibold mt-1">
                {supplier.support_phone ? (
                  <a href={`tel:${supplier.support_phone}`} className="text-primary hover:underline">
                    {supplier.support_phone}
                  </a>
                ) : (
                  <span className="text-muted-foreground">N/A</span>
                )}
              </p>
            </div>
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
                {new Date(supplier.created_at).toLocaleString('en-US', {
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
                {new Date(supplier.updated_at).toLocaleString('en-US', {
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

