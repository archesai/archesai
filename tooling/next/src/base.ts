import type { NextConfig } from 'next'

export default {
  eslint: {
    ignoreDuringBuilds: true
  },
  experimental: {
    optimizePackageImports: ['react-player']
  },
  transpilePackages: ['@archesai/ui', '@archesai/client', '@archesai/domain'],
  typescript: { tsconfigPath: 'tsconfig.lib.json', ignoreBuildErrors: true }
} satisfies NextConfig
