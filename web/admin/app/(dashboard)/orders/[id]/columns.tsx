import { Item } from "@/app/libs/types";
import { ColumnDef } from "@tanstack/react-table";

export const columns: ColumnDef<Item>[] = [
    {
        accessorKey: "ProductId",
        header: "Id",
    },
    {
        accessorKey: "ProductName",
        header: "Name",
    },
    {
        accessorKey: "Quantity",
        header: "Quantity",
    }
];

