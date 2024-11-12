import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { ContentEntity } from "../content/entities/content.entity";
import { PrismaService } from "../prisma/prisma.service";
import { ToolRunModel } from "./entities/tool-run.entity";

@Injectable()
export class ToolRunRepository extends BaseRepository<
  ToolRunModel,
  any,
  any,
  Prisma.ToolRunInclude,
  Prisma.ToolRunUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.toolRun);
  }

  async setOutputContent(toolRunId: string, contents: ContentEntity[]) {
    await this.prisma.toolRun.update({
      data: {
        inputs: {
          connect: contents.map((content) => ({ id: content.id })),
        },
      },
      where: { id: toolRunId },
    });
    return this.prisma.toolRun.findUnique({
      where: { id: toolRunId },
    });
  }
}
