import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { UsersService } from "../src/users/users.service";
import { createApp } from "./util";

describe("Organizations", () => {
  let app: INestApplication;
  let token: string;
  const credentials = {
    email: "organizations@archesai.com",
    password: "password",
    username: "organizations",
  };

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("should add a user to an organization when they create it", async () => {
    // Create user
    let res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials);
    expect(res.status).toBe(201);
    token = res.body.apiToken;

    // Verify e-mail
    const usersService = app.get<UsersService>(UsersService);
    await usersService.setEmailVerifiedByEmail(credentials.email);

    // Get organization
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(200);
    expect(res.body.defaultOrgname).toBeTruthy();
    const orgname = res.body.defaultOrgname;

    // Get organization members
    res = await request(app.getHttpServer())
      .get("/organizations/" + orgname + "/members")
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(200);
    expect(res.body.metadata.totalResults).toBe(1);

    res = await request(app.getHttpServer())
      .delete("/organizations/" + orgname)
      .send()
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(200);

    // Delete user
    res = await request(app.getHttpServer())
      .post("/auth/deactivate")
      .send(credentials);
    expect(res.status).toBe(201);
  });

  afterEach(async () => {
    await app.close();
  });
});
