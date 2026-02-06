"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { ArrowLeft, Save, Plus, Trash2 } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import { CreateProduct } from "@/app/actions/product";

export default function NewProductPage() {
  const router = useRouter();
  const [saving, setSaving] = useState(false);

  const [formData, setFormData] = useState({
    name: "",
    description: "",
    price: "",
    total_quantity: "",
    is_active: true,
    unit: "",
    external_id: "",
  });

  const [attributes, setAttributes] = useState<Array<{ key: string; value: string }>>([
    { key: "", value: "" },
  ]);

  const handleAddAttribute = () => {
    setAttributes([...attributes, { key: "", value: "" }]);
  };

  const handleRemoveAttribute = (index: number) => {
    if (attributes.length > 1) {
      setAttributes(attributes.filter((_, i) => i !== index));
    }
  };

  const handleAttributeChange = (index: number, field: "key" | "value", value: string) => {
    const updated = [...attributes];
    updated[index][field] = value;
    setAttributes(updated);
  };

  const buildAttributesObject = () => {
    const attrsObj: Record<string, string> = {};
    attributes.forEach((attr) => {
      if (attr.key.trim() && attr.value.trim()) {
        attrsObj[attr.key.trim()] = attr.value.trim();
      }
    });
    return Object.keys(attrsObj).length > 0 ? attrsObj : undefined;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setSaving(true);
      
      const attributesObj = buildAttributesObject();

      const createData: any = {
        name: formData.name,
        description: formData.description || undefined,
        price: formData.price ? parseFloat(formData.price) : undefined,
        total_quantity: formData.total_quantity ? parseInt(formData.total_quantity) : undefined,
        is_active: formData.is_active,
        unit: formData.unit || undefined,
        external_id: formData.external_id || undefined,
      };

      if (attributesObj) {
        createData.attributes = attributesObj;
      }

      await CreateProduct(createData);
      toast.success("Product created successfully");
      router.push("/supplier/products");
    } catch (err: any) {
      console.error("Error creating product:", err);
      toast.error(err.message || "Failed to create product");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="mb-8">
        <Button
          variant="ghost"
          asChild
          className="mb-4"
        >
          <Link href="/supplier/products">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Products
          </Link>
        </Button>
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-1 bg-primary rounded-full"></div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">New Product</h1>
        </div>
        <p className="text-gray-600 dark:text-gray-300 mt-2 ml-4">Create a new product</p>
      </div>

      <form onSubmit={handleSubmit}>
        <Card className="border-primary/20 hover:border-primary/40 transition-colors">
          <CardHeader>
            <CardTitle>Product Information</CardTitle>
            <CardDescription>Enter product details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid gap-6 md:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="name">Product Name *</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  required
                  placeholder="Enter product name"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="external_id">External ID</Label>
                <Input
                  id="external_id"
                  value={formData.external_id}
                  onChange={(e) => setFormData({ ...formData, external_id: e.target.value })}
                  placeholder="Optional external identifier"
                />
              </div>

              <div className="space-y-2 md:col-span-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  rows={4}
                  placeholder="Product description"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="price">Price (ETB)</Label>
                <Input
                  id="price"
                  type="number"
                  step="0.01"
                  min="0"
                  value={formData.price}
                  onChange={(e) => setFormData({ ...formData, price: e.target.value })}
                  placeholder="0.00"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="unit">Unit</Label>
                <Input
                  id="unit"
                  value={formData.unit}
                  onChange={(e) => setFormData({ ...formData, unit: e.target.value })}
                  placeholder="e.g., kg, pcs, box"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="total_quantity">Total Quantity</Label>
                <Input
                  id="total_quantity"
                  type="number"
                  min="0"
                  value={formData.total_quantity}
                  onChange={(e) => setFormData({ ...formData, total_quantity: e.target.value })}
                  placeholder="0"
                />
              </div>

              <div className="space-y-2 md:col-span-2">
                <div className="flex items-center justify-between">
                  <div>
                    <Label htmlFor="is_active">Active Status</Label>
                    <p className="text-sm text-muted-foreground">
                      Inactive products won't be visible to customers
                    </p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="is_active"
                      checked={formData.is_active}
                      onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked as boolean })}
                    />
                    <Label htmlFor="is_active" className="cursor-pointer">
                      {formData.is_active ? "Active" : "Inactive"}
                    </Label>
                  </div>
                </div>
              </div>
            </div>

            <div className="space-y-4 pt-4 border-t">
              <div className="flex items-center justify-between">
                <div>
                  <Label className="text-base font-semibold">Product Attributes</Label>
                  <p className="text-sm text-muted-foreground">
                    Add custom attributes as key-value pairs (e.g., Color: Red, Size: Large)
                  </p>
                </div>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={handleAddAttribute}
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add Attribute
                </Button>
              </div>

              <div className="space-y-3">
                {attributes.map((attr, index) => (
                  <div key={index} className="flex gap-3 items-start">
                    <div className="flex-1 space-y-2">
                      <Label htmlFor={`attr-key-${index}`} className="text-xs text-muted-foreground">
                        Key
                      </Label>
                      <Input
                        id={`attr-key-${index}`}
                        value={attr.key}
                        onChange={(e) => handleAttributeChange(index, "key", e.target.value)}
                        placeholder="e.g., Color, Size, Material"
                      />
                    </div>
                    <div className="flex-1 space-y-2">
                      <Label htmlFor={`attr-value-${index}`} className="text-xs text-muted-foreground">
                        Value
                      </Label>
                      <Input
                        id={`attr-value-${index}`}
                        value={attr.value}
                        onChange={(e) => handleAttributeChange(index, "value", e.target.value)}
                        placeholder="e.g., Red, Large, Cotton"
                      />
                    </div>
                    <div className="pt-7">
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        onClick={() => handleRemoveAttribute(index)}
                        disabled={attributes.length === 1}
                        className="text-destructive hover:text-destructive hover:bg-destructive/10"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>

              {attributes.length === 0 && (
                <div className="text-center py-4 text-sm text-muted-foreground">
                  No attributes added. Click "Add Attribute" to add custom product attributes.
                </div>
              )}
            </div>

            <div className="flex justify-end gap-4 pt-4">
              <Button
                type="button"
                variant="outline"
                asChild
              >
                <Link href="/supplier/products">Cancel</Link>
              </Button>
              <Button type="submit" disabled={saving}>
                <Save className="h-4 w-4 mr-2" />
                {saving ? "Creating..." : "Create Product"}
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  );
}




