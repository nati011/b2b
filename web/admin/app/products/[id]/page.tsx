'use client'
import useProductsStore from "@/app/libs/store/useProductStore";
import { useEffect, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import Heading from "@/app/components/breadcrumb";
import { MultiSelect } from "@/app/components/multiselect";
import { Badge } from "@/components/ui/badge"
import { Plus, X } from "lucide-react";
import ImageUpload from "@/app/components/image-upload";
import { useParams } from "next/navigation";


interface ConfigurableProductForm {
    name?: string;
    desc?: string;
    external_id?: string;
    images?: string[];
    attributes?: string[];
    distributor_id?: number;
    category_id?: number;
};


const EditProduct = () => {
    const [formData, setFormData] = useState<ConfigurableProductForm>()

    const updateFormData = <K extends keyof ConfigurableProductForm>(key: K, value: ConfigurableProductForm[K]) => {
        setFormData(prev => ({
            ...prev,
            [key]: value
        }));
    };

    const [errors, setErrors] = useState({
        name: "",
        price: "",
        attributeKeys: ""
    });

    const handleFormSubmit = () => {
        console.log("Submitted")
    }

    const routeParam = useParams<{ id: string }>();

    const {
        categories,
        success,
        loading,
        error,
        fetchProductDetail
    } = useProductsStore()

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

    const options = categories.map((c) => {
        return {
            value: c.id.toString(),
            label: c.name
        };
    });
    const [newKey, setNewKey] = useState("");
    const [newValue, setNewValue] = useState("");
    const [keyInput, setKeyInput] = useState("");
    const [e, setError] = useState("");


    const handleAddAttribute = () => {
        if (!keyInput.trim()) {
            setError("Attribute key is required");
            return;
        }

        if (formData?.attributes?.includes(keyInput.trim())) {
            setError("This attribute key already exists");
            return;
        }
        if (formData?.attributes) {
            updateFormData("attributes", [...formData.attributes, keyInput.trim()])
        }
        setKeyInput("");
        setError("");
    };


    useEffect(() => {
        fetchProductDetail(parseInt(routeParam.id));
    }, []);

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
                                    value={formData?.name}
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
                                    value={formData?.desc}
                                    // onChange={handleInputChange}
                                    placeholder="Describe your product"
                                    className="min-h-[120px]"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="category">
                                    Category <span className="text-destructive">*</span>
                                </Label>
                                <div className="relative">

                                    <MultiSelect options={options} onValueChange={() => { console.log("updated") }} />
                                </div>
                                {errors.price && (
                                    <p className="text-sm text-destructive">{errors.price}</p>
                                )}
                            </div>
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
                                            className="flex items-center gap-1"
                                            onClick={handleAddAttribute}
                                        >
                                            <Plus className="w-4 h-4" /> Add
                                        </Button>
                                    </div>
                                </div>

                                {error && <p className="text-sm text-destructive">{error}</p>}

                                <div className="border rounded-md p-4">
                                    <Label>Variation Keys</Label>
                                    <div className="flex flex-wrap gap-2 mt-2">
                                        {formData?.attributes?.length === 0 ? (
                                            <p className="text-sm text-muted-foreground italic">No attribute keys added yet</p>
                                        ) : (
                                            formData?.attributes?.map((key, index) => (
                                                <Badge key={index} variant="outline" className="px-3 py-1">
                                                    {key}
                                                    <X
                                                        className="ml-2 w-3 h-3 cursor-pointer"

                                                    />
                                                </Badge>
                                            ))
                                        )}
                                    </div>
                                    <p className="text-sm text-muted-foreground mt-2">
                                        These keys will be used as variation attributes across all product variants.
                                    </p>
                                </div>
                                <Label>Product Images</Label>
                                <ImageUpload onChange={(images) => updateFormData("images", images)}
                                    value={formData?.images} />
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
