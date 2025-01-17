import { Module } from '@nestjs/common'

import { PrismaModule } from '../prisma/prisma.module'
import { MemberRepository } from './member.repository'
import { MembersController } from './members.controller'
import { MembersService } from './members.service'
import { UsersModule } from '../users/users.module'
import { AuthModule } from '../auth/auth.module'

@Module({
  controllers: [MembersController],
  exports: [MembersService],
  imports: [PrismaModule, UsersModule, AuthModule],
  providers: [MembersService, MemberRepository]
})
export class MembersModule {}
