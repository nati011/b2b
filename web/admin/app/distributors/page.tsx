'use client'
import TableLayout from "@/app/components/table-layout";
import { columns } from "@/app/distributors/column"
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { useEffect } from "react";
import Heading from "../components/breadcrumb";
import { DataTable } from "@/components/ui/datatable";

export default function Distributors() {
  const {
    distributors,
    loading,
    error,
    fetchDistributors
  } = useDistributorsStore()

  useEffect(() => {
    fetchDistributors();
  }, []);

  const pages: any[] = [
    {
      "title": "Distributors",
      "href": `/distributors`,
    }
  ]



  return (
    <>
      <Heading page={pages} heading="Distributors" subheading="List of Registered Distributors" />
      <DataTable
        columns={columns}
        data={distributors}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Distributor", url: "/distributors/form" }}
        search="name"
        searchPlaceholder="Search distributor..."
      />
    </>

  );
}
