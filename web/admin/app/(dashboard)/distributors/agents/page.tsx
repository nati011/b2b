'use client'
import { DataTableLayout } from "@/components/ui/datatablelayout";
import { columns } from "@/app/(dashboard)/distributors/agents/column"
import useDistributorsStore from "@/app/libs/store/useDistributorStore"
import { useEffect } from "react";
import Heading from "@/components/breadcrumb";

export default function Distributors() {
  const {
    distributorUser,
    userCount,
    loading,
    error,
    fetchDistributorUser
  } = useDistributorsStore()

  useEffect(() => {
    fetchDistributorUser();
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
      <DataTableLayout
        columns={columns}
        data={distributorUser}
        loading={loading}
        button={true}
        buttonObj={{ name: "Register Distributor Agents", url: "/distributors/agents/form" }}
        search="name"
        searchPlaceholder="Search distributor agents..."
      />
    </>

  );
}
