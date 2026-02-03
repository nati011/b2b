"use client";
import { usePathname } from "next/navigation";
import { ThemeProvider, useTheme } from "next-themes";
import { useEffect } from "react";

function ThemeController({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const { theme, setTheme } = useTheme();
  const isAdminPage = pathname?.startsWith("/admin");

  useEffect(() => {
    // Force light mode for non-admin pages
    if (!isAdminPage) {
      if (theme !== "light") {
        setTheme("light");
      }
      // Ensure HTML element doesn't have dark class
      if (typeof document !== "undefined") {
        document.documentElement.classList.remove("dark");
      }
    }
  }, [isAdminPage, theme, setTheme]);

  return <>{children}</>;
}

export function ThemeWrapper({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider attribute="class" defaultTheme="light" enableSystem={false}>
      <ThemeController>{children}</ThemeController>
    </ThemeProvider>
  );
}

