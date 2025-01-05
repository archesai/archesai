import { Expose } from 'class-transformer'
import { IsString } from 'class-validator'

export class WriteUrlDto {
  /**
   * A write-only url that you can use to upload a file to secure storage
   * @example 'www.example.com?token=write-token'
   */
  @IsString()
  @Expose()
  write: string
}
