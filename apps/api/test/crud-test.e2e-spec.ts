import { INestApplication } from '@nestjs/common'
import request from 'supertest'
import { runCrudTestSuite } from './crud-test-suite'
import { LabelEntity } from '@/src/labels/entities/label.entity'
import { CreateLabelDto } from '@/src/labels/dto/create-label.dto'
import { UpdateLabelDto } from '@/src/labels/dto/update-label.dto'
import { createApp, getUser, registerUser } from './util'

describe('Labels E2E', () => {
  let app: INestApplication
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
    orgname = user.defaultOrgname
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
        },
        async delete(accessToken, orgname, id) {
          return request(app.getHttpServer())
            .delete(`/organizations/${orgname}/labels/${id}`)
            .set('Authorization', `Bearer ${accessToken}`)
        },
        async findAll(accessToken, orgname, searchQueryDto) {
          return request(app.getHttpServer())
            .get(`/organizations/${orgname}/labels`)
            .query(searchQueryDto)
            .set('Authorization', `Bearer ${accessToken}`)
        }
      },
      {
        create: [
          {
            name: 'valid label name',
            accessToken,
            orgname,
            createDto: { name: 'ValidName' },
            expectedStatus: 201
          },
          {
            name: 'empty label name',
            accessToken,
            orgname,
            createDto: { name: '' },
            expectedStatus: 400 // You might need custom logic in create() to check this
          }
        ],
        findOne: [
          {
            name: 'existing label',
            accessToken,
            orgname,
            createDto: { name: 'ToRead' },
            searchQueryDto: {},
            expectedStatus: 200
          },
          {
            name: 'non-existent label',
            accessToken,
            orgname,
            createDto: { name: 'temp' },
            searchQueryDto: {},
            expectedStatus: 404
          }
        ],
        update: [
          {
            name: 'update name successfully',
            accessToken,
            orgname,
            createDto: { name: 'Updatable' },
            updateDto: { name: 'NewName' },
            expectedStatus: 200
          },
          {
            name: 'update with invalid data',
            accessToken,
            orgname,
            createDto: { name: 'Original' },
            updateDto: { name: '' },
            expectedStatus: 400
          }
        ],
        findAll: [
          {
            name: 'default list scenario',
            accessToken,
            orgname,
            createDtos: [{ name: 'ListMe1' }, { name: 'ListMe2' }],
            searchQueryDto: {},
            expectedStatus: 200
          }
        ],
        delete: [
          {
            name: 'valid delete',
            accessToken,
            orgname,
            createDto: { name: 'ToDelete' },
            expectedStatus: 200
          },
          {
            name: 'delete nonexistent resource',
            accessToken,
            orgname,
            createDto: { name: 'Temp' },
            expectedStatus: 404
          }
        ]
      }
    )
  })
})
