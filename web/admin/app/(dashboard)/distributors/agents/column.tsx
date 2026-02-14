import { ColumnDef } from "@tanstack/react-table";
import { UserDetail } from '@/app/libs/types';

export const columns: ColumnDef<UserDetail>[] = [
    {
        accessorKey: "Id",
        header: "Id",
    },
    {
        accessorKey: "FirstName",
        header: "First Name",
    },
    {
        accessorKey: "LastName",
        header: "Last Name",
    },
    {
        accessorKey: "Email",
        header: "Email",
    },
    {
        accessorKey: "Phone",
        header: "Phone",
    },
    {
        accessorKey: "Username",
        header: "Username",
    },
];