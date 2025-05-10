'use client'
import ProductWizard from "@/app/products/form/components/productWizard";
import { toast } from "sonner"
import useProductsStore from "@/app/libs/store/useProductStore";
import { useEffect, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Product } from "@/app/libs/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import Heading from "@/app/components/breadcrumb";

const EditProduct = () => {
    const [formData, setFormData] = useState<Product>()
    const [errors, setErrors] = useState({
        name: "",
        price: "",
        attributeKeys: ""
    });

    const handleFormSubmit = () => {
        console.log("Submitted")
    }

    const pages = [
        {
            "title": "Products",
            "href": "/products"
        },
        {
            "title": "Product Form",
            "href": "/products/form"
        }
    ]

    return (
        <>
            <Heading page={pages} heading="Edit Selected Product" subheading="lorem ipsum lorem ipsum" />
            <Card className="shadow-none">
                <CardContent>
                    <form onSubmit={handleFormSubmit} className="space-y-6">
                        <div className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="name">
                                    Product Name <span className="text-destructive">*</span>
                                </Label>
                                <Input
                                    id="name"
                                    name="name"
                                    value={formData?.Name}
                                    // onChange={handleInputChange}
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
                                    value={formData?.Desc}
                                    // onChange={handleInputChange}
                                    placeholder="Describe your product"
                                    className="min-h-[120px]"
                                />
                            </div>
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
                                        value={formData?.Price || ""}
                                        // onChange={handleInputChange}
                                        placeholder="0.00"
                                        className={`pl-7 ${errors.price ? "border-destructive" : ""}`}
                                    />
                                </div>
                                {errors.price && (
                                    <p className="text-sm text-destructive">{errors.price}</p>
                                )}
                            </div>
                        </div>

                        <div className="flex justify-between">
                            <Button type="submit" className="flex items-center gap-2">
                                Submit
                            </Button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </>
    );
};

export default EditProduct;
