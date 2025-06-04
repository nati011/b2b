'use client'
import { columns } from "@/app/(dashboard)/orders/column"
import { useEffect, useState } from "react";
import Heading from "../../../components/breadcrumb";
import { DataTableLayout } from "@/components/ui/datatablelayout";
import { fetchOrders } from "@/actions/order";
import { Order } from "@/app/libs/types";

export default function Orders() {
  const [orders, setOrder] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const loadOrder = async () => {
      try {
        const data = await fetchOrders();
        setOrder(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    loadOrder();
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
