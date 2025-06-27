"use client";
import useRolesStore from "@/app/libs/store/useRoleStore";
import { useEffect, useState } from "react";
import Heading from "../../../components/breadcrumb";
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

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { ColumnDef } from "@tanstack/react-table";
import { LiaEdit } from "react-icons/lia";
import { Role } from "@/app/libs/types";
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/ui/datatable";
import { Search } from "lucide-react";
import { MdDeleteOutline } from "react-icons/md";
import { Textarea } from "@/components/ui/textarea";

export default function Roles() {
  const {
    success,
    loading,
    error,
    roles,
    fetchRoles,
    editRole,
    deleteRole,
    createRole,
  } = useRolesStore();

  const pages = [
    {
      title: "Roles",
      href: "/roles",
    },
    {
      title: "Role Categories",
      href: "/roles/role",
    },
  ];

  const [deleteModal, setDeleteModal] = useState(false);
  const [roleId, setRoleId] = useState(0);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [formData, setFormData] = useState<Role>({
    id: roleId,
    name: "",
    desc: "",
  });

  const filteredRoles = roles.filter((role) =>
    role.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleEditRole = async () => {
    await editRole(formData);
    setIsEditDialogOpen(false);
  };

  const handleAddRole = async () => {
    try {
      await createRole(formData);
      setIsAddDialogOpen(false);
      setFormData({ id: 0, name: "", desc: "" });
    } catch (error) {
      // Error handled by store
    }
  };

  const handleDeleteRole = async (id: number) => {
    deleteRole(id);
    setDeleteModal(false);
  };

  useEffect(() => {
    fetchRoles();
  }, []);

  useEffect(() => {
    if (error) toast.error(error);
    if (success) toast.success(success);
  }, [success, error]);

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
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <div className='flex items-center gap-2'>
            <LiaEdit
              className='text-gray-700 cursor-pointer'
              onClick={() => {
                setFormData({
                  id: row.original.id,
                  name: row.original.name,
                  desc: row.original.desc,
                });
                setIsEditDialogOpen(true);
              }}
            />
            <MdDeleteOutline
              className='text-red-900 cursor-pointer'
              onClick={() => {
                setRoleId(row.original.id);
                setDeleteModal(true);
              }}
            />
          </div>
        );
      },
    },
  ];

  return (
    <>
      <Heading
        page={pages}
        heading='Roles'
        subheading='List of registered roles'
      />

      <div className='w-full bg-card rounded-lg border border-border/50 shadow-sm'>
        {/* Header with Search and Add Button */}
        <div className='flex flex-col sm:flex-row justify-end gap-4 items-start sm:items-center p-6 border-b border-border/50 space-y-4 sm:space-y-0'>
          <div className='flex items-center space-x-4 w-full sm:w-auto'>
            <div className='relative flex-1 sm:flex-none'>
              <Search className='absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground' />
              <Input
                placeholder='Search roles...'
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className='pl-10 bg-background border-border focus:border-primary transition-colors rounded-sm'
              />
            </div>
          </div>

          <div className='flex items-center space-x-2'>
            <Button
              className='bg-primary hover:bg-primary/90 text-primary-foreground shadow-sm transition-all duration-200'
              onClick={() => setIsAddDialogOpen(true)}
            >
              + Add Role
            </Button>
          </div>
        </div>

        <DataTable columns={columns} data={filteredRoles} loading={loading} />
      </div>

      {/* Add Role Dialog */}
      <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add New Role</DialogTitle>
          </DialogHeader>
          <Input
            placeholder='Role Name'
            value={formData.name}
            onChange={(e) => {
              setFormData((prev) => ({
                ...prev,
                name: e.target.value,
              }));
            }}
            autoFocus
          />
          <Textarea
            placeholder='Role Description'
            value={formData.desc}
            onChange={(e) => {
              setFormData((prev) => ({
                ...prev,
                desc: e.target.value,
              }));
            }}
          />
          <DialogFooter>
            <Button variant='outline' onClick={() => setIsAddDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleAddRole}>Add Role</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteModal}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete role</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogDescription>
            Are you sure you want to delete this role?
          </AlertDialogDescription>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setDeleteModal(false)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              className='bg-red-900 hover:bg-red-800'
              onClick={() => handleDeleteRole(roleId)}
            >
              Confirm
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Edit Role Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Role</DialogTitle>
          </DialogHeader>
          <Input
            placeholder='Role Name'
            value={formData.name}
            onChange={(e) => {
              setFormData((prev) => ({
                ...prev,
                name: e.target.value,
              }));
            }}
            autoFocus
          />
          <Textarea
            placeholder='Role Description'
            value={formData.desc}
            onChange={(e) => {
              setFormData((prev) => ({
                ...prev,
                desc: e.target.value,
              }));
            }}
          />
          <DialogFooter>
            <Button
              variant='outline'
              onClick={() => setIsEditDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button onClick={handleEditRole}>Save Changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
