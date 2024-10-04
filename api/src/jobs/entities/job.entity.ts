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

const getResourceLink = (orgname: string, jobType: string, itemId: string) => {
  const mappings = {
    ANIMATION: "animations",
    DOCUMENT: "documents",
    IMAGE: "images",
  };

  return `/organizations/${orgname}/${mappings[jobType]}/${itemId}`;
};

@Exclude()
export class JobEntity extends BaseEntity implements Job {
  @ApiHideProperty()
  animationId: string;

  @ApiProperty({
    description: "The time that the job was completed",
    example: "2023-07-11T21:09:20.895Z",
  })
  @Expose()
  completedAt: Date;

  @ApiHideProperty()
  documentId: string;

  @ApiHideProperty()
  imageId: string;

  @ApiProperty({
    description: "The type of job that is being processed",
    example: "DOCUMENT",
  })
  @Expose()
  jobType: string;

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
    this.resourceLink = getResourceLink(job.orgname, job.jobType, job.id);
  }
}
