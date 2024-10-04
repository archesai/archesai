import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { ARTokensService } from "./ar-tokens.service";

@Module({
  exports: [ARTokensService],
  imports: [PrismaModule],
  providers: [ARTokensService],
})
export class ARTokensModule {}
