import path from 'node:path'

import tailwindcss from '@tailwindcss/vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact from '@vitejs/plugin-react'
import { visualizer } from 'rollup-plugin-visualizer'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    tailwindcss(),
    tanstackStart({
      customViteReactPlugin: true,
      tsr: {
        routesDirectory: 'src/app' // Defaults to "src/routes"
      }
    }),
    viteReact({
      babel: {
        plugins: [
          [
            'babel-plugin-react-compiler',
            {
              target: '19'
            }
          ]
        ]
      }
    }),
    visualizer({
      brotliSize: true,
      filename: 'dist/stats.html',
      gzipSize: true,
      open: true
    })
  ],
  preview: {
    allowedHosts: ['platform.archesai.dev'],
    host: '0.0.0.0',
    port: 3000
  },
  resolve: {
    alias: {
      '#': path.resolve(import.meta.dirname, './src')
    }
  },
  server: {
    allowedHosts: ['platform.archesai.dev'],
    host: '0.0.0.0',
    port: 3000
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
