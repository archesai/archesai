import { ApiProperty } from "@nestjs/swagger";
import { IsArray, IsOptional, IsString } from "class-validator";

export class CreateRunDto {
  @ApiProperty({
    description:
      "If using already created content, specify the content IDs to use as input for the run.",
    example: ["content-id-1", "content-id-2"],
    required: false,
    type: [String],
  })
  @IsArray()
  @IsOptional()
  contentIds?: string[];

  @ApiProperty({
    description:
      "If using direct text input, specify the text to use as input for the run. It will automatically be added as content.",
    example: "This is the text to use as input for the run.",
    required: false,
    type: String,
  })
  @IsString()
  @IsOptional()
  text?: string;

  @ApiProperty({
    description:
      "If using direct text input, specify the text to use as input for the run. It will automatically be added as content.",
    example: "This is a url to use as input for the run.",
    required: false,
    type: String,
  })
  @IsString()
  @IsOptional()
  url?: string;
}
