import { cookies } from "next/headers"
import { DM_Sans } from "next/font/google";

import { Toaster } from "@/components/ui/sonner"
import "@/app/globals.css";

import Topnav from "@/components/topnav";
import SessionProvider from "@/app/sessionprovider";
import { Provider } from "@/app/themeprovider";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"
import { Separator } from "@/components/ui/separator";

const font = DM_Sans({
  weight: ['100', '300', '400', '500', '600', '700'],
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
        <div className={`flex relative min-h-screen bg-gradient-to-br from-background via-background to-muted/20 ${font.className}`}>
          <SidebarProvider defaultOpen={true}>
            <AppSidebar />
            <SidebarInset className="flex-1">
              <div className="flex flex-col min-h-screen">
                <Topnav />
                <main className="flex-1 p-6 space-y-6">
                  <div className="animate-fade-in">
                    {children}
                  </div>
                </main>
              </div>
            </SidebarInset>
          </SidebarProvider>
          <Toaster 
            position="top-right"
            richColors
            closeButton
            duration={4000}
          />
        </div>
      </Provider>
    </SessionProvider>
  );
}
