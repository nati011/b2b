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

import { Input } from "@/components/ui/input";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";

import { Search } from "lucide-react";
import Link from "next/link";
import { PiSpinner } from "react-icons/pi";


interface buttonObj {
    name: string;
    url: string;
}

interface DataTableLayoutProps<TData, TValue> {
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
    total?: number;
    fetchProperties?: (url?: string) => void;
    onNext?: () => void;
    onPrevious?: () => void;
    canNext?: boolean;
    canPrevious?: boolean;
}

export function DataTableLayout<TData, TValue>({
    columns,
    data,
    button,
    search,
    buttonObj,
    searchPlaceholder,
    loading,
    total,
    onNext,
    onPrevious,
    canNext,
    canPrevious
}: DataTableLayoutProps<TData, TValue>) {
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
        <div className="w-full bg-card rounded-lg border border-border/50 shadow-sm">
            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-end gap-4 items-start sm:items-center p-6 border-b border-border/50 space-y-4 sm:space-y-0">
                <div className="flex items-center space-x-4 w-full sm:w-auto">
                    <div className="relative flex-1 sm:flex-none">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <Input
                            placeholder={searchPlaceholder || "Search..."}
                            value={
                                (table.getColumn(`${search}`)?.getFilterValue() as string) ?? ""
                            }
                            onChange={(event) =>
                                table.getColumn(`${search}`)?.setFilterValue(event.target.value)
                            }
                            className="pl-10 bg-background border-border focus:border-primary transition-colors rounded-sm"
                        />
                    </div>
                </div>
                
                <div className="flex items-center space-x-2">
                    {button && buttonObj?.name && (
                        <Link href={buttonObj.url} passHref>
                            <Button
                                className="bg-primary hover:bg-primary/90 text-primary-foreground shadow-sm transition-all duration-200"
                            >
                                + {buttonObj.name}
                            </Button>
                        </Link>
                    )}
                    
                </div>
            </div>

            {/* Table */}
            <div className="rounded-b-lg overflow-hidden">
                <Table>
                    <TableHeader className="bg-muted/30">
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow
                                key={headerGroup.id}
                                className="border-border hover:bg-muted/50 transition-colors"
                            >
                                {headerGroup.headers.map((header) => {
                                    return (
                                        <TableHead key={header.id} className="font-semibold text-foreground">
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
                        {table.getRowModel().rows?.length && !loading ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow
                                    key={row.id}
                                    data-state={row.getIsSelected() && "selected"}
                                    className="border-border hover:bg-muted/20 transition-colors"
                                >
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id} className="py-3">
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
                                        className="h-32 text-center"
                                    >
                                        <div className="flex flex-col items-center justify-center space-y-2">
                                            <PiSpinner className="h-6 w-6 animate-spin text-muted-foreground" />
                                            <p className="text-muted-foreground">Loading data...</p>
                                        </div>
                                    </TableCell>
                                ) : (
                                    <TableCell
                                        colSpan={columns.length}
                                        className="h-32 text-center"
                                    >
                                        <div className="flex flex-col items-center justify-center space-y-2">
                                            <div className="w-12 h-12 bg-muted rounded-full flex items-center justify-center">
                                                <Search className="h-6 w-6 text-muted-foreground" />
                                            </div>
                                            <p className="text-muted-foreground font-medium">No data available</p>
                                            <p className="text-sm text-muted-foreground">Try adjusting your search or filters</p>
                                        </div>
                                    </TableCell>
                                )}
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>

            {/* Pagination */}
            <div className="flex items-center justify-between px-6 py-4 border-t border-border/50 bg-muted/20">
                <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                    <p>
                        Showing {table.getFilteredRowModel().rows.length} of{" "}
                        {total==null ? table.getFilteredRowModel().rows.length : total} results
                    </p>
                </div>
                <div className="flex items-center space-x-2">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={onPrevious}
                        disabled={!canPrevious}
                        className="border-border"
                    >
                        Previous
                    </Button>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={onNext}
                        disabled={!canNext}
                        className="border-border"
                    >
                        Next
                    </Button>
                </div>
            </div>
        </div>
    );
}