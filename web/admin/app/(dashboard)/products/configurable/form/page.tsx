'use client'

import Heading from "@/app/components/breadcrumb";
import ConfigurableProductFormComponent from "@/app/components/configurable-product-form";

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
            <Heading page={pages} heading="Configurable Product" subheading="Add a new configurable product" />
            <ConfigurableProductFormComponent />
        </>
    );
}
