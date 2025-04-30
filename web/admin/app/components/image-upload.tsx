import { useState, useCallback, useEffect } from "react";

import Image from "next/image";

import { useDropzone } from "react-dropzone";
import { TbPhotoPlus } from "react-icons/tb";

import { Progress, Button } from "antd";

import axios from "axios";

export type Image = {
  image_url: string;
  is_cover: boolean;
  blur_hash?: string;
};

interface ImageUploadProps {
  onChange: (value: Image[]) => void;
  value?: Image[];
}

const ImageUpload: React.FC<ImageUploadProps> = ({ onChange, value }) => {
  const [previews, setPreviews] = useState<Image[]>([]);
  const [values, setValues] = useState<Image[]>([]);
  const [uploadProgress, setUploadProgress] = useState<{ [key: string]: number }>({});

  const api_key = "AmdeORpwsw7AJGbbjfwAYgPk1yQ";
  const cloud_name = "dnqkrebrb";

  const removeImage = (index: number) => () => {
    setPreviews((prevPreviews) => prevPreviews.filter((_, i) => i !== index));
    setValues((prevValues) => prevValues.filter((_, i) => i !== index));
    onChange(values.filter((_, i) => i !== index));
  };


  useEffect(() => {
    if (value) {
      setPreviews(value);
      setValues(value);
    } else {
      setPreviews([]);
      setValues([]);
    }
  }, [value, onChange]);


  const uploadImage = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("upload_preset", "uc2udcgh");
    formData.append("api_key", api_key);

    try {
      const response = await axios.post(
        `https://api.cloudinary.com/v1_1/${cloud_name}/image/upload`,
        formData,
        {
          onUploadProgress: (progressEvent) => {
            const progress = Math.round(
              (progressEvent.loaded * 100) / (progressEvent.total || 1)
            );
            setUploadProgress(prev => ({
              ...prev,
              [file.name]: progress
            }));
          },
        }
      );
      console.log(response.data)
      return {
        image_url: response.data.secure_url,
        is_cover: false,
        blur_hash: "eoG9Hrt5bbWAR*Y8s:ofoeofXVbcjYkCjss.ogWAWXWBaJWBa}j?oL",
      };
    } catch (error) {
      console.error("Error uploading image:", error);
      throw error;
    }
  };

  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      if (acceptedFiles.length > 0) {
        try {
          const uploadPromises = acceptedFiles.map(file => uploadImage(file));
          const uploadedImages = await Promise.all(uploadPromises);

          const newValues = [...values, ...uploadedImages];
          setValues(newValues);
          setPreviews(newValues);
          onChange(newValues);
          setUploadProgress({});
        } catch (error) {
          console.error("Error uploading images:", error);
          setUploadProgress({});
        }
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
    <div className="flex flex-col gap-4">
      <div
        className="
         relative
         cursor-pointer
         hover:opacity-70
         transition
         border-dashed
         border-2
         p-10
         border-neutral-300
         flex
         flex-col
         justify-center
         items-center
         gap-4
         text-neutral-600
         rounded-lg
        "
        {...getRootProps()}
      >
        <input {...getInputProps()} />
        <div className="flex flex-col items-center text-center">
          <TbPhotoPlus className="w-8 h-8 text-neutral-400" />
          {isDragActive ? (
            <p>Drop the files here ...</p>
          ) : (
            <p>Drag and drop images here, or click to select</p>
          )}
        </div>
      </div>

      {Object.keys(uploadProgress).length > 0 && (
        <div className="w-full space-y-2">
          {Object.entries(uploadProgress).map(([filename, progress]) => (
            <div key={filename} className="w-full">
              <div className="text-sm text-neutral-500 mb-1">{filename}</div>
              <Progress percent={progress} />
            </div>
          ))}
        </div>
      )}

      {previews.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
          {previews.map((preview, index) => (
            <div key={index} className="relative border rounded-lg p-4">
              <div className="relative aspect-square w-full overflow-hidden rounded-lg">
                <Image
                  src={preview.image_url}
                  alt={`Preview ${index + 1}`}
                  fill
                  className="object-cover"

                />
              </div>
              <div className="mt-2 flex items-center justify-between">

                <Button type="primary" danger onClick={removeImage(index)}>
                  Delete
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ImageUpload;