'use client'
import { columns } from "@/app/(dashboard)/orders/column"
import useOrdersStore from "@/app/libs/store/useOrderStore"
import { useEffect } from "react";
import Heading from "../../components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";

export default function Orders() {
  const {
    orders,
    loading,
    error,
    fetchOrders
  } = useOrdersStore()

  useEffect(() => {
    fetchOrders();
  }, []);

  const pages = [{
    "title": "Orders",
    "href": "/order"
  }]


  return (
    <>
      <Heading page={pages} heading="Orders" subheading="List of Registered Orders" />
      <DataTableLayout
        columns={columns}
        data={orders}
        loading={loading}
        button={false}
        search="DeliveryStatus"
        searchPlaceholder="Search orders..."
      />
    </>

  );
}
