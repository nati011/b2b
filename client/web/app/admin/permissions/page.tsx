"use client";
import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Plus, Edit, Trash2 } from "lucide-react";
import { toast } from "sonner";
import * as adminActions from "@/app/actions/admin";
import type { Permission } from "@/app/actions/admin";

export default function AdminPermissionsPage() {
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
  const [formData, setFormData] = useState({
    resource_code: "",
    action: "",
    description: "",
  });

  useEffect(() => {
    fetchPermissions();
  }, []);

  const fetchPermissions = async () => {
    try {
      setLoading(true);
      const data = await adminActions.getPermissions();
      setPermissions(data);
    } catch (error: any) {
      toast.error(error.message || "Failed to fetch permissions");
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingPermission) {
        await adminActions.updatePermission(editingPermission.id, formData);
        toast.success("Permission updated successfully");
      } else {
        await adminActions.createPermission(formData);
        toast.success("Permission created successfully");
      }
      setDialogOpen(false);
      resetForm();
      fetchPermissions();
    } catch (error: any) {
      toast.error(error.message || "Failed to save permission");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this permission?")) return;
    try {
      await adminActions.deletePermission(id);
      toast.success("Permission deleted successfully");
      fetchPermissions();
    } catch (error: any) {
      toast.error(error.message || "Failed to delete permission");
    }
  };

  const resetForm = () => {
    setFormData({
      resource_code: "",
      action: "",
      description: "",
    });
    setEditingPermission(null);
  };

  const openEditDialog = (permission: Permission) => {
    setEditingPermission(permission);
    setFormData({
      resource_code: permission.resource_code,
      action: permission.action,
      description: permission.description,
    });
    setDialogOpen(true);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-8 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Permissions</h1>
          <p className="text-gray-600 mt-2">Manage system permissions</p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={resetForm}>
              <Plus className="h-4 w-4 mr-2" />
              Add Permission
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {editingPermission ? "Edit Permission" : "Create Permission"}
              </DialogTitle>
              <DialogDescription>
                {editingPermission
                  ? "Update permission information"
                  : "Create a new permission"}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Label htmlFor="resource_code">Resource Code</Label>
                <Input
                  id="resource_code"
                  value={formData.resource_code}
                  onChange={(e) =>
                    setFormData({ ...formData, resource_code: e.target.value })
                  }
                  required
                />
              </div>
              <div>
                <Label htmlFor="action">Action</Label>
                <Input
                  id="action"
                  value={formData.action}
                  onChange={(e) => setFormData({ ...formData, action: e.target.value })}
                  required
                />
              </div>
              <div>
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  rows={3}
                />
              </div>
              <div className="flex justify-end gap-2">
                <Button type="button" variant="outline" onClick={() => setDialogOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit">Save</Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Permissions</CardTitle>
          <CardDescription>A list of all permissions in the system</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Resource</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Description</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {permissions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} className="text-center text-gray-500">
                    No permissions found
                  </TableCell>
                </TableRow>
              ) : (
                permissions.map((permission) => (
                  <TableRow key={permission.id}>
                    <TableCell>
                      <Badge variant="outline">{permission.resource_code}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant="secondary">{permission.action}</Badge>
                    </TableCell>
                    <TableCell>{permission.description || "-"}</TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => openEditDialog(permission)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDelete(permission.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}






