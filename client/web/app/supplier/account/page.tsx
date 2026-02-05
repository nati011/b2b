"use client";
import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Plus, Edit, Trash2, CreditCard } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import axios from "@/lib/axios";
import {
  listBankAccounts,
  createBankAccount,
  updateBankAccount,
  deleteBankAccount,
  BankAccountRequest,
  BankAccountResponse,
} from "@/app/actions/bankAccounts";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

export default function AccountManagementPage() {
  const { user } = useAuth();
  const [accounts, setAccounts] = useState<BankAccountResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [supplierId, setSupplierId] = useState<number | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<BankAccountResponse | null>(null);

  const [formData, setFormData] = useState<BankAccountRequest>({
    supplier_id: 0,
    bank_name: "",
    account_number: "",
    account_holder_name: "",
    branch_name: "",
    account_type: "checking",
    is_primary: false,
  });

  // Fetch supplier ID from user email
  useEffect(() => {
    const fetchSupplierId = async () => {
      if (!user?.email) {
        setLoading(false);
        return;
      }

      try {
        // Get supplier list - backend automatically filters by authenticated user if they're a supplier
        const response = await axios.get(`/supplier`);
        const suppliers = response.data;
        
        // Handle different response formats
        let supplier = null;
        if (suppliers?.items && Array.isArray(suppliers.items) && suppliers.items.length > 0) {
          supplier = suppliers.items[0];
        } else if (Array.isArray(suppliers) && suppliers.length > 0) {
          supplier = suppliers[0];
        } else if (suppliers?.id) {
          supplier = suppliers;
        }

        if (supplier?.id) {
          setSupplierId(supplier.id);
          setFormData(prev => ({ ...prev, supplier_id: supplier.id }));
        } else {
          toast.error("Supplier account not found. Please contact support.");
        }
      } catch (error: any) {
        console.error("Error fetching supplier:", error);
        toast.error(error.response?.data?.error || "Failed to fetch supplier information");
      } finally {
        setLoading(false);
      }
    };

    fetchSupplierId();
  }, [user, toast]);

  // Fetch bank accounts when supplier ID is available
  useEffect(() => {
    if (supplierId) {
      fetchAccounts();
    }
  }, [supplierId]);

  const fetchAccounts = async () => {
    if (!supplierId) return;

    try {
      setLoading(true);
      const accountsList = await listBankAccounts(supplierId);
      setAccounts(accountsList);
    } catch (error: any) {
      console.error("Error fetching bank accounts:", error);
      toast.error(error.message || "Failed to fetch bank accounts");
    } finally {
      setLoading(false);
    }
  };

  const handleOpenDialog = (account?: BankAccountResponse) => {
    if (account) {
      setEditingAccount(account);
      setFormData({
        supplier_id: account.supplier_id,
        bank_name: account.bank_name,
        account_number: account.account_number,
        account_holder_name: account.account_holder_name,
        branch_name: account.branch_name || "",
        account_type: account.account_type,
        is_primary: account.is_primary,
      });
    } else {
      setEditingAccount(null);
      setFormData({
        supplier_id: supplierId || 0,
        bank_name: "",
        account_number: "",
        account_holder_name: "",
        branch_name: "",
        account_type: "checking",
        is_primary: false,
      });
    }
    setIsDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    setEditingAccount(null);
    setFormData({
      supplier_id: supplierId || 0,
      bank_name: "",
      account_number: "",
      account_holder_name: "",
      branch_name: "",
      account_type: "checking",
      is_primary: false,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!supplierId) {
      toast.error("Supplier ID not found");
      return;
    }

    try {
      if (editingAccount) {
        await updateBankAccount(editingAccount.id, formData);
        toast.success("Bank account updated successfully");
      } else {
        await createBankAccount(formData);
        toast.success("Bank account created successfully");
      }
      handleCloseDialog();
      fetchAccounts();
    } catch (error: any) {
      toast.error(error.message || "Failed to save bank account");
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this bank account?")) {
      return;
    }

    try {
      await deleteBankAccount(id);
      toast.success("Bank account deleted successfully");
      fetchAccounts();
    } catch (error: any) {
      toast.error(error.message || "Failed to delete bank account");
    }
  };

  if (loading && !supplierId) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!supplierId) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Supplier Account Not Found</CardTitle>
            <CardDescription>
              Your supplier account could not be found. Please contact support.
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Account Management</h1>
          <p className="text-muted-foreground">
            Manage your bank accounts for receiving payments
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={() => handleOpenDialog()}>
              <Plus className="h-4 w-4 mr-2" />
              Add Bank Account
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>
                {editingAccount ? "Edit Bank Account" : "Add Bank Account"}
              </DialogTitle>
              <DialogDescription>
                {editingAccount
                  ? "Update your bank account information"
                  : "Add a new bank account to receive payments"}
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit}>
              <div className="space-y-4 py-4">
                <div className="space-y-2">
                  <Label htmlFor="bank_name">Bank Name *</Label>
                  <Input
                    id="bank_name"
                    value={formData.bank_name}
                    onChange={(e) =>
                      setFormData({ ...formData, bank_name: e.target.value })
                    }
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="account_number">Account Number *</Label>
                  <Input
                    id="account_number"
                    value={formData.account_number}
                    onChange={(e) =>
                      setFormData({ ...formData, account_number: e.target.value })
                    }
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="account_holder_name">Account Holder Name *</Label>
                  <Input
                    id="account_holder_name"
                    value={formData.account_holder_name}
                    onChange={(e) =>
                      setFormData({ ...formData, account_holder_name: e.target.value })
                    }
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="branch_name">Branch Name</Label>
                  <Input
                    id="branch_name"
                    value={formData.branch_name}
                    onChange={(e) =>
                      setFormData({ ...formData, branch_name: e.target.value })
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="account_type">Account Type *</Label>
                  <Select
                    value={formData.account_type}
                    onValueChange={(value) =>
                      setFormData({ ...formData, account_type: value })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select account type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="checking">Checking</SelectItem>
                      <SelectItem value="savings">Savings</SelectItem>
                      <SelectItem value="current">Current</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="is_primary"
                    checked={formData.is_primary}
                    onChange={(e) =>
                      setFormData({ ...formData, is_primary: e.target.checked })
                    }
                    className="h-4 w-4 rounded border-gray-300"
                  />
                  <Label htmlFor="is_primary" className="cursor-pointer">
                    Set as primary account
                  </Label>
                </div>
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={handleCloseDialog}>
                  Cancel
                </Button>
                <Button type="submit">Save</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <CreditCard className="h-5 w-5" />
            Bank Accounts
          </CardTitle>
          <CardDescription>
            Your registered bank accounts for receiving payments
          </CardDescription>
        </CardHeader>
        <CardContent>
          {accounts.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              <CreditCard className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>No bank accounts registered yet.</p>
              <p className="text-sm mt-2">Click "Add Bank Account" to get started.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Bank Name</TableHead>
                  <TableHead>Account Number</TableHead>
                  <TableHead>Account Holder</TableHead>
                  <TableHead>Branch</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {accounts.map((account) => (
                  <TableRow key={account.id}>
                    <TableCell className="font-medium">{account.bank_name}</TableCell>
                    <TableCell>{account.account_number}</TableCell>
                    <TableCell>{account.account_holder_name}</TableCell>
                    <TableCell>{account.branch_name || "N/A"}</TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {account.account_type.charAt(0).toUpperCase() + account.account_type.slice(1)}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {account.is_primary && (
                          <Badge className="bg-primary/10 text-primary border-primary/20">
                            Primary
                          </Badge>
                        )}
                        {account.is_active ? (
                          <Badge className="bg-emerald-500/10 text-emerald-600 border-emerald-500/20">
                            Active
                          </Badge>
                        ) : (
                          <Badge className="bg-gray-500/10 text-gray-600 border-gray-500/20">
                            Inactive
                          </Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleOpenDialog(account)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDelete(account.id)}
                          className="text-destructive hover:text-destructive"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

