import type { NextConfig } from 'next'

export default {
  transpilePackages: ['@archesai/ui', '@archesai/client', '@archesai/domain'],
  typescript: { tsconfigPath: 'tsconfig.lib.json' }
} satisfies NextConfig
