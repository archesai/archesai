import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { PrismaService } from "../prisma/prisma.service";
import { CreateToolDto } from "./dto/create-tool.dto";
import { UpdateToolDto } from "./dto/update-tool.dto";
import { ToolModel } from "./entities/tool.entity";

@Injectable()
export class ToolRepository extends BaseRepository<
  ToolModel,
  CreateToolDto,
  UpdateToolDto,
  Prisma.ToolInclude,
  Prisma.ToolSelect,
  Prisma.ToolUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.tool);
  }
}
