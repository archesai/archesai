import { bootstrap } from '#utils/bootstrap'

bootstrap().catch((err: unknown) => {
  console.error('Failed to start application:', err)
  process.exit(1)
})
