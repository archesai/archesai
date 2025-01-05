import { IntersectionType, PartialType, PickType } from '@nestjs/swagger'
import { RunEntity } from '@/src/runs/entities/run.entity'
import { IsArray, IsOptional, IsString } from 'class-validator'
import { Expose } from 'class-transformer'

export class CreateRunDto extends IntersectionType(
  PickType(RunEntity, ['runType'] as const),
  PartialType(PickType(RunEntity, ['pipelineId', 'toolId'] as const))
) {
  /**
   * If using already created content, specify the content IDs to use as input for the run.
   * example: ['content-id-1', 'content-id-2']
   */
  @IsOptional()
  @IsArray()
  @IsString({ each: true })
  @Expose()
  contentIds?: string[]

  /**
   * If using direct text input, specify the text to use as input for the run. It will automatically be added as content.
   * example: 'This is the text to use as input for the run.'
   */
  @IsOptional()
  @IsString()
  @Expose()
  text?: string

  /**
   * If using a URL as input, specify the URL to use as input for the run. It will automatically be added as content.
   * example: 'https://example.com'
   */
  @IsOptional()
  @IsString()
  @Expose()
  url?: string
}
