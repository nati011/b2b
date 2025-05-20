import { cookies } from "next/headers"
import { DM_Sans } from "next/font/google";

import { Toaster } from "@/components/ui/sonner"
import "@/app/globals.css";

import Topnav from "@/app/components/topnav";
import SessionProvider from "@/app/sessionprovider";
import { Provider } from "@/app/themeprovider";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"
import { Separator } from "@/components/ui/separator";



const font = DM_Sans({
  weight: ['100', '300', '700'],
  subsets: ['latin']
});


export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const cookieStore = await cookies()
  return (
    <SessionProvider>
      <Provider>
        <div className="flex relative">
          <SidebarProvider defaultOpen={true}>
            <AppSidebar />
            <Separator orientation="vertical" className="h-4" />
            <div className=" p-4 w-full bg-slate-50/50">
              <Topnav />
              {children}
            </div>

          </SidebarProvider>
        </div>
      </Provider>
    </SessionProvider>

  );
}
