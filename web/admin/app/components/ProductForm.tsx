'use client'

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";
import ProductAttributeForm from "./AttributeForm";
import ProductImagesForm from "./ProductImageForm";
import { MultiSelect } from "@/app/components/multiselect";
import useProductsStore from "@/app/libs/store/useProductStore";
import useCategoryStore from "@/app/libs/store/useCategories";
import { useRouter } from "next/navigation";

interface Product {
    Id?: number;
    Name: string;
    ExternalID: string;
    Desc: string;
    Price: number;
    Attributes: Array<{ key: string; value: string }>;
    Images: string[];
    CategoryId: number[];
}

interface ProductFormProps {
    initialData?: Partial<Product>;
    isEdit?: boolean;
    onSuccess?: () => void;
}

const ProductForm: React.FC<ProductFormProps> = ({
    initialData,
    isEdit = false,
    onSuccess,
}) => {
    const defaultProduct: Product = {
        Name: '',
        ExternalID: '',
        Desc: '',
        Price: 0,
        Attributes: [],
        Images: [],
        CategoryId: [],
        ...initialData
    };

    const [product, setProduct] = useState<Product>(defaultProduct);
    const router = useRouter()
    const [isSubmitting, setIsSubmitting] = useState(false);
    const { error, loading, createProduct, updateProduct, fetchProductDetail } = useProductsStore();
    const { categories, fetchCategories } = useCategoryStore();
    const [submitted, setSubmitted] = useState(false)

    useEffect(() => {
        fetchCategories();
        if (isEdit && initialData?.Id) {
            fetchProductDetail(initialData.Id);
        }
    }, [isEdit, initialData?.Id, fetchCategories, fetchProductDetail]);

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
    ) => {
        const { name, value } = e.target;
        setProduct(prev => ({
            ...prev,
            [name]: name === "Price" ? parseFloat(value) || 0 : value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);

        try {
            const attributesObject = (product.Attributes || []).reduce((acc, attr) => {
                acc[attr.key] = attr.value;
                return acc;
            }, {} as Record<string, string>);

            const productData = {
                ...product,
                Name: product.Name.trim(),
                ExternalID: product.ExternalID.trim(),
                Desc: product.Desc.trim(),
                CategoryId: product.CategoryId.filter(id => !isNaN(id)),
                Attributes: attributesObject
            };

            if (isEdit) {
                await updateProduct(productData);
            } else {
                await createProduct(productData);
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
        setProduct(defaultProduct);
        toast.info("Form reset");
    };

    const categoryOptions = categories.map(c => ({
        value: c.id.toString(),
        label: c.name,
    }));

    return (
        <form onSubmit={handleSubmit}>
            <div className="space-y-6">
                <BasicInfoCard
                    product={product}
                    handleChange={handleChange}
                    categoryOptions={categoryOptions}
                    setProduct={setProduct}
                />

                <ProductAttributeForm
                    attributes={product.Attributes}
                    onChange={attributes => setProduct(prev => ({ ...prev, Attributes: attributes }))}
                />

                <ProductImagesForm
                    images={product.Images}
                    onChange={images => setProduct(prev => ({ ...prev, Images: images }))}
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

const BasicInfoCard: React.FC<{
    product: Product;
    handleChange: React.ChangeEventHandler<HTMLInputElement | HTMLTextAreaElement>;
    categoryOptions: Array<{ value: string; label: string }>;
    setProduct: React.Dispatch<React.SetStateAction<Product>>;
}> = ({ product, handleChange, categoryOptions, setProduct }) => (
    <Card className="rounded-sm border-2 border-gray-200 shadow-none">
        <CardContent>
            <div className="mb-4">
                <h3 className="text-xl font-bold">Basic Information</h3>
                <p className="text-sm text-muted-foreground">Product details</p>
            </div>

            <div className="space-y-4">
                <div className="grid grid-cols-1 gap-6">
                    <div className="space-y-4">
                        <InputField
                            label="External ID"
                            id="ExternalID"
                            name="ExternalID"
                            value={product.ExternalID}
                            onChange={handleChange}
                            placeholder="e.g. 12345-abcde"
                        />

                        <InputField
                            label="Product Name"
                            id="Name"
                            name="Name"
                            value={product.Name}
                            onChange={handleChange}
                            placeholder="e.g. POLO black - xl"
                            required
                        />

                        <div className="grid gap-2">
                            <Label htmlFor="Desc">Description</Label>
                            <Textarea
                                id="Desc"
                                name="Desc"
                                value={product.Desc}
                                onChange={handleChange}
                                placeholder="Enter product description"
                                rows={5}
                            />
                        </div>

                        <InputField
                            label="Price"
                            id="Price"
                            name="Price"
                            type="number"
                            step="0.01"
                            value={product.Price.toString()}
                            onChange={handleChange}
                            placeholder="29.99"
                            required
                        />

                        <div className="grid gap-2">
                            <Label>Categories</Label>
                            <MultiSelect
                                options={categoryOptions}
                                onValueChange={values => setProduct(prev => ({
                                    ...prev,
                                    CategoryId: values.map(v => Number(v))
                                }))}
                            />
                        </div>
                    </div>
                </div>
            </div>
        </CardContent>
    </Card>
);

const InputField: React.FC<{
    label: string;
    id: string;
    name: string;
    value: string;
    onChange: React.ChangeEventHandler<HTMLInputElement>;
    placeholder?: string;
    type?: string;
    required?: boolean;
    step?: string;
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

export default ProductForm;