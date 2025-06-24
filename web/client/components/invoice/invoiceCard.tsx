"use client"
import Image from "next/image";
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/components/invoice/column";
import { Invoice } from "@/lib/types";

interface Props {
    loading: boolean
    invoice: Invoice
}


// @ts-ignore
export const InvoiceCard: React.FC<Props> = ({
    loading,
    invoice
}) => {

    return (
        <div className="gap-4 print:w-full important">
            <Card
                className="rounded-sm px-4 shadow-none mx-auto max-w-6xl"
            >
                <div className="flex items-center sm:px-4 border-gray-200 border-b-[1px] py-4 mb-2 bg-gray-100/[0.5]">
                    <div className="flex items-center">
                        <Image src='/logo.png' width={150} height={100} alt="logo" />
                        <div className="">
                            <p className="font-semibold text-lg text-blue-900">
                                Efoyeta Store
                            </p>
                            <p className="text-gray-600 text-sm">
                                +2519234567132
                            </p>
                            <p className="text-gray-600 text-sm">
                                Addis Ababa, Ethiopia
                            </p>
                        </div>


                    </div>
                    <div className="ml-auto ">
                        <div className="flex">
                            <div className="font-semibold text-md text-neutral-700">
                                <span className="dark:text-white"> Invoice:</span>
                            </div>
                            <div className="font-light sm:text-md  text-sm text-neutral-700 ml-2">
                                <span className="dark:text-white">#{invoice?.Id}</span>
                            </div>
                        </div>
                        <div className="flex">
                            <div className="font-semibold text-md text-neutral-700">
                                <span className="dark:text-white">Date Issued:</span>
                            </div>
                            <div className="font-light sm:text-md text-sm text-neutral-700 ml-2">
                                <span className="dark:text-white">{new Date(invoice.Created_Date).toDateString()}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <DataTable
                    columns={columns}
                    data={invoice.LineItems}
                    loading={loading}
                />
                <div className="flex flex-col w-full items-end gap-4 text-gray-600 text-left">
                    <p>
                        <strong>Subtotal:</strong>{invoice.SubTotal}
                    </p>
                    <p>
                        <strong>Tax Amount:</strong>{invoice.TaxAmount}
                    </p>
                    <p>
                        <strong>Total:</strong>{invoice.SubTotal + invoice.TaxAmount}
                    </p>
                </div>
            </Card>
        </div>


    );
}