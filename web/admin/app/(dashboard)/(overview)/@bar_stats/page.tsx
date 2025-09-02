'use client';

import * as React from 'react';
import { Bar, BarChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from 'recharts';

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card';
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent
} from '@/components/ui/chart';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { TrendingUp, Calendar } from 'lucide-react';
import useOrdersStore from '@/app/libs/store/useOrderStore';
import { Order } from '@/app/libs/types';


type ChartData = {
  date: string;
  orders: number;
  revenue: number;
}

const chartConfig = {
  orders: {
    label: 'Orders',
    color: '#179FDB'
  },
  revenue: {
    label: 'Revenue',
    color: '#65C4BC'
  }
} satisfies ChartConfig;


export function mapOrdersToChartDataComplete(orders: Order[]): ChartData[] {
  const monthlyData = orders.reduce((acc, order) => {
    const date = new Date(order.CreatedAt);
    const monthKey = date.toLocaleDateString('en-US', { month: 'short' });
    
    if (!acc[monthKey]) {
      acc[monthKey] = {
        date: monthKey,
        orders: 0,
        revenue: 0
      };
    }
    
    acc[monthKey].orders += 1;
    if (order.PaymentStatus=='ACCEPTED'){
      acc[monthKey].revenue += order.Total;
    }
    
    return acc;
  }, {} as Record<string, ChartData>);
  
  const monthOrder = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 
                      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  
  return monthOrder.map(month => monthlyData[month] || {
    date: month,
    orders: 0,
    revenue: 0
  });
}

export default function BarGraph() {
  const [activeChart, setActiveChart] = React.useState<keyof typeof chartConfig>('orders');

  const [timeRange, setTimeRange] = React.useState('12m');

  const {
    orders,
    loading
  } = useOrdersStore()

  const chartData = React.useMemo(() => {
    if (!orders || orders.length === 0) {
      return [];
    }
    return mapOrdersToChartDataComplete(orders);
  }, [orders]);

  const total = React.useMemo(
    () => ({
      orders: chartData.reduce((acc, curr) => acc + curr.orders, 0),
      revenue: chartData.reduce((acc, curr) => acc + curr.revenue, 0)
    }),
    [chartData]
  );


  const [isClient, setIsClient] = React.useState(false);

  React.useEffect(() => {
    setIsClient(true);
  }, []);

  if (!isClient) {
    return (
      <Card className="rounded-lg border-0 shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div className="space-y-1">
            <CardTitle className="text-lg font-semibold">Order Analytics</CardTitle>
            <CardDescription>Track your order performance over time</CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center">
            <div className="animate-pulse">Loading chart...</div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="rounded-lg border-0 shadow-sm bg-gradient-to-br from-card to-card/50">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div className="space-y-1">
          <CardTitle className="text-lg font-semibold flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-primary" />
            Order Analytics
          </CardTitle>
          <CardDescription>Track your order performance over time</CardDescription>
        </div>
      </CardHeader>
      <CardContent className="p-6">
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="space-y-1">
            <p className="text-sm text-muted-foreground">Total Orders</p>
            <p className="text-2xl font-bold text-primary">{total.orders.toLocaleString()}</p>
          </div>
          <div className="space-y-1">
            <p className="text-sm text-muted-foreground">Total Revenue</p>
            <p className="text-2xl text-primary font-bold">${total.revenue.toLocaleString()}</p>
          </div>
        </div>
        
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[300px] w-full"
        >
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={chartData}
              margin={{
                top: 20,
                right: 20,
                left: 20,
                bottom: 20
              }}
            >
              <defs>
                <linearGradient id="fillBar" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="0%"
                    stopColor="#179FDB"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="100%"
                    stopColor="#179FDB"
                    stopOpacity={0.2}
                  />
                </linearGradient>
                <linearGradient id="fillRevenue" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="0%"
                    stopColor="#65C4BC"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="100%"
                    stopColor="#65C4BC"
                    stopOpacity={0.2}
                  />
                </linearGradient>
              </defs>
              <CartesianGrid 
                strokeDasharray="3 3" 
                vertical={false}
                stroke="hsl(var(--border))"
                opacity={0.5}
              />
              <XAxis
                dataKey="date"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                tick={{ fontSize: 12, fill: 'hsl(var(--muted-foreground))' }}
              />
              <YAxis
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                tick={{ fontSize: 12, fill: 'hsl(var(--muted-foreground))' }}
              />
              <ChartTooltip
                cursor={{ fill: 'hsl(var(--accent))', opacity: 0.1 }}
                content={
                  <ChartTooltipContent
                    className="w-[200px]"
                    nameKey={activeChart}
                  />
                }
              />
              <Bar
                dataKey={activeChart}
                fill={activeChart === 'orders' ? 'url(#fillBar)' : 'url(#fillRevenue)'}
                radius={[4, 4, 0, 0]}
                maxBarSize={50}
              />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>

      </CardContent>
    </Card>
  );
}
