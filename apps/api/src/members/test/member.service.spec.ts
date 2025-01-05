import { Test, TestingModule } from '@nestjs/testing'

import { CreateMemberDto } from '../dto/create-member.dto'
import { MemberRepository } from '../member.repository'
import { RoleTypeEnum } from '../entities/member.entity'
import { MembersService } from '../members.service'
import { createMock, DeepMocked } from '@golevelup/ts-jest'
import { UsersService } from '@/src/users/users.service'

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
      const createMemberDto: CreateMemberDto = {
        inviteEmail: 'test@example.com',
        role: RoleTypeEnum.ADMIN
      }
      const existingUser = { username: 'testuser' }

      jest
        .spyOn(mockedUsersService, 'findOneByEmail')
        .mockResolvedValue(existingUser as any)
      jest.spyOn(mockedMemberRepository, 'create').mockResolvedValue({} as any)

      await memberService.create({
        ...createMemberDto,
        orgname
      })

      expect(mockedUsersService.findOneByEmail).toHaveBeenCalledWith(
        createMemberDto.inviteEmail
      )
      expect(mockedMemberRepository.create).toHaveBeenCalledWith({
        inviteEmail: createMemberDto.inviteEmail,
        orgname,
        role: createMemberDto.role
      })
    })

    it('should create a new member without existing user', async () => {
      const orgname = 'test-org'
      const createMemberDto: CreateMemberDto = {
        inviteEmail: 'test@example.com',
        role: RoleTypeEnum.ADMIN
      }

      jest
        .spyOn(mockedUsersService, 'findOneByEmail')
        .mockResolvedValue({} as any)
      jest.spyOn(mockedMemberRepository, 'create').mockResolvedValue({} as any)

      await memberService.create({
        ...createMemberDto,
        orgname
      })

      expect(mockedUsersService.findOneByEmail).toHaveBeenCalledWith(
        createMemberDto.inviteEmail
      )
      expect(mockedMemberRepository.create).toHaveBeenCalledWith({
        inviteEmail: createMemberDto.inviteEmail,
        role: createMemberDto.role,
        orgname
      })
    })
  })

  describe('join', () => {
    it('should update member to join organization', async () => {
      const orgname = 'test-org'
      const inviteEmail = 'test@example.com'
      const username = 'testuser'

      jest.spyOn(mockedMemberRepository, 'update').mockResolvedValue({} as any)

      await memberService.join(orgname, inviteEmail, username)

      expect(mockedMemberRepository.join).toHaveBeenCalledWith(
        orgname,
        inviteEmail,
        username
      )
    })
  })
})
