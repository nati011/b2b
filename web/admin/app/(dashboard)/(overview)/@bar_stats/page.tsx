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

export const description = 'An interactive bar chart showing order analytics';

const chartData = [
  { date: 'Jan', orders: 120, revenue: 2400 },
  { date: 'Feb', orders: 180, revenue: 3600 },
  { date: 'Mar', orders: 150, revenue: 3000 },
  { date: 'Apr', orders: 220, revenue: 4400 },
  { date: 'May', orders: 280, revenue: 5600 },
  { date: 'Jun', orders: 320, revenue: 6400 },
  { date: 'Jul', orders: 290, revenue: 5800 },
  { date: 'Aug', orders: 350, revenue: 7000 },
  { date: 'Sep', orders: 380, revenue: 7600 },
  { date: 'Oct', orders: 420, revenue: 8400 },
  { date: 'Nov', orders: 450, revenue: 9000 },
  { date: 'Dec', orders: 500, revenue: 10000 }
];

const chartConfig = {
  orders: {
    label: 'Orders',
    color: 'hsl(var(--primary))'
  },
  revenue: {
    label: 'Revenue',
    color: 'hsl(var(--chart-1))'
  }
} satisfies ChartConfig;

export default function BarGraph() {
  const [activeChart, setActiveChart] = React.useState<keyof typeof chartConfig>('orders');
  const [timeRange, setTimeRange] = React.useState('12m');

  const total = React.useMemo(
    () => ({
      orders: chartData.reduce((acc, curr) => acc + curr.orders, 0),
      revenue: chartData.reduce((acc, curr) => acc + curr.revenue, 0)
    }),
    []
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
        <div className="flex items-center space-x-2">
          <Select value={timeRange} onValueChange={setTimeRange}>
            <SelectTrigger className="w-[120px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="7d">Last 7 days</SelectItem>
              <SelectItem value="30d">Last 30 days</SelectItem>
              <SelectItem value="3m">Last 3 months</SelectItem>
              <SelectItem value="6m">Last 6 months</SelectItem>
              <SelectItem value="12m">Last 12 months</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent className="p-6">
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="space-y-1">
            <p className="text-sm text-muted-foreground">Total Orders</p>
            <p className="text-2xl font-bold">{total.orders.toLocaleString()}</p>
          </div>
          <div className="space-y-1">
            <p className="text-sm text-muted-foreground">Total Revenue</p>
            <p className="text-2xl font-bold">${total.revenue.toLocaleString()}</p>
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
                    stopColor="hsl(var(--primary))"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="100%"
                    stopColor="hsl(var(--primary))"
                    stopOpacity={0.2}
                  />
                </linearGradient>
                <linearGradient id="fillRevenue" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="0%"
                    stopColor="hsl(var(--chart-1))"
                    stopOpacity={0.8}
                  />
                  <stop
                    offset="100%"
                    stopColor="hsl(var(--chart-1))"
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
