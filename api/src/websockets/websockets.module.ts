import { Global, Module } from "@nestjs/common";

import { AuthModule } from "../auth/auth.module";
import { UsersModule } from "../users/users.module";
import { WebsocketsGateway } from "./websockets.gateway";
import { WebsocketsService } from "./websockets.service";

@Global()
@Module({
  exports: [WebsocketsService],
  imports: [AuthModule, UsersModule],
  providers: [WebsocketsGateway, WebsocketsService],
})
export class WebsocketsModule {}
