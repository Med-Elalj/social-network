import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/uploads/:path*',
        destination: 'http://backend:8080/uploads/:path*',
      },
    ];
  },
  reactStrictMode: true
};

export default nextConfig;
