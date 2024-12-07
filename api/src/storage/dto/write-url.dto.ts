import { ApiProperty } from '@nestjs/swagger'

export class WriteUrlDto {
  @ApiProperty({
    description:
      'A write-only url that you can use to upload a file to secure storage',
    example: 'www.example.com?token=write-token'
  })
  write!: string
}
