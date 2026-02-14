/// <reference types="node" />
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  compress: true,
  poweredByHeader: false,

  images: {
    remotePatterns: [
      { protocol: "https", hostname: "*" },
      { protocol: "http", hostname: "*" },
    ],
    formats: ["image/avif", "image/webp"],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    minimumCacheTTL: 60,
  },

  experimental: {
    optimizePackageImports: [
      "lucide-react",
      "@radix-ui/react-dropdown-menu",
      "@radix-ui/react-dialog",
      "@radix-ui/react-select",
      "@radix-ui/react-popover",
    ],
    ppr: false,
  },

  compiler: {
    removeConsole: process.env.NODE_ENV === "production",
  },

  webpack: (config: Record<string, unknown>, { isServer }: { isServer: boolean }) => {
    if (!isServer && config.optimization) {
      config.optimization = {
        ...(config.optimization as Record<string, unknown>),
        moduleIds: "deterministic",
      };
    }
    return config;
  },
};

export default nextConfig;
