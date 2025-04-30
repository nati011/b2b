'use client'
import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ArrowLeft, ArrowRight, Plus, X } from "lucide-react";
import { Badge } from "@/components/ui/badge";

export interface Attribute {
    key: string;
    value: string;
}

interface AttributesFormProps {
    attributes: Attribute[];
    onUpdateAttributes: (attributes: Attribute[]) => void;
    attributeKeys: string[];
    onUpdateAttributeKeys: (keys: string[]) => void;
    productType: "simple" | "configurable";
    onBack: () => void;
    onNext: () => void;
}

const AttributesForm: React.FC<AttributesFormProps> = ({
    attributes,
    onUpdateAttributes,
    attributeKeys,
    onUpdateAttributeKeys,
    productType,
    onBack,
    onNext,
}) => {
    const [newKey, setNewKey] = useState("");
    const [newValue, setNewValue] = useState("");
    const [keyInput, setKeyInput] = useState("");
    const [error, setError] = useState("");

    const handleAddAttribute = () => {
        if (!newKey.trim()) {
            setError("Attribute key is required");
            return;
        }

        if (!newValue.trim() && productType === "simple") {
            setError("Attribute value is required");
            return;
        }

        if (attributes.some(attr => attr.key === newKey)) {
            setError("This attribute key already exists");
            return;
        }

        const newAttribute = { key: newKey.trim(), value: newValue.trim() };
        onUpdateAttributes([...attributes, newAttribute]);
        setNewKey("");
        setNewValue("");
        setError("");
    };

    const removeAttribute = (key: string) => {
        onUpdateAttributes(attributes.filter(attr => attr.key !== key));
    };

    const handleAddAttributeKey = () => {
        if (!keyInput.trim()) {
            setError("Attribute key is required");
            return;
        }

        if (attributeKeys.includes(keyInput.trim())) {
            setError("This attribute key already exists");
            return;
        }

        onUpdateAttributeKeys([...attributeKeys, keyInput.trim()]);
        setKeyInput("");
        setError("");
    };

    const removeAttributeKey = (key: string) => {
        onUpdateAttributeKeys(attributeKeys.filter(k => k !== key));
    };

    return (
        <div className="space-y-6">
            <div>
                <h3 className="text-lg font-medium mb-4">
                    {productType === "simple" ? "Product Attributes" : "Variation Attributes"}
                </h3>

                {productType === "simple" ? (
                    <div className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label htmlFor="attributeKey">Attribute Key</Label>
                                <Input
                                    id="attributeKey"
                                    placeholder="e.g., Color, Size, Material"
                                    value={newKey}
                                    onChange={(e) => setNewKey(e.target.value)}
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="attributeValue">Attribute Value</Label>
                                <Input
                                    id="attributeValue"
                                    placeholder="e.g., Red, XL, Cotton"
                                    value={newValue}
                                    onChange={(e) => setNewValue(e.target.value)}
                                />
                            </div>
                        </div>

                        {error && <p className="text-sm text-destructive">{error}</p>}

                        <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={handleAddAttribute}
                            className="flex items-center gap-1"
                        >
                            <Plus className="w-4 h-4" /> Add Attribute
                        </Button>

                        <div className="border rounded-md p-4 space-y-2">
                            <Label>Current Attributes</Label>
                            {attributes.length === 0 ? (
                                <p className="text-sm text-muted-foreground italic">No attributes added yet</p>
                            ) : (
                                <div className="space-y-2">
                                    {attributes.map((attr, index) => (
                                        <div key={index} className="flex justify-between items-center p-2 bg-accent/50 rounded-md">
                                            <div>
                                                <span className="font-medium">{attr.key}:</span> {attr.value}
                                            </div>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => removeAttribute(attr.key)}
                                            >
                                                <X className="w-4 h-4" />
                                            </Button>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>
                ) : (
                    <div className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="attributeKeyInput">Variation Attribute Keys</Label>
                            <div className="flex gap-2">
                                <Input
                                    id="attributeKeyInput"
                                    placeholder="e.g., Color, Size"
                                    value={keyInput}
                                    onChange={(e) => setKeyInput(e.target.value)}
                                />
                                <Button
                                    type="button"
                                    onClick={handleAddAttributeKey}
                                    className="flex items-center gap-1"
                                >
                                    <Plus className="w-4 h-4" /> Add
                                </Button>
                            </div>
                        </div>

                        {error && <p className="text-sm text-destructive">{error}</p>}

                        <div className="border rounded-md p-4">
                            <Label>Variation Keys</Label>
                            <div className="flex flex-wrap gap-2 mt-2">
                                {attributeKeys.length === 0 ? (
                                    <p className="text-sm text-muted-foreground italic">No attribute keys added yet</p>
                                ) : (
                                    attributeKeys.map((key, index) => (
                                        <Badge key={index} variant="outline" className="px-3 py-1">
                                            {key}
                                            <X
                                                className="ml-2 w-3 h-3 cursor-pointer"
                                                onClick={() => removeAttributeKey(key)}
                                            />
                                        </Badge>
                                    ))
                                )}
                            </div>
                            <p className="text-sm text-muted-foreground mt-2">
                                These keys will be used as variation attributes across all product variants.
                            </p>
                        </div>
                    </div>
                )}
            </div>

            <div className="flex justify-between">
                <Button type="button" variant="outline" onClick={onBack} className="flex items-center gap-2">
                    <ArrowLeft className="w-4 h-4" /> Back
                </Button>
                <Button onClick={onNext} className="flex items-center gap-2">
                    Continue <ArrowRight className="w-4 h-4" />
                </Button>
            </div>
        </div>
    );
};

export default AttributesForm;