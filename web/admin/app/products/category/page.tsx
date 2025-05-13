'use client'
import useProductsStore from "@/app/libs/store/useProductStore"
import { useEffect, useState } from "react";
import Heading from "../../components/breadcrumb";
import { DataTable } from "./datatable";
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

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { ColumnDef } from "@tanstack/react-table";
import Link from "next/link";
import { LiaEdit } from "react-icons/lia";
import { Category } from '@/app/libs/types';
import { MdDeleteOutline } from "react-icons/md";
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";




export default function Products() {
  const {
    categories,
    success,
    loading,
    error,
    fetchCategories,
    editCategory,
    deleteCategory
  } = useProductsStore()

  const pages = [
    {
      "title": "Products",
      "href": "/products"
    },
    {
      "title": "Product Categories",
      "href": "/products/category"
    },
  ]

  const [deleteModal, setDeleteModal] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [newCategoryName, setNewCategoryName] = useState("")
  const [productId, setProductId] = useState(0)
  const handleEditCategory = async () => {
    editCategory(productId, newCategoryName)
    setIsEditDialogOpen(false)
  };


  const handleDeleteCategory = async (id: number) => {
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        toast.success("Category Deleted", {
          description: "The category has been deleted successfully.",
          position: "top-right"
        });

        deleteCategory(id)

        resolve();
      }, 500);
    });
  };

  useEffect(() => {
    fetchCategories();
  }, []);


  const columns: ColumnDef<Category>[] = [
    {
      accessorKey: "id",
      header: "Id",
    },
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <div className="flex items-center gap-2">
            <LiaEdit className="text-gray-700 cursor-pointer" onClick={() => { setProductId(row.original.id); setIsEditDialogOpen(true); setNewCategoryName(row.original.name) }} />
            <MdDeleteOutline className="text-red-900 cursor-pointer" onClick={() => { setProductId(row.original.id); setDeleteModal(true) }} />
          </div>
        );
      },
    },
  ];


  return (
    <>
      <Heading page={pages} heading="Categories" subheading="List of Registered categories" />
      <DataTable
        columns={columns}
        data={categories}
        loading={loading}
        search="name"
        searchPlaceholder="Search categories..."
      />

      <AlertDialog open={deleteModal}>
        <AlertDialogTrigger asChild>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete category</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogDescription>Are you sure you want to delete the category?</AlertDialogDescription>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => {
              setDeleteModal(false)
              setProductId(0)
            }}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction className="bg-red-900"
              onClick={() => {
                setDeleteModal(false)
                setProductId(0)
                handleDeleteCategory(productId)
              }
              }
            >
              Confirm
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Selected Category</DialogTitle>
          </DialogHeader>
          <Input
            placeholder="Category Name"
            value={newCategoryName}
            onChange={(e) => setNewCategoryName(e.target.value)}
            autoFocus
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleEditCategory}>Edit Category</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

    </>
  );
}
