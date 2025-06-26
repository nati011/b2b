import { Skeleton } from "@/components/ui/skeleton"

export const CardSkeleton = () => {
    return (
        <div className="flex flex-col space-y-3">
            <Skeleton className="h-[300px] w-[300px]" />
            <div className="space-y-2">
                <Skeleton className="h-4 w-[250px]" />
                <Skeleton className="h-4 w-[200px]" />
            </div>
        </div>
    )
}