import { Injectable } from "@nestjs/common";
import { Prisma } from "@prisma/client";

import { BaseRepository } from "../common/base.repository";
import { ContentEntity } from "../content/entities/content.entity";
import { PrismaService } from "../prisma/prisma.service";
import { TransformationModel } from "./entities/transformation.entity";

@Injectable()
export class TransformationRepository extends BaseRepository<
  TransformationModel,
  any,
  any,
  Prisma.TransformationInclude,
  Prisma.TransformationSelect,
  Prisma.TransformationUpdateInput
> {
  constructor(private prisma: PrismaService) {
    super(prisma.transformation);
  }

  async setOutputContent(transformationId: string, contents: ContentEntity[]) {
    await this.prisma.transformation.update({
      data: {
        inputs: {
          connect: contents.map((content) => ({ id: content.id })),
        },
      },
      where: { id: transformationId },
    });
    return this.prisma.transformation.findUnique({
      where: { id: transformationId },
    });
  }
}
