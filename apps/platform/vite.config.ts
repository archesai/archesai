import fs from 'node:fs'
import path from 'node:path'

import tailwindcss from '@tailwindcss/vite'
import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 3000,
    proxy: {},
    https: {
      key: fs.readFileSync(
        path.resolve(
          '../../deploy/kubernetes/overlays/development/certs/localhost-key.pem'
        )
      ),
      cert: fs.readFileSync(
        path.resolve(
          '../../deploy/kubernetes/overlays/development/certs/localhost.pem'
        )
      )
    }
  },
  plugins: [
    tailwindcss(),
    // Enables Vite to resolve imports using path aliases.
    tanstackStart({
      tsr: {
        // Specifies the directory TanStack Router uses for your routes.
        routesDirectory: 'src/app' // Defaults to "src/routes"
      }
    })
  ]
})
