import { ApiProperty } from '@nestjs/swagger'

export class ReadUrlDto {
  @ApiProperty({
    description:
      'A read-only url that you can use to download the file from secure storage',
    example: 'www.example.com?token=read-token'
  })
  read!: string
}
