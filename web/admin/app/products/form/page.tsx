'use client'
import ProductWizard from "@/app/products/form/components/productWizard";
import { toast } from "sonner"
import useProductsStore from "@/app/libs/store/useProductStore";
import { useEffect } from "react";

const Index = () => {
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
        createConfigurableProduct,
    } = useProductsStore();
    const handleCreateCategory = async (name: string) => {
        return new Promise<void>((resolve) => {
            setTimeout(() => {
                createCategory(name)
                toast("Category Created", {
                    description: `${name} has been added to categories.`
                });
                resolve();
            }, 500);
        });
    };

    const handleDeleteCategory = async (id: number) => {
        return new Promise<void>((resolve) => {
            setTimeout(() => {
                toast("Category Deleted", {
                    description: "The category has been deleted successfully."
                });

                deleteCategory(id)

                resolve();
            }, 500);
        });
    };
    const handleCreateProduct = async (productData: any) => {
        return new Promise<void>((resolve) => {
            setTimeout(() => {
                createProduct(productData)

                toast.success("Product Created", {
                    description: `${productData.name} has been created successfully.`
                });

                resolve();
            }, 1000);
        });
    };

    const handleCreateConfigurableProduct = async (productData: any) => {
        return new Promise<void>((resolve) => {
            setTimeout(() => {
                createConfigurableProduct(productData)
                toast("Configurable Product Created", {
                    description: `${productData.name} has been created with ${productData.products.length} variants.`
                });

                resolve();
            }, 1000);
        });
    };

    useEffect(() => {
        fetchCategories()
        fetchProducts()
    }, [])

    return (
        <ProductWizard
            categories={categories}
            products={products}
            onCreateCategory={handleCreateCategory}
            onDeleteCategory={handleDeleteCategory}
            onCreateProduct={handleCreateProduct}
            onCreateConfigurableProduct={handleCreateConfigurableProduct}
        />
    );
};

export default Index;
