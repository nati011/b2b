'use client'
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/app/distributors/agents/column"
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { useEffect } from "react";
import Heading from "@/app/components/breadcrumb";

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
      "title": "Distributor Agents",
      "href": `/distributors/agents`,
    }
  ]


  return (
    <>
      <Heading page={pages} heading="Distributor Agents" subheading="List of Registered Distributor Agents" />
      <DataTable
        columns={columns}
        data={distributors}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Distributor Agents", url: "/distributors/agents/form" }}
        search="name"
        searchPlaceholder="Search distributor agents..."
      />
    </>

  );
}
