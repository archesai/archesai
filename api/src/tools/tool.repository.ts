import { Injectable } from "@nestjs/common";
import { Prisma, Tool } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";

@Injectable()
export class ToolRepository extends BaseRepository<
  Tool,
  CreateToolDto,
  UpdateToolDto,
  Prisma.ToolInclude,
  Prisma.ToolSelect
> {
  constructor(private prisma: PrismaService) {
    super(prisma.tool);
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ToolUpdateInput) {
    return this.prisma.tool.update({
      data: raw,
      where: {
        id,
      },
    });
  }
}
