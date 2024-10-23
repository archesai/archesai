import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Job } from "@prisma/client";
import { JobStatus } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";

import { BaseEntity } from "../../common/base-entity.dto";

export const jobMap = {
  ANIMATION: "animationId",
  DOCUMENT: "documentId",
  IMAGE: "imagesId",
};

@Exclude()
export class JobEntity extends BaseEntity implements Job {
  @ApiProperty({
    description: "The time that the job was completed",
    example: "2023-07-11T21:09:20.895Z",
  })
  @Expose()
  completedAt: Date;

  @ApiProperty({
    description: "The error message if the job failed",
    example: "Could not process the document",
  })
  @Expose()
  error: string;

  // Private Properties
  @ApiHideProperty()
  orgname: string;

  @ApiProperty({
    description: "The percent progress of the current job",
    example: 0.9,
  })
  @Expose()
  progress: number;

  @ApiProperty({
    description: "The link to the resource that is being processed",
    example: "/organizations/archesai/documents/documentId",
  })
  @Expose()
  resourceLink: string;

  // Public Properties
  @ApiProperty({
    description: "The time that the job was started",
    example: "2023-07-11T21:09:20.895Z",
  })
  @Expose()
  startedAt: Date;

  @ApiProperty({
    description: "The status of the current animation processing",
    enum: JobStatus,
    example: JobStatus.COMPLETE,
  })
  @Expose()
  status: JobStatus;

  constructor(job: Job) {
    super();
    Object.assign(this, job);
    this.resourceLink = `/organizations/${job.orgname}/content/${job.id}`;
  }
}
