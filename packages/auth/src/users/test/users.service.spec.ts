import type { DeepMocked } from '@golevelup/ts-jest'
import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import type { BaseInsertion, UserEntity } from '@archesai/schemas'

import { OrganizationsService } from '#organizations/organizations.service'
import { createRandomUser } from '#users/factories/user.factory'
import { UserRepository } from '#users/user.repository'
import { UsersService } from '#users/users.service'

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
      const createUserDto = {
        deactivated: false,
        email: 'test@test.com',
        image: 'testphoto',
        name: 'test',
        orgname: 'test-org'
      } satisfies BaseInsertion<UserEntity>
      const mockedUser = createRandomUser(createUserDto)

      mockedUserRepository.create.mockResolvedValue(mockedUser)
      mockedUserRepository.findOne.mockResolvedValue(mockedUser)

      await service.create(createUserDto)

      expect(mockedUserRepository.create).toHaveBeenCalledWith(createUserDto)
      expect(mockedOrganizationsService.create).toHaveBeenCalled()
    })
  })

  describe('deactivate', () => {
    it('should deactivate user', async () => {
      const mockedUser = createRandomUser()
      mockedUserRepository.delete.mockResolvedValue(mockedUser)
      await service.deactivate(mockedUser.id)

      expect(mockedUserRepository.delete).toHaveBeenCalledWith(mockedUser.id)
    })
  })
})
