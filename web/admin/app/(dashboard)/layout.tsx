import { DM_Sans } from "next/font/google";
import { Toaster } from "@/components/ui/sonner"
import { Provider } from "@/app/themeprovider";
import "@/app/globals.css";

const font = DM_Sans({
    weight: ['100', '300', '700'],
    subsets: ['latin']
});

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
                    <Toaster />
                </Provider>
            </body>
        </html>

    );
}
