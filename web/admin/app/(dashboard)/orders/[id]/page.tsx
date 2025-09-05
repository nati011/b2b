"use client"
import { useParams} from "next/navigation";
import { useEffect, useState } from "react";
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/app/(dashboard)/orders/[id]/columns";
import useOrdersStore from "@/app/libs/store/useOrderStore";
import useRetailersStore from "@/app/libs/store/useRetailerStore";
import StatusBadge from "@/components/status-badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { DELIVERED_STATUS, CONFIRMED_STATUS, CONFIRM_COMMAND,DELIVERED_COMMAND,REJECTED_COMMAND } from '@/app/libs/enums';
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";

// @ts-ignore
export default function OrderDetail() {
    const routeParam = useParams<{ id: string }>();

    const {
        loading,
        fetchOrder,
        order,
        updateOrderStatus
    } = useOrdersStore()

    const {
        retailer,
        fetchRetailer
    } = useRetailersStore()

    const [confirmApproveDialog, setConfirmDialog] = useState(false)
    const [confirmRejectDialog, setRejectDialog] = useState(false)


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
                        <div className="flex flex-col gap-4">
                            {
                                (order.ConfirmationStatus == CONFIRMED_STATUS && order.DeliveryStatus == DELIVERED_STATUS) ? (
                                    <div className="flex gap-2">
                                        <StatusBadge status={order.DeliveryStatus} />
                                        <StatusBadge status={order.PaymentStatus} />
                                    </div>
        
                                ) : (
                                    <>
                                        {
                                            (order.ConfirmationStatus == CONFIRMED_STATUS) ? (
                                            <div className="flex gap-2">
                                                <Button variant={"outline"} className="border-emerald-400 text-emerald-800 hover:text-emerald-600" onClick={()=>setConfirmDialog(true)}>Approve</Button>
                                                <Button variant={"outline"} className="border-emerald-400 text-red-800 hover:text-red-600" onClick={()=>setRejectDialog(true)}>Reject</Button>
                                            </div>
                                            ) : (
                                                <div className="flex gap-2">
                                                    <Button variant={"outline"} className="border-emerald-400 text-emerald-800 hover:text-emerald-600" onClick={()=>setConfirmDialog(true)}>Confirm Delivery</Button>
                                                    <Button variant={"outline"} className="border-emerald-400 text-red-800 hover:text-red-600" onClick={()=>setRejectDialog(true)}>Reject</Button>
                                                </div>
                                            )
                                        }
                                    </>
                                )
                            }
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
            <Dialog open={confirmApproveDialog} onOpenChange={setConfirmDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Confirm Order</DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <p>
                        Are you sure you want to confirm order?
                    </p>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setConfirmDialog(false)}>Cancel</Button>
                        <Button onClick={() => {
                            setConfirmDialog(false)
                            updateOrderStatus(CONFIRM_COMMAND, parseInt(routeParam.id))
                            }}>Confirm</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
            <Dialog open={confirmRejectDialog} onOpenChange={setRejectDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Reject Order</DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <p>
                        Are you sure you want to reject order?
                    </p>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setRejectDialog(false)}>Cancel</Button>
                        <Button onClick={() => {
                            setRejectDialog(false)
                            updateOrderStatus(REJECTED_COMMAND, parseInt(routeParam.id))
                            }}>Confirm</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <Dialog open={confirmApproveDialog} onOpenChange={setConfirmDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Confirm Order Delivery</DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <p>
                        Are you sure you want to confirm order delivery?
                    </p>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setConfirmDialog(false)}>Cancel</Button>
                        <Button onClick={() => {
                            setConfirmDialog(false)
                            updateOrderStatus(CONFIRM_COMMAND, parseInt(routeParam.id))
                            }}>Confirm</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
            <Dialog open={confirmRejectDialog} onOpenChange={setRejectDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Reject Order Delivery</DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <p>
                        Are you sure you want to reject order delivery?
                    </p>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setRejectDialog(false)}>Cancel</Button>
                        <Button onClick={() => {
                            setRejectDialog(false)
                            updateOrderStatus(REJECTED_COMMAND, parseInt(routeParam.id))
                            }}>Confirm</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}