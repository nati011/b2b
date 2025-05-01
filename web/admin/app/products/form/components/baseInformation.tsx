'use client'
import React from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { ArrowLeft, ArrowRight } from "lucide-react";

interface ProductBasicInfo {
    name: string;
    description: string;
    price?: string;
}

interface BasicInfoFormProps {
    productType: string;
    formData: ProductBasicInfo;
    onChange: (data: Partial<ProductBasicInfo>) => void;
    onBack: () => void;

    onNext: () => void;
    errors: {
        name?: string;
        price?: string;
    };
}

const BasicInfoForm: React.FC<BasicInfoFormProps> = ({
    productType,
    formData,
    onChange,
    onBack,
    onNext,
    errors,
}) => {
    const handleFormSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onNext();
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        onChange({ [name]: value });
    };

    return (
        <form onSubmit={handleFormSubmit} className="space-y-6">
            <div className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="name">
                        Product Name <span className="text-destructive">*</span>
                    </Label>
                    <Input
                        id="name"
                        name="name"
                        value={formData.name}
                        onChange={handleInputChange}
                        placeholder="Enter product name"
                        className={errors.name ? "border-destructive" : ""}
                    />
                    {errors.name && (
                        <p className="text-sm text-destructive">{errors.name}</p>
                    )}
                </div>

                <div className="space-y-2">
                    <Label htmlFor="description">Product Description</Label>
                    <Textarea
                        id="description"
                        name="description"
                        value={formData.description}
                        onChange={handleInputChange}
                        placeholder="Describe your product"
                        className="min-h-[120px]"
                    />
                </div>

                {productType === "simple" && (
                    <div className="space-y-2">
                        <Label htmlFor="price">
                            Price <span className="text-destructive">*</span>
                        </Label>
                        <div className="relative">
                            <span className="absolute left-3 top-2.5 text-muted-foreground">
                                $
                            </span>
                            <Input
                                id="price"
                                name="price"
                                type="text"
                                value={formData.price || ""}
                                onChange={handleInputChange}
                                placeholder="0.00"
                                className={`pl-7 ${errors.price ? "border-destructive" : ""}`}
                            />
                        </div>
                        {errors.price && (
                            <p className="text-sm text-destructive">{errors.price}</p>
                        )}
                    </div>
                )}
            </div>

            <div className="flex justify-between">
                <Button type="button" variant="outline" onClick={onBack} className="flex items-center gap-2">
                    <ArrowLeft className="w-4 h-4" /> Back
                </Button>
                <Button type="submit" className="flex items-center gap-2">
                    Continue <ArrowRight className="w-4 h-4" />
                </Button>
            </div>
        </form>
    );
};

export default BasicInfoForm;