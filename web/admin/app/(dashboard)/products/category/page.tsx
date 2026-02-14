'use client'
import { useEffect, useState } from "react";
import Heading from "../../../../components/breadcrumb";
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
import { LiaEdit } from "react-icons/lia";
import { Category } from '@/app/libs/types';
import { MdDeleteOutline } from "react-icons/md";
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import useCategoryStore from "@/app/libs/store/useCategories";
import { Search } from "lucide-react";
import { DataTable } from "@/components/ui/datatable";




export default function Products() {

  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [deleteModal, setDeleteModal] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [newCategoryName, setNewCategoryName] = useState("")
  const [productId, setProductId] = useState(0)

  const {
    success,
    categoriesError,
    categoriesLoading,
    categories,
    fetchCategories,
    editCategory,
    deleteCategory,
    createCategory
  } = useCategoryStore()

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
      const handleAddCategory = () => {
          if (newCategoryName.trim()) {
              createCategory(newCategoryName.trim())
              setNewCategoryName("");
              setIsAddDialogOpen(false);
          }
      };
  
  

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(()=>{
    if(success !=null){
      toast.success(success)
    }

    if(categoriesError !=null){
      toast.error(categoriesError)
    }

  },[success, categoriesError])


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
      <div className="w-full bg-card rounded-lg border border-border/50 shadow-sm">
            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-end gap-4 items-start sm:items-center p-6 border-b border-border/50 space-y-4 sm:space-y-0">
                <div className="flex items-center space-x-4 w-full sm:w-auto">
                    <div className="relative flex-1 sm:flex-none">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <Input
                            placeholder={"Search categories..."}
                          
                            className="pl-10 bg-background border-border focus:border-primary transition-colors rounded-sm"
                        />
                    </div>
                </div>
                
                <div className="flex items-center space-x-2">
                            <Button
                                className="bg-primary hover:bg-primary/90 text-primary-foreground shadow-sm transition-all duration-200"
                                onClick={() => setIsAddDialogOpen(true)}
                            >
                                + Register Categories
                            </Button>
                </div>
            </div>
            <DataTable
          columns={columns}
          data={categories}
          loading={categoriesLoading}
                          />
            </div>


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
                deleteCategory(productId)
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
            <Button onClick={()=>{
              editCategory(productId, newCategoryName)
              setIsEditDialogOpen(false)
            }}>Edit Category</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

            {/* Add Category Dialog */}
            <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Add New Category</DialogTitle>
                    </DialogHeader>
                    <Input
                        placeholder="Category Name"
                        value={newCategoryName}
                        onChange={(e) => setNewCategoryName(e.target.value)}
                        autoFocus
                    />
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsAddDialogOpen(false)}>Cancel</Button>
                        <Button onClick={handleAddCategory}>Add Category</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
    </>

  );
}
