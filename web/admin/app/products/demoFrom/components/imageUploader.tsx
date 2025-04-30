
import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight, Plus, X } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import ImageUpload from "@/app/components/image-upload";



interface ImageUploaderProps {
    images: string[];
    onUpdateImages: (images: string[]) => void;
    onBack: () => void;
    onNext: () => void;
}

const ImageUploader: React.FC<ImageUploaderProps> = ({
    images,
    onUpdateImages,
    onBack,
    onNext,
}) => {
    return (
        <div className="space-y-6">
            <div className="space-y-4">
                <h3 className="text-lg font-medium">Product Images</h3>

                <ImageUpload onChange={onUpdateImages}
                    value={images} />

                <div className="flex justify-between">
                    <Button type="button" variant="outline" onClick={onBack} className="flex items-center gap-2">
                        <ArrowLeft className="w-4 h-4" /> Back
                    </Button>
                    <Button onClick={onNext} className="flex items-center gap-2">
                        Continue <ArrowRight className="w-4 h-4" />
                    </Button>
                </div>
            </div>
        </div>
    );
};

export default ImageUploader;