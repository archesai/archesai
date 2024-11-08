import { PipelineEntity } from "@/src/pipelines/entities/pipeline.entity";
import { ToolEntity } from "@/src/tools/entities/tool.entity";
import { ApiProperty } from "@nestjs/swagger";
import { Exclude, Expose } from "class-transformer";

import { RunEntity } from "./run.entity";

@Exclude()
export class RunDetailedEntity extends RunEntity {
  @ApiProperty({
    description: "The child runs associated with the run",
    type: [RunEntity],
  })
  @Expose()
  childRuns: RunEntity[];

  @ApiProperty({
    description: "The parent run associated with the run",
    type: RunEntity,
  })
  @Expose()
  parentRun: null | RunEntity;

  @ApiProperty({
    description: "The pipeline associated with the run",
    type: PipelineEntity,
  })
  @Expose()
  pipeline: PipelineEntity;

  @ApiProperty({
    description: "The input contents associated with the run",
    type: [String],
  })
  @Expose()
  runInputContentIds: string[];

  @ApiProperty({
    description: "The output contents associated with the run",
    type: [String],
  })
  @Expose()
  runOutputContentIds: string[];

  @ApiProperty({
    description: "The tool associated with the run",
    type: ToolEntity,
  })
  @Expose()
  tool: ToolEntity;

  constructor(
    run: {
      childRuns: RunEntity[];
      inputContents: { contentId: string }[];
      outputContents: { contentId: string }[];
      parentRun: null | RunEntity;
      pipeline: PipelineEntity;
      tool: ToolEntity;
    } & RunEntity
  ) {
    super(run);
    this.runInputContentIds = run.inputContents.map(
      (content) => content.contentId
    );
    this.runOutputContentIds = run.outputContents.map(
      (content) => content.contentId
    );
    Object.assign(this, run);
  }
}
