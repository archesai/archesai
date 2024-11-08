import { BaseEntity } from "@/src/common/dto/base.entity.dto";
import { ApiProperty } from "@nestjs/swagger";
import {
  AuthProvider as _PrismaAuthProvider,
  AuthProviderType,
} from "@prisma/client";
import { Expose } from "class-transformer";
import { IsEnum, IsString } from "class-validator";

export type AuthProviderModel = _PrismaAuthProvider;

export class AuthProviderEntity
  extends BaseEntity
  implements AuthProviderModel
{
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

  constructor(authProvider: AuthProviderModel) {
    super();
    Object.assign(this, authProvider);
  }
}
