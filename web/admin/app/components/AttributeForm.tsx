import React, { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { X, Plus } from "lucide-react";



interface ProductAttributeFormProps {
    attributes?: any[];
    onChange: (attributes: any[]) => void;
}

const ProductAttributeForm: React.FC<ProductAttributeFormProps> = ({
    attributes = [],
    onChange,
}) => {
    const [newKey, setNewKey] = useState("");
    const [newValue, setNewValue] = useState("");

    const handleAddAttribute = () => {
        if (newKey.trim() && newValue.trim()) {
            if (attributes.some(attr => attr.key === newKey.trim())) {
                alert("Attribute key already exists!");
                return;
            }

            onChange([...attributes, {
                key: newKey.trim(),
                value: newValue.trim()
            }]);
            setNewKey("");
            setNewValue("");
        }
    };

    const handleRemoveAttribute = (index: number) => {
        const newAttributes = attributes.filter((_, i) => i !== index);
        onChange(newAttributes);
    };

    const handleValueChange = (index: number, newValue: string) => {
        const newAttributes = attributes.map((attr, i) =>
            i === index ? { ...attr, value: newValue } : attr
        );
        onChange(newAttributes);
    };

    return (
        <Card className="rounded-sm border-2 border-gray-200 shadow-none">
            <CardContent>
                <div className="mb-4">
                    <h3 className="text-xl font-bold">Product Attributes</h3>
                    <p className="text-sm text-muted-foreground">
                        Add custom attributes for this product.
                    </p>
                </div>

                <div className="space-y-4">
                    {attributes.map((attribute, index) => (
                        <div key={`${attribute.key}-${index}`} className="flex items-center gap-2">
                            <div className="grid flex-1 gap-2">
                                <div className="grid grid-cols-2 gap-2">
                                    <div>
                                        <Input
                                            value={attribute.key}
                                            disabled
                                            className="bg-muted"
                                        />
                                    </div>
                                    <div>
                                        <Input
                                            value={attribute.value}
                                            onChange={(e) =>
                                                handleValueChange(index, e.target.value)
                                            }
                                        />
                                    </div>
                                </div>
                            </div>
                            <Button
                                size="icon"
                                variant="ghost"
                                className="h-8 w-8"
                                onClick={() => handleRemoveAttribute(index)}
                            >
                                <X className="h-4 w-4" />
                            </Button>
                        </div>
                    ))}

                    <div className="flex items-center gap-2">
                        <div className="grid flex-1 gap-2">
                            <div className="grid grid-cols-2 gap-2">
                                <div>
                                    <Input
                                        placeholder="Attribute name"
                                        value={newKey}
                                        onChange={(e) => setNewKey(e.target.value)}
                                    />
                                </div>
                                <div>
                                    <Input
                                        placeholder="Attribute value"
                                        value={newValue}
                                        onChange={(e) => setNewValue(e.target.value)}
                                    />
                                </div>
                            </div>
                        </div>
                        <Button
                            size="icon"
                            className="h-8 w-8"
                            onClick={handleAddAttribute}
                            type="button"
                            disabled={!newKey.trim() || !newValue.trim()}
                        >
                            <Plus className="h-4 w-4" />
                        </Button>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

export default ProductAttributeForm;