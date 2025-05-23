'use client'

import Heading from "@/components/breadcrumb";
import ConfigurableProductFormComponent from "@/components/configurable-product-form";
import useProductStore from "@/app/libs/store/useProductStore";
import { useParams } from "next/navigation";
import { useEffect } from "react";

export default function Retailers() {
    const {
        configurable_product,
        fetchConfigurableProductDetail
    } = useProductStore()
    const routeParam = useParams<{ id: string }>();

    const pages = [
        {
            "title": "Products",
            "href": "/product"
        },
        {
            "title": "Edit",
            "href": "/product/form"
        }
    ]

    useEffect(() => {
        fetchConfigurableProductDetail(parseInt(routeParam.id));
    }, []);
    return (
        <>
            <Heading page={pages} heading="Configurable Product" subheading="Edit configurable product" />
            <ConfigurableProductFormComponent
                isEdit
                loading
                initialData={configurable_product}
            />
        </>
    );
}
