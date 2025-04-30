'use client'
import ImageUpload, { Image } from "@/app/components/image-upload";
import useProductsStore from "@/app/libs/store/useProductStore"
import { Button, Card, Form, Input } from "antd";
import { useState } from "react";


export default function ProductsForm() {
  const {
    products,
    loading,
    error,
    fetchProducts
  } = useProductsStore()

  const [form] = Form.useForm();

  const [image, setImage] = useState<Image[]>([])

  const { TextArea } = Input;

  return (
    <div className="grid grid-cols-1 gap-4">
      <p className="font-semibold text-3xl text-black">
        Add New Product
      </p>
      <Card title="Product Information">
        <Form
          layout={'vertical'}
          form={form}
        >
          <div className="grid grid-cols-2 gap-2">
            <Form.Item label="External Id">
              <Input placeholder="External Id" />
            </Form.Item>
            <Form.Item label="Product Name">
              <Input placeholder="Product Name" />
            </Form.Item>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <Form.Item label="Price">
              <Input placeholder="Price" />
            </Form.Item>
            <Form.Item label="Category">
              <Input placeholder="Category" />
            </Form.Item>
          </div>
          <Form.Item label="Description">
            <TextArea placeholder="Product Description" autoSize={{ minRows: 5 }} />
          </Form.Item>
          <Form.Item label="Image">
            <ImageUpload onChange={setImage} />
          </Form.Item>

        </Form>
      </Card>
    </div>
  );
}
