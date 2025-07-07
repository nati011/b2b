"use client";
import { useEffect, useState } from "react";
import { Package } from "lucide-react";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import useOrdersStore from "@/lib/store/useOrderStore";
import { IoWarning } from "react-icons/io5";
import { OrderCard } from "@/components/OrderCard";
import { CANCELED_STATUS, COMPLETED_STATUS, PENDING_STATUS } from "@/lib/enums";

const Orders = () => {
  const {
    totalOrder,
    orders,
    loading,
    error,
    fetchOrders,
    completePayment,
    cancelOrder,
  } = useOrdersStore();
  const [currentPage, setCurrentPage] = useState(0);
  const [selectedStatus, setSelectedStatus] = useState("ALL");

  useEffect(() => {
    fetchOrders(currentPage, selectedStatus);
    console.log(orders);
    console.log(orders?.length)
  }, [currentPage, selectedStatus]);

  if (error) {
    return (
      <div className="w-full flex flex-col gap-6 items-center justify-center text-center min-h-screen">
        <div className="bg-red-50 rounded-full p-8 mb-4">
          <IoWarning className="w-16 h-16 text-red-900" />
        </div>
        <div>
          <h4 className="text-2xl font-semibold text-gray-900 mb-2">
            Something went wrong while fetching orders
          </h4>
        </div>
        <Button
          size="lg"
          className="px-8"
          onClick={() => {
            fetchOrders(0, selectedStatus);
          }}
        >
          Retry
        </Button>
      </div>
    );
  }

  const totalPages = Math.ceil(totalOrder / 10);

  return (
    <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div>
        <div className="mb-8">
          <h3 className="font-semibold text-xl text-primary mb-2">
            Order History
          </h3>
          <p className="text-md text-gray-600">
            Check the status of recent orders, manage returns, and discover
            similar products.
          </p>
        </div>

        {orders?.length === 0 && !loading ? (
          <div className="w-full flex flex-col gap-6 items-center justify-center text-center py-16">
            <div className="bg-gray-50 rounded-full p-8 mb-4">
              <Package className="w-16 h-16 text-gray-400" />
            </div>
            <div>
              <h4 className="text-2xl font-semibold text-gray-900 mb-2">
                No orders yet
              </h4>
              <p className="text-gray-600 mb-6 max-w-md">
                When you place your first order, it will appear here. Start
                shopping to see your order history.
              </p>
            </div>
            <Button size="lg" className="px-8">
              Start Shopping
            </Button>
          </div>
        ) : (
          <Tabs defaultValue="ALL" onValueChange={setSelectedStatus}>
            <TabsList className="w-full">
              <TabsTrigger value="ALL" className="border-none">
                ALL
              </TabsTrigger>
              <TabsTrigger value={PENDING_STATUS}>{PENDING_STATUS}</TabsTrigger>
              <TabsTrigger value={COMPLETED_STATUS}>
                {COMPLETED_STATUS}
              </TabsTrigger>
              <TabsTrigger value={CANCELED_STATUS}>
                {CANCELED_STATUS}
              </TabsTrigger>
            </TabsList>

            <div className="space-y-6">
              {loading ? (
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
              ) : (
                <>
                  {orders?.map((order, i) => (
                    <OrderCard
                      key={i}
                      completePayment={completePayment}
                      cancelOrder={cancelOrder}
                      loading={loading}
                      order={order}
                    />
                  ))}
                  {totalPages > 0 && (
                    <Pagination>
                      <PaginationContent>
                        <PaginationItem
                          onClick={() =>
                            currentPage > 0 && setCurrentPage(currentPage - 1)
                          }
                          className={`hover:cursor-pointer ${
                            currentPage === 0 ? "opacity-50 pointer-events-none" : ""
                          }`}
                          aria-disabled={currentPage === 0}
                        >
                          <PaginationPrevious />
                        </PaginationItem>
                        {[...Array(totalPages)].map((_, i) => (
                          <PaginationItem
                            key={i}
                            onClick={() => setCurrentPage(i)}
                            className={`hover:cursor-pointer ${
                              currentPage === i ? "font-bold underline" : ""
                            }`}
                          >
                            <PaginationLink>{i + 1}</PaginationLink>
                          </PaginationItem>
                        ))}
                        <PaginationItem
                          onClick={() =>
                            currentPage < totalPages - 1 &&
                            setCurrentPage(currentPage + 1)
                          }
                          className={`hover:cursor-pointer ${
                            currentPage === totalPages - 1 || totalPages === 0
                              ? "opacity-50 pointer-events-none"
                              : ""
                          }`}
                          aria-disabled={
                            currentPage === totalPages - 1 || totalPages === 0
                          }
                        >
                          <PaginationNext />
                        </PaginationItem>
                      </PaginationContent>
                    </Pagination>
                  )}
                </>
              )}
            </div>
          </Tabs>
        )}
      </div>
    </div>
  );
};

export default Orders;
