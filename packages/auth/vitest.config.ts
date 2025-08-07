import path from 'path'

import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    watch: false,
    globals: true,
    environment: 'node',
    include: ['src/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    reporters: ['default'],
    coverage: {
      reportsDirectory: '.coverage',
      provider: 'v8' as const
    }
  },
  resolve: {
    alias: {
      '#': path.resolve(import.meta.dirname, './src')
    }
  }
})
