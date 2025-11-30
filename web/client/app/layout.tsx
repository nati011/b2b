import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { Toaster } from "@/components/ui/sonner";
import SessionProvider from "@/app/SessionProvider";
import TawkChat from "@/components/TawkChat";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
  fallback: ["-apple-system", "BlinkMacSystemFont", "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "sans-serif"],
});

export const metadata: Metadata = {
  title: "Efoyeta Store",
  viewport: {
    width: 'device-width',
    initialScale: 1,
    maximumScale: 5,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${inter.variable} font-sans antialiased min-h-screen`}
      >
        <SessionProvider>
          <Navbar />
          {children}
          <Footer />
          <Toaster richColors />
          <TawkChat />
        </SessionProvider>
      </body>
    </html>
  );
}
