import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { ToolRepository } from "./tool.repository";
import { ToolsController } from "./tools.controller";
import { ToolsService } from "./tools.service";

@Module({
  controllers: [ToolsController],
  exports: [ToolsService],
  imports: [PrismaModule],
  providers: [ToolsService, ToolRepository],
})
export class ToolsModule {}
