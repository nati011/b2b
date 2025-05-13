import { ColumnDef } from "@tanstack/react-table";
import { Distributor } from '@/app/libs/types';

export const columns: ColumnDef<Distributor>[] = [
    {
        accessorKey: "id",
        header: "Id",
    },
    {
        accessorKey: "name",
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