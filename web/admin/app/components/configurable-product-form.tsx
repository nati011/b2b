'use client'

import React, { useEffect, useState, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";
import ProductImagesForm from "./ProductImageForm";
import { MultiSelect } from "@/app/components/multiselect";
import useProductsStore from "@/app/libs/store/useProductStore";
import ConfigurableProductAttributeForm from "./ConfigurableAttributeForm";
import { ConfigurableProduct, ConfigurableProductRequest } from "../libs/types";
import { useRouter } from "next/navigation";

interface ProductFormProps {
    initialData?: Partial<ConfigurableProduct>;
    isEdit?: boolean;
    loading?: boolean;
    onSuccess?: () => void;
}

const defaultProduct: ConfigurableProductRequest = {
    id: 0,
    name: "",
    desc: "",
    external_id: "",
    images: [],
    attribute_keys: [],
    products: []
};

const ConfigurableProductForm: React.FC<ProductFormProps> = ({
    initialData,
    isEdit = false,
    onSuccess,
}) => {
    const router = useRouter()
    const [product, setProduct] = useState<ConfigurableProductRequest>(defaultProduct);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { error, loading, products, createConfigurableProduct, fetchProducts, fetchProductDetail } = useProductsStore();
    const [submitted, setSubmitted] = useState(false)

    const initializeProduct = useCallback(() => {
        const initialProduct = initialData ? {
            id: initialData.Id || 0,
            name: initialData.Name || "",
            desc: initialData.Desc || "",
            external_id: initialData.ExternalId || "",
            images: initialData.Images?.map(i => i.ImageUrl) || [],
            attribute_keys: initialData.Attributes || [],
            products: initialData.Products || []
        } : defaultProduct;

        setProduct(initialProduct);
    }, [initialData]);

    useEffect(() => {
        fetchProducts();
        if (isEdit && initialData?.Id) {
            fetchProductDetail(initialData.Id);
        }
    }, [isEdit, initialData?.Id, fetchProducts, fetchProductDetail]);

    useEffect(() => {
        if (!loading) initializeProduct();
    }, [initialData, loading, initializeProduct]);

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
    ) => {
        const { name, value } = e.target;
        setProduct(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);

        try {
            if (isEdit) {
                // TODO: Implement update
                // await updateConfigurableProduct(product);
            } else {
                await createConfigurableProduct(product);

            }
            setSubmitted(true)
        } catch (error: unknown) {
            toast.error(error instanceof Error ? error.message : "An error occurred");
        } finally {
            setIsSubmitting(false);
        }
    };

    useEffect(() => {
        if (submitted && !loading) {
            if (error) {
                toast.error(error)
            } else {
                toast.success("Success!")
                handleReset()
                router.push('/products')
            }
            setSubmitted(false)
        }
    }, [loading, error, submitted])


    const handleReset = () => {
        initializeProduct();
        toast.info("Form reset");
    };

    const productOptions = products.map(product => ({
        value: product.Id.toString(),
        label: product.Name
    }));

    return (
        <form onSubmit={handleSubmit}>
            <div className="space-y-6">
                <Card className="rounded-sm border-2 border-gray-200 shadow-none">
                    <CardContent>
                        <div className="mb-4">
                            <h3 className="text-xl font-bold">Basic Information</h3>
                            <p className="text-sm text-muted-foreground">
                                {isEdit ? "Update product details" : "Add a new product"}
                            </p>
                        </div>

                        <div className="space-y-4">
                            <div className="grid grid-cols-1 gap-6">
                                <div className="space-y-4">
                                    <InputField
                                        label="External ID"
                                        id="external_id"
                                        name="external_id"
                                        value={product.external_id}
                                        onChange={handleChange}
                                        placeholder="e.g. 12345-abcde"
                                    />

                                    <InputField
                                        label="Product Name"
                                        id="name"
                                        name="name"
                                        value={product.name}
                                        onChange={handleChange}
                                        placeholder="e.g. POLO black - xl"
                                        required
                                    />

                                    <div className="grid gap-2">
                                        <Label htmlFor="desc">Description</Label>
                                        <Textarea
                                            id="desc"
                                            name="desc"
                                            value={product.desc}
                                            onChange={handleChange}
                                            placeholder="Enter product description"
                                            rows={5}
                                        />
                                    </div>

                                    <div className="grid gap-2">
                                        <Label>Associated Products</Label>
                                        <MultiSelect
                                            options={productOptions}
                                            onValueChange={values => setProduct(prev => ({
                                                ...prev,
                                                products: values.map(Number)
                                            }))}
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <ConfigurableProductAttributeForm
                    attributes={product.attribute_keys}
                    onChange={attributes => setProduct(prev => ({
                        ...prev,
                        attribute_keys: attributes
                    }))}
                />

                <ProductImagesForm
                    images={product.images}
                    onChange={images => setProduct(prev => ({
                        ...prev,
                        images
                    }))}
                />

                <FormActions
                    isSubmitting={isSubmitting}
                    isEdit={isEdit}
                    onReset={handleReset}
                />
            </div>
        </form>
    );
};

// Sub-components for better readability
const InputField: React.FC<{
    label: string;
    id: string;
    name: string;
    value: string;
    onChange: React.ChangeEventHandler<HTMLInputElement>;
    placeholder?: string;
    required?: boolean;
}> = ({ label, ...props }) => (
    <div className="grid gap-2">
        <Label htmlFor={props.id}>{label}</Label>
        <Input {...props} />
    </div>
);

const FormActions: React.FC<{
    isSubmitting: boolean;
    isEdit: boolean;
    onReset: () => void;
}> = ({ isSubmitting, isEdit, onReset }) => (
    <div className="flex items-center justify-end space-x-4">
        <Button
            type="button"
            variant="outline"
            onClick={onReset}
            disabled={isSubmitting}
        >
            Reset
        </Button>
        <Button type="submit" disabled={isSubmitting}>
            {isSubmitting
                ? isEdit ? "Updating..." : "Creating..."
                : isEdit ? "Update Product" : "Create Product"}
        </Button>
    </div>
);

export default ConfigurableProductForm;