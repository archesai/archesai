import type { NextConfig } from 'next'

export default {
  eslint: { ignoreDuringBuilds: true },
  transpilePackages: ['@archesai/ui', '@archesai/client'],
  typescript: { ignoreBuildErrors: true, tsconfigPath: 'tsconfig.lib.json' }
} satisfies NextConfig
