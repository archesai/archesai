import { Module } from "@nestjs/common";

import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { UserRepository } from "./user.repository";
import { UsersController } from "./users.controller";
import { UsersService } from "./users.service";

@Module({
  controllers: [UsersController],
  exports: [UsersService],
  imports: [PrismaModule, OrganizationsModule],
  providers: [UsersService, UserRepository],
})
export class UsersModule {}
