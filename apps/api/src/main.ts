import { ConfigService } from '@archesai/core'

import { setup } from '#utils/setup'

async function bootstrap() {
  const app = await setup()

  const configService = app.get(ConfigService)
  await app.listen({
    host: configService.get('server.host'),
    port: configService.get('server.port')
  })
}

bootstrap().catch((err: unknown) => {
  console.error(err)
  process.exit(1)
})
