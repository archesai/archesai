import { Injectable, Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { AuthProviderType } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { CreateUserDto } from "./dto/create-user.dto";
import { UpdateUserDto } from "./dto/update-user.dto";
import {
  UserEntity,
  UserWithMembershipsAndAuthProvidersModel,
} from "./entities/user.entity";
import { UserRepository } from "./user.repository";

@Injectable()
export class UsersService extends BaseService<
  UserEntity,
  undefined,
  UpdateUserDto,
  UserRepository,
  UserWithMembershipsAndAuthProvidersModel
> {
  private readonly logger: Logger = new Logger(UsersService.name);
  constructor(
    private userRepository: UserRepository,
    private organizationsService: OrganizationsService,
    private configService: ConfigService
  ) {
    super(userRepository);
  }

  async create(orgname: string, createUserDto: CreateUserDto) {
    const user = await this.userRepository.create("", {
      emailVerified:
        this.configService.get("FEATURE_EMAIL") === true
          ? createUserDto.emailVerified
          : true,
      ...createUserDto,
    });
    await this.organizationsService.createAndInitialize(user, {
      billingEmail: user.email,
      orgname: user.username,
    });
    return this.toEntity(await this.userRepository.findOne("", user.id));
  }

  async deactivate(id: string) {
    await this.userRepository.deactivate(id);
  }

  async findOneByEmail(email: string) {
    return this.toEntity(await this.userRepository.findOneByEmail(email));
  }

  async setEmail(id: string, newEmail: string) {
    return this.toEntity(
      await this.userRepository.updateRaw(null, id, {
        email: newEmail,
      })
    );
  }

  async setEmailVerified(id: string) {
    return this.toEntity(
      await this.userRepository.updateRaw(null, id, {
        emailVerified: true,
      })
    );
  }

  async setEmailVerifiedByEmail(email: string) {
    const user = await this.findOneByEmail(email);
    return this.toEntity(await this.setEmailVerified(user.id));
  }

  async setRefreshToken(id: string, refreshToken: string) {
    return this.toEntity(
      await this.userRepository.updateRaw(null, id, {
        refreshToken,
      })
    );
  }

  async syncAuthProvider(
    email: string,
    provider: AuthProviderType,
    providerId: string
  ) {
    const user = await this.userRepository.findOneByEmail(email);
    // if it does not have this provider, add it
    if (!user.authProviders.some((p) => p.provider === provider)) {
      return this.userRepository.addAuthProvider(email, provider, providerId);
    }
    return this.toEntity(user);
  }

  protected toEntity(
    model: UserWithMembershipsAndAuthProvidersModel
  ): UserEntity {
    return new UserEntity(model);
  }
}
