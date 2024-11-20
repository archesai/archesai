import { ApiProperty, ApiPropertyOptional } from "@nestjs/swagger";
import { Type } from "class-transformer";
import {
  IsArray,
  IsBoolean,
  IsNumber,
  IsObject,
  IsOptional,
  IsString,
  ValidateNested,
} from "class-validator";

export class CreateChatCompletionDto {
  @ApiPropertyOptional()
  @IsNumber()
  @IsOptional()
  best_of?: number;

  @ApiPropertyOptional({ default: 0.0 })
  @IsNumber()
  @IsOptional()
  frequency_penalty?: number;

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  ignore_eos?: boolean;

  @ApiPropertyOptional()
  @IsObject()
  @IsOptional()
  logit_bias?: Record<string, number>;

  @ApiPropertyOptional()
  @IsNumber()
  @IsOptional()
  max_tokens?: number;

  @ApiProperty({
    oneOf: [{ items: { type: "object" }, type: "array" }],
  })
  messages: MessageDto[];

  @ApiProperty()
  @IsString()
  model?: string;

  @ApiProperty({ default: 1 })
  @IsNumber()
  n?: number = 1;

  name?: string;

  @ApiPropertyOptional({ default: 0.0 })
  @IsNumber()
  @IsOptional()
  presence_penalty?: number;

  @ApiPropertyOptional({ default: true })
  @IsBoolean()
  @IsOptional()
  skip_special_tokens?: boolean;

  @ApiPropertyOptional({ type: [String] })
  @IsArray()
  @IsOptional()
  @IsString({ each: true })
  stop?: string[];

  @ApiPropertyOptional({ type: [Number] })
  @IsArray()
  @IsOptional()
  @Type(() => Number)
  @ValidateNested({ each: true })
  stop_token_ids?: number[];

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  stream?: boolean;

  @ApiProperty({ default: 0.7 })
  @IsNumber()
  temperature?: number = 0.7;

  @ApiProperty({ default: -1 })
  @IsNumber()
  top_k?: number = -1;

  @ApiProperty({ default: 1.0 })
  @IsNumber()
  top_p?: number = 1;

  @ApiPropertyOptional({ default: false })
  @IsBoolean()
  @IsOptional()
  use_beam_search?: boolean;

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  user?: string;
}

export class MessageDto {
  @ApiProperty()
  @IsString()
  content: string;

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  name?: string;

  @ApiProperty()
  @IsString()
  role: "assistant" | "function" | "system" | "user";
}
