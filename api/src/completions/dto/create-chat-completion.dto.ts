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

export class CreateChatCompletionDto {
  @ApiPropertyOptional()
  @IsOptional()
  @IsNumber()
  best_of?: number;

  @ApiPropertyOptional({ default: 0.0 })
  @IsOptional()
  @IsNumber()
  frequency_penalty?: number;

  @ApiPropertyOptional({ default: false })
  @IsOptional()
  @IsBoolean()
  ignore_eos?: boolean;

  @ApiPropertyOptional()
  @IsOptional()
  @IsObject()
  logit_bias?: Record<string, number>;

  @ApiPropertyOptional()
  @IsOptional()
  @IsNumber()
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
  @IsOptional()
  @IsNumber()
  presence_penalty?: number;

  @ApiPropertyOptional({ default: true })
  @IsOptional()
  @IsBoolean()
  skip_special_tokens?: boolean;

  @ApiPropertyOptional({ type: [String] })
  @IsOptional()
  @IsArray()
  @IsString({ each: true })
  stop?: string[];

  @ApiPropertyOptional({ type: [Number] })
  @IsOptional()
  @IsArray()
  @ValidateNested({ each: true })
  @Type(() => Number)
  stop_token_ids?: number[];

  @ApiPropertyOptional({ default: false })
  @IsOptional()
  @IsBoolean()
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
  @IsOptional()
  @IsBoolean()
  use_beam_search?: boolean;

  @ApiPropertyOptional()
  @IsOptional()
  @IsString()
  user?: string;
}
