"use client";
import useRolesStore from "@/app/libs/store/useRoleStore";
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ColumnDef } from "@tanstack/react-table";
import { LiaEdit } from "react-icons/lia";
import { Role, Resource } from "@/app/libs/types";
import { toast } from "sonner";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/ui/datatable";
import { Search, MoreHorizontal, Settings } from "lucide-react";
import { MdDeleteOutline } from "react-icons/md";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import useResourceStore from "@/app/libs/store/useResourceStore";

export default function Roles() {
  const {
    success,
    rolesLoading: loading,
    rolesError: error,
    roles,
    fetchRoles,
    editRole,
    deleteRole,
    createRole,
  } = useRolesStore();

  const {
    resources,
    resourcesLoading,
    fetchResources,
    roleResources,
      addMultipleResourcesToRole
  } = useResourceStore();

  const pages = [
    {
      title: "Roles",
      href: "/roles", 
    },
    {
      title: "Role Categories",
      href: "/roles/role",
      title: "Role Categories",
      href: "/roles/role",
    },
  ];
  ];

  const [deleteModal, setDeleteModal] = useState(false);
  const [roleId, setRoleId] = useState(0);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [isAssignResourceDialogOpen, setIsAssignResourceDialogOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedResources, setSelectedResources] = useState<number[]>([]);
  const [formData, setFormData] = useState<Role>({
    id: roleId,
    name: "",
    desc: "",
  });

  const filteredRoles = roles.filter((role) =>
    role.name.toLowerCase().includes(searchTerm.toLowerCase())
  );
    desc: "",
  });

  const filteredRoles = roles.filter((role) =>
    role.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleEditRole = async () => {
    await editRole(formData);
    setIsEditDialogOpen(false);
    await editRole(formData);
    setIsEditDialogOpen(false);
  };

  const handleAddRole = async () => {
    try {
      await createRole(formData);
      setIsAddDialogOpen(false);
      setFormData({ id: 0, name: "", desc: "" });
    } catch (error) {
    }
  };

  const handleDeleteRole = async (id: number) => {
    deleteRole(id);
    setDeleteModal(false);
  };

  const handleAssignResource = async () => {
    try {
      console.log(roleId, selectedResources)
      await addMultipleResourcesToRole(roleId, selectedResources);
      setIsAssignResourceDialogOpen(false);
      setSelectedResources([]);
      toast.success("Resources assigned successfully");
    } catch (error) {
      toast.error("Failed to assign resources");
    }
  };

  const handleOpenAssignResource = async (id: number) => {
    setRoleId(id);
    await fetchResources();
    setIsAssignResourceDialogOpen(true);
  };

  useEffect(() => {
    fetchRoles();
    fetchResources();
  }, []);

  useEffect(() => {
    if (error) toast.error(error);
    if (success) toast.success(success);
  }, [success, error]);
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
      accessorKey:"desc",
      header:"Description"
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) => {
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                onClick={() => {
                  setFormData({
                    id: row.original.id,
                    name: row.original.name,
                    desc: row.original.desc,
                  });
                  setIsEditDialogOpen(true);
                }}
              >
                <LiaEdit className="mr-2 h-4 w-4" />
                Edit
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => handleOpenAssignResource(row.original.id)}
              >
                <Settings className="mr-2 h-4 w-4" />
                Assign Resources
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => {
                  setRoleId(row.original.id);
                  setDeleteModal(true);
                }}
                className="text-red-600"
              >
                <MdDeleteOutline className="mr-2 h-4 w-4" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
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
          <AlertDialogDescription>
            Are you sure you want to delete this role?
          </AlertDialogDescription>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setDeleteModal(false)}>
            <AlertDialogCancel onClick={() => setDeleteModal(false)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              className='bg-red-900 hover:bg-red-800'
              onClick={() => handleDeleteRole(roleId)}
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
      {/* Edit Role Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Role</DialogTitle>
            <DialogTitle>Edit Role</DialogTitle>
          </DialogHeader>
          <Input
            placeholder='Role Name'
            placeholder='Role Name'
            value={formData.name}
            onChange={(e) => {
              setFormData((prev) => ({
              setFormData((prev) => ({
                ...prev,
                name: e.target.value,
                name: e.target.value,
              }));
            }}
            autoFocus
          />
          <Textarea
            placeholder='Role Description'
            placeholder='Role Description'
            value={formData.desc}
            onChange={(e) => {
              setFormData((prev) => ({
              setFormData((prev) => ({
                ...prev,
                desc: e.target.value,
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

      {/* Assign Resources Dialog */}
      <Dialog open={isAssignResourceDialogOpen} onOpenChange={setIsAssignResourceDialogOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Assign Resources to Role</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="text-sm text-muted-foreground">
              Select the resources you want to assign to this role:
            </div>
            <div className="space-y-2 max-h-60 overflow-y-auto">
              {resources.map((resource) => (
                <div key={resource.id} className="flex items-center space-x-2">
                  <Checkbox
                    id={`resource-${resource.id}`}
                    checked={selectedResources.includes(resource.id)}
                    onCheckedChange={(checked) => {
                      if (checked) {
                        setSelectedResources([...selectedResources, resource.id]);
                      } else {
                        setSelectedResources(selectedResources.filter(id => id !== resource.id));
                      }
                    }}
                  />
                  <label
                    htmlFor={`resource-${resource.id}`}
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  >
                    {resource.resource} - {resource.action}
                  </label>
                </div>
              ))}
            </div>
          </div>
          <DialogFooter>
            <Button
              variant='outline'
              onClick={() => {
                setIsAssignResourceDialogOpen(false);
                setSelectedResources([]);
              }}
            >
              Cancel
            </Button>
            <Button 
              onClick={handleAssignResource}
              disabled={selectedResources.length === 0}
            >
              Assign Resources
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}