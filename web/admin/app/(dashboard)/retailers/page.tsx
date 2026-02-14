'use client'
import { columns } from "@/app/(dashboard)/retailers/column"
import useRetailersStore from "@/app/libs/store/useRetailerStore"
import { useEffect, useState } from "react";
import Heading from "../../../components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";

export default function Retailers() {
  const [offset, setOffset] = useState(0);
  const limit = 10;

  const {
    retailers,
    loading,
    error,
    total,
    fetchRetailers
  } = useRetailersStore()



  useEffect(() => {
    fetchRetailers(offset)
  }, [offset]);

  const handleNext = () => {
    if (offset + limit < (total || 0)) {
      setOffset(offset + limit);
    }
  };
  const handlePrevious = () => {
    if (offset - limit >= 0) {
      setOffset(offset - limit);
    }
  };
  const canNext = offset + limit < (total || 0);
  const canPrevious = offset > 0;
  const pages = [
    {
      "title": "Retailers",
      "href": "/retailers"
    },
  ]

  return (
    <>
      <Heading page={pages} heading="Retailers" subheading="List of Registered Retailers" />
      {error ? (
        <div className="flex flex-col items-center justify-center h-full">
          <h1 className="text-2xl font-bold">Error</h1>
          <p className="text-gray-500">{error}</p>
        </div>
      ):(
        <DataTableLayout
        columns={columns}
        data={retailers}
        loading={loading}
        search="name"
        searchPlaceholder="Search Retailers..."
        onNext={handleNext}
        onPrevious={handlePrevious}
        canNext={canNext}
        canPrevious={canPrevious}
        total={total || 0}
      />
      )}


    </>
  );
}
