import { BaseEntity } from "@/src/common/base-entity.dto";
import { ApiProperty } from "@nestjs/swagger";
import { AuthProvider, AuthProviderType } from "@prisma/client";
import { Expose } from "class-transformer";
import { IsEnum, IsString } from "class-validator";

export class AuthProviderEntity extends BaseEntity implements AuthProvider {
  @ApiProperty({
    description: "The auth provider's provider",
    enum: AuthProviderType,
  })
  @Expose()
  @IsString()
  @IsEnum(AuthProviderType)
  provider: AuthProviderType;

  @ApiProperty({
    description: "The auth provider's provider ID",
  })
  @Expose()
  @IsString()
  providerId: string;

  @ApiProperty({
    description: "The auth provider's user ID",
  })
  @Expose()
  @IsString()
  userId: string;

  constructor(authProvider: AuthProvider) {
    super();
    Object.assign(this, authProvider);
  }
}
