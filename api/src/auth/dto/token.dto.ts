import { ApiProperty } from "@nestjs/swagger";

export class TokenDto {
  @ApiProperty({
    description: "The authorization token that can be used to access Arches AI",
    example: "supersecretauthorizationtoken",
  })
  accessToken: string;

  @ApiProperty({
    description: "The refresh token that can be used to get a new access token",
    example: "supersecretauthorizationtoken",
  })
  refreshToken: string;

  constructor({
    accessToken,
    refreshToken,
  }: {
    accessToken: string;
    refreshToken: string;
  }) {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
  }
}
