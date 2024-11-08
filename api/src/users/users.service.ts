import { Injectable, Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { AuthProvider, AuthProviderType, Member, User } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateUserDto } from "./dto/create-user.dto";
import { UpdateUserDto } from "./dto/update-user.dto";
import { UserEntity } from "./entities/user.entity";
import { UserRepository } from "./user.repository";

@Injectable()
export class UsersService extends BaseService<
  UserEntity,
  undefined,
  UpdateUserDto,
  UserRepository,
  { memberships: Member[] } & User
> {
  private readonly logger: Logger = new Logger(UsersService.name);
  constructor(
    private userRepository: UserRepository,
    private organizationsService: OrganizationsService,
    private configService: ConfigService,
    private websocketsService: WebsocketsService
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

  async setEmailVerified(id: string) {
    return this.toEntity(await this.userRepository.setEmailVerified(id));
  }

  async setEmailVerifiedByEmail(email: string) {
    return this.toEntity(
      await this.userRepository.setEmailVerifiedByEmail(email)
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
    model: { authProviders: AuthProvider[]; memberships: Member[] } & User
  ): UserEntity {
    return new UserEntity(model);
  }

  async updateEmail(id: string, email: string) {
    return this.userRepository.updateEmail(id, email);
  }

  async updateRefreshToken(id: string, refreshToken: string) {
    return this.toEntity(
      await this.userRepository.updateRefreshToken(id, refreshToken)
    );
  }
}
