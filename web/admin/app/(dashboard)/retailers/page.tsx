'use client'
import { columns } from "@/app/(dashboard)/retailers/column"
import useRetailersStore from "@/app/libs/store/useRetailerStore"
import { useEffect } from "react";
import Heading from "../../../components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";

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
  const pages = [
    {
      "title": "Retailers",
      "href": "/retailers"
    },
  ]

  return (
    <>
      <Heading page={pages} heading="Retailers" subheading="List of Registered Retailers" />
      <DataTableLayout
        columns={columns}
        data={retailers}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Retailers", url: "/Retailers/form" }}
        search="name"
        searchPlaceholder="Search Retailers..."
      />

    </>
  );
}
