import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser } from '#utils/helpers'

describe('Members', () => {
  let app: HttpInstance
  let accessToken: string
  let orgname: string

  const credentials = {
    email: 'admin@archesai.com',
    password: 'password'
  }

  const invitedUser = {
    email: 'invitedUser@archesai.com',
    password: 'password2'
  }

  const uninvitedUser = {
    email: 'uninvitedUser@archesai.com',
    password: 'password'
  }

  beforeAll(async () => {
    app = await createApp()
    await app.init()

    usersService = app.get<UsersService>(UsersService)
    organizationsService = app.get<OrganizationsService>(OrganizationsService)

    accessToken = (await registerUser(app, credentials)).accessToken

    orgname = (await getUser(app, accessToken)).defaultOrgname

    const userEntity = await getUser(app, accessToken)
    await usersService.setEmailVerified(userEntity.id)
    await organizationsService.setPlan(orgname, PlanTypeEnum.UNLIMITED)
  })

  afterAll(async () => {
    await app.close()
  })

  it('Should allow admins to add members', async () => {
    // Invite user with invalid role
    await inviteUser(invitedUser.email, 'BADROLE', 400)

    // Invite user with valid role
    await inviteUser(invitedUser.email, 'ADMIN', 201)

    // Register uninvited user
    const uninvitedRes = await registerUser(app, uninvitedUser)
    const uninvitedUserToken = uninvitedRes.accessToken

    // Register invited user
    const invitedRes = await registerUser(app, invitedUser)
    const invitedUserToken = invitedRes.accessToken

    // Attempt to join with various scenarios
    await joinOrganization(invitedUserToken, 403) // Not verified
    await joinOrganization(uninvitedUserToken, 403) // Not verified

    const uninvitedUserEntity = await getUser(app, uninvitedUserToken)
    await usersService.setEmailVerified(uninvitedUserEntity.id)
    await joinOrganization(uninvitedUserToken, 404) // Uninvited

    const invitedUserEntity = await getUser(app, invitedUserToken)
    await usersService.setEmailVerified(invitedUserEntity.id)
    await joinOrganization(invitedUserToken, 201) // Verified and invited

    // Verify invited user added as member
    const memberRes = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .set('Authorization', 'Bearer ' + invitedUserToken)
    expect(memberRes.status).toBe(200)
    expect(memberRes.body.metadata.totalResults).toBe(2)

    // Verify uninvited user not a member
    const nonMemberRes = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .set('Authorization', 'Bearer ' + uninvitedUserToken)
    expect(nonMemberRes.status).toBe(404)

    // Cleanup: delete users and organization
    await cleanupOrganizationAndUsers()
  })

  const inviteUser = async (
    email: string,
    role: string,
    expectedStatus: number
  ) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .send({ inviteEmail: email, role })
      .set('Authorization', 'Bearer ' + accessToken)
    if (res.status != 400) {
      expect(res).toSatisfyApiSpec()
    }
    expect(res.status).toBe(expectedStatus)
  }

  const joinOrganization = async (
    userToken: string,
    expectedStatus: number
  ) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members/join`)
      .send({})
      .set('Authorization', 'Bearer ' + userToken)
    expect(res.status).toBe(expectedStatus)
  }

  const cleanupOrganizationAndUsers = async () => {
    await request(app.getHttpServer())
      .delete(`/organizations/${orgname}`)
      .set('Authorization', 'Bearer ' + accessToken)

    await deactivateUser(credentials)
    await deactivateUser(invitedUser)
    await deactivateUser(uninvitedUser)
  }

  const deactivateUser = (user: { email: string; password: string }) =>
    request(app.getHttpServer()).post('/auth/deactivate').send(user)
})
