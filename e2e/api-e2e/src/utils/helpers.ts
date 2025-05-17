import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

// Helper function to register a user and return the API token
export const registerUser = async (
  app: HttpInstance,
  createAccountRequest: {
    email: string
    password: string
  }
): Promise<AccessTokenEntity> => {
  const res = await request(app.getHttpServer())
    .post('/auth/accounts')
    .send(createAccountRequest)
  if (res.status !== 201) {
    throw new Error(`Failed to register user: ${res.status} ${res.text}`)
  }

  return res.body
}

export const setEmailVerified = async (app: HttpInstance, id: string) => {
  const userService = app.get<UsersService>(UsersService)
  await userService.setEmailVerified(id)
}

// Helper function to get user data
export const getUser = async (
  app: HttpInstance,
  accessToken: string
): Promise<UserEntity> => {
  const res = await request(app.getHttpServer())
    .get('/user')
    .set('Authorization', `Bearer ${accessToken}`)
  expect(res.status).toBe(200)
  expect(res.body.defaultOrgname).toBeTruthy()
  expect(res).toSatisfyApiSpec()
  return res.body
}

// Helper function to check organization data
export const getOrganization = async (
  app: HttpInstance,
  orgname: string,
  accessToken: string
): Promise<OrganizationEntity> => {
  const res = await request(app.getHttpServer())
    .get(`/organizations/${orgname}`)
    .set('Authorization', `Bearer ${accessToken}`)
  expect(res.status).toBe(200)
  expect(res).toSatisfyApiSpec()
  return res.body
}

// Helper function to deactivate a user
export const deactivateUser = async (
  app: HttpInstance,
  accessToken: string
) => {
  const res = await request(app.getHttpServer())
    .post('/user/deactivate')
    .set('Authorization', `Bearer ${accessToken}`)
  expect(res.status).toBe(201)
}
