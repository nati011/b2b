'use client'
import useRolesStore from "@/app/libs/store/useRoleStore"
import { useEffect, useState } from "react";
import Heading from "../../components/breadcrumb";
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
import { Role } from '@/app/libs/types';
import { MdDeleteOutline } from "react-icons/md";
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { DataTableLayout } from "@/components/ui/datatablelayout";
import { Textarea } from "@/components/ui/textarea";




export default function Roles() {
  const {
    success,
    loading,
    error,
    roles,
    fetchRoles,
    editRole,
    deleteRole
  } = useRolesStore()

  const pages = [
    {
      "title": "Roles",
      "href": "/roles"
    },
    {
      "title": "Role Categories",
      "href": "/roles/role"
    },
  ]

  const [deleteModal, setDeleteModal] = useState(false)
  const [roleId, setRoleId] = useState(0)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [formData, setFormData] = useState<Role>({
    id: roleId,
    name: "",
    desc: ""
  })

  const handleEditRole = async () => {
    editRole(formData)
    setIsEditDialogOpen(false)
  };


  const handleDeleteRole = async (id: number) => {
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        toast.success("Role Deleted", {
          description: "The role has been deleted successfully.",
          position: "top-right"
        });

        deleteRole(id)

        resolve();
      }, 500);
    });
  };

  useEffect(() => {
    fetchRoles();
  }, []);


  const columns: ColumnDef<Role>[] = [
    {
      accessorKey: "id",
      header: "Id",
    },
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "desc",
      header: "Description",
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <div className="flex items-center gap-2">
            <LiaEdit className="text-gray-700 cursor-pointer" onClick={() => { setIsEditDialogOpen(true); setFormData(prev => ({ ...prev, id: row.original.id, name: row.original.name, desc: row.original.desc })) }} />
          </div>
        );
      },
    },
  ];


  return (
    <>
      <Heading page={pages} heading="Roles" subheading="List of registered roles" />
      <DataTableLayout
        columns={columns}
        data={roles}
        loading={loading}
        search="name"
        searchPlaceholder="Search roles..."
      />

      <AlertDialog open={deleteModal}>
        <AlertDialogTrigger asChild>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete role</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogDescription>Are you sure you want to delete the role?</AlertDialogDescription>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => {
              setDeleteModal(false)
              setRoleId(0)
            }}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction className="bg-red-900"
              onClick={() => {
                setDeleteModal(false)
                setRoleId(0)
                handleDeleteRole(roleId)
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
            <DialogTitle>Edit Selected Role</DialogTitle>
          </DialogHeader>
          <Input
            placeholder="Role Name"
            value={formData.name}
            onChange={(e) => {
              setFormData(prev => ({
                ...prev,
                name: e.target.value
              }));
            }}
            autoFocus
          />
          <Textarea

            placeholder="Role Description"
            value={formData.desc}
            onChange={(e) => {
              setFormData(prev => ({
                ...prev,
                desc: e.target.value
              }));
            }}
            autoFocus
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleEditRole}>Edit Role</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

    </>
  );
}
