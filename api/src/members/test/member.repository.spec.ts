import { Test, TestingModule } from '@nestjs/testing'

import { PrismaService } from '../../prisma/prisma.service'
import { CreateMemberDto } from '../dto/create-member.dto'
import { MemberRepository } from '../member.repository'

describe('MemberRepository', () => {
  let memberRepository: MemberRepository
  let prismaService: PrismaService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        MemberRepository,
        {
          provide: PrismaService,
          useValue: {
            member: {
              create: jest.fn(),
              update: jest.fn()
            },
            organization: {
              findUniqueOrThrow: jest.fn()
            },
            user: {
              findUnique: jest.fn()
            }
          }
        }
      ]
    }).compile()

    memberRepository = module.get<MemberRepository>(MemberRepository)
    prismaService = module.get<PrismaService>(PrismaService)
  })

  describe('create', () => {
    it('should create a new member', async () => {
      const orgname = 'test-org'
      const createMemberDto: CreateMemberDto = {
        inviteEmail: 'test@example.com',
        role: 'ADMIN'
      }
      const existingUser = { username: 'testuser' }

      jest.spyOn(prismaService.user, 'findUnique').mockResolvedValue(existingUser as any)
      jest.spyOn(prismaService.member, 'create').mockResolvedValue({} as any)

      await memberRepository.create(orgname, createMemberDto)

      expect(prismaService.user.findUnique).toHaveBeenCalledWith({
        where: { email: createMemberDto.inviteEmail }
      })
      expect(prismaService.member.create).toHaveBeenCalledWith({
        data: {
          inviteEmail: createMemberDto.inviteEmail,
          organization: { connect: { orgname } },
          role: createMemberDto.role,
          user: { connect: { username: existingUser.username } }
        }
      })
    })

    it('should create a new member without existing user', async () => {
      const orgname = 'test-org'
      const createMemberDto: CreateMemberDto = {
        inviteEmail: 'test@example.com',
        role: 'ADMIN'
      }

      jest.spyOn(prismaService.user, 'findUnique').mockResolvedValue(null)
      jest.spyOn(prismaService.member, 'create').mockResolvedValue({} as any)

      await memberRepository.create(orgname, createMemberDto)

      expect(prismaService.user.findUnique).toHaveBeenCalledWith({
        where: { email: createMemberDto.inviteEmail }
      })
      expect(prismaService.member.create).toHaveBeenCalledWith({
        data: {
          inviteEmail: createMemberDto.inviteEmail,
          organization: { connect: { orgname } },
          role: createMemberDto.role
        }
      })
    })
  })

  describe('join', () => {
    it('should update member to join organization', async () => {
      const orgname = 'test-org'
      const inviteEmail = 'test@example.com'
      const username = 'testuser'

      jest.spyOn(prismaService.organization, 'findUniqueOrThrow').mockResolvedValue({} as any)
      jest.spyOn(prismaService.member, 'update').mockResolvedValue({} as any)

      await memberRepository.join(orgname, inviteEmail, username)

      expect(prismaService.organization.findUniqueOrThrow).toHaveBeenCalledWith({
        where: { orgname }
      })
      expect(prismaService.member.update).toHaveBeenCalledWith({
        data: {
          inviteAccepted: true,
          user: { connect: { username } }
        },
        where: {
          inviteEmail_orgname: {
            inviteEmail,
            orgname
          }
        }
      })
    })
  })
})
