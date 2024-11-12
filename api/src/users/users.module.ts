import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { UserController } from "./user.controller";
import { UserRepository } from "./user.repository";
import { UsersService } from "./users.service";

@Module({
  controllers: [UserController],
  exports: [UsersService],
  imports: [PrismaModule, OrganizationsModule, ConfigModule],
  providers: [UsersService, UserRepository],
})
export class UsersModule {}
