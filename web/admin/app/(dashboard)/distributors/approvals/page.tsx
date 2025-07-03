'use client'
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { useEffect, useState } from "react";
import Heading from "@/components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";
import { ColumnDef } from "@tanstack/react-table";
import { Distributor, DistributorVerdict } from '@/app/libs/types';
import Link from "next/link";
import { GoEye } from "react-icons/go";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CheckCircle, XCircle, Clock } from "lucide-react";
import { toast } from "sonner";

export default function DistributorApprovals() {
  const {
    distributors,
    loading,
    error,
    fetchDistributors,
    approveDistributor,
    rejectDistributor
  } = useDistributorsStore()

  const [pendingDistributors, setPendingDistributors] = useState<Distributor[]>([])

  useEffect(() => {
    fetchDistributors();
  }, []);

  useEffect(() => {
    // Filter distributors that are pending approval
    const pending = distributors.filter(dist => dist.verdict === DistributorVerdict.PENDING);
    setPendingDistributors(pending);
  }, [distributors]);

  useEffect(() => {
    if (error) {
      toast.error(error);
    }
  }, [error]);

  const pages: any[] = [
    {
      "title": "Distributors",
      "href": "/distributors",
    },
    {
      "title": "Approvals",
      "href": "/distributors/approvals",
    }
  ]

  const handleQuickApprove = async (id: number) => {
    try {
      await approveDistributor(id);
      toast.success("Distributor approved successfully");
    } catch (error: any) {
      toast.error(error.message || "Failed to approve distributor");
    }
  };

  const handleQuickReject = async (id: number) => {
    try {
      // For quick reject, we'll use a default message
      await rejectDistributor(id, "Application rejected by admin");
      toast.success("Distributor rejected successfully");
    } catch (error: any) {
      toast.error(error.message || "Failed to reject distributor");
    }
  };

  const columns: ColumnDef<Distributor>[] = [
    {
      accessorKey: "id",
      header: "ID",
      cell: ({ row }) => <span className="font-mono text-sm">#{row.original.id}</span>,
    },
    {
      accessorKey: "name",
      header: "Business Name",
    },
    {
      accessorKey: "tin",
      header: "TIN",
      cell: ({ row }) => <span className="font-mono text-sm">{row.original.tin}</span>,
    },
    {
      accessorKey: "region",
      header: "Region",
    },
    {
      accessorKey: "general_zone",
      header: "Zone",
    },
    {
      id: "status",
      header: "Status",
      cell: ({ row }) => {
        if (row.original.verdict === DistributorVerdict.PENDING) {
          return (
            <Badge variant="secondary" className="flex items-center gap-1">
              <Clock className="w-3 h-3" />
              Pending Approval
            </Badge>
          );
        }
        if (row.original.verdict === DistributorVerdict.REJECTED) {
          return (
            <Badge variant="destructive" className="flex items-center gap-1">
              <XCircle className="w-3 h-3" />
              Rejected
            </Badge>
          );
        }
        if (row.original.verdict === DistributorVerdict.APPROVED) {
          return (
            <Badge variant="default" className="flex items-center gap-1 bg-green-100 text-green-800 hover:bg-green-100">
              <CheckCircle className="w-3 h-3" />
              Approved
            </Badge>
          );
        }
        return null;
      },
    },
    {
      id: "actions",
      enableHiding: false,
      header: "Actions",
      cell: ({ row }: any) => {
        return (
          <div className="flex items-center gap-2">
            <Link 
              href={`/distributors/detail/${row.original.id}`}
              className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-8 w-8"
            >
              <GoEye className="h-4 w-4" />
            </Link>
            <Button
              size="sm"
              onClick={() => handleQuickApprove(row.original.id)}
              disabled={loading}
              className="bg-transparent border border-green-600 text-green-700 h-8 px-2 hover:bg-green-100"
            >
              <CheckCircle className="w-3 h-3 mr-1" />
              Approve
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={() => handleQuickReject(row.original.id)}
              disabled={loading}
              className="h-8 px-2 bg-transparent border border-red-900 text-red-900 hover:bg-red-100"
            >
              <XCircle className="w-3 h-3 mr-1" />
              Reject
            </Button>
          </div>
        );
      },
    },
  ];

  return (
    <>
      <Heading 
        page={pages} 
        heading="Distributor Approvals" 
        subheading={`${pendingDistributors.length} pending approval${pendingDistributors.length !== 1 ? 's' : ''}`} 
      />
      
      {pendingDistributors.length === 0 ? (
        <div className="text-center py-12">
          <Clock className="w-16 h-16 mx-auto mb-4 text-gray-400" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No Pending Approvals</h3>
          <p className="text-gray-600 mb-4">All distributor applications have been reviewed.</p>
          <Link href="/distributors">
            <Button variant="outline">View All Distributors</Button>
          </Link>
        </div>
      ) : (
        <DataTableLayout
          columns={columns}
          data={pendingDistributors}
          loading={loading}
          button={false}
          search="name"
          searchPlaceholder="Search pending distributors..."
        />
      )}
    </>
  );
} 