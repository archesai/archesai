import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { OrganizationsService } from "../src/organizations/organizations.service";
import { UsersService } from "../src/users/users.service";
import { createApp, getUser, registerUser } from "./util";

describe("Members", () => {
  let app: INestApplication;
  let usersService: UsersService;
  let organizationsService: OrganizationsService;
  let accessToken: string;
  let orgname: string;

  const credentials = {
    email: "admin@archesai.com",
    password: "password",
  };

  const invitedUser = {
    email: "invitedUser@archesai.com",
    password: "password2",
  };

  const uninvitedUser = {
    email: "uninvitedUser@archesai.com",
    password: "password",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    usersService = app.get<UsersService>(UsersService);
    organizationsService = app.get<OrganizationsService>(OrganizationsService);

    accessToken = (await registerUser(app, credentials)).accessToken;

    orgname = (await getUser(app, accessToken)).defaultOrgname;

    await usersService.setEmailVerified(credentials.email);
    await organizationsService.setPlan(orgname, "UNLIMITED");
  });

  afterAll(async () => {
    await app.close();
  });

  it("Should allow admins to add members", async () => {
    // Invite user with invalid role
    await inviteUser(invitedUser.email, "BADROLE", 400);

    // Invite user with valid role
    await inviteUser(invitedUser.email, "ADMIN", 201);

    // Register uninvited user
    const uninvitedRes = await registerUser(app, uninvitedUser);
    const uninvitedUserToken = uninvitedRes.accessToken;

    // Register invited user
    const invitedRes = await registerUser(app, invitedUser);
    const invitedUserToken = invitedRes.accessToken;

    // Attempt to join with various scenarios
    await joinOrganization(invitedUserToken, 403); // Not verified
    await joinOrganization(uninvitedUserToken, 403); // Not verified

    await usersService.setEmailVerified(uninvitedUser.email);
    await joinOrganization(uninvitedUserToken, 404); // Uninvited

    await usersService.setEmailVerified(invitedUser.email);
    await joinOrganization(invitedUserToken, 201); // Verified and invited

    // Verify invited user added as member
    const memberRes = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .set("Authorization", "Bearer " + invitedUserToken);
    expect(memberRes.status).toBe(200);
    expect(memberRes.body.metadata.totalResults).toBe(2);

    // Verify uninvited user not a member
    const nonMemberRes = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .set("Authorization", "Bearer " + uninvitedUserToken);
    expect(nonMemberRes.status).toBe(404);

    // Cleanup: delete users and organization
    await cleanupOrganizationAndUsers();
  });

  const inviteUser = async (email, role, expectedStatus) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .send({ inviteEmail: email, role })
      .set("Authorization", "Bearer " + accessToken);
    if (res.status != 400) {
      expect(res).toSatisfyApiSpec();
    }
    expect(res.status).toBe(expectedStatus);
  };

  const joinOrganization = async (userToken, expectedStatus) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members/join`)
      .send({})
      .set("Authorization", "Bearer " + userToken);
    expect(res.status).toBe(expectedStatus);
  };

  const cleanupOrganizationAndUsers = async () => {
    await request(app.getHttpServer())
      .delete(`/organizations/${orgname}`)
      .set("Authorization", "Bearer " + accessToken);

    await deactivateUser(credentials);
    await deactivateUser(invitedUser);
    await deactivateUser(uninvitedUser);
  };

  const deactivateUser = (user) =>
    request(app.getHttpServer()).post("/auth/deactivate").send(user);
});
