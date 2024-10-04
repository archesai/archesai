import { ApiProperty } from "@nestjs/swagger";
import { Content, Job } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsBoolean, IsNumber, IsOptional, IsString } from "class-validator";

import { ContentEntity } from "../../content/entities/content.entity";
import { JobEntity } from "../../jobs/entities/job.entity";

@Exclude()
export class AnimationEntity extends ContentEntity implements Content {
  @ApiProperty({
    description: "The animation prompts of the animation",
    example: "3D",
  })
  @Expose()
  @IsString()
  animationPrompts: string;

  @ApiProperty({
    default: 0,
    description: "The audio start time in seconds",
    required: false,
    type: Number,
  })
  @Expose()
  @IsNumber()
  @IsOptional()
  audioStart = 0;

  @ApiProperty({
    default: "",
    description: "The file path or url of the audio",
    example: "/animations/1125162/init.png",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  @IsOptional()
  audioUrl = "";

  @ApiProperty({
    description: "The frames used per second",
    example: 30,
  })
  @Expose()
  @IsNumber()
  fps: number;

  @ApiProperty({
    description: "The height of the animation",
    example: 768,
  })
  @Expose()
  @IsNumber()
  height: number;

  @ApiProperty({
    description: "The length of the animation",
    example: 234,
  })
  @Expose()
  @IsNumber()
  length: number;

  @ApiProperty({
    description: "The number of frames used in this animation",
    example: 60,
  })
  @Expose()
  @IsNumber()
  maxFrames: number;

  @ApiProperty({
    default: false,
    description: "Whether or not to use the audio",
    required: false,
    type: Boolean,
  })
  @Expose()
  @IsBoolean()
  @IsOptional()
  useAudio = false;

  @ApiProperty({
    description: "The width of the animation",
    example: 768,
  })
  @Expose()
  @IsNumber()
  width: number;

  constructor(content: { job: Job } & Content) {
    super(content);
    Object.assign(this, content);
    this.length = this.maxFrames / this.fps;
    this.job = new JobEntity(content.job);
  }
}
