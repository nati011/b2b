'use client'

import Heading from "@/components/breadcrumb";
import ProductForm from "@/components/product-form";

export default function Retailers() {

    const pages = [
        {
            "title": "Products",
            "href": "/product"
        },
        {
            "title": "Form",
            "href": "/product/form"
        }
    ]

    return (
        <>
            <Heading page={pages} heading="Product" subheading="Add a new Product" />
            <ProductForm />
        </>
    );
}
