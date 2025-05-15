'use client'
import useProductsStore from "@/app/libs/store/useProductStore"
import { useEffect, useState } from "react";
import Heading from "@/app/components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";
import { ColumnDef } from "@tanstack/react-table";
import { ConfigurableProduct, Product } from "@/app/libs/types";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal } from "lucide-react";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import Image from "next/image";
import Link from "next/link";
import { LiaEdit } from "react-icons/lia";
import { Badge } from "@/components/ui/badge";

type statusProduct = {
  Id: number;
  Name: string
  Status: boolean
}

export default function Products() {
  const {
    configurable_products,
    success,
    loading,
    error,
    fetchConfigurableProducts,
    updateConfigurableProductStatus
  } = useProductsStore()

  const pages = [
    {
      "title": "Products",
      "href": "/products"
    },
  ]

  const [statusModal, setStatusModal] = useState(false)
  const [statusProduct, setStatusProduct] = useState<statusProduct>()
  const handleActivateProduct = () => {
    updateConfigurableProductStatus(statusProduct!.Id, statusProduct!.Status)
    setStatusModal(false)
  }



  useEffect(() => {
    fetchConfigurableProducts();
  }, []);


  const columns: ColumnDef<ConfigurableProduct>[] = [
    {
      accessorKey: "Images",
      header: "",
      cell: ({ row }) => {
        const image = row.original.Images[0].ImageUrl
        console.log(image)
        return <div className="border rounded w-fit">
          <Image src={image} alt="product-image" width={100} height={100} />
        </div>
      },
    },
    {
      accessorKey: "Id",
      header: "Id",
    },
    {
      accessorKey: "Name",
      header: "Name",
    },
    {
      accessorKey: " Attributes",
      header: "Attributes",
      cell: ({ row }) => {
        return (
          <div className="flex gap-2">
            {row.original.Attributes.map((a, index) => (
              <Badge key={index} variant={"outline"} className="border-blue-900 text-blue-950">
                {a}
              </Badge>
            ))}
          </div>
        );
      },
    },
    {
      accessorKey: "IsAvailable",
      header: () => <div className="text-left">Status</div>,
      cell: ({ row }) => {
        const status = row.getValue("IsAvailable")
        return <div className={!status ? "border border-amber-500 py-1 mx-auto rounded-md text-amber-500 font-medium text-center text-xs" : "border border-emerald-500  py-1 mx-auto rounded-md  text-emerald-500 font-medium text-center text-xs"}>
          {status ? "Available" : "Unavailable"}
        </div>
      },
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <div className="flex items-center gap-2">
            <Link href={`/products/configurable/${row.original.Id}`}>
              <LiaEdit className="text-gray-700 text-xl" />
            </Link>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                <DropdownMenuItem onClick={() => {
                  setStatusProduct({ Id: row.original.Id, Name: row.original.Name, Status: row.original.IsAvailable })
                  setStatusModal(true)
                }
                }>
                  {

                    !row.original.IsAvailable ? "Activate Product" : "Deactivate Product"
                  }
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        );
      },
    },

  ];
  return (
    <>
      <Heading page={pages} heading="Configurable Products" subheading="List of Registered cofigurable products" />
      <DataTableLayout
        columns={columns}
        data={configurable_products}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Products", url: "/products/configurable/form" }}
        search="Name"
        searchPlaceholder="Search products..."
      />
      <AlertDialog open={statusModal} onOpenChange={setStatusModal}>
        <AlertDialogTrigger asChild>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Activate Product</AlertDialogTitle>
          </AlertDialogHeader>
          <p>
            Are you sure you want to activate <strong>{statusProduct?.Name}</strong>?
          </p>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setStatusModal(false)}>Cancel</Button>
            <Button onClick={handleActivateProduct}>
              {
                statusProduct?.Status ? "Deactivate" : "Activate"
              }
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

    </>
  );
}
