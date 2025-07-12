// import type { DeepMocked } from '@golevelup/ts-jest'
// import type { TestingModule } from '@nestjs/testing'

// import { createMock } from '@golevelup/ts-jest'
// import { Test } from '@nestjs/testing'

// import type { ArchesApiRequest, ExecutionContext } from '@archesai/core'

// import { ConfigModule, Logger } from '@archesai/core'

// import { ApiTokensController } from '#api-tokens/api-tokens.controller'
// import { ApiTokensService } from '#api-tokens/api-tokens.service'
// import { createRandomApiToken } from '#api-tokens/factories/api-token.factory'
// import { AuthenticatedGuard } from '#auth/guards/authenticated.guard'

// describe('ApiTokensController', () => {
//   // eslint-disable-next-line @typescript-eslint/no-explicit-any
//   let app: any
//   let mockedApiTokensService: DeepMocked<ApiTokensService>
//   let organizationId: string
//   let username: string

//   beforeAll(async () => {
//     const moduleRef: TestingModule = await Test.createTestingModule({
//       controllers: [ApiTokensController],
//       imports: [ConfigModule],
//       providers: [
//         {
//           provide: ApiTokensService,
//           useValue: createMock<ApiTokensService>()
//         }
//       ]
//     })
//       .overrideGuard(AuthenticatedGuard)
//       .useValue({
//         canActivate(ctx: ExecutionContext) {
//           const request = ctx.switchToHttp().getRequest<ArchesApiRequest>()
//           request.user = mockUserEntity
//           return true
//         }
//       })
//       .compile()
//     app = moduleRef.createNestApplication()
//     app.useLogger(app.get(Logger))

//     await app.init()

//     mockedApiTokensService = moduleRef.get(ApiTokensService)

//     const mockUserEntity = createRandomUser()
//     organizationId = mockUserEntity.organizationId
//     username = mockUserEntity.id
//   })

//   afterAll(async () => {
//     await app.close()
//   })

//   it('should be defined', () => {
//     expect(app.get(ApiTokensController)).toBeDefined()
//     expect(mockedApiTokensService).toBeDefined()
//   })

//   it('POST /organizations/:organizationId/api-tokens should call service.create', async () => {
//     const createApiTokenDto = {
//       name: 'testToken',
//       role: 'ADMIN' as const
//     }
//     const mockedApiToken = createRandomApiToken(createApiTokenDto)
//     mockedApiTokensService.create.mockResolvedValue(mockedApiToken)

//     const response = await request(app.getHttpServer())
//       .post(`/organizations/${organizationId}/api-tokens`)
//       .send(createApiTokenDto)

//     expect(response.status).toBe(201)
//     expect(response.body).toEqual({
//       ...mockedApiToken,
//       createdAt: mockedApiToken.createdAt,
//       updatedAt: undefined
//     })

//     expect(mockedApiTokensService.create).toHaveBeenCalledWith({
//       domains: '*',
//       name: 'testToken',
//       organizationId,
//       role: 'ADMIN' as const,
//       username
//     })
//   })

//   it('GET /organizations/:organizationId/api-tokens should call service.findMany', async () => {
//     const mockedApiToken = createRandomApiToken()
//     const mockedPaginatedApiTokens = {
//       metadata: {
//         limit: 10,
//         offset: 0,
//         totalResults: 1
//       },
//       results: [mockedApiToken]
//     }

//     const response = await request(app.getHttpServer())
//       .get(`/organizations/${organizationId}/api-tokens`)
//       .query({})

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({
//       ...mockedPaginatedApiTokens,
//       results: [
//         {
//           ...mockedApiToken,
//           createdAt: mockedApiToken.createdAt,
//           updatedAt: undefined
//         }
//       ]
//     })
//     expect(mockedApiTokensService.findMany).toHaveBeenCalledWith({
//       endDate: undefined,
//       filters: [
//         {
//           field: 'organizationId',
//           operator: 'equals',
//           value: organizationId
//         }
//       ],
//       limit: 10,
//       offset: 0,
//       sortBy: 'createdAt',
//       sortDirection: 'desc',
//       startDate: undefined
//     })
//   })

//   it('GET /organizations/:organizationId/api-tokens/:id should call service.findOne', async () => {
//     const mockedApiToken = createRandomApiToken()
//     mockedApiTokensService.findOne.mockResolvedValue(mockedApiToken)

//     const response = await request(app.getHttpServer()).get(
//       `/organizations/${organizationId}/api-tokens/1`
//     )

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({
//       ...mockedApiToken,
//       createdAt: mockedApiToken.createdAt,
//       updatedAt: undefined
//     })
//     expect(mockedApiTokensService.findOne).toHaveBeenCalledWith('1')
//   })

//   it('PATCH /organizations/:organizationId/api-tokens/:id should call service.update', async () => {
//     const mockedApiToken = createRandomApiToken()
//     mockedApiTokensService.update.mockResolvedValue(mockedApiToken)

//     const response = await request(app.getHttpServer())
//       .patch(`/organizations/${organizationId}/api-tokens/1`)
//       .send({ name: 'updatedToken' })
//       .set('Authorization', 'Bearer token')

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({
//       ...mockedApiToken,
//       createdAt: mockedApiToken.createdAt,
//       updatedAt: undefined
//     })
//     expect(mockedApiTokensService.update).toHaveBeenCalledWith('1', {
//       name: 'updatedToken'
//     })
//   })

//   it('DELETE /organizations/:organizationId/api-tokens/:id should call service.remove', async () => {
//     const response = await request(app.getHttpServer())
//       .delete(`/organizations/${organizationId}/api-tokens/1`)
//       .set('Authorization', 'Bearer token')

//     expect(response.status).toBe(200)
//     expect(response.body).toEqual({})
//     expect(mockedApiTokensService.delete).toHaveBeenCalledWith('1')
//   })
// })
