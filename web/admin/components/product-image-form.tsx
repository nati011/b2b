import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { X, Images } from "lucide-react";
import ImageUpload from "@/components/image-upload";

interface ProductImagesFormProps {
    images: string[];
    onChange: (images: string[]) => void;
}

const ProductImagesForm: React.FC<ProductImagesFormProps> = ({
    images,
    onChange,
}) => {
    const [newImageUrl, setNewImageUrl] = React.useState("");

    const addImage = () => {
        if (newImageUrl.trim()) {
            onChange([...images, newImageUrl.trim()]);
            setNewImageUrl("");
        }
    };

    const removeImage = (index: number) => {
        const updatedImages = [...images];
        updatedImages.splice(index, 1);
        onChange(updatedImages);
    };

    return (
        <Card className="rounded-sm border-2 border-gray-200 shadow-none">
            <CardContent>
                <div className="mb-4">
                    <h3 className="text-xl font-bold">Product Images</h3>
                    <p className="text-sm text-muted-foreground">
                        Add images for this product.
                    </p>
                </div>

                <div className="space-y-4">
                    {images.length > 0 && (
                        <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
                            {images.map((image, index) => (
                                <div
                                    key={index}
                                    className="group relative aspect-square rounded-md border bg-muted"
                                >
                                    <img
                                        alt={`Product image ${index + 1}`}
                                        className="object-cover h-full w-full rounded-md"
                                        src={image}
                                        onError={(e) => {
                                            (e.target as HTMLImageElement).src = "/placeholder.svg";
                                        }}
                                    />
                                    <Button
                                        size="icon"
                                        variant="secondary"
                                        className="absolute right-1 top-1 h-6 w-6 opacity-0 group-hover:opacity-100"
                                        onClick={() => removeImage(index)}
                                    >
                                        <X className="h-4 w-4" />
                                    </Button>
                                </div>
                            ))}
                        </div>
                    )}

                    <div className="flex items-center gap-2">
                        <Input
                            placeholder="Enter image URL"
                            value={newImageUrl}
                            onChange={(e) => setNewImageUrl(e.target.value)}
                            className="flex-1"
                        />
                        <Button onClick={addImage}>Add Image</Button>
                    </div>
                    <ImageUpload value={images} onChange={onChange} />
                </div>
            </CardContent>
        </Card>
    );
};

export default ProductImagesForm;