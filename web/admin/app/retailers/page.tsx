'use client'
import TableLayout from "@/app/components/table-layout";
import { column } from "@/app/products/column"
import useRetailersStore from "@/app/libs/store/useRetailerStore"
import { useEffect } from "react";

export default function Retailers() {
    const {
        retailers,
        loading,
        error,
        fetchRetailers
     } = useRetailersStore()

     useEffect(() => {
        fetchRetailers();
      }, []);
    
  
  return (
    <TableLayout
       column={column}
       data={retailers}
       loading={loading}
        heading={"Retailers"}
         pages={[
          {
           "title":"Retailers",
            "href":"/retailers"
           },
         ]}
         buttonURL="/retailers/form"
      />
  );
}
