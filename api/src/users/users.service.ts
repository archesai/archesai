import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { User } from "@prisma/client";
import { AuthProviderType } from "@prisma/client";

import { OrganizationsService } from "../organizations/organizations.service";
import { UpdateUserDto } from "./dto/update-user.dto";
import { CreateUserInput } from "./types/create-user.type";
import { UserRepository } from "./user.repository";

@Injectable()
export class UsersService {
  constructor(
    private userRepository: UserRepository,
    private organizationsService: OrganizationsService,
    private configService: ConfigService
  ) {}

  async create(createUser: CreateUserInput) {
    const user = await this.userRepository.create({
      emailVerified:
        this.configService.get("FEATURE_EMAIL") === true
          ? createUser.emailVerified
          : true,
      ...createUser,
    });
    await this.organizationsService.createAndInitialize(user, {
      billingEmail: user.email,
      orgname: user.username,
    });
    return this.findOne(user.id);
  }

  async deactivate(id: string) {
    await this.userRepository.deactivate(id);
  }

  async findAll(): Promise<User[]> {
    return this.userRepository.findAll();
  }

  async findOne(id: string) {
    return this.userRepository.findOne(id);
  }

  async findOneByEmail(email: string) {
    return this.userRepository.findOneByEmail(email);
  }

  async remove(id: string) {
    await this.userRepository.remove(id);
  }

  async setEmailVerified(id: string) {
    return this.userRepository.setEmailVerified(id);
  }

  async setEmailVerifiedByEmail(email: string) {
    return this.userRepository.setEmailVerifiedByEmail(email);
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
    return user;
  }

  async update(id: string, updateUserDto: UpdateUserDto) {
    return this.userRepository.update(id, updateUserDto);
  }

  async updateEmail(id: string, email: string) {
    return this.userRepository.updateEmail(id, email);
  }

  async updateRefreshToken(id: string, refreshToken: string) {
    return this.userRepository.updateRefreshToken(id, refreshToken);
  }
}
