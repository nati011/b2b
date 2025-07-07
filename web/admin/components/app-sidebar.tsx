import * as React from "react"
import { NavUser } from "@/components/nav-user"

import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarMenuSub,
    SidebarMenuSubButton,
    SidebarMenuSubItem,
    SidebarRail,
} from "@/components/ui/sidebar"

import {
    Package,
    Users,
    ShoppingCart,
    Settings,
    BarChart3,
    Store,
} from "lucide-react"
import Image from "next/image"
import Link from "next/link"

export async function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
    const data = {
        navMain: [
            {
                title: "Dashboard",
                url: "/",
                icon: BarChart3,
                isActive: true,
            },
            {
                title: "Products",
                url: "/products",
                icon: Package,
                items: [
                    {
                        title: "Simple Products",
                        url: "/products",
                    },
                    {
                        title: "Configurable Products",
                        url: "/products/configurable",
                    },
                    {
                        title: "Categories",
                        url: "/products/category",
                    },
                ],
            },
            {
                title: "Retailers",
                url: "/retailers",
                icon: Users,
            },
            {
                title: "Distributors",
                url: "/distributors",
                icon: Store,
                items: [
                    {
                        title: "All Distributors",
                        url: "/distributors",
                    },
                    {
                        title: "Distributor Agents",
                        url: "/distributors/agents",
                    },
                ],
            },
            {
                title: "Orders",
                url: "/orders",
                icon: ShoppingCart,
                items: [
                    {
                        title: "All Orders",
                        url: "/orders",
                    },
                    {
                        title: "Transactions",
                        url: "/transactions",
                    }
                ],
            },
            {
                title: "Settings",
                url: "/role",
                icon: Settings,
                items: [
                    {
                        title: "Roles & Permissions",
                        url: "/role",
                    }
                ],
            }
        ]
    }

    return (
        <Sidebar {...props} className="border-r border-border/40">
            <SidebarHeader className="border-b border-border/40">
                <SidebarMenu>
                    <SidebarMenuItem className="py-[0.15rem]">
                        <Link href="/">
                            <SidebarMenuButton
                                size="lg"
                                className="data-[state=open]:bg-accent data-[state=open]:text-accent-foreground hover:bg-accent/50 transition-colors"
                            >
                                <div className="flex  size-10 items-center justify-center  text-primary-foreground">
                                    <Image src='/logo.svg' width={32} height={32} alt="logo" className="rounded-lg" />
                                </div>
                                <div className="grid flex-1 text-left text-lg leading-tight">
                                    <span className="truncate font-bold">
                                        Efoyeta Store
                                    </span>
                                    <span className="truncate text-xs text-muted-foreground">
                                        Admin Dashboard
                                    </span>
                                </div>
                            </SidebarMenuButton>
                        </Link>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent className="px-2 py-4">
                <SidebarGroup>
                    <SidebarMenu>
                        {data.navMain.map((item) => (
                            <SidebarMenuItem key={item.title}>
                                <Link href={item.url}>
                                <SidebarMenuButton
                                    tooltip={item.title}
                                    className="hover:bg-accent/50 transition-all duration-200 group"
                                >
                                    <item.icon className="h-5 w-5 transition-transform group-hover:scale-110 text-primary" />
                                    <span className="font-medium">{item.title}</span>
                                </SidebarMenuButton>
                                </Link>
                                {item.items?.length ? (
                                    <SidebarMenuSub>
                                        {item.items.map((subItem) => (
                                            <SidebarMenuSubItem key={subItem.title}>
                                                <SidebarMenuSubButton asChild>
                                                    <Link href={subItem.url} className="hover:bg-accent/30 transition-colors">
                                                        {subItem.title}
                                                    </Link>
                                                </SidebarMenuSubButton>
                                            </SidebarMenuSubItem>
                                        ))}
                                    </SidebarMenuSub>
                                ) : null}
                            </SidebarMenuItem>
                        ))}
                    </SidebarMenu>
                </SidebarGroup>
            </SidebarContent>
            <SidebarFooter className="border-t border-border/40 p-2 active:bg-neutral-100">
                <NavUser />
            </SidebarFooter>
            <SidebarRail />
        </Sidebar>
    )
}

