"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";
import ProductAttributeForm from "./attribute-form";
import ProductImagesForm from "@/components/product-image-form";
import { MultiSelect } from "@/components/multiselect";
import useProductsStore from "@/app/libs/store/useProductStore";
import useCategoryStore from "@/app/libs/store/useCategories";
import useDistributorsStore from "@/app/libs/store/useDistributorStore";

interface ProductFormProps {
  initialData?: Partial<any>;
  isEdit?: boolean;
  onSuccess?: () => void;
}

const ProductForm: React.FC<ProductFormProps> = ({
  initialData,
  isEdit = false,
  onSuccess,
}) => {
  console.log("Initial Data_________");
  console.log(initialData);
  const [product, setProduct] = useState(initialData || {});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const {
    loading,
    success,
    error,
    createProduct,
    updateProduct,
    fetchProductDetail,
  } = useProductsStore();
  const { categories, fetchCategories } = useCategoryStore();

  const { distributors, fetchDistributors } = useDistributorsStore();

  useEffect(() => {
    fetchDistributors();
  }, []);
  useEffect(() => {
    fetchCategories();
    if (isEdit && initialData?.Id) {
      fetchProductDetail(initialData.Id);
    }
  }, [isEdit, initialData?.Id]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setProduct((prev) => ({
      ...prev,
      [name]: name === "Price" ? parseFloat(value) || 0 : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const attributesObject =
        product.Attributes?.reduce(
          (
            acc: { [x: string]: any },
            attr: { key: string | number; value: any }
          ) => {
            acc[attr.key] = attr.value;
            return acc;
          },
          {} as Record<string, string>
        ) || {};
      const productData = {
        ...product,
        Id: product.Id,
        Name: product.Name || "",
        ExternalID: product.ExternalID || "",
        Desc: product.Desc || "",
        Price: product.Price || 0,
        Attributes: attributesObject || [],
        Images: product.Images || [],
        CategoryId: product.CategoryId?.map(Number) || [],
      };

      if (isEdit) {
        await updateProduct(productData);
      } else {
        await createProduct(productData);
      }

      onSuccess?.();
    } catch (error: any) {
      toast.error(error.message || "An error occurred");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleReset = () => {
    setProduct(initialData || {});
    toast.info("Form reset");
  };

  const categoryOptions = categories.map((c) => ({
    value: c.id.toString(),
    label: c.name,
  }));

  useEffect(() => {
    if (initialData?.Images) {
      const images = initialData?.Images.map((i: { ImageUrl: string }) => {
        return i.ImageUrl;
      });
      initialData.Images = images;
    }
    setProduct(initialData || {});
  }, [initialData]);
  useEffect(() => {
    if (success != null) {
      toast.success(success);
    }
  }, [success]);

  useEffect(() => {
    if (error != null) {
      toast.error(error);
    }
  }, [error]);
  return (
    <form onSubmit={handleSubmit}>
      <div className="space-y-6">
        <Card className="rounded-sm border-2 border-gray-200 shadow-none">
          <CardContent className="">
            <div className="mb-4">
              <h3 className="text-xl font-bold">Basic Information</h3>
              <p className="text-sm text-muted-foreground">
                {isEdit ? "Update product details" : "Add a new product"}
              </p>
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-1 gap-6">
                <div className="space-y-4">
                  <div className="grid gap-2">
                    <Label htmlFor="ExternalID">External ID</Label>
                    <Input
                      id="ExternalID"
                      name="ExternalID"
                      value={product.ExternalID || ""}
                      onChange={handleChange}
                      placeholder="e.g. 12345-abcde"
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="Name">Product Name</Label>
                    <Input
                      id="Name"
                      name="Name"
                      value={product.Name}
                      onChange={handleChange}
                      placeholder="e.g. POLO black - xl"
                      required
                    />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="Desc">Description</Label>
                    <Textarea
                      id="Desc"
                      name="Desc"
                      value={product.Desc || ""}
                      onChange={handleChange}
                      placeholder="Enter product description"
                      rows={5}
                    />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="Price">Price</Label>
                    <Input
                      id="Price"
                      name="Price"
                      type="number"
                      step="0.01"
                      value={product.Price || ""}
                      onChange={handleChange}
                      placeholder="29.99"
                      required
                    />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="Distributor">Distributor</Label>
                    <Select>
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Select a distributor" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          <SelectLabel>Distributor</SelectLabel>
                          {distributors?.map((d) => (
                            <SelectItem value={d.id.toString()}>
                              {d.name}
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="grid gap-2">
                    <Label>Categories</Label>
                    <MultiSelect
                      options={categoryOptions}
                      onValueChange={(values) =>
                        setProduct((prev) => ({
                          ...prev,
                          CategoryId: values.map(Number),
                        }))
                      }
                    />
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <ProductAttributeForm
          attributes={product.Attributes || []}
          onChange={(attributes) =>
            setProduct((prev) => ({ ...prev, Attributes: attributes }))
          }
        />

        <ProductImagesForm
          images={product.Images || []}
          onChange={(images) =>
            setProduct((prev) => ({
              ...prev,
              Images: images,
            }))
          }
        />

        <div className="flex items-center justify-end space-x-4">
          <Button
            type="button"
            variant="outline"
            onClick={handleReset}
            disabled={isSubmitting}
          >
            Reset
          </Button>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting
              ? isEdit
                ? "Updating..."
                : "Creating..."
              : isEdit
              ? "Update Product"
              : "Create Product"}
          </Button>
        </div>
      </div>
    </form>
  );
};

export default ProductForm;
