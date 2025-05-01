import { ColumnDef } from "@tanstack/react-table";
import { Order } from '@/app/libs/types';

export const columns: ColumnDef<Order>[] = [
    {
        accessorKey: "Id",
        header: "Id",
    },
    {
        accessorKey: "Status",
        header: "Status",
    },
    {
        accessorKey: "DeliveryStatus",
        header: "Delivery Status",
    },
    {
        accessorKey: "PaymentStatus",
        header: "Payment Status",
    },

];