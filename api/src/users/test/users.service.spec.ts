import { createRandomUser } from '@/prisma/factories/user.factory'
import { createMock, DeepMocked } from '@golevelup/ts-jest'
import { Test, TestingModule } from '@nestjs/testing'
import { AuthProviderType } from '@prisma/client'

import { OrganizationsService } from '../../organizations/organizations.service'
import { CreateUserDto } from '../dto/create-user.dto'
import { UserRepository } from '../user.repository'
import { UsersService } from '../users.service'

describe('UsersService', () => {
  let service: UsersService
  let mockedUserRepository: DeepMocked<UserRepository>
  let mockedOrganizationsService: DeepMocked<OrganizationsService>

  beforeEach(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      providers: [UsersService]
    })
      .useMocker(createMock)
      .compile()

    service = moduleRef.get(UsersService)
    mockedUserRepository = moduleRef.get(UserRepository)
    mockedOrganizationsService = moduleRef.get(OrganizationsService)
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  describe('create', () => {
    it('should create a user and associated organization', async () => {
      const createUserDto: CreateUserDto = {
        email: 'test@test.com',
        emailVerified: false,
        photoUrl: 'testphoto',
        username: 'testuser'
      }
      const mockedUser = createRandomUser(createUserDto)

      mockedUserRepository.create.mockResolvedValue(mockedUser)
      mockedUserRepository.findOne.mockResolvedValue(mockedUser)

      await service.create('orgname', createUserDto)

      expect(mockedUserRepository.create).toHaveBeenCalledWith('', createUserDto)
      expect(mockedOrganizationsService.create).toHaveBeenCalled()
    })
  })

  describe('syncAuthProvider', () => {
    it('should add auth provider if not exists', async () => {
      const email = 'test@test.com'
      const provider = AuthProviderType.LOCAL
      const providerId = '123'

      const mockedUser = createRandomUser()

      mockedUserRepository.findOneByEmail.mockResolvedValue(mockedUser)
      mockedUserRepository.addAuthProvider.mockResolvedValue({
        authProviders: [
          {
            createdAt: new Date(),
            id: '1',
            provider,
            providerId,
            userId: mockedUser.id
          }
        ],
        ...mockedUser
      })

      await service.syncAuthProvider(email, provider, providerId)

      expect(mockedUserRepository.addAuthProvider).toHaveBeenCalledWith(email, provider, providerId)
    })
  })

  describe('setEmailVerified', () => {
    it('should set email as verified', async () => {
      const mockedUser = createRandomUser()

      mockedUserRepository.updateRaw.mockResolvedValue(mockedUser)
      await service.setEmailVerified(mockedUser.id)

      expect(mockedUserRepository.updateRaw).toHaveBeenCalledWith(null, mockedUser.id, {
        emailVerified: true
      })
    })
  })

  describe('deactivate', () => {
    it('should deactivate user', async () => {
      const mockedUser = createRandomUser()
      mockedUserRepository.deactivate.mockResolvedValue()
      await service.deactivate(mockedUser.id)

      expect(mockedUserRepository.deactivate).toHaveBeenCalledWith(mockedUser.id)
    })
  })
})
