import { SubItemEntity } from "@/src/common/entities/base-sub-item.entity";
import { BaseEntity } from "@/src/common/entities/base.entity";
import { ToolEntity } from "@/src/tools/entities/tool.entity";
import {
  PipelineStep as _PrismaPipelineStep,
  Tool as _PrismaTool,
} from "@prisma/client";
import { Exclude, Expose } from "class-transformer";

type PipelineStepModel = _PrismaPipelineStep & {
  tool: _PrismaTool;
};

@Exclude()
export class PipelineStepEntity
  extends BaseEntity
  implements PipelineStepModel
{
  /**
   * The order of the step in the pipeline
   */
  @Expose()
  dependents: SubItemEntity[];

  /**
   * These are the steps that this step depends on.
   */
  @Expose()
  dependsOn: SubItemEntity[];

  /**
   * The name of the step in the pipeline. It must be unique within the pipeline.
   */
  @Expose()
  name: string;

  /**
   * The ID of the pipelin that this step belongs to
   * @example 'pipeline-id'
   */
  @Expose()
  pipelineId: string;

  /**
   * The name of the tool that this step uses.
   */
  @Expose()
  tool: ToolEntity;

  /**
   * This is the ID of the tool that this step uses.
   * @example 'tool-id'
   */
  @Expose()
  toolId: string;

  constructor(pipelineStep: PipelineStepModel) {
    super();
    Object.assign(this, pipelineStep);
  }
}
