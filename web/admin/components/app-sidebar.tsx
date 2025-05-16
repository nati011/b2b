"use client"

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

import { AiOutlineProduct } from "react-icons/ai"
import { PiUsersThreeLight } from "react-icons/pi"
import { GoGear } from "react-icons/go"
import Image from "next/image"
import Link from "next/link"

const data = {
    user: {
        name: "shadcn",
        email: "m@example.com",
        avatar: "/avatars/shadcn.jpg",
    },
    singular: [
        {
            title: "Admin Users",
            url: "/admin",
            icon: PiUsersThreeLight,
        }
    ],
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
            title: "Customers",
            url: "/retailers",
            icon: PiUsersThreeLight,
        },
        {
            title: "Distributors",
            url: "/distributors",
            icon: PiUsersThreeLight
        },
        {
            title: "Order",
            url: "/Order",
            icon: PiUsersThreeLight,
            items: [
                {
                    title: "Orders",
                    url: "/orders",
                },
                {
                    title: "Transactions",
                    url: "/transactions",
                }
            ],
        },
        {
            title: "Configurations",
            url: "/settings",
            icon: GoGear,
            items: [
                {
                    title: "Role",
                    url: "/role",
                }
            ],
        }
    ]
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
    return (
        <Sidebar {...props}>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <Link href="/">
                            <SidebarMenuButton
                                size="lg"
                                className="data-[state=open]:bg-gray-800 text-white-accent data-[state=open]:text-sidebar-accent-foreground"
                            >

                                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-gray-800 text-sidebar-primary-foreground">
                                    <Image src='/logo.png' width={150} height={100} alt="logo" />
                                </div>
                                <div className="grid flex-1 text-left text-lg leading-tight">
                                    <span className="truncate font-semibold">
                                        Efoyeta Store
                                    </span>
                                </div>
                            </SidebarMenuButton>

                        </Link>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarMenu>
                        {data.navMain.map((item) => (
                            <SidebarMenuItem key={item.title}>
                                <SidebarMenuButton tooltip={item.title}>
                                    {item.icon && <item.icon />}
                                    <a href={item.url}>{item.title}</a>
                                </SidebarMenuButton>
                                {item.items?.length ? (
                                    <SidebarMenuSub>
                                        {item.items.map((item) => (
                                            <SidebarMenuSubItem key={item.title}>
                                                <SidebarMenuSubButton asChild>
                                                    <a href={item.url}>{item.title}</a>
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
            <SidebarFooter>
                <NavUser user={data.user} />
            </SidebarFooter>
            <SidebarRail />
        </Sidebar>
    )
}
