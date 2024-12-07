import { Module } from '@nestjs/common'

import { PrismaModule } from '../prisma/prisma.module'
import { MemberRepository } from './member.repository'
import { MembersController } from './members.controller'
import { MembersService } from './members.service'

@Module({
  controllers: [MembersController],
  exports: [MembersService],
  imports: [PrismaModule],
  providers: [MembersService, MemberRepository]
})
export class MembersModule {}
