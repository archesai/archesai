import { Module } from "@nestjs/common";

import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { UserController } from "./user.controller";
import { UserRepository } from "./user.repository";
import { UsersService } from "./users.service";

@Module({
  controllers: [UserController],
  exports: [UsersService],
  imports: [PrismaModule, OrganizationsModule],
  providers: [UsersService, UserRepository],
})
export class UsersModule {}
