import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { JwtModule } from "@nestjs/jwt";

import { PrismaModule } from "../prisma/prisma.module";
import { ApiTokenRepository } from "./api-token.repository";
import { ApiTokensController } from "./api-tokens.controller";
import { ApiTokensService } from "./api-tokens.service";

@Module({
  controllers: [ApiTokensController],
  exports: [ApiTokensService],
  imports: [PrismaModule, ConfigModule, JwtModule],
  providers: [ApiTokensService, ApiTokenRepository],
})
export class ApiTokensModule {}
