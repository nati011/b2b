'use client'
import React, { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner"
import BasicInfoForm from "./baseInformation";
import CategoriesForm, { Category } from "./categoriesForm";
import AttributesForm from "./attributesForm";
import ImageUploader from './imageUploader';
import ConfigurableProducts from "./configurableProduct";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Product } from "@/app/libs/types";

interface ProductFormData {
    productType: string;
    basicInfo: {
        name: string;
        description: string;
        price?: string;
    };
    categories: number[];
    attributes: { key: string; value: string }[];
    attributeKeys: string[];
    images: string[];
    selectedProducts: number[];
}

interface ProductWizardProps {
    categories: Category[];
    products: Product[];
    onCreateCategory: (name: string) => Promise<void>;
    onDeleteCategory: (id: number) => Promise<void>;
    onCreateProduct: (data: any) => Promise<void>;
    onCreateConfigurableProduct: (data: any) => Promise<void>;
}

const ProductWizard: React.FC<ProductWizardProps> = ({
    categories,
    products,
    onCreateCategory,
    onDeleteCategory,
    onCreateProduct,
    onCreateConfigurableProduct
}) => {
    const [formData, setFormData] = useState<ProductFormData>({
        productType: 'simple',
        basicInfo: {
            name: "",
            description: "",
            price: ""
        },
        categories: [],
        attributes: [],
        attributeKeys: [],
        images: [],
        selectedProducts: []
    });

    const [errors, setErrors] = useState({
        name: "",
        price: "",
        attributeKeys: ""
    });

    const [currentStep, setCurrentStep] = useState(0);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (formData.productType === "simple") {
            setFormData(prev => ({
                ...prev,
                attributeKeys: [],
                selectedProducts: []
            }));
        } else if (formData.productType === "configurable") {
            setFormData(prev => ({
                ...prev,
                attributes: [],
                basicInfo: {
                    ...prev.basicInfo,
                    price: ""
                }
            }));
        }
    }, [formData.productType]);

    const buildSteps = () => {
        if (!formData.productType) return ["Product Type"];

        if (formData.productType === "simple") {
            return [
                "Basic Info",
                "Categories",
                "Attributes",
                "Images",
                "Review"
            ];
        } else {
            return [
                "Basic Info",
                "Categories",
                "Variation Keys",
                "Products",
                "Images",
                "Review"
            ];
        }
    };

    const steps = buildSteps();

    const validateCurrentStep = (): boolean => {
        let isValid = true;
        const newErrors = { name: "", price: "", attributeKeys: "" };

        if (currentStep === 1) {
            if (!formData.basicInfo.name.trim()) {
                newErrors.name = "Product name is required";
                isValid = false;
            }

            if (formData.productType === "simple" && !formData.basicInfo.price?.trim()) {
                newErrors.price = "Price is required";
                isValid = false;
            }
        }

        if (currentStep === 3 && formData.productType === "configurable" && formData.attributeKeys.length === 0) {
            newErrors.attributeKeys = "At least one variation key is required";
            isValid = false;
        }

        setErrors(newErrors);
        return isValid;
    };

    const goToNext = () => {
        if (validateCurrentStep()) {
            setCurrentStep(prev => Math.min(prev + 1, steps.length - 1));
        }
    };

    const goToPrevious = () => {
        setCurrentStep(prev => Math.max(prev - 1, 0));
    };

    const updateFormData = <K extends keyof ProductFormData>(key: K, value: ProductFormData[K]) => {
        setFormData(prev => ({
            ...prev,
            [key]: value
        }));
    };

    const updateBasicInfo = (data: Partial<ProductFormData["basicInfo"]>) => {
        setFormData(prev => ({
            ...prev,
            basicInfo: {
                ...prev.basicInfo,
                ...data
            }
        }));

        if (data.name) setErrors(prev => ({ ...prev, name: "" }));
        if (data.price) setErrors(prev => ({ ...prev, price: "" }));
    };

    const handleSubmit = async () => {
        setIsSubmitting(true);

        try {
            if (formData.productType === "simple") {
                await onCreateProduct({
                    name: formData.basicInfo.name,
                    desc: formData.basicInfo.description,
                    price: parseFloat(formData.basicInfo.price || "0"),
                    images: formData.images,
                    attributes: Object.fromEntries(
                        formData.attributes.map(attr => [attr.key, attr.value])
                    ),
                    category_id: formData.categories,
                    external_id: `product-${Date.now()}`,
                    distributor_id: 1,
                });
            } else {
                await onCreateConfigurableProduct({
                    name: formData.basicInfo.name,
                    desc: formData.basicInfo.description,
                    attribute_keys: formData.attributeKeys,
                    products: formData.selectedProducts,
                    images: formData.images,
                    external_id: `configurable-${Date.now()}`
                });
            }

            toast("Product created successfully", {
                description: `Your ${formData.productType} product has been created.`,
            });

            // // Reset form
            // setFormData({
            //     productType: 'simple',
            //     basicInfo: {
            //         name: "",
            //         description: "",
            //         price: ""
            //     },
            //     categories: [],
            //     attributes: [],
            //     attributeKeys: [],
            //     images: [],
            //     selectedProducts: []
            // });
            // setCurrentStep(0);

        } catch (error) {
            console.error("Error creating product:", error);
            toast("Error creating product", {
                description: "There was a problem creating your product.",
            });
        } finally {
            setIsSubmitting(false);
        }
    };

    const renderStepContent = () => {
        switch (steps[currentStep]) {
            case "Basic Info":
                return (
                    <BasicInfoForm
                        productType={formData.productType}
                        formData={formData.basicInfo}
                        onChange={updateBasicInfo}
                        onBack={goToPrevious}
                        onNext={goToNext}
                        errors={errors}
                    />
                );
            case "Categories":
                return (
                    <CategoriesForm
                        categories={categories}
                        selectedCategories={formData.categories}
                        onSelectCategories={(categoryIds) => updateFormData("categories", categoryIds)}
                        onAddCategory={onCreateCategory}
                        onDeleteCategory={onDeleteCategory}
                        onBack={goToPrevious}
                        onNext={goToNext}
                    />
                );
            case "Attributes":
                return (
                    <AttributesForm
                        productType="simple"
                        attributes={formData.attributes}
                        onUpdateAttributes={(attrs) => updateFormData("attributes", attrs)}
                        attributeKeys={formData.attributeKeys}
                        onUpdateAttributeKeys={(keys) => updateFormData("attributeKeys", keys)}
                        onBack={goToPrevious}
                        onNext={goToNext}
                    />
                );
            case "Variation Keys":
                return (
                    <AttributesForm
                        productType="configurable"
                        attributes={formData.attributes}
                        onUpdateAttributes={(attrs) => updateFormData("attributes", attrs)}
                        attributeKeys={formData.attributeKeys}
                        onUpdateAttributeKeys={(keys) => updateFormData("attributeKeys", keys)}
                        onBack={goToPrevious}
                        onNext={goToNext}
                    />
                );
            case "Products":
                return (
                    <ConfigurableProducts
                        products={products}
                        selectedProductIds={formData.selectedProducts}
                        onUpdateSelectedProducts={(productIds) => updateFormData("selectedProducts", productIds)}
                        attributeKeys={formData.attributeKeys}
                        onBack={goToPrevious}
                        onNext={goToNext}
                    />
                );
            case "Images":
                return (
                    <ImageUploader
                        images={formData.images}
                        onUpdateImages={(images) => updateFormData("images", images)}
                        onBack={goToPrevious}
                        onNext={goToNext}
                    />
                );
            case "Review":
                return (
                    <div className="space-y-4">
                        <div>
                            <h3 className="font-medium">Product Type</h3>
                            <p>{formData.productType}</p>
                        </div>
                        <div>
                            <h3 className="font-medium">Basic Information</h3>
                            <p>Name: {formData.basicInfo.name}</p>
                            <p>Description: {formData.basicInfo.description}</p>
                            {formData.productType === "simple" && (
                                <p>Price: {formData.basicInfo.price}</p>
                            )}
                        </div>
                        <div>
                            <h3 className="font-medium">Categories</h3>
                            <ul>
                                {formData.categories.map(id => (
                                    <li key={id}>
                                        {categories.find(c => c.id === id)?.name || `Category ${id}`}
                                    </li>
                                ))}
                            </ul>
                        </div>
                        {formData.productType === "simple" ? (
                            <div>
                                <h3 className="font-medium">Attributes</h3>
                                <ul>
                                    {formData.attributes.map((attr, index) => (
                                        <li key={index}>
                                            {attr.key}: {attr.value}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        ) : (
                            <div>
                                <h3 className="font-medium">Variation Keys</h3>
                                <ul>
                                    {formData.attributeKeys.map((key, index) => (
                                        <li key={index}>{key}</li>
                                    ))}
                                </ul>
                                <h3 className="font-medium mt-4">Selected Products</h3>
                                <ul>
                                    {formData.selectedProducts.map(id => (
                                        <li key={id}>
                                            {products.find(p => p.Id === id)?.Name || `Product ${id}`}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        )}
                        <div>
                            <h3 className="font-medium">Images</h3>
                            <p>{formData.images.length} images selected</p>
                        </div>
                        <div className="flex justify-between mt-6">
                            <button
                                onClick={goToPrevious}
                                className="px-4 py-2 border rounded-md"
                            >
                                Back
                            </button>
                            <button
                                onClick={handleSubmit}
                                disabled={isSubmitting}
                                className="px-4 py-2 bg-primary text-primary-foreground rounded-md disabled:opacity-50"
                            >
                                {isSubmitting ? "Creating..." : "Create Product"}
                            </button>
                        </div>
                    </div>
                );
            default:
                return null;
        }
    };

    return (
        <div className="container py-8">
            <Tabs onValueChange={(type) => updateFormData("productType", type)} value={formData.productType} className="bg-transparent">
                <TabsList className="w-full">
                    <TabsTrigger value="simple">Simple Product</TabsTrigger>
                    <TabsTrigger value="configurable-product">Configurable Product</TabsTrigger>
                </TabsList>
                <TabsContent value="simple">
                    <Card>
                        <CardHeader>
                            <div className="">
                                <h1 className="text-3xl font-bold tracking-tight">Register Product</h1>
                                <p className="text-muted-foreground">
                                    Create and manage products
                                </p>
                            </div>
                            <hr className="mb-4" />
                            {/* <div className="space-y-2">
                                <div className="flex justify-between text-sm mb-1">
                                    <span>Progress</span>
                                    <span className="font-medium">{Math.round((currentStep / (steps.length - 1)) * 100)}%</span>
                                </div>
                                <div className="w-full bg-muted rounded-full h-2.5">
                                    <div
                                        className="bg-primary h-2.5 rounded-full transition-all"
                                        style={{ width: `${(currentStep / (steps.length - 1)) * 100}%` }}
                                    ></div>
                                </div>
                                <p className="text-xs text-muted-foreground mt-2">
                                    {
                                        `Step ${currentStep + 1} of ${steps.length}:  `}<strong className="text-black text-md">{steps[currentStep]}</strong>
                                </p>
                            </div> */}

                        </CardHeader>
                        <CardContent>
                            {renderStepContent()}</CardContent>
                    </Card>
                </TabsContent>
                <TabsContent value='configurable-product'>
                    <Card>
                        <CardHeader>
                            <div className="">
                                <h1 className="text-3xl font-bold tracking-tight">Register Product</h1>
                                <p className="text-muted-foreground">
                                    Create and manage products
                                </p>
                            </div>
                            <hr className="mb-4" />


                        </CardHeader>
                        <CardContent>
                            {renderStepContent()}</CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
};

export default ProductWizard;