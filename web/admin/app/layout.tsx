"use client"
import { DM_Sans } from "next/font/google";

import { Toaster } from "@/components/ui/sonner"
import "./globals.css";

import Topnav from "./components/topnav";
import SessionProvider from "./sessionprovider";
import { Provider } from "./themeprovider";
import SideBar from "@/app/components/sidebar";


const font = DM_Sans({
  weight: ['100', '300', '700'],
  subsets: ['latin']
});


export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${font} antialiased`}
      >
        <SessionProvider>
          <Provider>
            <div className="flex relative bg-slate-50">
              <div className="lg:w-[15%] fixed z-50">
                <SideBar />
              </div>
              <div className="dark:bg-neutral-900 p-4  sm:px-10 lg:ml-[15%] w-full min-h-screen">
                <div className="mb-10">
                  <Topnav />
                </div>
                {children}
                <Toaster />
              </div>
            </div>
          </Provider>
        </SessionProvider>
      </body>
    </html>

  );
}
