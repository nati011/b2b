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
            const status = row.getValue("status") == 'COMPLETED'
            return <div className={!status ? "font-semibold py-1 rounded-md text-amber-500 text-center text-xs w-fit flex" : "border border-emerald-500  py-1 rounded-md  text-emerald-500 text-center text-xs"}>
                {row.original.Status}
            </div>
        },
    },
    {
        accessorKey: "DeliveryStatus",
        header: "Delivery Status",
        cell: ({ row }) => {
            const status = row.getValue("status") == 'COMPLETED'
            return <p className={!status ? "font-semibold py-1 rounded-md text-amber-500 text-center text-xs w-fit flex" : "border border-emerald-500  py-1 rounded-md  text-emerald-500 text-center text-xs"}>
                {row.original.DeliveryStatus}
            </p>
        },

    },
    {
        accessorKey: "PaymentStatus",
        header: "Payment Status",
        cell: ({ row }) => {
            const status = row.getValue("status") == 'COMPLETED'
            return <div className={!status ? "font-semibold py-1 rounded-md text-amber-500 text-center text-xs w-fit flex" : "border border-emerald-500  py-1 rounded-md  text-emerald-500 text-center text-xs"}>
                {row.original.PaymentStatus}
            </div>
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