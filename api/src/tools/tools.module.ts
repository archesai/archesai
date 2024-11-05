import { Module } from "@nestjs/common";

import { PrismaModule } from "../prisma/prisma.module";
import { RunsModule } from "../runs/runs.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ToolRepository } from "./tool.repository";
import { ToolsController } from "./tools.controller";
import { ToolsService } from "./tools.service";

@Module({
  controllers: [ToolsController],
  exports: [ToolsService],
  imports: [PrismaModule, WebsocketsModule, RunsModule],
  providers: [ToolsService, ToolRepository],
})
export class ToolsModule {}
