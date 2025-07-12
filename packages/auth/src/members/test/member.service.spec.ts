import type { DeepMocked } from '@golevelup/ts-jest'
import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import type { BaseInsertion, MemberEntity } from '@archesai/schemas'

import { createRandomMember } from '#members/factories/member.factory'
import { MemberRepository } from '#members/member.repository'
import { MembersService } from '#members/members.service'
import { createRandomUser } from '#users/factories/user.factory'
import { UsersService } from '#users/users.service'

describe('MembersService', () => {
  let mockedMemberRepository: DeepMocked<MemberRepository>
  let mockedUsersService: DeepMocked<UsersService>
  let memberService: MembersService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [MembersService]
    })
      .useMocker(createMock)
      .compile()

    memberService = module.get(MembersService)
    mockedUsersService = module.get(UsersService)
    mockedMemberRepository = module.get(MemberRepository)
  })

  describe('create', () => {
    it('should create a new smember', async () => {
      const orgname = 'test-org'
      const createMemberRequest: BaseInsertion<MemberEntity> = {
        createdAt: new Date().toISOString(),
        invitationId: 'test-invitation-id',
        orgname,
        role: 'ADMIN',
        updatedAt: new Date().toISOString(),
        userId: 'test-user-id'
      }
      const existingUser = createRandomUser()

      jest
        .spyOn(mockedUsersService, 'findOneByEmail')
        .mockResolvedValue(existingUser)
      jest
        .spyOn(mockedMemberRepository, 'create')
        .mockResolvedValue(createRandomMember())

      await memberService.create(createMemberRequest)
    })

    it('should create a new member without existing user', async () => {
      const orgname = 'test-org'
      const createMemberRequest: BaseInsertion<MemberEntity> = {
        createdAt: new Date().toISOString(),
        invitationId: 'test-invitation-id',
        orgname,
        role: 'ADMIN',
        updatedAt: new Date().toISOString(),
        userId: 'test-user-id'
      }

      jest
        .spyOn(mockedUsersService, 'findOneByEmail')
        .mockResolvedValue(createRandomUser())
      jest
        .spyOn(mockedMemberRepository, 'create')
        .mockResolvedValue(createRandomMember())

      await memberService.create({
        ...createMemberRequest,
        orgname
      })
    })
  })

  describe('join', () => {
    it('should update member to join organization', () => {
      jest
        .spyOn(mockedMemberRepository, 'update')
        .mockResolvedValue(createRandomMember())
    })
  })
})
