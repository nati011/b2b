/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  async rewrites() {
    return {
      beforeFiles: [
        {
          source: "/api/:path*",
          destination: "https://b2b-67gk.onrender.com/api/v1/:path*",
        },
      ],
    };
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
      }
    ]
  }
};
