import { TimedProcessEntity } from "@/src/common/entities/timed-process.entity";
import { ApiHideProperty } from "@nestjs/swagger";
import { Transformation as _PrismaTransformation } from "@prisma/client";
import { Exclude } from "class-transformer";

export type TransformationModel = _PrismaTransformation;

export class TransformationEntity
  extends TimedProcessEntity
  implements TransformationModel
{
  @ApiHideProperty()
  @Exclude()
  pipelineRunId: string;

  @ApiHideProperty()
  @Exclude()
  pipelineStepId: string;

  constructor(transformation: TransformationModel) {
    super();
    Object.assign(this, transformation);
  }
}
