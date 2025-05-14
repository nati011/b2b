
import React, { useState, useCallback, useEffect } from "react";
import { useDropzone } from "react-dropzone";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { X, Upload, Loader } from "lucide-react";
import { toast } from "sonner";

interface ImageUploadProps {
  onChange: (value: string[]) => void;
  value?: string[];
}

const ImageUpload: React.FC<ImageUploadProps> = ({ onChange, value = [] }) => {
  const [previews, setPreviews] = useState<string[]>([]);
  const [values, setValues] = useState<string[]>([]);
  const [uploadProgress, setUploadProgress] = useState<{ [key: string]: number }>({});
  const [isUploading, setIsUploading] = useState(false);

  const api_key = "AmdeORpwsw7AJGbbjfwAYgPk1yQ";
  const cloud_name = "dnqkrebrb";

  const removeImage = (index: number) => {
    const newPreviews = [...previews];
    const newValues = [...values];
    newPreviews.splice(index, 1);
    newValues.splice(index, 1);

    setPreviews(newPreviews);
    setValues(newValues);
    onChange(newValues);
  };

  useEffect(() => {
    if (value) {
      setPreviews(value);
      setValues(value);
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
          setUploadProgress(prev => ({
            ...prev,
            [file.name]: progress
          }));
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

      setIsUploading(true);
      try {
        const uploadPromises = acceptedFiles.map(file => uploadImage(file));
        toast("Uploading images...");

        const uploadedImages = await Promise.all(uploadPromises);

        const newValues = [...values, ...uploadedImages];
        setValues(newValues);
        setPreviews(newValues);
        onChange(newValues);
        toast.success("Images uploaded successfully!");
      } catch (error) {
        console.error("Error uploading images:", error);
        toast.error("Failed to upload some images");
      } finally {
        setUploadProgress({});
        setIsUploading(false);
      }
    },
    [values, onChange]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    multiple: true,
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
            {isDragActive ? "Drop images here" : "Drag & drop images here"}
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

      {/* {previews.length > 0 && (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
          {previews.map((preview, index) => (
            <div key={index} className="group relative aspect-square rounded-md border bg-muted">
              <img
                src={preview}
                alt={`Preview ${index + 1}`}
                className="h-full w-full rounded-md object-cover"
                onError={(e) => {
                  (e.target as HTMLImageElement).src = "/placeholder.svg";
                }}
              />
              <Button
                size="icon"
                variant="destructive"
                className="absolute right-1 top-1 h-6 w-6 opacity-0 group-hover:opacity-100 transition"
                onClick={() => removeImage(index)}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      )} */}
    </div>
  );
};

export default ImageUpload;
