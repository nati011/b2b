"use client";
import React from "react";
import { useTheme } from "next-themes";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
    Bell,
    Search,
    Settings,
    LogOut,
    User,
    Moon,
    Sun,
    Monitor,
    ChevronDown
} from "lucide-react";
import { useAuth } from "@/hooks/useAuth";

const Topnav = () => {
    const { setTheme } = useTheme();
    const { user, logout, isLoading } = useAuth();
    // @ts-ignore
    const user_initials = (user?.name)
        ?.split(" ")
        .map((n: string) => n[0])
        .join("");

    if (isLoading) {
        return (
            <div className="flex justify-end items-center mb-4 animate-pulse">
                <div className="h-10 w-10 bg-muted rounded-full"></div>
            </div>
        );
    }

    return (
        <div className="flex items-center justify-end p-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 print:hidden">
            <div className="flex items-center space-x-4">
                {/* Theme Toggle */}
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                            <Sun className="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                            <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                            <span className="sr-only">Toggle theme</span>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => setTheme("light")}>
                            <Sun className="mr-2 h-4 w-4" />
                            <span>Light</span>
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => setTheme("dark")}>
                            <Moon className="mr-2 h-4 w-4" />
                            <span>Dark</span>
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => setTheme("system")}>
                            <Monitor className="mr-2 h-4 w-4" />
                            <span>System</span>
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>

                {/* User Menu */}
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="flex items-center space-x-2 px-3 py-2">
                            <Avatar className="h-8 w-8">
                                <AvatarImage src="" alt={user?.first_name ?? ""} />
                                <AvatarFallback className="bg-gradient-to-br from-primary to-primary/80 text-primary-foreground text-sm font-medium">
                                    {user_initials || <User className="h-4 w-4" />}
                                </AvatarFallback>
                            </Avatar>
                            <div className="hidden md:flex flex-col items-start">
                                {/* @ts-ignore */}
                                <span className="text-sm font-medium">{user?.name}</span>
                                <span className="text-xs text-muted-foreground">{user?.email}</span>
                            </div>
                            <ChevronDown className="h-4 w-4 text-muted-foreground" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-56">
                        <DropdownMenuLabel className="flex items-center space-x-2 p-3">
                            <Avatar className="h-10 w-10">
                                <AvatarImage src="" alt={user?.first_name ?? ""} />
                                <AvatarFallback className="bg-gradient-to-br from-primary to-primary/80 text-primary-foreground">
                                    {user_initials || <User className="h-5 w-5" />}
                                </AvatarFallback>
                            </Avatar>
                            <div className="flex flex-col">
                                <span className="font-medium">{user?.first_name} {user?.last_name}</span>
                                <span className="text-sm text-muted-foreground">{user?.email}</span>
                                {/* {user?.roles && user.roles.length > 0 && (
                                    <div className="flex flex-wrap gap-1 mt-1">
                                        {user.roles.slice(0, 2).map((role, index) => (
                                            <Badge key={index} variant="secondary" className="text-xs">
                                                {role}
                                            </Badge>
                                        ))}
                                        {user.roles.length > 2 && (
                                            <Badge variant="outline" className="text-xs">
                                                +{user.roles.length - 2}
                                            </Badge>
                                        )}
                                    </div>
                                )} */}
                            </div>
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuGroup>
                            <DropdownMenuItem className="flex items-center">
                                <User className="h-4 w-4 mr-2" />
                                Profile
                            </DropdownMenuItem>
                            <DropdownMenuItem className="flex items-center">
                                <Settings className="h-4 w-4 mr-2" />
                                Settings
                            </DropdownMenuItem>
                        </DropdownMenuGroup>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                            onClick={() => logout()}
                            className="flex items-center text-destructive focus:text-destructive"
                        >
                            <LogOut className="h-4 w-4 mr-2" />
                            Sign out
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
    );
};

export default Topnav;