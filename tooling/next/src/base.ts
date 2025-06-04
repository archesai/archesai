import type { NextConfig } from 'next'

export default {
  eslint: {
    ignoreDuringBuilds: true
  },
  transpilePackages: ['@archesai/ui', '@archesai/client', '@archesai/domain'],
  typescript: { tsconfigPath: 'tsconfig.lib.json', ignoreBuildErrors: true }
} satisfies NextConfig
