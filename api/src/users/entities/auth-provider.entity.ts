import { BaseEntity } from "@/src/common/entities/base.entity";
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
  /**
   * The auth provider's provider
   * @example LOCAL
   */
  @Expose()
  @IsEnum(AuthProviderType)
  @IsString()
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
