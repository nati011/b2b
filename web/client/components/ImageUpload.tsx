
import React, { useState, useCallback, useEffect } from "react";
import { useDropzone } from "react-dropzone";
import { Button } from "@/components/ui/button";
import { X, Upload, Loader } from "lucide-react";
import { toast } from "sonner";
import { Progress } from "@/components/ui/progress";

interface ImageUploadProps {
  onChange: (value: string[]) => void;
  value?: string[];
}

const ImageUpload: React.FC<ImageUploadProps> = ({ onChange, value = [] }) => {
  const [preview, setPreview] = useState<string | null>(null);
  const [imageValue, setImageValue] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState<{ [key: string]: number }>({});
  const [isUploading, setIsUploading] = useState(false);

  const api_key = "AmdeORpwsw7AJGbbjfwAYgPk1yQ";
  const cloud_name = "dnqkrebrb";

  const removeImage = () => {
    setPreview(null);
    setImageValue(null);
    onChange([]);
  };

  useEffect(() => {
    if (value && value.length > 0) {
      setPreview(value[0]);
      setImageValue(value[0]);
    } else {
      setPreview(null);
      setImageValue(null);
    }
  }, [value]);

  const uploadImage = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("upload_preset", "uc2udcgh");
    formData.append("api_key", api_key);

    try {
      const xhr = new XMLHttpRequest();

      xhr.open("POST", `https://api.cloudinary.com/v1_1/${cloud_name}/image/upload`);

      xhr.upload.addEventListener("progress", (event) => {
        if (event.lengthComputable) {
          const progress = Math.round((event.loaded * 100) / event.total);
          setUploadProgress({ [file.name]: progress });
        }
      });

      return new Promise<string>((resolve, reject) => {
        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            const response = JSON.parse(xhr.responseText);
            resolve(response.secure_url);
          } else {
            reject(new Error(`Upload failed with status ${xhr.status}`));
          }
        };

        xhr.onerror = () => {
          reject(new Error("Network error during upload"));
        };

        xhr.send(formData);
      });
    } catch (error) {
      console.error("Error uploading image:", error);
      throw error;
    }
  };

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      if (acceptedFiles.length === 0) return;
      const file = acceptedFiles[0];
      setIsUploading(true);
      try {
        toast("Uploading image...");
        const uploadedImage = await uploadImage(file);
        setImageValue(uploadedImage);
        setPreview(uploadedImage);
        onChange([uploadedImage]);
        toast.success("Image uploaded successfully!");
      } catch (error) {
        console.error("Error uploading image:", error);
        toast.error("Failed to upload image");
      } finally {
        setUploadProgress({});
        setIsUploading(false);
      }
    },
    [onChange]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    multiple: false,
    accept: {
      'image/*': ['.png', '.jpg', '.jpeg', '.gif']
    }
  });

  return (
    <div className="space-y-4">
      <div
        {...getRootProps()}
        className={`
          relative
          cursor-pointer
          transition
          border-dashed
          border-2
          p-6
          border-input
          flex
          flex-col
          justify-center
          items-center
          gap-4
          text-muted-foreground
          rounded-md
          ${isDragActive ? 'bg-muted/50' : ''}
          ${isUploading ? 'opacity-50 pointer-events-none' : 'hover:bg-muted/50'}
        `}
      >
        <input {...getInputProps()} />
        <div className="flex flex-col items-center text-center">
          {isUploading ? (
            <Loader className="h-8 w-8 animate-spin text-muted-foreground" />
          ) : (
            <Upload className="h-8 w-8 text-muted-foreground" />
          )}
          <h3 className="mt-2 text-sm font-semibold">
            {isDragActive ? "Drop image here" : "Drag & drop image here"}
          </h3>
          <p className="mt-1 text-xs text-muted-foreground">
            or click to browse
          </p>
        </div>
      </div>

      {Object.keys(uploadProgress).length > 0 && (
        <div className="w-full space-y-3">
          {Object.entries(uploadProgress).map(([filename, progress]) => (
            <div key={filename} className="w-full">
              <div className="flex justify-between text-sm text-muted-foreground mb-1">
                <span className="truncate max-w-[80%]">{filename}</span>
                <span>{progress}%</span>
              </div>
              <Progress value={progress} className="h-2" />
            </div>
          ))}
        </div>
      )}

      {preview && (
        <div className="grid grid-cols-1 gap-4">
          <div className="group relative aspect-square rounded-md border bg-muted">
            <img
              src={preview}
              alt="Preview"
              className="h-full w-full rounded-md object-cover"
              onError={(e) => {
                (e.target as HTMLImageElement).src = "/placeholder.svg";
              }}
            />
            <Button
              size="icon"
              variant="destructive"
              className="absolute right-1 top-1 h-6 w-6 opacity-0 group-hover:opacity-100 transition"
              onClick={removeImage}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ImageUpload;
