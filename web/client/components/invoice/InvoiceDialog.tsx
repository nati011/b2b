"use client"
import Image from "next/image";
import { IoPrintOutline } from "react-icons/io5";

import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogFooter,
    DialogTrigger,
} from "@/components/ui/dialog"
import { Card } from "@/components/ui/card";
import { DataTable } from "@/components/ui/datatable";
import { columns } from "@/components/invoice/column";
import { Button } from "@/components/ui/button";
import { InvoiceSkeleton } from "@/components/InvoiceSkeleton";
import useInvoiceStore from "@/lib/store/useInvoiceStore";
import { InvoiceCard } from "./invoiceCard";

interface Props {
    id: number
}


// @ts-ignore
export const InvoiceDetail: React.FC<Props> = ({
    id
}) => {
    const handlePrint = () => {
        window.print();
    };
    const {
        invoice,
        loading,
        fetchInvoice,
    } = useInvoiceStore()
    return (

        <Dialog>
            <DialogTrigger onClick={() => { fetchInvoice(id) }}><Button>View Invoice</Button></DialogTrigger>
            <DialogContent className="print:border-0 print:shadow-none">
                {
                    !loading ? (
                        <div className="grid grid-cols-1 gap-4 p-4">
                            <InvoiceCard invoice={invoice} loading={loading} />
                            <DialogFooter className="print:hidden">
                                <DialogClose asChild>
                                    <Button variant="outline">Cancel</Button>
                                </DialogClose>
                                <Button onClick={handlePrint}><IoPrintOutline /><span>Print Invoice</span></Button>
                            </DialogFooter>
                        </div>

                    ) : (
                        <InvoiceSkeleton />
                    )
                }

            </DialogContent>
        </Dialog>


    );
}