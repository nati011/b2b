'use client'
import useResourceStore from "@/app/libs/store/useResourceStore"
import { useEffect, useState } from "react";
import Heading from "../components/breadcrumb";
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
import { Resource } from '@/app/libs/types';
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/ui/datatable";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";




export default function ResourceList() {
  const {
    success,
    loading,
    error,
    resource,
    fetchResources,
    editResource,
    deleteResource

  } = useResourceStore()

  const pages = [
    {
      "title": "Resource",
      "href": "/resources"
    }
  ]

  const [deleteModal, setDeleteModal] = useState(false)
  const [resourceId, setResourceId] = useState(0)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [formData, setFormData] = useState<Resource>({
    id: resourceId,
    name: "",
    action: ""
  })

  const handleEditResource = async () => {
    editResource(formData)
    setIsEditDialogOpen(false)
  };


  const handleDeleteResource = async (id: number) => {
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        toast.success("Resource Deleted", {
          description: "The resource has been deleted successfully.",
          position: "top-right"
        });

        deleteResource(id)

        resolve();
      }, 500);
    });
  };

  useEffect(() => {
    fetchResources();
  }, []);


  const columns: ColumnDef<Resource>[] = [
    {
      accessorKey: "id",
      header: "Id",
    },
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "action",
      header: "Action",
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <div className="flex items-center gap-2">
            <LiaEdit className="text-gray-700 cursor-pointer" onClick={() => { setIsEditDialogOpen(true); setFormData(prev => ({ ...prev, id: row.original.id, name: row.original.name, action: row.original.action })) }} />
          </div>
        );
      },
    },
  ];


  return (
    <>
      <Heading page={pages} heading="Resource" subheading="List of registered resources" />
      <DataTable
        columns={columns}
        data={resource}
        loading={loading}
        search="name"
        searchPlaceholder="Search resources..."
      />

      <AlertDialog open={deleteModal}>
        <AlertDialogTrigger asChild>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete resource</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogDescription>Are you sure you want to delete the resource?</AlertDialogDescription>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => {
              setDeleteModal(false)
              setResourceId(0)
            }}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction className="bg-red-900"
              onClick={() => {
                setDeleteModal(false)
                setResourceId(0)
                handleDeleteResource(resourceId)
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
            <DialogTitle>Edit Selected Resource</DialogTitle>
          </DialogHeader>
          <Separator orientation="horizontal" />
          <div className="grid gap-2">
            <Label className="font-semibold">
              Resource Name
            </Label>
            <Input
              placeholder="Resource Name"
              value={formData.name}
              onChange={(e) => {
                setFormData(prev => ({
                  ...prev,
                  name: e.target.value
                }));
              }}
              autoFocus
            />
          </div>

          <div className="grid gap-2">
            <Label className="font-semibold">
              Resource Description
            </Label>
            <Textarea

              placeholder="Resource Description"
              value={formData.action}
              onChange={(e) => {
                setFormData(prev => ({
                  ...prev,
                  action: e.target.value
                }));
              }}
              autoFocus
            />
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleEditResource}>Edit Resource</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

    </>
  );
}
