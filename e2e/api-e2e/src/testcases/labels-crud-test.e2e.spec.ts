import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser } from '#utils/helpers'
import { runCrudTestSuite } from '#utils/test-suite.spec'

describe('Labels E2E', () => {
  let app: HttpInstance
  let accessToken: string
  let orgname: string
  const credentials = {
    email: 'chatbots-test@archesai.com',
    password: 'password'
  }

  beforeAll(async () => {
    app = await createApp()
    await app.init()
  })

  afterAll(async () => {
    await app.close()
  })

  beforeEach(async () => {
    accessToken = (await registerUser(app, credentials)).accessToken
    const user = await getUser(app, accessToken)
    orgname = user.orgname
  })

  afterEach(async () => {
    await request(app.getHttpServer())
      .delete(`/user`)
      .set('Authorization', `Bearer ${accessToken}`)
  })

  describe('CRUD', () => {
    runCrudTestSuite<LabelEntity, CreateLabelDto, UpdateLabelDto>(
      {
        async create(accessToken, orgname, createDto) {
          return request(app.getHttpServer())
            .post(`/organizations/${orgname}/labels`)
            .set('Authorization', `Bearer ${accessToken}`)
            .send(createDto)
        },
        async delete(accessToken, orgname, id) {
          return request(app.getHttpServer())
            .delete(`/organizations/${orgname}/labels/${id}`)
            .set('Authorization', `Bearer ${accessToken}`)
        },
        async findMany(accessToken, orgname, query) {
          return request(app.getHttpServer())
            .get(`/organizations/${orgname}/labels`)
            .query(query)
            .set('Authorization', `Bearer ${accessToken}`)
        },
        async findOne(accessToken, orgname, id) {
          return request(app.getHttpServer())
            .get(`/organizations/${orgname}/labels/${id}`)
            .set('Authorization', `Bearer ${accessToken}`)
        },
        async update(accessToken, orgname, id, updateDto) {
          return request(app.getHttpServer())
            .patch(`/organizations/${orgname}/labels/${id}`)
            .set('Authorization', `Bearer ${accessToken}`)
            .send(updateDto)
        }
      },
      {
        create: [
          {
            accessToken,
            createDto: { name: 'ValidName' },
            expectedStatus: 201,
            name: 'valid label name',
            orgname
          },
          {
            accessToken,
            createDto: { name: '' },
            expectedStatus: 400, // You might need custom logic in create() to check this
            name: 'empty label name',
            orgname
          }
        ],
        delete: [
          {
            accessToken,
            createDto: { name: 'ToDelete' },
            expectedStatus: 200,
            name: 'valid delete',
            orgname
          },
          {
            accessToken,
            createDto: { name: 'Temp' },
            expectedStatus: 404,
            name: 'delete nonexistent resource',
            orgname
          }
        ],
        findMany: [
          {
            accessToken,
            createDtos: [{ name: 'ListMe1' }, { name: 'ListMe2' }],
            expectedStatus: 200,
            name: 'default list scenario',
            orgname,
            query: {}
          }
        ],
        findOne: [
          {
            accessToken,
            createDto: { name: 'ToRead' },
            expectedStatus: 200,
            name: 'existing label',
            orgname,
            query: {}
          },
          {
            accessToken,
            createDto: { name: 'temp' },
            expectedStatus: 404,
            name: 'non-existent label',
            orgname,
            query: {}
          }
        ],
        update: [
          {
            accessToken,
            createDto: { name: 'Updatable' },
            expectedStatus: 200,
            name: 'update name successfully',
            orgname,
            updateDto: { name: 'NewName' }
          },
          {
            accessToken,
            createDto: { name: 'Original' },
            expectedStatus: 400,
            name: 'update with invalid data',
            orgname,
            updateDto: { name: '' }
          }
        ]
      }
    )
  })
})
