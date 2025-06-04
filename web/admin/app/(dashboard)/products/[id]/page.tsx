'use client'
import useProductsStore from "@/app/libs/store/useProductStore";
import { useEffect, useState } from "react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import Heading from "@/components/breadcrumb";
import { MultiSelect } from "@/components/multiselect";
import { Badge } from "@/components/ui/badge"
import { Plus, X } from "lucide-react";
import ImageUpload from "@/components/image-upload";
import { useParams } from "next/navigation";
import ProductForm from "@/components/product-form";



const EditProduct = () => {

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
        product,
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




    useEffect(() => {
        fetchProductDetail(parseInt(routeParam.id));
    }, []);

    return (
        <>
            <Heading page={pages} heading="Edit Selected Product" />
            <ProductForm
                isEdit
                initialData={product}
            />
        </>
    );
};

export default EditProduct;
