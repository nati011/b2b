/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  async rewrites() {
    return {
      beforeFiles: [
        {
          source: "/api/:path*",
          destination: "http://localhost:8084/api/v1/:path*",
        },
      ],
    };
  },
  images: {
    domains: ["res.cloudinary.com"],
  },
};
