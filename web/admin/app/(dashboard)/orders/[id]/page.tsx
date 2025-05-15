"use client"
import { useParams, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import Heading from "@/app/components/breadcrumb";
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/app/(dashboard)/orders/[id]/columns";
import useOrdersStore from "@/app/libs/store/useOrderStore";
import useRetailersStore from "@/app/libs/store/useRetailerStore";


// @ts-ignore
export default function OrderDetail({ params: { locale } }) {
    const routeParam = useParams<{ id: string }>();

    const {
        loading,
        error,
        fetchOrder,
        order
    } = useOrdersStore()

    const {
        retailer,
        fetchRetailer
    } = useRetailersStore()


    const handlePrint = () => {
        window.print();
    };

    const pages = [
        {
            name: "Orders",
            href: "/Order",
        },

    ];

    useEffect(() => {
        fetchOrder(parseInt(routeParam.id))
    }, [])

    useEffect(() => {
        if (order) {
            fetchRetailer(order.RetailerId)
        }
    }, [order])
    return (
        <div className="px-20">
            <div className="flex flex-col sm:px-4 border-gray-200 border-b-[1px] py-4 mb-5">
                <div className="font-semibold">
                    <p>Order #{order.Id}</p>
                    <div className="bg-green-100[0.5] text-green-900 rounded-full">
                        {order.DeliveryStatus}
                    </div>
                    <div className="bg-green-100[0.5] text-green-900 rounded-full">
                        {order.PaymentStatus}
                    </div>
                </div>
                <div className="text-gray-800">
                    {new Date(order.CreatedAt).toUTCString()}
                </div>
            </div>
            <div className="flex gap-2">
                <Card
                    className="rounded-sm w-full px-4 shadow-none"
                >
                    <div className="font-semibold">
                        <p>Order Details</p>
                    </div>
                    <DataTable
                        columns={columns}
                        data={order.Items}
                        loading={loading}
                        button={false}
                    />
                    <div className="w-full flex flex-col items-end">
                        <p>
                            <span className="font-semibold">Total:</span>  {order.Total}
                        </p>
                    </div>
                </Card>
                <Card
                    className="rounded-sm px-4 shadow-none"
                >
                    <div className="font-semibold">
                        <p>Customer Information</p>
                    </div>
                    <div className="">
                        <p className="">
                            {retailer.user.first_name} {retailer.user.last_name}
                        </p>
                        <p className="text-gray-700">
                            {retailer.user.email}
                        </p>
                    </div>
                </Card>

            </div>
        </div>
    );
}