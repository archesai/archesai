import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser } from '#utils/helpers'

describe('Labels', () => {
  let app: HttpInstance
  let accessToken: string
  let organizationId: string
  let labelId: string

  const credentials = {
    email: 'chatbots-test@archesai.com',
    password: 'password'
  }

  beforeAll(async () => {
    app = await createApp()
    await app.init()
    accessToken = (await registerUser(app, credentials)).accessToken
    const usersService = app.get<UsersService>(UsersService)
    const user = await getUser(app, accessToken)
    organizationId = user.defaultOrgname
    await usersService.setEmailVerified(user.id)
  })

  afterAll(async () => {
    await app.close()
  })

  it('CREATE - should be able to create a label', async () => {
    const labelRes = await request(app.getHttpServer())
      .post(`/organizations/${organizationId}/labels`)
      .send({
        name: 'name'
      })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(labelRes).toMatchObject({
      status: 201
    })
    expect(labelRes.status).toBe(201)
    expect(labelRes).toSatisfyApiSpec()
    labelId = labelRes.body.id
  })

  it('GET - should be able to get a label', async () => {
    const getRes = await request(app.getHttpServer())
      .get(`/organizations/${organizationId}/labels/${labelId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(getRes.status).toBe(200)
    expect(getRes).toSatisfyApiSpec()
  })

  it('UPDATE - should be able to update a label', async () => {
    const updateRes = await request(app.getHttpServer())
      .patch(`/organizations/${organizationId}/labels/${labelId}`)
      .send({
        name: 'new-label-name'
      })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(updateRes.status).toBe(200)
    expect(updateRes).toSatisfyApiSpec()
    expect(updateRes.body.name).toEqual('new-label-name')
  })

  it('FIND_MANY - should be able to get all labels', async () => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${organizationId}/labels`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(200)
    expect(res).toSatisfyApiSpec()
    const labels = res.body
    expect(labels.results.length).toBeGreaterThan(0)
  })

  it('DELETE - should be able to delete a label', async () => {
    // Delete the chatbot
    const deleteRes = await request(app.getHttpServer())
      .delete(`/organizations/${organizationId}/labels/${labelId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(deleteRes.status).toBe(200)
  })
})
