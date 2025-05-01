"use client";
import * as React from "react";

import {
    ColumnDef,
    ColumnFiltersState,
    SortingState,
    VisibilityState,
    flexRender,
    getCoreRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    useReactTable,
} from "@tanstack/react-table";

import { Button } from "@/components/ui/button";

import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";

import { CiFilter } from "react-icons/ci";
import Link from "next/link";
import { PiSpinner } from "react-icons/pi";

interface buttonObj {
    name: string;
    url: string;
}

interface DataTableProps<TData, TValue> {
    columns: ColumnDef<TData, TValue>[];
    data: TData[];
    button?: boolean;
    heading?: string;
    subheading?: string;
    search?: string;
    buttonObj?: buttonObj;
    title?: string;
    searchPlaceholder?: string;
    loading?: boolean
    previous?: string | null
    next?: string | null
    fetchProperties?: (url?: string) => void;
}

export function DataTable<TData, TValue>({
    columns,
    data,
    button,
    search,
    buttonObj,
    title,
    searchPlaceholder,
    loading,
    previous,
    next,
    fetchProperties
}: DataTableProps<TData, TValue>) {
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
        []
    );
    const [columnVisibility, setColumnVisibility] =
        React.useState<VisibilityState>({});
    const [sorting, setSorting] = React.useState<SortingState>([]);
    const [rowSelection, setRowSelection] = React.useState({});

    const table = useReactTable({
        data,
        columns,
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        onColumnVisibilityChange: setColumnVisibility,
        onRowSelectionChange: setRowSelection,
        state: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
        },
    });


    return (
        <div className="w-full bg-white dark:bg-black p-4 rounded-md mt-4 print:hidden">

            <div className="sm:flex customers-center py-4 gap-2 items-center">
                {button && buttonObj?.name && (
                    <div className="flex gap-2">
                        <Link href={buttonObj.url} passHref>
                            <Button
                                variant="outline"
                                size="sm"
                                className="border-blue-900 bg-slate-100 dark:bg-black hover:text-blue-900 text-blue-900 px-6 py-4 sm:mb-0 mb-2"
                            >
                                + {buttonObj.name}
                            </Button>
                        </Link>

                    </div>
                )}
                <div className={`flex ml-auto items-center ${!button ? 'w-full' : ''}`}>
                    <Input
                        placeholder={searchPlaceholder}
                        value={
                            (table.getColumn(`${search}`)?.getFilterValue() as string) ?? ""
                        }
                        onChange={(event) =>
                            table.getColumn(`${search}`)?.setFilterValue(event.target.value)
                        }
                        className="bg-transparent"
                    />
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <button className="focus:outline-none ring-none">
                                <CiFilter className="ml-2 h-4 w-4" />
                            </button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            {table
                                .getAllColumns()
                                .filter((column) => column.getCanHide())
                                .map((column) => {
                                    return (
                                        <DropdownMenuCheckboxItem
                                            key={column.id}
                                            className="capitalize"
                                            checked={column.getIsVisible()}
                                            onCheckedChange={(value: any) =>
                                                column.toggleVisibility(!!value)
                                            }
                                        >
                                            {column.id}
                                        </DropdownMenuCheckboxItem>
                                    );
                                })}
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>

            {!button && title && (
                <p className="font-semibold text-md my-4">{title}</p>
            )}
            <div className="rounded-md border">
                <Table>
                    <TableHeader className="bg-transparent hover:bg-transparent">
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow
                                key={headerGroup.id}
                                className="bg-transparent hover:bg-transparent"
                            >
                                {headerGroup.headers.map((header) => {
                                    return (
                                        <TableHead key={header.id}>
                                            {header.isPlaceholder
                                                ? null
                                                : flexRender(
                                                    header.column.columnDef.header,
                                                    header.getContext()
                                                )}
                                        </TableHead>
                                    );
                                })}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody>
                        {table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow
                                    key={row.id}
                                    data-state={row.getIsSelected() && "selected"}
                                >
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id}>
                                            {flexRender(
                                                cell.column.columnDef.cell,
                                                cell.getContext()
                                            )}
                                        </TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                {loading ? (
                                    <TableCell
                                        colSpan={columns.length}
                                        className="h-24 text-center"
                                    >
                                        <div className="flex items-center justify-center">
                                            <PiSpinner className="h-4 w-4 mr-2 animate-spin" />
                                            <p>Loading...</p>
                                        </div>
                                    </TableCell>
                                ) : (
                                    <TableCell
                                        colSpan={columns.length}
                                        className="h-24 text-center"
                                    >
                                        No data available.
                                    </TableCell>
                                )}
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            <div className="flex items-center justify-end space-x-2 py-4">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => fetchProperties ? previous && fetchProperties(previous) : table.previousPage()}
                    disabled={loading}
                    className="bg-transparent font-medium"
                >
                    Previous
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => fetchProperties ? (next && fetchProperties(next)) : table.nextPage()}
                    disabled={loading}
                    className="bg-transparent font-medium"
                >
                    Next
                </Button>
            </div>
        </div>
    );
}