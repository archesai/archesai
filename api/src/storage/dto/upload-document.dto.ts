import { ApiProperty } from "@nestjs/swagger";
import { IsEmpty } from "class-validator";

export class UploadDocumentDto {
  @ApiProperty({
    description: "The file to upload",
    format: "binary",
    type: "string",
  })
  @IsEmpty()
  file: Express.Multer.File;
}
