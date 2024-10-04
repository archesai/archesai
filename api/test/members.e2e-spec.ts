import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { OrganizationsService } from "../src/organizations/organizations.service";
import { UsersService } from "../src/users/users.service";
import { createApp } from "./util";

describe("Members", () => {
  let app: INestApplication;
  let token: string;

  const credentials = {
    email: "member-test-admin@archesai.com",
    password: "password",
    username: "member-test-admin",
  };

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("Should allow admins to add members", async () => {
    // Create user
    let res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials);
    expect(res.status).toBe(201);
    token = res.body.apiToken;

    // Verify email
    const usersService = app.get<UsersService>(UsersService);
    const organizationsService =
      app.get<OrganizationsService>(OrganizationsService);

    await usersService.setEmailVerifiedByEmail(credentials.email);

    // Get organization name
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(200);
    const orgname = res.body.defaultOrgname;

    // Invite user without API plan
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .send({
        inviteEmail: "jonathankingfc@archesai.com",
        role: "BADROLE",
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(403);
    await organizationsService.setPlan(orgname, "API");

    // Invite user with bad role
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .send({
        inviteEmail: "jonathankingfc@archesai.com",
        role: "BADROLE",
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(400);

    // Invite user
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .send({
        firstName: "jonathan",
        inviteEmail: "jonathankingfc@archesai.com",
        role: "ADMIN",
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(201);

    // Create user
    res = await request(app.getHttpServer()).post("/auth/register").send({
      email: "jonathankingfc@archesai.com",
      password: "password",
      username: "jonathankingfc",
    });
    expect(res.status).toBe(201);
    const invitedUserToken = res.body.apiToken;

    // Create user
    res = await request(app.getHttpServer()).post("/auth/register").send({
      email: "uninviteduser@archesai.com",
      password: "password",
      username: "uninviteduser",
    });
    expect(res.status).toBe(201);
    const uninvitedUserToken = res.body.apiToken;

    // Try to join as invited but not e-mail verified and fail
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members/join`)
      .send({})
      .set("Authorization", "Bearer " + invitedUserToken);
    expect(res.status).toBe(403);

    // Try to join as uninvited and not e-mail verified and fail
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members/join`)
      .send({})
      .set("Authorization", "Bearer " + uninvitedUserToken);
    expect(res.status).toBe(403);

    // // Verify email
    await usersService.setEmailVerifiedByEmail("uninviteduser@archesai.com");

    // Try to join as uninvited fail
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members/join`)
      .send({})
      .set("Authorization", "Bearer " + uninvitedUserToken);
    expect(res.status).toBe(404);

    // Try to join as invited as verified and succeed
    await usersService.setEmailVerifiedByEmail("jonathankingfc@archesai.com");
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members/join`)
      .send({})
      .set("Authorization", "Bearer " + invitedUserToken);
    expect(res.status).toBe(201);

    // Check and see added as member
    res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .set("Authorization", "Bearer " + invitedUserToken);
    expect(res.status).toBe(200);
    expect(res.body.metadata.totalResults).toBe(2);

    // Uninvited user should get 404 since not a member
    res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .set("Authorization", "Bearer " + uninvitedUserToken);
    expect(res.status).toBe(404);

    // Delete users
    res = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}`)
      .set("Authorization", "Bearer " + token);

    res = await request(app.getHttpServer())
      .post("/auth/deactivate")
      .send(credentials);
    expect(res.status).toBe(201);
    res = await request(app.getHttpServer())
      .post("/auth/deactivate")
      .send({ email: "jonathankingfc@archesai.com", password: "password" });
    expect(res.status).toBe(201);
    res = await request(app.getHttpServer())
      .post("/auth/deactivate")
      .send({ email: "uninviteduser@archesai.com", password: "password" });
    expect(res.status).toBe(201);
  });

  afterEach(async () => {
    await app.close();
  });
});
