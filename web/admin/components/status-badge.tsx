import React from 'react';
import { ScrollArea } from '@/components/ui/scroll-area';
import { CANCELED_STATUS, COMPLETED_STATUS, PENDING_STATUS } from '@/app/libs/enums';


export default function StatusBadge({
    status
}: {
    status: string;
}) {
    const statusConfig = {
        [CANCELED_STATUS]: {
            bg: "bg-blue-100/50",
            text: "text-blue-800",
            border: "border-blue-200",
            label: "Canceled"
        },
        [PENDING_STATUS]: {
            bg: "bg-amber-100/50",
            text: "text-amber-800",
            border: "border-amber-200",
            label: "Pending"
        },
        [COMPLETED_STATUS]: {
            bg: "bg-emerald-100/50",
            text: "text-emerald-800",
            border: "border-emerald-200",
            label: "Completed"
        },
    };
    //@ts-ignore
    const config = statusConfig[status as string] || {
        bg: "bg-gray-100",
        text: "text-gray-800",
        border: "border-gray-200",
        label: "Unknown"
    };
    return (
        <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.bg} ${config.text} ${config.border}`}>
            {config.label}
        </div>
    );
}
