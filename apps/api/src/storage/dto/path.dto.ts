import { Expose } from 'class-transformer'
import { IsBoolean, IsOptional, IsString } from 'class-validator'

export class PathDto {
  /**
   * Whether or not this path points to a directory
   * @example false
   */
  @IsOptional()
  @IsBoolean()
  @Expose()
  isDir: boolean = false

  /**
   * The path that the file should upload to
   * @example '/location/in/storage'
   */
  @IsString()
  @Expose()
  path: string
}
