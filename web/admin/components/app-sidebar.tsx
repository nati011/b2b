"use client"

import * as React from "react"

import { NavMain } from "@/components/nav-main"
import { NavUser } from "@/components/nav-user"
import {
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    useSidebar,
} from "@/components/ui/sidebar"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarRail,
} from "@/components/ui/sidebar"
import { CiShoppingCart } from "react-icons/ci"
import { AiOutlineProduct } from "react-icons/ai"
import { PiUsersThreeLight } from "react-icons/pi"
import Link from "next/link"

// This is sample data.
const data = {
    user: {
        name: "shadcn",
        email: "m@example.com",
        avatar: "/avatars/shadcn.jpg",
    },
    navMain: [
        {
            title: "Product",
            url: "/products",
            icon: AiOutlineProduct,
            isActive: true,
            items: [
                {
                    title: "Simple Product",
                    url: "/products",
                },
                {
                    title: "Configurable Products",
                    url: "/products/configurable",
                },
                {
                    title: "Product Categories",
                    url: "/products/category",
                },
            ],
        },
        {
            title: "Users",
            url: "#",
            icon: PiUsersThreeLight,
            items: [
                {
                    title: "Retailers",
                    url: "/retailers",
                },
                {
                    title: "Distributor",
                    url: "/distributor",
                },
                {
                    title: "Admin Users",
                    url: "/admin",
                },
            ],
        }
    ]
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
    return (
        <Sidebar collapsible="icon" {...props}>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <Link href="/">
                            <SidebarMenuButton
                                size="lg"
                                className="data-[state=open]:bg-gray-800 text-white-accent data-[state=open]:text-sidebar-accent-foreground"
                            >

                                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-gray-800 text-sidebar-primary-foreground">
                                    <CiShoppingCart className="text-2xl" />
                                </div>
                                <div className="grid flex-1 text-left text-lg leading-tight">
                                    <span className="truncate font-semibold">
                                        Efoyeta Market
                                    </span>
                                </div>
                            </SidebarMenuButton>

                        </Link>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <NavMain items={data.navMain} />
            </SidebarContent>
            <SidebarFooter>
                <NavUser user={data.user} />
            </SidebarFooter>
            <SidebarRail />
        </Sidebar>
    )
}
