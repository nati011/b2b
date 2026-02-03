"use client";
import { usePathname } from "next/navigation";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";

export function ConditionalLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isAdminPage = pathname?.startsWith("/admin");
  const isSupplierPage = pathname?.startsWith("/supplier");
  const isAdminOrSupplierPage = isAdminPage || isSupplierPage;

  return (
    <>
      {!isAdminOrSupplierPage && <Navbar />}
      {children}
      {!isAdminOrSupplierPage && <Footer />}
    </>
  );
}






