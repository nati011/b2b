import React from "react";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

export interface ProductFormData {
    productType: string;
    basicInfo: {
        name: string;
        description: string;
        price?: string;
    };
    categories: number[];
    categoryNames?: string[];
    attributes: { key: string; value: string }[];
    attributeKeys: string[];
    images: string[];
    selectedProducts: number[];
    selectedProductNames?: string[];
}

interface ReviewScreenProps {
    formData: ProductFormData;
    onBack: () => void;
    onSubmit: () => void;
    isSubmitting: boolean;
}

const ReviewScreen: React.FC<ReviewScreenProps> = ({
    formData,
    onBack,
    onSubmit,
    isSubmitting,
}) => {
    const displayPrice = formData.basicInfo.price
        ? `$${parseFloat(formData.basicInfo.price).toFixed(2)}`
        : "N/A";

    return (
        <div className="space-y-6">
            <div>
                <h3 className="text-lg font-medium">Review Product Information</h3>
                <p className="text-muted-foreground">
                    Please review the information before creating the product
                </p>
            </div>

            <Tabs defaultValue="summary">
                <TabsList className="grid w-full grid-cols-4">
                    <TabsTrigger value="summary">Summary</TabsTrigger>
                    <TabsTrigger value="details">Details</TabsTrigger>
                    <TabsTrigger value="attributes">Attributes</TabsTrigger>
                    <TabsTrigger value="images">Images</TabsTrigger>
                </TabsList>

                <TabsContent value="summary" className="space-y-4 pt-4">
                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-lg">Product Type</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <Badge className="capitalize">
                                {formData.productType === "simple" ? "Simple Product" : "Configurable Product"}
                            </Badge>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-lg">Basic Information</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-2">
                            <div>
                                <span className="font-medium">Name:</span> {formData.basicInfo.name}
                            </div>
                            {formData.productType === "simple" && (
                                <div>
                                    <span className="font-medium">Price:</span> {displayPrice}
                                </div>
                            )}
                            {formData.basicInfo.description && (
                                <div>
                                    <span className="font-medium">Description:</span>{" "}
                                    <span className="text-muted-foreground">
                                        {formData.basicInfo.description.length > 100
                                            ? `${formData.basicInfo.description.substring(0, 100)}...`
                                            : formData.basicInfo.description}
                                    </span>
                                </div>
                            )}
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader className="pb-2">
                            <CardTitle className="text-lg">
                                {formData.productType === "simple" ? "Categories & Attributes" : "Configurations"}
                            </CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-2">
                            <div className="flex flex-wrap gap-2">
                                <span className="font-medium mr-2">Categories:</span>
                                {formData.categoryNames && formData.categoryNames.length > 0 ? (
                                    formData.categoryNames.map((name, index) => (
                                        <Badge key={index} variant="outline">
                                            {name}
                                        </Badge>
                                    ))
                                ) : (
                                    <span className="text-muted-foreground">None selected</span>
                                )}
                            </div>

                            {formData.productType === "simple" ? (
                                <div>
                                    <span className="font-medium">Attributes:</span>{" "}
                                    {formData.attributes.length > 0 ? (
                                        <span>{formData.attributes.length} defined</span>
                                    ) : (
                                        <span className="text-muted-foreground">None defined</span>
                                    )}
                                </div>
                            ) : (
                                <>
                                    <div className="flex flex-wrap gap-2">
                                        <span className="font-medium mr-2">Variation Keys:</span>
                                        {formData.attributeKeys.length > 0 ? (
                                            formData.attributeKeys.map((key, index) => (
                                                <Badge key={index} variant="outline">
                                                    {key}
                                                </Badge>
                                            ))
                                        ) : (
                                            <span className="text-muted-foreground">None defined</span>
                                        )}
                                    </div>
                                    <div>
                                        <span className="font-medium">Selected Products:</span>{" "}
                                        {formData.selectedProductNames ? (
                                            <span>{formData.selectedProductNames.length} products</span>
                                        ) : (
                                            <span className="text-muted-foreground">None selected</span>
                                        )}
                                    </div>
                                </>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="details" className="space-y-4 pt-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Product Details</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div>
                                <h4 className="font-medium mb-1">Product Name</h4>
                                <p>{formData.basicInfo.name}</p>
                            </div>

                            {formData.productType === "simple" && (
                                <div>
                                    <h4 className="font-medium mb-1">Price</h4>
                                    <p>{displayPrice}</p>
                                </div>
                            )}

                            <div>
                                <h4 className="font-medium mb-1">Description</h4>
                                <p className="whitespace-pre-line">
                                    {formData.basicInfo.description || "No description provided"}
                                </p>
                            </div>

                            <div>
                                <h4 className="font-medium mb-1">Categories</h4>
                                <div className="flex flex-wrap gap-2">
                                    {formData.categoryNames && formData.categoryNames.length > 0 ? (
                                        formData.categoryNames.map((name, index) => (
                                            <Badge key={index} variant="outline">
                                                {name}
                                            </Badge>
                                        ))
                                    ) : (
                                        <p className="text-muted-foreground">No categories selected</p>
                                    )}
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="attributes" className="space-y-4 pt-4">
                    {formData.productType === "simple" ? (
                        <Card>
                            <CardHeader>
                                <CardTitle>Product Attributes</CardTitle>
                            </CardHeader>
                            <CardContent>
                                {formData.attributes.length > 0 ? (
                                    <div className="space-y-2">
                                        {formData.attributes.map((attr, index) => (
                                            <div key={index} className="flex p-2 border rounded-md">
                                                <div className="font-medium min-w-[100px]">{attr.key}:</div>
                                                <div>{attr.value}</div>
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <p className="text-muted-foreground">No attributes defined for this product</p>
                                )}
                            </CardContent>
                        </Card>
                    ) : (
                        <>
                            <Card>
                                <CardHeader>
                                    <CardTitle>Variation Attributes</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <div className="space-y-2">
                                        <h4 className="font-medium">Attribute Keys</h4>
                                        <div className="flex flex-wrap gap-2">
                                            {formData.attributeKeys.length > 0 ? (
                                                formData.attributeKeys.map((key, index) => (
                                                    <Badge key={index} variant="outline">
                                                        {key}
                                                    </Badge>
                                                ))
                                            ) : (
                                                <p className="text-muted-foreground">No attribute keys defined</p>
                                            )}
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader>
                                    <CardTitle>Selected Products</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    {formData.selectedProductNames && formData.selectedProductNames.length > 0 ? (
                                        <div className="space-y-1">
                                            {formData.selectedProductNames.map((name, index) => (
                                                <div key={index} className="p-2 border rounded-md">
                                                    {name}
                                                </div>
                                            ))}
                                        </div>
                                    ) : (
                                        <p className="text-muted-foreground">No products selected</p>
                                    )}
                                </CardContent>
                            </Card>
                        </>
                    )}
                </TabsContent>

                <TabsContent value="images" className="pt-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Product Images</CardTitle>
                        </CardHeader>
                        <CardContent>
                            {formData.images.length > 0 ? (
                                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                                    {formData.images.map((image, index) => (
                                        <div key={index} className="aspect-square overflow-hidden rounded-md border">
                                            <img
                                                src={image}
                                                alt="Product"
                                                className="h-full w-full object-cover"
                                            />
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-muted-foreground">No images uploaded</p>
                            )}
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>

            <div className="flex justify-between">
                <Button type="button" variant="outline" onClick={onBack} className="flex items-center gap-2">
                    <ArrowLeft className="w-4 h-4" /> Back
                </Button>
                <Button
                    onClick={onSubmit}
                    disabled={isSubmitting}
                    className="min-w-[120px]"
                >
                    {isSubmitting ? "Creating..." : "Create Product"}
                </Button>
            </div>
        </div>
    );
};

export default ReviewScreen;