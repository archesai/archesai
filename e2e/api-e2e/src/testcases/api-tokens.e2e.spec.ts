import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser, setEmailVerified } from '#utils/helpers'

describe('Access Tokens', () => {
  let app: HttpInstance
  let accessToken: string
  let organizationId: string
  const credentials = {
    email: 'api-tokens-test@archesai.com',
    password: 'password'
  }

  beforeAll(async () => {
    app = await createApp()
    await app.init()
    accessToken = (await registerUser(app, credentials)).accessToken
    const user = await getUser(app, accessToken)
    organizationId = user.defaultOrgname
    await setEmailVerified(app, user.id)
  })

  afterAll(async () => {
    await app.close()
  })

  it('Should let users create and delete scoped api tokens', async () => {
    // Create USER token
    const apiToken = await createToken(organizationId, accessToken, 'USER')

    // Verify token exists
    await verifyTokenExists(organizationId, accessToken)

    // Verify USER token role
    await verifyUserRoleTypeEnum(apiToken.key, 'USER')

    // Verify restricted actions with USER token
    await verifyRestrictedActions(organizationId, apiToken.key)

    // Delete token
    await deleteToken(organizationId, accessToken, apiToken.id)

    // Ensure token is no longer valid
    await verifyTokenRevocation(organizationId, apiToken.key)

    // Create a token with restricted domain
    const badUserToken = await createToken(organizationId, accessToken, 'USER')

    // Verify restricted access due to bad domain
    await verifyRestrictedDomainAccess(organizationId, badUserToken.key)

    // Verify allowed access due to good domain
    await verifyAllowedDomainAccess(organizationId, badUserToken.key)

    expect(true).toBe(true)
  })

  const createToken = async (
    organizationId: string,
    accessToken: string,
    role: RoleType
  ): Promise<ApiTokenEntity> => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${organizationId}/api-tokens`)
      .send({ name: `${role}-token`, role })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(201)
    expect(res).toSatisfyApiSpec()
    return res.body
  }

  const verifyTokenExists = async (
    organizationId: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${organizationId}/api-tokens`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(200)
    expect(res.body.results.length).toBe(1)
    expect(res).toSatisfyApiSpec()
    return res.body.results[0].id
  }

  const verifyUserRoleTypeEnum = async (
    accessToken: string,
    expectedRoleTypeEnum: RoleType
  ) => {
    const res = await request(app.getHttpServer())
      .get('/user')
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(200)
    expect(res).toSatisfyApiSpec()
    expect(res.body.memberships.length).toBe(1)
    expect(res.body.memberships[0].role).toBe(expectedRoleTypeEnum)
  }

  const verifyRestrictedActions = async (
    organizationId: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .delete(`/organizations/${organizationId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(403)
  }

  const deleteToken = async (
    organizationId: string,
    accessToken: string,
    tokenId: string
  ) => {
    const res = await request(app.getHttpServer())
      .delete(`/organizations/${organizationId}/api-tokens/${tokenId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(200)
  }

  const verifyTokenRevocation = async (
    organizationId: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .delete(`/organizations/${organizationId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(401)
  }

  const verifyRestrictedDomainAccess = async (
    organizationId: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${organizationId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(401)
  }

  const verifyAllowedDomainAccess = async (
    organizationId: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${organizationId}`)
      .set('Authorization', `Bearer ${accessToken}`)
      .set('Origin', 'localhost')
    expect(res.status).toBe(200)
  }
})
