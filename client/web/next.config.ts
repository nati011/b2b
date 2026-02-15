/** @type {import('next').NextConfig} */
const backendApiUrl = process.env.BACKEND_API_URL || 'http://185.222.240.66';

module.exports = {
  reactStrictMode: true,
  // Performance optimizations
  swcMinify: true,
  compress: true,
  poweredByHeader: false,

  // Proxy to backend so browser never calls HTTP (avoids mixed content when site is HTTPS on Vercel)
  async rewrites() {
    return [
      { source: '/api-backend/:path*', destination: `${backendApiUrl}/api/:path*` },
    ];
  },

  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: 'res.cloudinary.com',
      },
      {
        protocol: "https",
        hostname: 'images.unsplash.com',
      },
      {
        protocol: "https",
        hostname: 'unsplash.com',
      },
    ],
    formats: ['image/avif', 'image/webp'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    minimumCacheTTL: 60,
  },
  
  // Optimize package imports
  experimental: {
    optimizePackageImports: [
      'next/font',
      'lucide-react',
      '@radix-ui/react-dropdown-menu',
      '@radix-ui/react-dialog',
    ],
  },
  
  // Compiler optimizations
  compiler: {
    removeConsole: process.env.NODE_ENV === 'production',
  },
  
  // Performance optimizations for faster navigation
  onDemandEntries: {
    maxInactiveAge: 25 * 1000,
    pagesBufferLength: 2,
  },
  
  // Enable static page generation optimizations
  generateBuildId: async () => {
    return 'build-' + Date.now();
  },
};
