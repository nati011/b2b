'use client'
import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Product } from "@/app/libs/types";

interface ConfigurableProductsProps {
    products: Product[];
    selectedProductIds: number[];
    onUpdateSelectedProducts: (ids: number[]) => void;
    attributeKeys: string[];
    onBack: () => void;
    onNext: () => void;
}

const ConfigurableProducts: React.FC<ConfigurableProductsProps> = ({
    products,
    selectedProductIds,
    onUpdateSelectedProducts,
    attributeKeys,
    onBack,
    onNext,
}) => {
    const [searchTerm, setSearchTerm] = useState("");

    const toggleProduct = (productId: number) => {
        if (selectedProductIds.includes(productId)) {
            onUpdateSelectedProducts(selectedProductIds.filter(id => id !== productId));
        } else {
            onUpdateSelectedProducts([...selectedProductIds, productId]);
        }
    };

    const filteredProducts = products.filter(product =>
        product.Name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="space-y-6">
            <div className="space-y-4">
                <h3 className="text-lg font-medium">Select Products for Configuration</h3>

                <div className="relative">
                    <input
                        type="text"
                        placeholder="Search products..."
                        className="w-full p-2 pl-8 border rounded-md"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-5 w-5 absolute top-2.5 left-2 text-gray-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                        />
                    </svg>
                </div>

                <ScrollArea className="h-[320px] border rounded-md">
                    {filteredProducts.length === 0 ? (
                        <div className="p-4 text-center text-muted-foreground">
                            No products found matching your search
                        </div>
                    ) : (
                        <div className="p-1">
                            {filteredProducts.map((product) => (
                                <Card
                                    key={product.Id}
                                    className={`mb-2 cursor-pointer transition-colors hover:bg-accent/50 ${selectedProductIds.includes(product.Id) ? "border-primary" : ""
                                        }`}
                                    onClick={() => toggleProduct(product.Id)}
                                >
                                    <CardContent className="p-3 flex items-center space-x-3">
                                        <Checkbox
                                            checked={selectedProductIds.includes(product.Id)}
                                            onCheckedChange={() => toggleProduct(product.Id)}
                                            onClick={(e) => e.stopPropagation()}
                                        />
                                        <div className="flex-1">
                                            <p className="font-medium">{product.Name}</p>
                                            <div className="flex flex-wrap gap-2 mt-1">
                                                {product.Attributes?.filter((attr: { key: string; }) => attributeKeys.includes(attr.key))
                                                    .map((attr: any, idx: any) => (
                                                        <span key={idx} className="text-xs bg-accent rounded px-2 py-0.5">
                                                            {attr.key}: {attr.value}
                                                        </span>
                                                    ))
                                                }
                                            </div>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}
                        </div>
                    )}
                </ScrollArea>

                <div className="mt-4">
                    <p className="text-sm font-medium">Selected Products: {selectedProductIds.length}</p>
                    {attributeKeys.length > 0 && (
                        <p className="text-sm text-muted-foreground mt-1">
                            These products will be grouped as variations based on: {attributeKeys.join(", ")}
                        </p>
                    )}
                </div>
            </div>

            <div className="flex justify-between">
                <Button type="button" variant="outline" onClick={onBack} className="flex items-center gap-2">
                    <ArrowLeft className="w-4 h-4" /> Back
                </Button>
                <Button
                    onClick={onNext}
                    disabled={selectedProductIds.length === 0}
                    className="flex items-center gap-2"
                >
                    Continue <ArrowRight className="w-4 h-4" />
                </Button>
            </div>
        </div>
    );
};

export default ConfigurableProducts;