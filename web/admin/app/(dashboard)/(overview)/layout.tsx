"use client"
import { PENDING_STATUS } from '@/app/libs/enums';
import useOrdersStore from '@/app/libs/store/useOrderStore';
import PageContainer from '@/components/page-container';
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card';
import { 
  Users, 
  Package, 
  ShoppingCart, 
  DollarSign,
} from 'lucide-react';
import React from 'react';

export default function OverViewLayout({
  sales,
  bar_stats,
}: {
  sales: React.ReactNode;
  bar_stats: React.ReactNode;
}) {
  const {
    orders,
  } = useOrdersStore()
  const totalRevenue = orders.reduce((sum, order) => sum + (order.Total || 0), 0).toLocaleString();
  const activeOrders = orders.filter(order => order.Status === PENDING_STATUS).length.toLocaleString();
  
  const stats = [
    {
      title: "Total Revenue",
      value: `$ ${totalRevenue}`,
      icon: DollarSign,
      color: "text-emerald-600"
    },
    {
      title: "Active Orders",
      value: `${activeOrders}`,
      icon: ShoppingCart,
      color: "text-blue-600"
    },
    {
      title: "Total Customers",
      value: "12,234",
      icon: Users,
      color: "text-purple-600"
    },
    {
      title: "Products Sold",
      value: "573",
      icon: Package,
      color: "text-orange-600"
    }
  ];

  return (
    <PageContainer>
      <div className='flex flex-1 flex-col space-y-6 animate-fade-in'>
        {/* Header */}
        <div className='flex items-center justify-between'>
          <div>
            <h1 className='text-3xl font-bold tracking-tight text-primary'>
              Dashboard Overview
            </h1>
            <p className='text-muted-foreground mt-1'>
              Welcome back! Here's what's happening with your business today.
            </p>
          </div>
        </div>

        {/* Stats Cards */}
        <div className='grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4'>
          {stats.map((stat, index) => (
            <Card key={index} className='card-hover border-0 shadow-sm bg-gradient-to-br from-card to-card/50'>
              <CardHeader className='flex flex-row items-center justify-between space-y-0'>
                <CardTitle className='text-sm font-medium text-muted-'>
                  {stat.title}
                </CardTitle>
                <stat.icon className={`h-4 w-4 ${stat.color}`} />
              </CardHeader>
              <CardContent>
                <div className='text-2xl font-bold text-primary'>{stat.value}</div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Charts Section */}
        <div className='grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-7'>
          <div className='col-span-4'>
            {bar_stats}
          </div>
          <div className='col-span-4 md:col-span-3'>
            {sales}
          </div>
        </div>

       
      </div>
    </PageContainer>
  );
}
