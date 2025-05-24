// import type { DeepMocked } from '@golevelup/ts-jest'
// import type { TestingModule } from '@nestjs/testing'

// import { createMock } from '@golevelup/ts-jest'
// import { Test } from '@nestjs/testing'

// import type { CreateMemberRequest } from '#members/dto/create-member.req.dto'

// import { createRandomMember } from '#members/factories/member.factory'
// import { MembersController } from '#members/members.controller'
// import { MembersService } from '#members/members.service'
// import { createRandomUser } from '#users/factories/user.factory'

// describe('MembersController', () => {
//   // eslint-disable-next-line @typescript-eslint/no-explicit-any
//   let app: any
//   let mockedMembersService: DeepMocked<MembersService>

//   beforeAll(async () => {
//     const moduleRef: TestingModule = await Test.createTestingModule({
//       controllers: [MembersController],
//       providers: [
//         {
//           provide: 'APP_GUARD',
//           useValue: createMock({
//             canActivate: jest.fn().mockImplementation((context) => {
//               const request = context.switchToHttp().getRequest()
//               request.user = createRandomUser({ id: 'testUser' })
//               return true
//             })
//           })
//         },
//         {
//           provide: MembersService,
//           useValue: createMock<MembersService>({
//             create: jest.fn(),
//             delete: jest.fn(),
//             findMany: jest.fn(),
//             findOne: jest.fn(),
//             update: jest.fn()
//           })
//         }
//       ]
//     }).compile()

//     app = moduleRef.createNestApplication()
//     await app.init()

//     mockedMembersService = moduleRef.get(MembersService)
//   })

//   afterAll(async () => {
//     await app.close()
//   })

//   it('should be defined', () => {
//     expect(app.get(MembersController)).toBeDefined()
//     expect(mockedMembersService).toBeDefined()
//   })

//   it('POST /organizations/:orgname/members should validate role', async () => {
//     const orgname = 'testOrg'
//     const createMemberRequest = {
//       inviteEmail: 'jonathan@gmail.com',
//       role: 'BADROLE'
//     }
//     const response = await request(app.getHttpServer())
//       .post(`/organizations/${orgname}/members`)
//       .set('Content-Type', 'application/json')
//       .send(createMemberRequest)

//     expect(response.status).toBe(400)
//   })

//   it('POST /organizations/:orgname/members should call service.create', async () => {
//     const orgname = 'testOrg'
//     const createMemberRequest: CreateMemberRequest = {
//       name: 'Jonathan',
//       role: 'ADMIN'
//     }
//     const mockedMember = createRandomMember(createMemberRequest)
//     mockedMembersService.create.mockResolvedValue(mockedMember)

//     const response = await request(app.getHttpServer())
//       .post(`/organizations/${orgname}/members`)
//       .send(createMemberRequest)

//     expect(response.status).toBe(201)
//     expect(response.body).toEqual({
//       ...mockedMember,
//       createdAt: mockedMember.createdAt,
//       updatedAt: undefined
//     })

//     expect(mockedMembersService.create).toHaveBeenCalledWith(
//       createMemberRequest
//     )
//   })

//   it('GET /organizations/:orgname/members should call service.findMany', async () => {
//     const orgname = 'testOrg'
//     const mockedMember = createRandomMember()
//     const mockedPaginatedMembers = {
//       count: 100,
//       data: [mockedMember]
//     }
//     mockedMembersService.findMany.mockResolvedValue(mockedPaginatedMembers)

//     const response = await request(app.getHttpServer())
//       .get(`/organizations/${orgname}/members`)
//       .query({})

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({
//       ...mockedPaginatedMembers,
//       data: [
//         {
//           ...mockedMember,
//           createdAt: mockedMember.createdAt,
//           updatedAt: undefined
//         }
//       ]
//     })
//     expect(mockedMembersService.findMany).toHaveBeenCalledWith({
//       endDate: undefined,
//       filters: [
//         {
//           field: 'orgname',
//           operator: 'equals',
//           value: orgname
//         }
//       ],
//       limit: 10,
//       offset: 0,
//       sortBy: 'createdAt',
//       sortDirection: 'desc',
//       startDate: undefined
//     })
//   })

//   it('GET /organizations/:orgname/members/:id should call service.findOne', async () => {
//     const orgname = 'testOrg'
//     const mockedMember = createRandomMember()
//     mockedMembersService.findOne.mockResolvedValue(mockedMember)

//     const response = await request(app.getHttpServer()).get(
//       `/organizations/${orgname}/members/1`
//     )

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({
//       ...mockedMember,
//       createdAt: mockedMember.createdAt,
//       updatedAt: undefined
//     })
//     expect(mockedMembersService.findOne).toHaveBeenCalledWith('1')
//   })

//   it('PATCH /organizations/:orgname/members/:id should call service.update', async () => {
//     const orgname = 'testOrg'
//     const mockedMember = createRandomMember({
//       orgname
//     })
//     mockedMembersService.update.mockResolvedValue(mockedMember)

//     const response = await request(app.getHttpServer())
//       .patch(`/organizations/${orgname}/members/1`)
//       .send({ role: 'ADMIN' })
//       .set('Authorization', 'Bearer token')

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({
//       ...mockedMember,
//       createdAt: mockedMember.createdAt,
//       updatedAt: undefined
//     })
//     expect(mockedMembersService.update).toHaveBeenCalledWith('1', {
//       role: 'ADMIN'
//     })
//   })

//   it('DELETE /organizations/:orgname/members/:id should call service.remove', async () => {
//     const response = await request(app.getHttpServer())
//       .delete('/organizations/testOrg/members/1')
//       .set('Authorization', 'Bearer token')

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({})
//     expect(mockedMembersService.delete).toHaveBeenCalledWith('1')
//   })
// })
