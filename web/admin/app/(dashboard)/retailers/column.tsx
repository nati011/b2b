import { Retailer } from "../../libs/types";
import { ColumnDef } from "@tanstack/react-table";

export const columns: ColumnDef<Retailer>[] = [
    {
        accessorKey: "id",
        header: "Id",
    },
    {
        accessorKey: "name",
        header: "Name",
    },
    {
        accessorKey: "phone",
        header: ({ column }) => {
            return (
                <p>Phone</p>
            );
    },
    cell: ({ row }) => {
        return (
            <div className="flex items-center gap-2 text-gray-900">
                {row.original.user.phone}
            </div>
        );
    },
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
    {
        accessorKey: "is_active",
        header: () => <div className="text-left">Status</div>,
        cell: ({ row }) => {
          const status = row.original.user.is_active
          return <div className={!status ? "border border-amber-500 py-1 mx-auto rounded-md text-amber-500 font-medium text-center text-xs" : "border border-emerald-500  py-1 mx-auto rounded-md  text-emerald-500 font-medium text-center text-xs"}>
            {status ? "Active" : "Inactive"}
          </div>
        },
      },
];

