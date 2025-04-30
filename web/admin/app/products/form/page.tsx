"use client";
import ImageUpload, { Image } from "@/app/components/image-upload";
import useProductsStore from "@/app/libs/store/useProductStore";
import {
  Button,
  Card,
  Form,
  Input,
  Select,
  Modal,
  Space,
  Typography,
  Row,
  Col,
} from "antd";
import { useState, useEffect } from "react";
import {
  StyledManageCategoriesButton,
  StyledModal,
  StyledCategoryItem,
  StyledAddCategoryButton,
} from "./styles";

export default function ProductsForm() {
  const {
    products,
    loading,
    error,
    fetchProducts,
    categories,
    categoriesLoading,
    categoriesError,
    fetchCategories,
    createCategory,
    deleteCategory,
    createProduct,
  } = useProductsStore();
  const [form] = Form.useForm();
  const { TextArea } = Input;
  const [isAddCategoryModalVisible, setIsAddCategoryModalVisible] =
    useState(false);
  const [newCategoryName, setNewCategoryName] = useState("");
  const [isManageCategoriesModalVisible, setIsManageCategoriesModalVisible] =
    useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState<number | null>(null);
  const [isDeleteConfirmationVisible, setIsDeleteConfirmationVisible] =
    useState(false);
  const [attributes, setAttributes] = useState<{ [key: string]: string }>({});
  const [newAttributeKey, setNewAttributeKey] = useState("");
  const [newAttributeValue, setNewAttributeValue] = useState("");
  const [selectedCategoryIds, setSelectedCategoryIds] = useState<number[]>([]);
  const [productImages, setProductImages] = useState<Image[]>([]);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  const handleAddCategory = () => {
    setIsAddCategoryModalVisible(true);
  };

  const handleCreateCategory = async () => {
    try {
      await createCategory(newCategoryName);
      setIsAddCategoryModalVisible(false);
      setNewCategoryName("");
    } catch (error) {
      console.error("Error creating category:", error);
    }
  };

  const categoryOptions = categories.map((category) => ({
    value: category?.id,
    label: category?.name,
  }));

  const showManageCategoriesModal = () => {
    setIsManageCategoriesModalVisible(true);
  };

  const handleCancelManageCategories = () => {
    setIsManageCategoriesModalVisible(false);
  };

  const handleDeleteCategory = (categoryId: number) => {
    setCategoryToDelete(categoryId);
    setIsDeleteConfirmationVisible(true);
  };

  const handleConfirmDelete = async () => {
    try {
      if (categoryToDelete !== null) {
        await deleteCategory(categoryToDelete);
        setIsDeleteConfirmationVisible(false);
        setCategoryToDelete(null);
      }
    } catch (error) {
      console.error("Error deleting category:", error);
    }
  };

  const handleCancelDelete = () => {
    setIsDeleteConfirmationVisible(false);
    setCategoryToDelete(null);
  };

  const closeAddCategoryModal = () => {
    setIsAddCategoryModalVisible(false);
  };

  const handleAddAttribute = () => {
    if (newAttributeKey && newAttributeValue) {
      setAttributes({ ...attributes, [newAttributeKey]: newAttributeValue });
      setNewAttributeKey("");
      setNewAttributeValue("");
    }
  };

  const handleDeleteAttribute = (key: string) => {
    const newAttributes = { ...attributes };
    delete newAttributes[key];
    setAttributes(newAttributes);
  };

  const handleCategoryChange = (value: number[]) => {
    setSelectedCategoryIds(value);
  };

  const generateExternalId = () => {
    const randomNumber = Math.floor(Math.random() * 100000);
    const randomString = Math.random().toString(36).substring(2, 7);
    return `${randomNumber}-${randomString}`;
  };

  const handleCreateProduct = async () => {
    try {
      const values = await form.validateFields();
      const externalId = generateExternalId();

      const productData = {
        name: values.productName,
        desc: values.description,
        external_id: externalId,
        images: productImages.map((img) => img.image_url),
        price: parseFloat(values.price),
        attributes: attributes,
        distributor_id: 1,
        category_id: selectedCategoryIds,
      };

      await createProduct(productData);
      form.resetFields();
      setAttributes({});
      setProductImages([]);
      setSelectedCategoryIds([]);
      console.log("Product created successfully!");
    } catch (error: any) {
      console.error("Error creating product:", error);
    }
  };

  const handleImageChange = (images: Image[]) => {
    setProductImages(images);
  };

  if (categoriesLoading) {
    return <div>Loading categories...</div>;
  }

  if (categoriesError) {
    return <div>Error loading categories: {categoriesError}</div>;
  }

  return (
    <div className='grid grid-cols-1 gap-4'>
      <p className='font-semibold text-3xl text-black'>Add New Product</p>
      <Card title='Product Information'>
        <Form layout={"vertical"} form={form} name='productForm'>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item
                label='Product Name'
                name='productName'
                rules={[
                  { required: true, message: "Please enter product name!" },
                ]}
              >
                <Input placeholder='Product Name' />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label='Price'
                name='price'
                rules={[{ required: true, message: "Please enter price!" }]}
              >
                <Input placeholder='Price' />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label='Category'
                name='category'
                rules={[
                  { required: true, message: "Please select a category!" },
                ]}
              >
                <Select
                  mode='multiple'
                  placeholder='Select a category'
                  options={categoryOptions}
                  onChange={handleCategoryChange}
                  dropdownRender={(menu) => (
                    <div>
                      {menu}
                      <StyledManageCategoriesButton
                        onClick={showManageCategoriesModal}
                      >
                        Manage Categories
                      </StyledManageCategoriesButton>
                    </div>
                  )}
                />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item label='Attributes'>
            <Input
              placeholder='Attribute Key'
              value={newAttributeKey}
              onChange={(e) => setNewAttributeKey(e.target.value)}
              style={{ marginBottom: "8px" }}
            />
            <Input
              placeholder='Attribute Value'
              value={newAttributeValue}
              onChange={(e) => setNewAttributeValue(e.target.value)}
              style={{ marginBottom: "8px" }}
            />
            <Button
              type='primary'
              onClick={handleAddAttribute}
              style={{ marginBottom: "8px" }}
            >
              Add Attribute
            </Button>

            {Object.entries(attributes).map(([key, value]) => (
              <div
                key={key}
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  marginBottom: "4px",
                }}
              >
                <span>
                  {key}: {value}
                </span>
                <Button
                  size='small'
                  danger
                  onClick={() => handleDeleteAttribute(key)}
                >
                  Delete
                </Button>
              </div>
            ))}
          </Form.Item>

          <Form.Item label='Description' name='description'>
            <TextArea
              placeholder='Product Description'
              autoSize={{ minRows: 5 }}
            />
          </Form.Item>

          {/* Image Upload */}
          <Form.Item label='Images'>
            <ImageUpload onChange={handleImageChange} value={productImages} />
          </Form.Item>

          <div style={{ textAlign: "right" }}>
            <Button type='primary' onClick={handleCreateProduct}>
              Add Product
            </Button>
          </div>
        </Form>
      </Card>

      <Modal
        title='Add New Category'
        visible={isAddCategoryModalVisible}
        onOk={handleCreateCategory}
        onCancel={closeAddCategoryModal}
      >
        <Input
          placeholder='Category Name'
          value={newCategoryName}
          onChange={(e) => setNewCategoryName(e.target.value)}
        />
      </Modal>

      <StyledModal
        title='Manage Categories'
        visible={isManageCategoriesModalVisible}
        onCancel={handleCancelManageCategories}
        footer={null}
      >
        <StyledAddCategoryButton
          type='primary'
          size='small'
          onClick={handleAddCategory}
        >
          Add New Category
        </StyledAddCategoryButton>
        <Space direction='vertical' style={{ width: "100%" }}>
          {categories.map((category) => (
            <StyledCategoryItem key={category.id}>
              <Typography.Text>{category.name}</Typography.Text>
              <Button
                danger
                size='small'
                onClick={() => handleDeleteCategory(category.id)}
              >
                Delete
              </Button>
            </StyledCategoryItem>
          ))}
        </Space>
      </StyledModal>

      <Modal
        title='Confirm Delete'
        visible={isDeleteConfirmationVisible}
        onOk={handleConfirmDelete}
        onCancel={handleCancelDelete}
      >
        <p>Are you sure you want to delete this category?</p>
      </Modal>
    </div>
  );
}
