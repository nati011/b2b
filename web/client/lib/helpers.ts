
import { CANCELED_STATUS, COMPLETED_STATUS, PENDING_STATUS } from '@/lib/enums';

export const statusConfig = {
    [CANCELED_STATUS]: {
        bg: "bg-red-800",
        text: "text-red-800",
        border: "border-red-200",
        label: "Canceled"
    },
    [PENDING_STATUS]: {
        bg: "bg-amber-500",
        text: "text-amber-500",
        border: "border-amber-200",
        label: "Pending"
    },
    [COMPLETED_STATUS]: {
        bg: "bg-emerald-800",
        text: "text-emerald-800",
        border: "border-emerald-200",
        label: "Completed"
    },
};


export const getStatusConfig = (status: any) => {
    // @ts-ignore
    const s = statusConfig[status] || {
        bg: "bg-gray-100",
        text: "text-gray-800",
        border: "border-gray-200",
        label: "Unknown"
    }

    return s
};