'use client'
import { useState } from "react";
import { MailCheckIcon } from "lucide-react";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import useOrdersStore from "@/lib/store/useOrderStore";
import Link from "next/link";

const CheckYourEmail = () => {
    const {
        loading,
        error
    } = useOrdersStore()
    const [currentPage, setCurrentPage] = useState(0)

    if (loading) {
        return (
            <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
                <div className="space-y-6">
                    <div>
                        <Skeleton className="h-8 w-48" />
                        <Skeleton className="h-4 w-96 mt-2" />
                    </div>
                    {[...Array(3)].map((_, i) => (
                        <Card key={i}>
                            <CardHeader>
                                <Skeleton className="h-6 w-full" />
                            </CardHeader>
                            <CardContent>
                                <Skeleton className="h-4 w-full" />
                            </CardContent>
                        </Card>
                    ))}
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
                <div className="text-center">
                    <h3 className="font-semibold text-2xl mb-4">Something went wrong</h3>
                    <p className="text-gray-600 mb-6">{error}</p>
                    <Link href={"/login"}>
                        <Button>Try again</Button>
                    </Link>
                </div>
            </div>
        );
    }

    return (
        <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            <div>
                <div className="w-full flex flex-col gap-6 items-center justify-center text-center py-16">
                    <div className="bg-emerald-50 rounded-full p-8 mb-4">
                        <MailCheckIcon className="w-16 h-16 text-emerald-400" />
                    </div>
                    <div>
                        <h4 className="text-2xl font-semibold text-gray-900 mb-2">
                            Email Sent Successfully
                        </h4>
                        <p className="text-gray-600 mb-6 max-w-md">
                            Lorem Ipsum and so on. Placeholder text
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default CheckYourEmail;