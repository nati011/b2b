'use client'
import { columns } from "@/app/(dashboard)/orders/column"
import { useEffect, useState } from "react";
import Heading from "../../../components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";
import useOrdersStore from "@/app/libs/store/useOrderStore";

export default function Orders() {
  const [offset, setOffset] = useState(0);
  const limit = 10;

  const {
    orders,
    total,
    loading,
    error,
    fetchOrders,
  } = useOrdersStore()

  useEffect(() => {
    fetchOrders(offset)
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

  const pages = [{
    "title": "Orders",
    "href": "/order"
  }]

  return (
    <>
      <Heading page={pages} heading="Orders" subheading="List of Registered Orders" />
      {error && (
        <div className="flex flex-col items-center justify-center h-full">
          <h1 className="text-2xl font-bold">Error</h1>
          <p className="text-gray-500">{error}</p>
        </div>
      )}
      <DataTableLayout
        columns={columns}
        data={orders}
        total={total || 0}
        loading={loading}
        button={false}
        search="DeliveryStatus"
        searchPlaceholder="Search orders..."
        onNext={handleNext}
        onPrevious={handlePrevious}
        canNext={canNext}
        canPrevious={canPrevious}
      />
    </>
  );
}
