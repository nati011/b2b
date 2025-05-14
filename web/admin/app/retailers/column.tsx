import { Retailer } from "../libs/types";
import { ColumnDef } from "@tanstack/react-table";

export const columns: ColumnDef<Retailer>[] = [
    {
        accessorKey: "Id",
        header: "Id",
    },
    {
        accessorKey: "Name",
        header: "Name",
    },
    {
        accessorKey: "tin",
        header: "Tin",
    },
    {
        accessorKey: "general_zone",
        header: "General Zone",
    },
    {
        accessorKey: "region",
        header: "Region",
    },
    {
        accessorKey: "woreda",
        header: "Woreda",
    },
];

