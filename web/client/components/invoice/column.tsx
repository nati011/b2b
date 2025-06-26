import { InvoiceItem } from "@/lib/types";
import { ColumnDef } from "@tanstack/react-table";

export const columns: ColumnDef<InvoiceItem>[] = [
    {
        accessorKey: "ProductId",
        header: "Id",
    },
    {
        accessorKey: "ProductName",
        header: "Name",
    },
    {
        accessorKey: "ProductQuantity",
        header: "Quantity",
    },
    {
        accessorKey: "ProductPrice",
        header: "Price",
    }
];

