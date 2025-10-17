import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  reactStrictMode: true,
  experimental: {
    serverActions: {
      bodySizeLimit: '1mb'
    }
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'images.unsplash.com'
      },
      {
        protocol: 'https',
        hostname: 'bitdash-a.akamaihd.net'
      },
      {
        protocol: 'https',
        hostname: 'test-streams.mux.dev'
      },
      {
        protocol: 'https',
        hostname: 'www.shutterstock.com'
      }
    ]
  },
  headers: async () => [
    {
      source: '/:path*',
      headers: [
        {
          key: 'X-Accel-Buffering',
          value: 'no'
        }
      ]
    }
  ]
};

export default nextConfig;
