'use client'

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";
import ProductImagesForm from "./ProductImageForm";
import { MultiSelect } from "@/components/multiselect";
import useProductsStore from "@/app/libs/store/useProductStore";
import ConfigurableProductAttributeForm from "./ConfigurableAttributeForm";
import { ConfigurableProduct, ConfigurableProductForm } from "../app/libs/types";

interface ProductFormProps {
    initialData?: Partial<ConfigurableProduct>;
    isEdit?: boolean;
    loading?: boolean;
    onSuccess?: () => void;
}

const ConfigurableProductFormComponent: React.FC<ProductFormProps> = ({
    initialData,
    loading = false,
    isEdit = false,
    onSuccess,
}) => {
    const [product, setProduct] = useState<Partial<ConfigurableProductForm>>();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const {
        products,
        createConfigurableProduct,
        fetchProducts,
        fetchProductDetail,
    } = useProductsStore();

    useEffect(() => {
        fetchProducts();
        if (isEdit && initialData?.Id) {
            fetchProductDetail(initialData.Id);
        }
    }, [isEdit, initialData?.Id]);

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
    ) => {
        const { name, value } = e.target;
        setProduct((prev: any) => ({
            ...prev,
            [name]: value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);

        try {

            if (isEdit) {
                // await updateCon(productData);
                toast.success("Product updated successfully!");
            } else {
                await createConfigurableProduct(product);
                toast.success("Product created successfully!");
            }

            onSuccess?.();
        } catch (error: any) {
            toast.error(error.message || "An error occurred");
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleReset = () => {
        const p: ConfigurableProductForm = {
            id: initialData?.Id || 0,
            name: initialData?.Name || "",
            desc: initialData?.Desc || "",
            external_id: initialData?.Desc || "",
            images: initialData?.Images?.map((i: { ImageUrl: string; }) => { return i.ImageUrl }) || [],
            attribute_keys: initialData?.Attributes || [],
            products: initialData?.Products || []
        }
        setProduct(p);
        toast.info("Form reset");
    };

    const productOptions = products.map((c) => ({
        value: c.Id.toString(),
        label: c.Name,
    }));

    useEffect(() => {
        if (!loading) {
            const p: ConfigurableProductForm = {
                id: initialData?.Id || 0,
                name: initialData?.Name || "",
                desc: initialData?.Desc || "",
                external_id: initialData?.Desc || "",
                images: initialData?.Images?.map((i: { ImageUrl: string; }) => { return i.ImageUrl }) || [],
                attribute_keys: initialData?.Attributes || [],
                products: initialData?.Products || []
            }
            setProduct(p)
        }
    }, [initialData, loading])

    useEffect(() => {
        fetchProducts()
    }, [])

    return (
        <form onSubmit={handleSubmit}>
            <div className="space-y-6">
                <Card className="rounded-sm border-2 border-gray-200 shadow-none">
                    <CardContent className="">
                        <div className="mb-4">
                            <h3 className="text-xl font-bold">
                                Basic Information
                            </h3>
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
                                            value={product?.external_id}
                                            onChange={(e) => {
                                                setProduct((prev: any) => ({
                                                    ...prev,
                                                    ExternalId: e.target.value,
                                                }));
                                            }}
                                            placeholder="e.g. 12345-abcde"
                                        />
                                    </div>
                                    <div className="grid gap-2">
                                        <Label htmlFor="Name">Product Name</Label>
                                        <Input
                                            id="Name"
                                            name="Name"
                                            value={product?.name}
                                            onChange={(e) => {
                                                setProduct((prev: any) => ({
                                                    ...prev,
                                                    Name: e.target.value,
                                                }));
                                            }}
                                            placeholder="e.g. POLO black - xl"
                                            required
                                        />
                                    </div>

                                    <div className="grid gap-2">
                                        <Label htmlFor="Desc">Description</Label>
                                        <Textarea
                                            id="Desc"
                                            name="Desc"
                                            value={product?.desc || ""}
                                            onChange={(e) => {
                                                setProduct((prev: any) => ({
                                                    ...prev,
                                                    desc: e.target.value,
                                                }));
                                            }}
                                            placeholder="Enter product description"
                                            rows={5}
                                        />
                                    </div>


                                    <div className="grid gap-2">
                                        <Label>Products</Label>
                                        <MultiSelect
                                            options={productOptions}
                                            onValueChange={(values) =>
                                                setProduct((prev: any) => ({
                                                    ...prev,
                                                    Products: values.map(Number),
                                                }))
                                            }
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <ConfigurableProductAttributeForm
                    attributes={product?.attribute_keys || []}
                    onChange={(attributes) =>
                        setProduct((prev: any) => ({ ...prev, attribute_keys: attributes }))
                    }
                />

                <ProductImagesForm
                    images={product?.images || []}
                    onChange={(images) => setProduct((prev: any) => ({
                        ...prev,
                        images: images
                    }))}
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

export default ConfigurableProductFormComponent;