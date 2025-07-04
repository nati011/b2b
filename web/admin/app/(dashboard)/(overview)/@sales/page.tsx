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
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useEffect } from 'react';
import { CANCELED_STATUS, COMPLETED_STATUS, PENDING_STATUS } from "@/app/libs/enums";
import Loading from './loading';
import { Separator } from '@/components/ui/separator';
import Link from 'next/link';
import { 
  Clock, 
  CheckCircle, 
  XCircle, 
  ArrowRight, 
  TrendingUp,
  Calendar,
  DollarSign,
  ShoppingCart
} from 'lucide-react';

const statusConfig = {
  [CANCELED_STATUS]: {
    bg: "bg-red-100/50 dark:bg-red-900/20",
    text: "text-red-800 dark:text-red-400",
    border: "border-red-200 dark:border-red-800",
    label: "Canceled",
    icon: XCircle
  },
  [PENDING_STATUS]: {
    bg: "bg-amber-100/50 dark:bg-amber-900/20",
    text: "text-amber-800 dark:text-amber-400",
    border: "border-amber-200 dark:border-amber-800",
    label: "Pending",
    icon: Clock
  },
  [COMPLETED_STATUS]: {
    bg: "bg-emerald-100/50 dark:bg-emerald-900/20",
    text: "text-emerald-800 dark:text-emerald-400",
    border: "border-emerald-200 dark:border-emerald-800",
    label: "Completed",
    icon: CheckCircle
  },
};

export default function RecentSales() {
  const {
    loading,
    orders,
    fetchOrders,
  } = useOrdersStore()
  
  useEffect(() => {
    fetchOrders(0)
  }, [])

  if (loading) {
    return <Loading />
  }

  const recentOrders = orders.slice(0, 4);
  const totalRevenue = orders.reduce((sum, order) => sum + (order.Total || 0), 0);
  const completedOrders = orders.filter(order => order.Status === COMPLETED_STATUS).length;
  const pendingOrders = orders.filter(order => order.Status === PENDING_STATUS).length;

  return (
    <Card className='h-full rounded-lg border-0 shadow-sm bg-gradient-to-br from-card to-card/50'>
      <CardHeader className="">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-lg font-semibold flex items-center gap-2">
              <TrendingUp className="h-5 w-5 text-primary" />
              Recent Orders
            </CardTitle>
            <CardDescription className="mt-1">
              Latest order activity and performance metrics
            </CardDescription>
          </div>
          <Link href="/orders">
            <Button variant="outline" size="sm" className="text-xs">
              View All
              <ArrowRight className="h-3 w-3 ml-1" />
            </Button>
          </Link>
        </div>
        
        {/* Quick Stats */}
        <div className="grid grid-cols-3 gap-4 mt-4">
          <div className="text-center p-3 bg-muted/50 rounded-lg">
            <div className="flex items-center justify-center mb-1">
              <DollarSign className="h-4 w-4 text-emerald-600" />
              <p className="text-sm font-medium">{totalRevenue.toLocaleString()}</p>
            </div>
           
            <p className="text-xs text-muted-foreground font-semibold">Total Revenue</p>
          </div>
          <div className="text-center p-3 bg-muted/50 rounded-lg">
            <div className="flex items-center justify-center mb-1 gap-1">
              <CheckCircle className="h-4 w-4 text-emerald-600" />
              <p className="text-sm font-medium">{completedOrders}</p>
            </div>
          
            <p className="text-xs text-muted-foreground font-semibold">Completed</p>
          </div>
          <div className="text-center p-3 bg-muted/50 rounded-lg">
            <div className="flex items-center justify-center mb-1 gap-1">
              <Clock className="h-4 w-4 text-amber-600" />
              <p className="text-sm font-medium">{pendingOrders}</p>
            </div>
            
            <p className="text-xs text-muted-foreground font-semibold">Pending</p>
          </div>
        </div>
      </CardHeader>
      
      <Separator className="mb-4" />
      
      <CardContent className="p-0">
        <div className='space-y-3 px-6'>
          {recentOrders.length > 0 ? (
            recentOrders.map((order, index) => {
              const status = statusConfig[order.Status as keyof typeof statusConfig];
              const StatusIcon = status?.icon || Clock;
              
              return (
                <Link key={index} href={`/orders/${order.Id}`} className='block'>
                  <div className='flex items-center p-3 rounded-lg hover:bg-muted/50 transition-colors group'>
                    <Avatar className='h-10 w-10 mr-3'>
                      <AvatarFallback className="bg-gradient-to-br from-primary to-primary/80 text-primary-foreground font-medium">
                        {order.RetailerName?.split("")[0]?.slice(0, 2) || "OR"}
                      </AvatarFallback>
                    </Avatar>
                    
                    <div className='flex-1 min-w-0'>
                      <div className="flex items-center justify-between">
                        <p className='text-sm font-medium truncate'>{order.RetailerName}</p>
                        <div className="flex items-center space-x-2">
                          <Badge 
                            variant="secondary" 
                            className={`text-xs ${status?.bg} ${status?.text} ${status?.border}`}
                          >
                            <StatusIcon className="h-3 w-3 mr-1" />
                            {status?.label}
                          </Badge>
                        </div>
                      </div>
                      
                      <div className="flex items-center justify-between mt-1">
                        <div className="flex items-center space-x-4 text-xs text-muted-foreground">
                          <span className="flex items-center">
                            <Calendar className="h-3 w-3 mr-1" />
                            {new Date(order.CreatedAt).toLocaleDateString()}
                          </span>
                          <span>{order.Items?.length || 0} items</span>
                        </div>
                        <p className="text-sm font-medium">
                          ${order.Total?.toLocaleString() || '0'}
                        </p>
                      </div>
                    </div>
                    
                    <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors ml-2" />
                  </div>
                </Link>
              );
            })
          ) : (
            <div className="text-center py-8">
              <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mx-auto mb-4">
                <ShoppingCart className="h-8 w-8 text-muted-foreground" />
              </div>
              <p className="text-muted-foreground">No orders yet</p>
              <p className="text-sm text-muted-foreground mt-1">Orders will appear here when they're created</p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
