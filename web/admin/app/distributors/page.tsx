'use client'
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { useEffect } from "react";
import Heading from "../components/breadcrumb";
import { DataTable } from "@/components/ui/datatable";
import { ColumnDef } from "@tanstack/react-table";
import { Distributor } from '@/app/libs/types';
import Link from "next/link";
import { GoEye } from "react-icons/go";

export default function Distributors() {
  const {
    distributors,
    loading,
    error,
    fetchDistributors
  } = useDistributorsStore()

  useEffect(() => {
    fetchDistributors();
  }, []);

  const pages: any[] = [
    {
      "title": "Distributors",
      "href": `/distributors`,
    }
  ]
  const columns: ColumnDef<Distributor>[] = [
    {
      accessorKey: "id",
      header: "Id",
    },
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "tin",
      header: "Tin",
    },
    {
      accessorKey: "general_zone",
      header: "General Zone",
    },
    {
      accessorKey: "region",
      header: "Region",
    },
    {
      accessorKey: "woreda",
      header: "Woreda",
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <div className="flex items-center gap-2 text-gray-900">
            <Link href={`/distributors/detail/${row.original.id}`}>
              <GoEye />
            </Link>
          </div>
        );
      },
    },
  ];


  return (
    <>
      <Heading page={pages} heading="Distributors" subheading="List of Registered Distributors" />
      <DataTable
        columns={columns}
        data={distributors}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Distributor", url: "/distributors/form" }}
        search="name"
        searchPlaceholder="Search distributor..."
      />
    </>

  );
}
