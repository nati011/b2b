'use client'

import Heading from "@/app/components/breadcrumb";
import ConfigurableProductFormComponent from "@/app/components/configurable-product-form";
import useConfigurableProductStore from "@/app/libs/store/useConfigurableProduct";
import { useParams } from "next/navigation";
import { useEffect } from "react";

export default function Retailers() {
    const {
        product,
        fetchConfigurableProductDetail
    } = useConfigurableProductStore()
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
                initialData={product}
            />
        </>
    );
}
