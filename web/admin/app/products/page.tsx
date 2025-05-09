'use client'
import useProductsStore from "@/app/libs/store/useProductStore"
import { useEffect, useState } from "react";
import Heading from "../components/breadcrumb";
import { DataTable } from "@/components/ui/datatable";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { ColumnDef } from "@tanstack/react-table";
import { Product } from "../libs/types";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import Image from "next/image";


export default function Products() {
  const {
    products,
    loading,
    error,
    fetchProducts,
    addStock,
    depleteStock
  } = useProductsStore()

  const pages = [
    {
      "title": "Products",
      "href": "/products"
    },
  ]

  const [stockModal, setStockModal] = useState(false)
  const [depleteStockModal, setDepleteStockModal] = useState(false)
  const [stockQuantity, setStockQuantity] = useState(0)
  const [productId, setProductId] = useState(0)


  useEffect(() => {
    fetchProducts();
  }, []);


  const columns: ColumnDef<Product>[] = [
    {
      accessorKey: "Images",
      header: "",
      cell: ({ row }) => {
        const image = row.original.Images[0].ImageUrl
        console.log(image)
        return <div className="border rounded">
          <Image src={image} alt="product-image" width={40} height={40} />
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
      accessorKey: "Price",
      header: "Price",
    },
    {
      accessorKey: "AvailableStock",
      header: "Available Stock",
    },
    {
      accessorKey: "ReservedStock",
      header: "Reserved Stock",
    },
    {
      accessorKey: "Stock",
      header: "Stock",
    },
    {
      accessorKey: "IsActive",
      header: () => <div className="text-left">Status</div>,
      cell: ({ row }) => {
        const status = row.getValue("IsActive")
        return <div className={status ? "border border-amber-500 py-1 mx-auto rounded-md text-amber-500 font-medium text-center text-xs" : "border border-emerald-500  py-1 mx-auto rounded-md  text-emerald-500 font-medium text-center text-xs"}>
          {status ? "Inactive" : "Active"}
        </div>
      },
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        const item = row.original;
        return (
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
                setProductId(row.original.Id)
                setStockModal(true)
              }
              }>
                Add Stock
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => {
                setProductId(row.original.Id)
                setDepleteStockModal(true)
              }
              }>
                Deplete Stock
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },

  ];

  const handleAddStock = () => {
    console.log(stockQuantity, productId)
    addStock(stockQuantity, productId)
  };


  const handlDepleteStock = () => {
    depleteStock(stockQuantity, productId)

  };
  return (
    <>
      <Heading page={pages} heading="Products" subheading="List of Registered products" />
      <DataTable
        columns={columns}
        data={products}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Products", url: "/products/form" }}
        search="Name"
        searchPlaceholder="Search products..."
      />
      <AlertDialog open={stockModal}>
        <AlertDialogTrigger asChild>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Add Stock</AlertDialogTitle>
          </AlertDialogHeader>
          <Input
            placeholder="Quantity"
            value={stockQuantity}
            onChange={(e) => setStockQuantity(parseFloat(e.target.value))}
            autoFocus
          />
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => { setStockModal(false); setStockQuantity(0) }}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                handleAddStock();
                setStockModal(false)
              }
              }
            >
              Add Stock
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      <AlertDialog open={depleteStockModal}>
        <AlertDialogTrigger asChild>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Deplete Stock</AlertDialogTitle>
          </AlertDialogHeader>
          <Input
            placeholder="Quantity"
            value={stockQuantity}
            onChange={(e) => setStockQuantity(parseFloat(e.target.value))}
            autoFocus
          />
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => { setDepleteStockModal(false); setStockQuantity(0) }}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                handlDepleteStock();
                setDepleteStockModal(false)
              }
              }
            >
              Deplete Stock
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
