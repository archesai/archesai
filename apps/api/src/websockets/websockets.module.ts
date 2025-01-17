import { Global, Module } from '@nestjs/common'

import { AuthModule } from '../auth/auth.module'
import { WebsocketsGateway } from './websockets.gateway'
import { WebsocketsService } from './websockets.service'

@Global()
@Module({
  exports: [WebsocketsService],
  imports: [AuthModule],
  providers: [WebsocketsGateway, WebsocketsService]
})
export class WebsocketsModule {}
