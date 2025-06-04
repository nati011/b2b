'use client'
import useOrdersStore from '@/app/libs/store/useOrderStore';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  Card,
  CardHeader,
  CardContent,
  CardTitle,
  CardDescription
} from '@/components/ui/card';
import { useEffect } from 'react';
import { CANCELED_STATUS, COMPLETED_STATUS, PENDING_STATUS } from "@/app/libs/enums";
import Loading from './loading';
import { Separator } from '@/components/ui/separator';
import Link from 'next/link';



const statusConfig = {
  [CANCELED_STATUS]: {
    bg: "bg-red-100/50",
    text: "text-red-800",
    border: "border-red-200",
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

export default function RecentSales() {
  const {
    loading,
    error,
    orders,
    fetchOrders,
    invoice
  } = useOrdersStore()

  useEffect(() => {
    fetchOrders()
  }, [])

  return (
    loading ? <Loading /> :
      <Card className='h-full  rounded-sm shadow-none'>
        <CardHeader>
          <CardTitle>Recent orders</CardTitle>
          <CardDescription>You have {orders.length} orders.</CardDescription>
          <Separator orientation='horizontal' />
        </CardHeader>
        <CardContent>
          <div className='space-y-8'>
            {orders.slice(0, 4).map((sale, index) => (
              <Link key={index} href={`/orders/${sale.Id}`} className='flex items-center'>
                <Avatar className='h-10 w-10'>
                  <AvatarFallback>{sale.RetailerName.split("")[0].slice(0, 2)}</AvatarFallback>
                </Avatar>
                <div className='ml-4 space-y-1'>
                  <p className='text-sm leading-none font-medium'>{sale.RetailerName}</p>
                  <div className='ml-auto font-medium text-sm text-gray-700'>{sale.Items.length} Items</div>
                </div>
                {/* @ts-expect-error */}
                <div className={`ml-auto inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusConfig[sale.Status as string].bg} ${statusConfig[sale.Status as string].text} ${statusConfig[sale.Status as string].border}`}>
                  {/* @ts-expect-error */}
                  {statusConfig[sale.Status as string].label}
                </div>
              </Link>

            ))}
          </div>
        </CardContent>
      </Card>


  );
}
