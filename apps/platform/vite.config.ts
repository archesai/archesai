import path from 'node:path'

import tailwindcss from '@tailwindcss/vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact from '@vitejs/plugin-react'
// import { visualizer } from 'rollup-plugin-visualizer'
import { defineConfig } from 'vite'

export default defineConfig({
  //   root: __dirname,
  // cacheDir: '../../node_modules/.vite/apps/my-app',
  server: {
    port: 3000,
    host: '0.0.0.0',
    allowedHosts: ['platform.archesai.dev']
  },
  preview: {
    port: 3000,
    host: '0.0.0.0',
    allowedHosts: ['platform.archesai.dev']
  },
  build: {
    outDir: 'dist',
    lib: {
      entry: 'src/main.tsx',
      name: 'platform',
      fileName: 'index',
      formats: ['es' as const]
    }
  },
  resolve: {
    alias: {
      '#': path.resolve(__dirname, './src')
    }
  },
  plugins: [
    tailwindcss(),
    tanstackStart({
      tsr: {
        routesDirectory: 'src/app' // Defaults to "src/routes"
      },
      customViteReactPlugin: true
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
    })
  ],
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
  }
})

// optimizeDeps: {
//   include: ['@archesai/schemas', '@archesai/client'] // your monorepo packages
// },

// visualizer({
//   filename: 'dist/stats.html',
//   open: true,
//   gzipSize: true,
//   brotliSize: true
// })

// export default defineConfig(() => ({
//   root: __dirname,
//   cacheDir: '../../node_modules/.vite/libs/mylib2',
//   plugins: [],
//   // Uncomment this if you are using workers.
//   // worker: {
//   //  plugins: [ nxViteTsPaths() ],
//   // },
//   test: {
//   },
// }));

// /// <reference types='vitest' />
// import { defineConfig } from 'vite';
// import dts from 'vite-plugin-dts';
// import * as path from 'path';

// export default defineConfig(() => ({
//   root: __dirname,
//   cacheDir: '../../node_modules/.vite/libs/mylib',
//   plugins: [
//     dts({
//       entryRoot: 'src',
//       tsconfigPath: path.join(__dirname, 'tsconfig.lib.json'),
//     }),
//   ],
//   // Uncomment this if you are using workers.
//   // worker: {
//   //  plugins: [ nxViteTsPaths() ],
//   // },
//   // Configuration for building your library.
//   // See: https://vitejs.dev/guide/build.html#library-mode
//   build: {
//     outDir: './dist',
//     emptyOutDir: true,
//     reportCompressedSize: true,
//     commonjsOptions: {
//       transformMixedEsModules: true,
//     },
//     lib: {
//       // Could also be a dictionary or array of multiple entry points.
//       entry: 'src/index.ts',
//       name: '@org/mylib',
//       fileName: 'index',
//       // Change this to the formats you want to support.
//       // Don't forget to update your package.json as well.
//       formats: ['es' as const],
//     },
//     rollupOptions: {
//       // External packages that should not be bundled into your library.
//       external: [],
//     },
//   },
//   test: {
//     watch: false,
//     globals: true,
//     environment: 'node',
//     include: ['{src,tests}/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
//     reporters: ['default'],
//     coverage: {
//       reportsDirectory: './test-output/vitest/coverage',
//       provider: 'v8' as const,
//     },
//   },
// }));

// import react from '@vitejs/plugin-react';
// import dts from 'vite-plugin-dts';
// import { nxCopyAssetsPlugin } from '@nx/vite/plugins/nx-copy-assets.plugin';

// export default defineConfig({
//   // ...
//   plugins: [
//     // any needed plugins, but remove nxViteTsPaths
//     react(),
//     nxCopyAssetsPlugin(['*.md', 'package.json']),
//     dts({
//       entryRoot: 'src',
//       tsconfigPath: path.join(__dirname, 'tsconfig.lib.json'),
//     }),
//   ],
//   build: {
//     // ...
//     outDir: './dist',
//     // ...
//     lib: {
//       name: '@myorg/ui',
//       // ...
//     },
//   },
// });
