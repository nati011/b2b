
import { ColumnDef } from "@tanstack/react-table";
import { Product } from '@/app/libs/types';

export const columns: ColumnDef<Product>[] = [
    {
        accessorKey: "Id",
        header: "Id",
    },
    {
        accessorKey: "Name",
        header: "Name",
    },
    {
        accessorKey: "Price",
        header: "Price",
    },
    {
        accessorKey: "Stock",
        header: "Stock",
    },
    {
        accessorKey: "IsActive",
        header: () => <div className="text-left">Status</div>,
        cell: ({ row }) => {
            const status = row.getValue("IsActive")
            return <div className={status ? "border border-amber-500 py-1 mx-auto rounded-md text-amber-500 font-medium text-center text-xs" : "border border-emerald-500  py-1 mx-auto rounded-md  text-emerald-500 font-medium text-center text-xs"}>
                {status ? "Inactive" : "Active"}
            </div>
        },
    },
];
