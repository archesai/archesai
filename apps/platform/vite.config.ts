import tailwindcss from '@tailwindcss/vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 3000,
    host: '0.0.0.0',
    allowedHosts: ['platform.archesai.dev'],
    proxy: {}
  },
  plugins: [
    tailwindcss(),
    // Enables Vite to resolve imports using path aliases.
    tanstackStart({
      tsr: {
        // Specifies the directory TanStack Router uses for your routes.
        routesDirectory: 'src/app' // Defaults to "src/routes"
      },
      react: {
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
      }
    })
  ]
})
