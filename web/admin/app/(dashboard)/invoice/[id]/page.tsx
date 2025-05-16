"use client"
import { IoPrintOutline } from "react-icons/io5";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import axios from "axios";
import Heading from "@/app/components/breadcrumb";
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/app/(dashboard)/invoice/[id]/columns";
import useOrdersStore from "@/app/libs/store/useOrderStore";
import Logo from "@/public/logo.png"
import Image from "next/image";
import { Button } from "@/components/ui/button";


// @ts-ignore
export default function InvoiceDetail({ params: { locale } }) {
    const routeParam = useParams<{ id: string }>();

    const {
        loading,
        error,
        fetchInvoice,
        invoice
    } = useOrdersStore()


    const handlePrint = () => {
        window.print();
    };

    const pages = [
        {
            name: "Invoices",
            href: "/invoice",
        },

    ];

    useEffect(() => {
        fetchInvoice(parseInt(routeParam.id))
    }, [])

    return (
        <div className="h-screen">
            <Heading page={pages} heading="Invoice" subheading="Invoice Details" />
            <div className="sm:flex gap-4 print:w-full">

                <Card
                    className="rounded-sm w-full px-4 shadow-none"
                >
                    <div className="flex items-center sm:px-4 border-gray-200 border-b-[1px] py-4 mb-2 bg-gray-100/[0.5]">
                        <div className="">
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

                    {/* <div className="customer-details">
                        <div className="logo">
                            <p className="text-lg">
                                Customer Details
                            </p>
                            <p className="text-gray-600 text-sm">
                                Abebe Kebede
                            </p>
                            <p className="text-gray-600 text-sm">
                                +2519234567132
                            </p>
                        </div>
                    </div> */}
                    <DataTable
                        columns={columns}
                        data={invoice.LineItems}
                        loading={loading}
                        button={false}
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
            <div className="w-full flex items-end p-4 flex-col gap-4 print:hidden">
                <Button variant={'outline'}
                    className="border-blue-900 text-blue-900 font-semibold px-4"
                    onClick={handlePrint}
                >
                    <IoPrintOutline />
                    <span>Print Invoice</span>
                </Button>
                {/* <Link
                        className=" flex items-center border-2 gap-1 rounded-sm border-[#1C40CA] font-medium text-[#1C40CA] text-md px-8 py-2 hover:bg-gray-200/[30%]"
                        href={`/${locale}/dashboard/invoices/edit/${routeParam?.id}/`}
                    >
                        <TbPencil />
                        <span>{dict?.editInvoice || "Edit Invoice"}</span>
                    </Link> */}
            </div>
        </div>
    );
}