import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger'
import { Type } from 'class-transformer'
import { IsArray, IsBoolean, IsNumber, IsOptional, IsString, ValidateNested } from 'class-validator'

export class CreateCompletionDto {
  @ApiPropertyOptional()
  @IsNumber()
  @IsOptional()
  best_of?: number

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  echo?: boolean = false

  @ApiPropertyOptional({ default: 0.0 })
  @IsNumber()
  @IsOptional()
  frequency_penalty?: number

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  ignore_eos?: boolean

  @ApiPropertyOptional()
  @IsOptional()
  logit_bias?: Record<string, number>

  @ApiPropertyOptional()
  @IsNumber()
  @IsOptional()
  logprobs?: number

  @ApiPropertyOptional({ default: 16 })
  @IsNumber()
  @IsOptional()
  max_tokens?: number = 16

  @ApiProperty()
  @IsString()
  model: string

  @ApiProperty({ default: 1 })
  @IsNumber()
  n: number = 1

  @ApiPropertyOptional({ default: 0.0 })
  @IsNumber()
  @IsOptional()
  presence_penalty?: number

  @ApiProperty({
    oneOf: [
      { type: 'string' },
      { items: { type: 'number' }, type: 'array' },
      { items: { items: { type: 'number' }, type: 'array' }, type: 'array' },
      { items: { type: 'string' }, type: 'array' }
    ]
  })
  prompt: number[] | number[][] | string | string[]

  @ApiPropertyOptional({ default: true })
  @IsBoolean()
  @IsOptional()
  skip_special_tokens?: boolean

  @ApiPropertyOptional({ type: [String] })
  @IsArray()
  @IsOptional()
  @IsString({ each: true })
  stop?: string[]

  @ApiPropertyOptional({ type: [Number] })
  @IsArray()
  @IsOptional()
  @Type(() => Number)
  @ValidateNested({ each: true })
  stop_token_ids?: number[]

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  stream?: boolean = false

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  suffix?: string

  @ApiProperty({ default: 1.0 })
  @IsNumber()
  temperature?: number = 1

  @ApiProperty({ default: -1 })
  @IsNumber()
  top_k: number

  @ApiProperty({ default: 1.0 })
  @IsNumber()
  top_p?: number = 1

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  use_beam_search?: boolean

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  user?: string
}
