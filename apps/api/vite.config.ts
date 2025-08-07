import path from 'path'
import type { FastifyInstance } from 'fastify'
import type { IncomingMessage, ServerResponse } from 'http'
import type { Connect, ViteDevServer } from 'vite'

import { defineConfig } from 'vite'

// FastifyHandler from the plugin code
const FastifyHandler = async ({
  app,
  req,
  res
}: {
  app: FastifyInstance
  req: IncomingMessage
  res: ServerResponse
}) => {
  await app.ready()
  app.routing(req, res)
}

export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 3001,
    allowedHosts: ['api.archesai.dev'],
    hmr: false
  },
  preview: {
    host: '0.0.0.0',
    port: 3001,
    allowedHosts: ['api.archesai.dev']
  },
  build: {
    ssr: './src/main.ts',
    rollupOptions: {
      input: './src/main.ts',
      output: {
        format: 'es'
      }
    },
    commonjsOptions: {
      transformMixedEsModules: true
    }
  },
  optimizeDeps: {
    noDiscovery: true,
    // Vite does not work well with optional dependencies,
    // mark them as ignored for now
    exclude: ['@swc/core']
  },
  resolve: {
    alias: {
      '#': path.resolve(__dirname, './src')
    }
  },
  plugins: [
    {
      name: 'custom-fastify-plugin',
      configureServer: async (server: ViteDevServer) => {
        const logger = server.config.logger
        const appPath = './src/main.ts'
        const exportName = 'app'

        async function loadApp() {
          try {
            const appModule = await server.ssrLoadModule(appPath)
            let app = appModule[exportName]

            if (!app) {
              logger.error(
                `Failed to find a named export ${exportName} from ${appPath}`
              )
              process.exit(1)
            }

            // Handle apps that return a promise
            app = await app
            return app
          } catch (error) {
            logger.error('Failed to load app:')
            logger.error(JSON.stringify(error, null, 2))
            throw error
          }
        }

        // Initialize app on boot
        server.httpServer?.once('listening', async () => {
          try {
            await loadApp()
            logger.info('Fastify app initialized')
          } catch (error) {
            logger.error('Failed to initialize app:')
            logger.error(JSON.stringify(error, null, 2))
          }
        })

        // Reload app on file changes (with debounce)
        let reloadTimeout: NodeJS.Timeout
        server.watcher.on('change', () => {
          clearTimeout(reloadTimeout)
          reloadTimeout = setTimeout(async () => {
            try {
              await loadApp()
              logger.info('App reloaded')
            } catch (error) {
              logger.error('Failed to reload app:')
              logger.error(JSON.stringify(error, null, 2))
            }
          }, 500)
        })

        // Add middleware to handle requests
        server.middlewares.use(
          async (
            req: IncomingMessage,
            res: ServerResponse,
            next: Connect.NextFunction
          ) => {
            try {
              const app = await loadApp()
              if (app) {
                await FastifyHandler({ app, req, res })
              } else {
                next()
              }
            } catch (error) {
              logger.error('Request handling error:')
              logger.error(JSON.stringify(error, null, 2))
              next(error)
            }
          }
        )
      }
    }
  ]
})

// import { builtinModules } from 'node:module'
// import path from 'node:path'

// import { defineConfig } from 'vite'

// import pkg from './package.json' with { type: 'json' }

// export default defineConfig({
//   build: {
//     ssr
//     emptyOutDir: true,
//     reportCompressedSize: true,
//     lib: {
//       entry: 'src/index.ts',
//       formats: ['es'] // optional: ensures ESM output only
//     },
//     rollupOptions: {
//       external: [
//         ...builtinModules,
//         ...builtinModules.map((m) => `node:${m}`),
//         ...Object.keys(pkg.dependencies || {})
//       ]
//     }
//   },

// })
