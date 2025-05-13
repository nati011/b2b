import { cookies } from "next/headers"
import { DM_Sans } from "next/font/google";

import { Toaster } from "@/components/ui/sonner"
import "./globals.css";

import Topnav from "./components/topnav";
import SessionProvider from "./sessionprovider";
import { Provider } from "./themeprovider";
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
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${font} antialiased`}
      >
        <SessionProvider>
          <Provider>
            <div className="flex relative">
              <SidebarProvider defaultOpen={true}>
                <AppSidebar />
                <SidebarInset>
                  <main>
                    <Separator orientation="vertical" className="mr-2 h-4" />
                    <div className=" p-4  sm:px-10 ml-[16%] min-h-screen  bg-slate-50">
                      <Topnav />
                      {children}
                      <Toaster />
                    </div>
                  </main>
                </SidebarInset>
              </SidebarProvider>
            </div>
          </Provider>
        </SessionProvider>
      </body>
    </html>

  );
}
