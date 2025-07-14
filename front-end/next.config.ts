import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  NEXT_PUBLIC_BACKEND_URL: process.env.NEXT_PUBLIC_BACKEND_URL,
  reactStrictMode: true
};

export default nextConfig;
