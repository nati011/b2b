
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

const ThankYouSkeleton = () => {
    return (
        <div className="min-h-screen">
            <main className="container mx-auto px-4 py-8">
                <div className="max-w-6xl mx-auto">
                    {/* Success Message Skeleton */}
                    <Card className="mb-8 border-green-200 bg-green-50 shadow-none rounded-sm">
                        <CardContent className="p-6 text-center">
                            <Skeleton className="w-16 h-16 mx-auto mb-4 rounded-full" />
                            <Skeleton className="h-8 w-80 mx-auto mb-2" />
                            <Skeleton className="h-5 w-60 mx-auto mb-2" />
                            <Skeleton className="h-4 w-48 mx-auto" />
                        </CardContent>
                    </Card>

                    {/* Invoice Skeleton */}
                    <Card className="shadow-none rounded-sm">
                        <CardHeader className="border-b">
                            <div className="flex justify-between items-start">
                                <div>
                                    <Skeleton className="h-7 w-24 mb-2" />
                                    <Skeleton className="h-5 w-32 mb-1" />
                                    <Skeleton className="h-4 w-28" />
                                </div>
                                <div className="text-right">
                                    <Skeleton className="h-6 w-28 mb-2" />
                                    <Skeleton className="h-4 w-36 mb-1" />
                                    <Skeleton className="h-4 w-32 mb-1" />
                                    <Skeleton className="h-4 w-40" />
                                </div>
                            </div>
                        </CardHeader>

                        <CardContent className="p-6">
                            {/* Order Items Skeleton */}
                            <div className="mb-6">
                                <Skeleton className="h-6 w-28 mb-4" />
                                <div className="space-y-3">
                                    {[1, 2, 3].map((item) => (
                                        <div key={item} className="flex justify-between items-start">
                                            <div className="flex-1">
                                                <Skeleton className="h-5 w-48 mb-1" />
                                                <Skeleton className="h-4 w-32" />
                                            </div>
                                            <div className="text-right">
                                                <Skeleton className="h-5 w-16" />
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="border-t my-6"></div>

                            {/* Order Summary Skeleton */}
                            <div className="space-y-3">
                                <div className="flex justify-between">
                                    <Skeleton className="h-4 w-16" />
                                    <Skeleton className="h-4 w-20" />
                                </div>
                                <div className="flex justify-between">
                                    <Skeleton className="h-4 w-12" />
                                    <Skeleton className="h-4 w-16" />
                                </div>

                                <div className="border-t my-3"></div>

                                <div className="flex justify-between">
                                    <Skeleton className="h-6 w-16" />
                                    <Skeleton className="h-6 w-24" />
                                </div>
                            </div>

                            {/* Footer Skeleton */}
                            <div className="mt-8 pt-6 text-center">
                                <Skeleton className="h-4 w-52 mx-auto mb-2" />
                                <Skeleton className="h-4 w-72 mx-auto" />
                            </div>
                        </CardContent>
                    </Card>

                    {/* Additional Actions Skeleton */}
                    <div className="mt-8 text-center">
                        <div className="flex flex-col sm:flex-row gap-4 justify-center">
                            <Skeleton className="h-10 w-40" />
                            <Skeleton className="h-10 w-32" />
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default ThankYouSkeleton;