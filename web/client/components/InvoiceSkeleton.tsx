import { Skeleton } from "@/components/ui/skeleton"

export const InvoiceSkeleton = () => {
    return (
        <div className="flex flex-col space-y-3">
            <Skeleton className="h-[150px] w-xl" />
            <div className="space-y-2 align-right">
                <Skeleton className="h-4 w-xl" />
                <Skeleton className="h-4 w-xl" />
            </div>
            {/* <Skeleton className="h-[100px] w-xl" /> */}
            <div className="space-y-2">
                <Skeleton className="h-4 w-xl" />
                <Skeleton className="h-4 w-xl" />
            </div>
        </div>
    )
}