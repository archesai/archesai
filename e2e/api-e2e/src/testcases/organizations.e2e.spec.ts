import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { deactivateUser, getUser, registerUser } from '#utils/helpers'

describe('Organizations', () => {
  let app: HttpInstance
  let accessToken: string

  const credentials = {
    email: 'organizations-test@archesai.com',
    password: 'password'
  }

  beforeAll(async () => {
    app = await createApp()
    await app.init()

    accessToken = (await registerUser(app, credentials)).accessToken

    const user = await getUser(app, accessToken)
  })

  afterAll(async () => {
    await app.close()
  })

  it('should add a user to an organization when they create it', async () => {
    // Get user
    const user = await getUser(app, accessToken)

    // Get members of default organization
    const res = await request(app.getHttpServer())
      .get('/organizations/' + user.organizationId + '/members')
      .set('Authorization', 'Bearer ' + accessToken)
    expect(res.status).toBe(200)
    expect(res.body.metadata.totalResults).toBe(1)

    // Delete the organization
    await request(app.getHttpServer())
      .delete('/organizations/' + user.organizationId)
      .send()
      .set('Authorization', 'Bearer ' + accessToken)
    expect(res.status).toBe(200)

    // Deactivate the user
    await deactivateUser(app, accessToken)
  })
})
