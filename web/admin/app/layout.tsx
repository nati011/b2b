import type { Metadata } from "next";
import { DM_Sans } from "next/font/google";
import { Toaster } from "@/components/ui/sonner"
import { Provider } from "@/app/themeprovider";
import "@/app/globals.css";
import { Footer } from "@/components/footer";

const font = DM_Sans({
    weight: ['100', '300', '700'],
    subsets: ['latin'],
    display: "swap",
    fallback: ["system-ui", "arial"],
});

export const metadata: Metadata = {
  title: "Efoyeta Store Admin",
  viewport: {
    width: 'device-width',
    initialScale: 1,
    maximumScale: 5,
  },
};

export default async function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" suppressHydrationWarning>
            <body
                className={`${font} antialiased`}
            >
                <Provider>
                    {children}
                    <Toaster richColors position="top-right" />
                </Provider>
            </body>
        </html>

    );
}
