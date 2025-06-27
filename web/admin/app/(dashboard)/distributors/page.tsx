"use client";
import useDistributorsStore from "@/app/libs/store/useDistributorStore";
import { useEffect, useState } from "react";
import Heading from "@/components/breadcrumb";
import { ColumnDef } from "@tanstack/react-table";
import { Distributor, DistributorVerdict } from "@/app/libs/types";
import Link from "next/link";
import { GoEye } from "react-icons/go";
import { Badge } from "@/components/ui/badge";
import { CheckCircle, XCircle, Clock, PauseCircle, Search } from "lucide-react";
import { DataTable } from "@/components/ui/datatable";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { CiFilter } from "react-icons/ci";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export default function Distributors() {
  const [offset, setOffset] = useState(0);
  const [status, setStatus] = useState("ALL");
  const limit = 10;
  const { distributors, loading, error, totalCount, fetchDistributors } =
    useDistributorsStore();

  useEffect(() => {
    fetchDistributors(status);
  }, [status]);

  const pages: any[] = [
    {
      title: "Distributors",
      href: `/distributors`,
    },
  ];
  const handleNext = () => {
    if (offset + limit < (totalCount || 0)) {
      setOffset(offset + limit);
    }
  };
  const handlePrevious = () => {
    if (offset - limit >= 0) {
      setOffset(offset - limit);
    }
  };
  const canNext = offset + limit < (totalCount || 0);
  const canPrevious = offset > 0;

  const getStatusBadge = (distributor: Distributor) => {
    if (distributor.verdict === DistributorVerdict.PENDING) {
      return (
        <Badge variant='secondary' className='flex items-center gap-1'>
          <Clock className='w-3 h-3' />
          Pending Approval
        </Badge>
      );
    }
    if (distributor.verdict === DistributorVerdict.REJECTED) {
      return (
        <Badge variant='destructive' className='flex items-center gap-1'>
          <XCircle className='w-3 h-3' />
          Rejected
        </Badge>
      );
    }
    if (
      distributor.verdict === DistributorVerdict.APPROVED &&
      distributor.is_active
    ) {
      return (
        <Badge
          variant='default'
          className='flex items-center gap-1 bg-green-100 text-green-800 hover:bg-green-100'
        >
          <CheckCircle className='w-3 h-3' />
          Approved
        </Badge>
      );
    }
    if (
      distributor.verdict === DistributorVerdict.APPROVED &&
      !distributor.is_active
    ) {
      return (
        <Badge variant='destructive' className='flex items-center gap-1'>
          <PauseCircle className='w-3 h-3' />
          Inactive
        </Badge>
      );
    }
    return null;
  };

  const columns: ColumnDef<Distributor>[] = [
    {
      accessorKey: "id",
      header: "ID",
      cell: ({ row }) => (
        <span className='font-mono text-sm'>#{row.original.id}</span>
      ),
    },
    {
      accessorKey: "name",
      header: "Business Name",
    },
    {
      accessorKey: "tin",
      header: "TIN",
      cell: ({ row }) => (
        <span className='font-mono text-sm'>{row.original.tin}</span>
      ),
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
      cell: ({ row }) => getStatusBadge(row.original),
    },
    {
      id: "actions",
      enableHiding: false,
      header: "Actions",
      cell: ({ row }: any) => {
        return (
          <div className='flex items-center gap-2'>
            <Link
              href={`/distributors/detail/${row.original.id}`}
              className='inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 border border-input bg-background hover:bg-accent hover:text-accent-foreground h-8 w-8'
            >
              <GoEye className='h-4 w-4' />
            </Link>
          </div>
        );
      },
    },
  ];

  return (
    <>
      <Heading
        page={pages}
        heading='Distributors'
        subheading='Manage distributor registrations and approvals'
      />

      <div className='w-full bg-card rounded-lg border border-border/50 shadow-sm'>
        {/* Header */}
        <div className='flex flex-col sm:flex-row justify-between gap-4 items-start sm:items-center p-6 border-b border-border/50 space-y-4 sm:space-y-0'>
          <div className='flex items-center space-x-4 w-full sm:w-auto'>
            <div className='relative flex-1 sm:flex-none'>
              <Search className='absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground' />
              <Input
                placeholder={"Search distributors..."}
                className='pl-10 bg-background border-border focus:border-primary transition-colors rounded-sm'
              />
            </div>
          </div>

          <div className='flex items-center space-x-2'>
            <Link href='/distributors/form' passHref>
              <Button className='bg-primary hover:bg-primary/90 text-primary-foreground shadow-sm transition-all duration-200'>
                + Register Distributor
              </Button>
            </Link>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant={"outline"} className='shadow-none'>
                  <CiFilter />
                  Filter
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align='end'>
                <DropdownMenuItem
                  onClick={() => {
                    setStatus("PENDING");
                  }}
                >
                  Pending
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => {
                    setStatus("APPROVED");
                  }}
                >
                  Approved
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => {
                    setStatus("REJECTED");
                  }}
                >
                  Rejected
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
        <DataTable
          columns={columns}
          data={distributors}
          loading={loading}
          total={totalCount || 0}
          onNext={handleNext}
          onPrevious={handlePrevious}
          canNext={canNext}
          canPrevious={canPrevious}
        />
      </div>
    </>
  );
}
