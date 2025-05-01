"use client";
import ImageUpload from "@/app/components/image-upload";
import useProductsStore from "@/app/libs/store/useProductStore";
import { Button, Card, Form, Input, Select, Modal, Space, Typography, Row, Col } from "antd"; // Import Row and Col
import { useState, useEffect } from "react";
import AttributeKeysInput from "@/app/components/attribute-keys-text-input";
import {
  StyledManageCategoriesButton,
  StyledModal,
  StyledCategoryItem,
  StyledAddCategoryButton,
} from "./styles"; // Import from the styles file

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
    createCategory, // This line is important!
    deleteCategory,   // And this one!
    createProduct,    // Add createProduct
    createConfigurableProduct, // Add createConfigurableProduct
  } = useProductsStore();
  const [form] = Form.useForm();
  const { TextArea } = Input;
  const [isAddCategoryModalVisible, setIsAddCategoryModalVisible] =
    useState(false);
  const [newCategoryName, setNewCategoryName] = useState("");
  const [isManageCategoriesModalVisible, setIsManageCategoriesModalVisible] =
    useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState<number | null>(null); // State to hold the ID of the category to delete
  const [isDeleteConfirmationVisible, setIsDeleteConfirmationVisible] =
    useState(false);
  const [attributes, setAttributes] = useState<{ [key: string]: string }>({}); // State for attributes
  const [newAttributeKey, setNewAttributeKey] = useState("");
  const [newAttributeValue, setNewAttributeValue] = useState("");
  const [selectedCategoryIds, setSelectedCategoryIds] = useState<number[]>([]); // State for selected category IDs
  const [productImages, setProductImages] = useState<string[]>([]); // State for product images using your Image type
  const [activeForm, setActiveForm] = useState<'product' | 'configurable'>('configurable'); // State for active form
  const [selectedProducts, setSelectedProducts] = useState<number[]>([]); // State for selected products in configurable product form

  useEffect(() => {
    fetchCategories(); // Fetch categories when the component mounts
    fetchProducts(); // Fetch products when the component mounts
  }, [fetchCategories, fetchProducts]); // Add fetchCategories and fetchProducts to the dependency array

  const handleAddCategory = () => {
    setIsAddCategoryModalVisible(true);
  };

  const handleCreateCategory = async () => {
    try {
      await createCategory(newCategoryName); // Use the createCategory function from the store
      setIsAddCategoryModalVisible(false);
      setNewCategoryName("");
    } catch (error) {
      console.error("Error creating category:", error);
      // Optionally, display an error message to the user
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
        await deleteCategory(categoryToDelete); // Use the deleteCategory function from the store
        setIsDeleteConfirmationVisible(false);
        setCategoryToDelete(null);
      }
    } catch (error) {
      console.error("Error deleting category:", error);
      // Optionally, display an error message to the user
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
      const values = await form.validateFields(); // Validate the form
      const externalId = generateExternalId();

      const productData = {
        name: values.productName, // Use the field name from the form
        desc: values.description, // Use the field name from the form
        external_id: externalId,
        images: productImages,
        price: parseFloat(values.price),
        attributes: attributes,
        distributor_id: 1,
        category_id: selectedCategoryIds,
      };

      await createProduct(productData);
      form.resetFields(); // Clear the form
      setAttributes({}); // Clear the attributes
      setProductImages([]); // Clear the product images
      setSelectedCategoryIds([]); // Clear selected categories
      console.log("Product created successfully!");
    } catch (error: any) {
      console.error("Error creating product:", error);
      // Optionally, display an error message to the user
    }
  };

  const handleImageChange = (images: string[]) => {
    setProductImages(images);
  };

  const handleFormChange = (formType: 'product' | 'configurable') => {
    setActiveForm(formType);
  };

  const handleProductSelect = (value: number[]) => {
    setSelectedProducts(value);
  };

  const handleCreateConfigurableProduct = async () => {
    try {
      const values = await form.validateFields();

      const configurableProductData = {
        name: values.configurableProductName,
        desc: values.configurableProductDescription,
        external_id: "", // Or generate a random one if needed
        attribute_keys: values.attributeKeys, // Use the array of attribute keys
        products: selectedProducts, // Send the selected product IDs
        images: productImages,
      };

      await createConfigurableProduct(configurableProductData);
      form.resetFields();
      setProductImages([]);
      setSelectedProducts([]);
      console.log("Configurable product created successfully!");
    } catch (error: any) {
      console.error("Error creating configurable product:", error);
    }
  };

  const productOptions = products?.map((product) => ({
    value: product?.Id,
    label: product?.Name,
  }));

  return (
    <div className='grid grid-cols-1 gap-4'>
      <p className='font-semibold text-3xl text-black'>Add New Product</p>

      {/* Form Type Buttons */}
      <div style={{ display: 'flex', justifyContent: 'space-around', marginBottom: '16px' }}>
        <Button type={activeForm === 'product' ? 'primary' : 'default'} onClick={() => handleFormChange('product')}>
          Product
        </Button>
        <Button type={activeForm === 'configurable' ? 'primary' : 'default'} onClick={() => handleFormChange('configurable')}>
          Configurable Product
        </Button>
      </div>

      {/* Product Form */}
      {activeForm === 'product' && (
        <Card title='Product Information'>
          <Form layout={"vertical"} form={form} name="productForm">
            {/* Product Name, Price, and Category on the same line */}
            <Row gutter={16}>
              <Col span={8}>
                <Form.Item label='Product Name' name="productName" rules={[{ required: true, message: 'Please enter product name!' }]}>
                  <Input placeholder='Product Name' />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item label='Price' name="price" rules={[{ required: true, message: 'Please enter price!' }]}>
                  <Input placeholder='Price' />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item label='Category' name="category" rules={[{ required: true, message: 'Please select a category!' }]}>
                  <Select
                    mode="multiple" // Allow multiple selections
                    placeholder='Select a category'
                    options={categoryOptions}
                    onChange={handleCategoryChange} // Handle category selection
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

            {/* Attributes below Product Name, Price, and Category */}
            <Form.Item label='Attributes'>
              {/* Input for new attribute key */}
              <Input
                placeholder='Attribute Key'
                value={newAttributeKey}
                onChange={(e) => setNewAttributeKey(e.target.value)}
                style={{ marginBottom: '8px' }}
              />
              {/* Input for new attribute value */}
              <Input
                placeholder='Attribute Value'
                value={newAttributeValue}
                onChange={(e) => setNewAttributeValue(e.target.value)}
                style={{ marginBottom: '8px' }}
              />
              {/* Button to add attribute */}
              <Button type='primary' onClick={handleAddAttribute} style={{ marginBottom: '8px' }}>
                Add Attribute
              </Button>

              {/* Display existing attributes */}
              {Object.entries(attributes).map(([key, value]) => (
                <div key={key} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '4px' }}>
                  <span>{key}: {value}</span>
                  <Button size='small' danger onClick={() => handleDeleteAttribute(key)}>
                    Delete
                  </Button>
                </div>
              ))}
            </Form.Item>

            {/* Description below Attributes */}
            <Form.Item label='Description' name="description">
              <TextArea
                placeholder='Product Description'
                autoSize={{ minRows: 5 }}
              />
            </Form.Item>

            {/* Image Upload */}
            <Form.Item label='Images'>
              <ImageUpload
                onChange={handleImageChange}
                value={productImages} // Pass the current images as the value
              />
            </Form.Item>

            {/* Add Product Button */}
            <div style={{ textAlign: 'right' }}>
              <Button type='primary' onClick={handleCreateProduct}>
                Add Product
              </Button>
            </div>
          </Form>
        </Card>
      )}

      {/* Configurable Product Form */}
      {activeForm === 'configurable' && (
        <Card title='Configurable Product Information'>
          <Form layout={"vertical"} form={form} name="configurableProductForm">
            {/* Name, Attribute Key, and Products on the same line */}
            <Row gutter={16}>
              <Col span={8}>
                <Form.Item label='Name' name="configurableProductName" rules={[{ required: true, message: 'Please enter configurable product name!' }]}>
                  <Input placeholder='Name' />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label='Attribute Keys'
                  name="attributeKeys"
                  rules={[{ required: true, message: 'Please enter attribute keys!' }]}
                >
                  <AttributeKeysInput />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item label='Products' name="products" rules={[{ required: true, message: 'Please select products!' }]}>
                  <Select
                    mode="multiple"
                    placeholder="Select Products"
                    options={productOptions}
                    onChange={handleProductSelect}
                  />
                </Form.Item>
              </Col>
            </Row>

            {/* Description below Name, Attribute Key, and Products */}
            <Form.Item label='Description' name="configurableProductDescription">
              <TextArea
                placeholder='Description'
                autoSize={{ minRows: 5 }}
              />
            </Form.Item>

            {/* Image Upload */}
            <Form.Item label='Images'>
              <ImageUpload
                onChange={handleImageChange}
                value={productImages}
              />
            </Form.Item>

            {/* Add Configurable Product Button */}
            <div style={{ textAlign: 'right' }}>
              <Button type='primary' onClick={handleCreateConfigurableProduct}>
                Add Configurable Product
              </Button>
            </div>
          </Form>
        </Card>
      )}

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
        <StyledAddCategoryButton type="primary" size="small" onClick={handleAddCategory}>
          Add New Category
        </StyledAddCategoryButton>
        <Space direction='vertical' style={{ width: "100%" }}>

          {categories.map((category) => (
            <StyledCategoryItem key={category.id}>
              <Typography.Text>{category.name}</Typography.Text>
              <Button danger size="small" onClick={() => handleDeleteCategory(category.id)}>
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