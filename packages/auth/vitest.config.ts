import path from 'node:path'

import { defineConfig } from 'vitest/config'

export default defineConfig({
  resolve: {
    alias: {
      '#': path.resolve(import.meta.dirname, './src')
    }
  },
  test: {
    coverage: {
      provider: 'v8' as const,
      reportsDirectory: '.coverage'
    },
    environment: 'node',
    globals: true,
    include: ['src/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    reporters: ['default'],
    watch: false
  }
})
