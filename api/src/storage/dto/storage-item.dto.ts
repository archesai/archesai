import { ApiProperty } from "@nestjs/swagger";
import { IsBoolean, IsNumber, IsString } from "class-validator";

export class StorageItemDto {
  @ApiProperty({
    description: "Whether or not this is a directory",
    example: true,
  })
  createdAt: Date;

  @ApiProperty({
    description: "The id of the storage item",
    example: "14",
  })
  @IsString()
  id: string;

  @ApiProperty({
    description: "Whether or not this is a directory",
    example: true,
  })
  @IsBoolean()
  isDir: boolean;

  @ApiProperty({
    description: "The path that the file is located in",
    example: "/location/in/storage",
  })
  @IsString()
  name: string;

  @ApiProperty({
    description: "The size of the item in bytes",
    example: 12341234,
  })
  @IsNumber()
  size: number;

  constructor(storageItemDto: StorageItemDto) {
    Object.assign(this, storageItemDto);
  }
}
