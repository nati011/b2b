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
    domains: ["res.cloudinary.com"],
  },
};
