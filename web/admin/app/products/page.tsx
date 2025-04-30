'use client'
import TableLayout from "@/app/components/table-layout";
import { column } from "@/app/products/column"
import useProductsStore from "@/app/libs/store/useProductStore"
import { useEffect } from "react";

export default function Products() {
    const {
        products,
        loading,
        error,
        fetchProducts
     } = useProductsStore()

     useEffect(() => {
        fetchProducts();
      }, []);
    
  
  return (
    <TableLayout
       column={column}
       data={products}
       loading={loading}
        heading={"Products"}
         pages={[
          {
           "title":"Products",
            "href":"/products"
           },
         ]}
         buttonURL="/products/form"
      />
  );
}
