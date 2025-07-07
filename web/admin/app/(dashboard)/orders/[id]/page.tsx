"use client"
import { useParams} from "next/navigation";
import { useEffect } from "react";
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/app/(dashboard)/orders/[id]/columns";
import useOrdersStore from "@/app/libs/store/useOrderStore";
import useRetailersStore from "@/app/libs/store/useRetailerStore";
import StatusBadge from "@/components/status-badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";


// @ts-ignore
export default function OrderDetail() {
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
            href: "/orders",
        },

    ];

    useEffect(() => {
        fetchOrder(parseInt(routeParam.id))
    }, [])

    useEffect(() => {
        if (order.RetailerId != 0) {
            console.log("HEERRRRRRRRRRRREEEEEEEEEEEEE__________________")
            fetchRetailer(order.RetailerId)
        }

    }, [order])
    return (
        <div className="px-20">
            <div className="flex flex-col sm:px-4 border-gray-200 border-b-[1px] py-4 mb-5">

                <div className="font-semibold">
                    <p>Order #{order.Id}</p>
                    <div className="flex justify-between">
                        <div className="text-gray-800">
                            {new Date(order.CreatedAt).toUTCString()}
                        </div>
                        <div className="">
                            <StatusBadge status={order.DeliveryStatus} />
                            <StatusBadge status={order.PaymentStatus} />
                        </div>
                    </div>

                </div>
            </div>
            <div className="flex gap-2 w-full">
                <Card
                    className="rounded-sm px-4 shadow-none w-3/4"
                >
                    <div className="font-semibold">
                        <p>Order Details</p>
                    </div>
                    <DataTable
                        columns={columns}
                        data={order.Items}
                        loading={loading}
                    />
                    <div className="w-full flex flex-col items-end">
                        <p>
                            <span className="font-semibold">Total:</span>  {order.Total}
                        </p>
                    </div>
                </Card>
                <Card
                    className="rounded-sm px-4 shadow-none w-1/4 h-fit"
                >
                    <div className="font-semibold ">
                        <p>Customer Information</p>
                    </div>
                    <Separator orientation="horizontal" />
                    <div className="flex gap-4">
                        <Avatar className='h-10 w-10'>
                            <AvatarFallback>{retailer.user.first_name.slice(0, 1).toUpperCase()}{retailer.user.last_name.slice(0, 1).toUpperCase()}</AvatarFallback>
                        </Avatar>
                        <div className="">
                            <p className="">
                                {retailer.user.first_name} {retailer.user.last_name}
                            </p>
                            <p className="text-gray-700">
                                {retailer.user.email}
                            </p>
                        </div>
                    </div>

                </Card>

            </div>
        </div>
    );
}