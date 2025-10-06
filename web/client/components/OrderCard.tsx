"use client"

import { Calendar, Receipt } from "lucide-react";
import { PiSpinner } from "react-icons/pi";
import React, { useEffect, useState } from "react";

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
    cancelOrder: (id: number) => Promise<void>
}


// @ts-ignore
export const OrderCard: React.FC<Props> = ({
    loading,
    order,
    completePayment,
    cancelOrder
}) => {
    console.log(order)
    // Countdown logic
    const [timeLeft, setTimeLeft] = useState<string>("");

    useEffect(() => {
        if (order.PaymentStatus === "PENDING" && order.ExpiresAt) {
            console.log(order.ExpiresAt)
            const interval = setInterval(() => {
                const now = new Date();
                const expires = new Date(order.ExpiresAt);
                const diff = expires.getTime() - now.getTime();
                if (diff <= 0) {
                    setTimeLeft("Expired");
                    clearInterval(interval);
                } else {
                    const hours = Math.floor(diff / (1000 * 60 * 60));
                    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
                    const seconds = Math.floor((diff % (1000 * 60)) / 1000);
                    setTimeLeft(
                        hours > 0
                            ? `${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`
                            : `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`
                    );
                }
            }, 1000);
            return () => clearInterval(interval);
        } else {
            setTimeLeft("");
        }
    }, [order.PaymentStatus, order.ExpiresAt]);

    return (
        <Card
            key={order.Id}
            className="rounded-sm shadow-none hover:shadow-md transition-shadow duration-200 border-accent border-2"
        >
            <CardHeader className="border-b px-6 border-red-100">
                <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
                    <div className="space-y-1">
                        <div className="flex items-center gap-2">
                            <Receipt className="w-4 h-4 text-primary" />
                            <span className="font-semibold">Order #{order.Id}</span>
                            {order.PaymentStatus === "PENDING" && timeLeft && (
                                <span className="ml-2 px-2 py-0.5 rounded text-xs font-semibold bg-amber-100/50 text-amber-600 border border-amber-200">
                                    Expires in {timeLeft}
                                </span>
                            )}
                        </div>
                        <div className="flex items-center gap-2 text-sm text-gray-600">
                            <Calendar className="w-4 h-4 text-primary" />
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
                                                    <div className="flex gap-4">
                                                        <Button variant={'outline'} className="border-red-900" onClick={() => { completePayment(order.Id) }} disabled={timeLeft==""}>
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
                                                        <Button variant={'outline'} className="border-gray-400" onClick={() => { cancelOrder(order.Id) }}>
                                                            {
                                                                loading ? (
                                                                    <>
                                                                        <PiSpinner className="animate-spin" />
                                                                        Loading
                                                                    </>
                                                                ) : (
                                                                    <span>
                                                                        Cancel Order
                                                                    </span>
                                                                )
                                                            }
                                                        </Button>
                                                    </div>
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
                                <h4 className="font-medium text-primary">{orderItem.ProductName}</h4>
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
                        <p className="font-bold text-lg text-primary">
                            ${order.Total.toFixed(2)}
                        </p>
                    </div>
                </div>
            </CardFooter>
        </Card>
    );
}