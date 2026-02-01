import type { Metadata } from "next";
import { Overpass } from "next/font/google";
import "./globals.css";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { Toaster } from "@/components/ui/sonner";
import SessionProvider from "@/app/SessionProvider";
import { AuthProvider } from "@/context/AuthContext";
import TawkChat from "@/components/TawkChat";

const overpass = Overpass({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-overpass",
  weight: ["400", "500", "600", "700"],
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
        className={`${overpass.variable} font-sans antialiased min-h-screen`}
      >
        <SessionProvider>
          <AuthProvider>
            <Navbar />
            {children}
            <Footer />
            <Toaster richColors />
            <TawkChat />
          </AuthProvider>
        </SessionProvider>
      </body>
    </html>
  );
}
