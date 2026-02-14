"use client";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/app/(dashboard)/orders/[id]/columns";
import useOrdersStore from "@/app/libs/store/useOrderStore";
import useRetailersStore from "@/app/libs/store/useRetailerStore";
import StatusBadge from "@/components/status-badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import {
  ACCEPTED_STATUS,
  CONFIRMED_STATUS,
  CONFIRM_COMMAND,
  DELIVERED_COMMAND,
  REJECTED_COMMAND,
} from '@/app/libs/enums';
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";

export default function OrderDetail() {
  const routeParam = useParams<{ id: string }>();

  const {
    loading,
    fetchOrder,
    order,
    updateOrderStatus,
  } = useOrdersStore();

  const {
    retailer,
    fetchRetailer,
  } = useRetailersStore();

  const [dialogType, setDialogType] = useState<"approve-order" | "reject-order" | "confirm-delivery" | "reject-delivery" | null>(null);

  const handlePrint = () => {
    window.print();
  };

  useEffect(() => {
    fetchOrder(parseInt(routeParam.id));
  }, [routeParam.id]);

  useEffect(() => {
    if (order.RetailerId) {
      fetchRetailer(order.RetailerId);
    }
  }, [order.RetailerId]);

  const isOrderConfirmed = order.ConfirmationStatus === CONFIRMED_STATUS;
  const isOrderDelivered = order.DeliveryStatus === ACCEPTED_STATUS;

  // Determine button visibility
  const showOrderActions = !isOrderConfirmed && !isOrderDelivered;
  const showDeliveryActions = isOrderConfirmed && !isOrderDelivered;

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
              {isOrderConfirmed && isOrderDelivered ? (
                <div className="flex gap-2">
                  <StatusBadge status={order.DeliveryStatus} />
                  <StatusBadge status={order.PaymentStatus} />
                </div>
              ) : (
                <div className="flex gap-2">
                  {showOrderActions && (
                    <>
                      <Button
                        variant={"outline"}
                        className="border-emerald-400 text-emerald-800 hover:text-emerald-600"
                        onClick={() => setDialogType("approve-order")}
                      >
                        Approve
                      </Button>
                      <Button
                        variant={"outline"}
                        className="border-red-400 text-red-800 hover:text-red-600"
                        onClick={() => setDialogType("reject-order")}
                      >
                        Reject
                      </Button>
                    </>
                  )}
                  {showDeliveryActions && (
                    <>
                      <Button
                        variant={"outline"}
                        className="border-emerald-400 text-emerald-800 hover:text-emerald-600"
                        onClick={() => setDialogType("confirm-delivery")}
                      >
                        Confirm Delivery
                      </Button>
                      <Button
                        variant={"outline"}
                        className="border-red-400 text-red-800 hover:text-red-600"
                        onClick={() => setDialogType("reject-delivery")}
                      >
                        Reject
                      </Button>
                    </>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="flex gap-2 w-full">
        <Card className="rounded-sm px-4 shadow-none w-3/4">
          <div className="font-semibold">
            <p>Order Details</p>
          </div>
          <DataTable columns={columns} data={order.Items || []} loading={loading} />
          <div className="w-full flex flex-col items-end mt-4">
            <p>
              <span className="font-semibold">Total:</span> {order.Total}
            </p>
          </div>
        </Card>

        <Card className="rounded-sm px-4 shadow-none w-1/4 h-fit">
          <div className="font-semibold">
            <p>Customer Information</p>
          </div>
          <Separator orientation="horizontal" className="my-2" />
          <div className="flex gap-4 items-center">
            <Avatar className="h-10 w-10">
              <AvatarFallback>
                {retailer.user?.first_name?.[0]?.toUpperCase()}
                {retailer.user?.last_name?.[0]?.toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div>
              <p className="font-medium">
                {retailer.user?.first_name} {retailer.user?.last_name}
              </p>
              <p className="text-gray-700 text-sm">{retailer.user?.email}</p>
            </div>
          </div>
        </Card>
      </div>

      {/* Confirm Order Dialog */}
      <Dialog open={dialogType === "approve-order"} onOpenChange={() => setDialogType(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Approve Order</DialogTitle>
          </DialogHeader>
          <Separator />
          <p>Are you sure you want to approve this order?</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogType(null)}>Cancel</Button>
            <Button
              onClick={() => {
                setDialogType(null);
                updateOrderStatus(CONFIRM_COMMAND, parseInt(routeParam.id));
              }}
            >
              Approve
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reject Order Dialog */}
      <Dialog open={dialogType === "reject-order"} onOpenChange={() => setDialogType(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Order</DialogTitle>
          </DialogHeader>
          <Separator />
          <p>Are you sure you want to reject this order?</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogType(null)}>Cancel</Button>
            <Button
              onClick={() => {
                setDialogType(null);
                updateOrderStatus(REJECTED_COMMAND, parseInt(routeParam.id));
              }}
            >
              Reject
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Confirm Delivery Dialog */}
      <Dialog open={dialogType === "confirm-delivery"} onOpenChange={() => setDialogType(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Delivery</DialogTitle>
          </DialogHeader>
          <Separator />
          <p>Are you sure you want to confirm delivery of this order?</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogType(null)}>Cancel</Button>
            <Button
              onClick={() => {
                setDialogType(null);
                updateOrderStatus(DELIVERED_COMMAND, parseInt(routeParam.id));
              }}
            >
              Confirm Delivery
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reject Delivery Dialog */}
      <Dialog open={dialogType === "reject-delivery"} onOpenChange={() => setDialogType(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Delivery</DialogTitle>
          </DialogHeader>
          <Separator />
          <p>Are you sure you want to reject delivery of this order?</p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogType(null)}>Cancel</Button>
            <Button
              onClick={() => {
                setDialogType(null);
                updateOrderStatus(REJECTED_COMMAND, parseInt(routeParam.id));
              }}
            >
              Reject
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}