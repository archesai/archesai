// import type { DeepMocked } from '@golevelup/ts-jest'
// import type { TestingModule } from '@nestjs/testing'

// import { createMock } from '@golevelup/ts-jest'
// import { Test } from '@nestjs/testing'

// import type { UpdateUserRequest } from '#users/dto/update-user.req.dto'

// import { createRandomUser } from '#users/factories/user.factory'
// import { UsersController } from '#users/users.controller'
// import { UsersService } from '#users/users.service'

// describe('UsersController', () => {
//   // eslint-disable-next-line @typescript-eslint/no-explicit-any
//   let app: any
//   let mockedUsersService: DeepMocked<UsersService>

//   beforeAll(async () => {
//     const moduleRef: TestingModule = await Test.createTestingModule({
//       controllers: [UsersController],
//       providers: [
//         {
//           provide: UsersService,
//           useValue: createMock<UsersService>({
//             create: jest.fn(),
//             deactivate: jest.fn(),
//             findOneByEmail: jest.fn()
//           })
//         }
//       ]
//     }).compile()

//     app = moduleRef.createNestApplication()
//     await app.init()

//     mockedUsersService = moduleRef.get(UsersService)
//   })

//   afterAll(async () => {
//     await app.close()
//   })

//   it('should be defined', () => {
//     expect(app.get(UsersController)).toBeDefined()
//     expect(mockedUsersService).toBeDefined()
//   })

//   describe('POST /user/deactivate', () => {
//     it('should deactivate a user', async () => {
//       const response = await request(app.getHttpServer())
//         .post('/user/deactivate')
//         .send()
//         .expect(201)

//       expect(response.body).toEqual({})
//       expect(mockedUsersService.deactivate).toHaveBeenCalledWith('test-id')
//     })
//   })

//   describe('GET /user', () => {
//     it('should return the current user', async () => {
//       const response = await request(app.getHttpServer())
//         .get('/user')
//         .expect(200)

//       expect(response.body.id).toEqual('test-id')
//     })
//   })

//   describe('PATCH /user', () => {
//     it('should update a user', async () => {
//       const updateUserDto: UpdateUserRequest = {
//         name: 'John'
//       }
//       const mockedUser = createRandomUser({
//         id: 'test-id',
//         organizationId: 'test-org'
//       })
//       mockedUsersService.update.mockResolvedValue(mockedUser)

//       const response = await request(app.getHttpServer())
//         .patch('/user')
//         .send(updateUserDto)
//         .expect(200)

//       expect(response.body.firstName).toEqual('John')
//       expect(response.body.lastName).toEqual('Doe')
//       expect(response.body.id).toEqual('test-id')

//       expect(mockedUsersService.update).toHaveBeenCalledWith(
//         mockedUser.id,
//         updateUserDto
//       )
//     })
//   })
// })
