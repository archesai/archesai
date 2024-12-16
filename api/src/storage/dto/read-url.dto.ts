import { Expose } from 'class-transformer'
import { IsString } from 'class-validator'

export class ReadUrlDto {
  /**
   * The read-only url that you can use to download the file from secure storage
   * @example 'www.example.com?token=read-token'
   */
  @IsString()
  @Expose()
  read: string
}
