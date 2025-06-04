import { ColumnDef } from "@tanstack/react-table";
import { Order } from '@/app/libs/types';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal } from "lucide-react";
import Link from "next/link";
import { ACCEPTED_STATUS, CANCELED_STATUS, COMPLETED_STATUS, PENDING_STATUS } from "@/app/libs/enums";


export const columns: ColumnDef<Order>[] = [
    {
        accessorKey: "Id",
        header: "Id",
    },

    {
        accessorKey: "RetailerName",
        header: "Retailer Name",
    },
    {
        accessorKey: "Status",
        header: "Status",
        cell: ({ row }) => {
            const status = row.getValue("Status");
            console.log(status)
            const statusConfig = {
                [CANCELED_STATUS]: {
                    bg: "bg-red-100/50",
                    text: "text-red-800",
                    border: "border-red-200",
                    label: "Canceled"
                },
                [PENDING_STATUS]: {
                    bg: "bg-amber-100/50",
                    text: "text-amber-800",
                    border: "border-amber-200",
                    label: "Pending"
                },
                [COMPLETED_STATUS]: {
                    bg: "bg-emerald-100/50",
                    text: "text-emerald-800",
                    border: "border-emerald-200",
                    label: "Completed"
                },
            };
            //   @ts-ignore
            const config = statusConfig[status as string] || {
                bg: "bg-gray-100",
                text: "text-gray-800",
                border: "border-gray-200",
                label: "Unknown"
            };

            return (
                <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.bg} ${config.text} ${config.border}`}>
                    {config.label}
                </div>
            );
        },
    },
    {
        accessorKey: "DeliveryStatus",
        header: "Delivery Status",
        cell: ({ row }) => {
            const status = row.getValue("DeliveryStatus")
            const statusConfig = {
                [CANCELED_STATUS]: {
                    bg: "bg-blue-100/50",
                    text: "text-blue-800",
                    border: "border-blue-200",
                    label: "Dispatched"
                },
                [PENDING_STATUS]: {
                    bg: "bg-amber-100/50",
                    text: "text-amber-800",
                    border: "border-amber-200",
                    label: "Pending"
                },
                [COMPLETED_STATUS]: {
                    bg: "bg-emerald-100/50",
                    text: "text-emerald-800",
                    border: "border-emerald-200",
                    label: "Completed"
                },
            };
            //   @ts-ignore
            const config = statusConfig[status as string] || {
                bg: "bg-gray-100",
                text: "text-gray-800",
                border: "border-gray-200",
                label: "Unknown"
            };

            return (
                <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.bg} ${config.text} ${config.border}`}>
                    {config.label}
                </div>
            );
        },

    },
    {
        accessorKey: "PaymentStatus",
        header: "Payment Status",
        cell: ({ row }) => {
            const status = row.getValue("PaymentStatus")
            const statusConfig = {
                [PENDING_STATUS]: {
                    bg: "bg-amber-100/50",
                    text: "text-amber-800",
                    border: "border-amber-200",
                    label: "Pending"
                },
                [ACCEPTED_STATUS]: {
                    bg: "bg-emerald-100/50",
                    text: "text-emerald-800",
                    border: "border-emerald-200",
                    label: "Accepted"
                },
            };
            //   @ts-ignore
            const config = statusConfig[status as string] || {
                bg: "bg-gray-100",
                text: "text-gray-800",
                border: "border-gray-200",
                label: "Unknown"
            };

            return (
                <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.bg} ${config.text} ${config.border}`}>
                    {config.label}
                </div>
            );
        },
    },

    {
        id: "actions",
        enableHiding: false,
        cell: ({ row }) => {
            return (
                <div className="flex items-center gap-2">

                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <button className="h-8 w-8 p-0">
                                <span className="sr-only">Open menu</span>
                                <MoreHorizontal className="h-4 w-4" />
                            </button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuLabel>Actions</DropdownMenuLabel>
                            <Link href={`/orders/${row.original.Id}`}>
                                <DropdownMenuItem>
                                    Order Detail
                                </DropdownMenuItem>
                            </Link>
                            <Link href={`/invoice/${row.original.Id}`}>
                                <DropdownMenuItem>
                                    Invoice
                                </DropdownMenuItem>
                            </Link>

                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            );
        },
    },

];