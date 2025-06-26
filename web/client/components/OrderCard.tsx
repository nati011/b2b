"use client"

import { Calendar, Receipt } from "lucide-react";
import { PiSpinner } from "react-icons/pi";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { InvoiceDetail } from "@/components/invoice/InvoiceDialog";
import { Separator } from "@/components/ui/separator";

import { getStatusConfig, statusConfig } from "@/lib/helpers";
import { Order } from "@/lib/types";


interface Props {
    loading: boolean
    order: Order
    completePayment: (id: number) => Promise<void>
}


// @ts-ignore
export const OrderCard: React.FC<Props> = ({
    loading,
    order,
    completePayment
}) => {

    return (
        <Card key={order.Id} className="rounded-sm shadow-none hover:shadow-md transition-shadow duration-200">
            <CardHeader className="border-b px-6">
                <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                    <div className="space-y-1">
                        <div className="flex items-center gap-2">
                            <Receipt className="w-4 h-4 text-gray-500" />
                            <span className="font-semibold text-gray-900">Order #{order.Id}</span>
                        </div>
                        <div className="flex items-center gap-2 text-sm text-gray-600">
                            <Calendar className="w-4 h-4" />
                            <span>{new Date(order.CreatedAt).toLocaleDateString('en-US', {
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric'
                            })}</span>
                        </div>
                    </div>
                    {
                        order.PaymentStatus == "ACCEPTED" ? (
                            <InvoiceDetail id={order.Id} />
                        ) : (

                            <>

                                {
                                    order.PaymentStatus == "ACCEPTED" ? (
                                        <InvoiceDetail id={order.Id} />
                                    ) : (
                                        <>
                                            {
                                                order.PaymentStatus == "PENDING" ? (
                                                    <Button onClick={() => { completePayment(order.Id) }}>
                                                        {
                                                            loading ? (
                                                                <>
                                                                    <PiSpinner className="animate-spin" />
                                                                    Loading
                                                                </>
                                                            ) : (
                                                                <span>
                                                                    Complete Payment
                                                                </span>
                                                            )
                                                        }

                                                    </Button>
                                                ) : (
                                                    <></>
                                                )
                                            }
                                        </>

                                    )
                                }
                            </>

                        )
                    }

                </div>
            </CardHeader>

            <CardContent className="px-6 py-4">
                <div className="space-y-3">
                    {order.Items.map((orderItem, index) => (
                        <div key={index} className="flex justify-between items-start">
                            <div className="flex-1">
                                <h4 className="font-medium text-gray-900">{orderItem.ProductName}</h4>
                                <p className="text-sm text-gray-600">
                                    ${orderItem.ProductPrice.toFixed(2)} × {orderItem.Quantity}
                                </p>
                            </div>
                            <div className="text-right">
                                <p className="font-medium text-gray-900">
                                    ${(orderItem.ProductPrice * orderItem.Quantity).toFixed(2)}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>

            <Separator />

            <CardFooter className="px-6">
                <div className="flex justify-between items-center w-full">
                    <div className="flex items-center gap-2">
                        {/* @ts-ignore */}
                        <div className={`w-2 h-2 ${getStatusConfig(order?.DeliveryStatus).bg} rounded-full`}></div>
                        {/* @ts-ignore */}
                        <span className={`text-sm font-medium  ${getStatusConfig(order?.DeliveryStatus).text}`}>Delivery {order ? order?.DeliveryStatus.toLocaleLowerCase() : ""}</span>
                    </div>
                    <div className="text-right">
                        <p className="text-sm text-gray-600">Total</p>
                        <p className="font-bold text-lg text-gray-900">
                            ${order.Total.toFixed(2)}
                        </p>
                    </div>
                </div>
            </CardFooter>
        </Card>
    );
}