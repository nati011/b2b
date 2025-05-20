'use client'
import { useEffect, useState } from "react";
import Heading from "../../components/breadcrumb";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { DataTableLayout } from "@/components/ui/datatablelayout";
import useOrdersStore from "@/app/libs/store/useOrderStore";
import { Separator } from "@radix-ui/react-dropdown-menu";
import { Transaction } from "@/app/libs/types";
import { MoreHorizontal } from "lucide-react";
import { RxCaretSort } from 'react-icons/rx'




export default function Transactions() {
    const {
        transactions,
        fetchTransactions
    } = useOrdersStore()

    const pages = [
        {
            "title": "Transactions",
            "href": "/transactions"
        }
    ]
    const [confirmTransactionOpen, setConfirmTransactionOpen] = useState(false)
    useEffect(() => {
        fetchTransactions();
    }, []);


    const columns: ColumnDef<Transaction>[] = [
        {
            accessorKey: "id",
            header: "Id",
        },
        {
            accessorKey: "tx_ref",
            header: "Transaction Ref",
        },
        {
            accessorKey: "date",
            header: ({ column }) => {
                return (
                    <Button
                        variant="ghost"
                        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    >
                        Date
                        <RxCaretSort className="ml-2 h-4 w-4" />
                    </Button>
                );
            },
            cell: ({ row }) => {
                const date = new Date(row.original.date)
                return (
                    <div className="flex items-center gap-2 text-gray-900">
                        {date.getDate()}-{date.getMonth()}-{date.getFullYear()}
                    </div>
                );
            },
        },
        {
            accessorKey: "amount",
            header: "Amount",
        },
        {
            accessorKey: "status",
            header: "Status",
            cell: ({ row }) => {
                const status = row.getValue("status") == 'COMPLETED'
                return <div className={!status ? "border border-amber-500 py-1 mx-auto rounded-md text-amber-500 font-medium text-center text-xs" : "border border-emerald-500  py-1 mx-auto rounded-md  text-emerald-500 font-medium text-center text-xs"}>
                    {row.getValue("status")}
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
                                <DropdownMenuItem onClick={() => {
                                    setConfirmTransactionOpen(true)
                                }} disabled={row.original.status == "COMPLETED"}>
                                    Confirm Payment
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                );
            },
        },
    ];


    return (
        <>
            <Heading page={pages} heading="Transactions" subheading="List of registered transactions" />
            <DataTableLayout
                columns={columns}
                data={transactions}
                search="tx_ref"
                searchPlaceholder="Search transactions..."
            />
            <Dialog open={confirmTransactionOpen} onOpenChange={setConfirmTransactionOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Confirm Transaction</DialogTitle>
                    </DialogHeader>
                    <Separator />
                    <p>
                        Are you sure you want to confirm payment?
                    </p>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setConfirmTransactionOpen(false)}>Cancel</Button>
                        <Button onClick={() => setConfirmTransactionOpen(false)}>Confirm</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

        </>
    );
}
