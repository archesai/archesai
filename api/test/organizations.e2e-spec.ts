import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { UsersService } from "../src/users/users.service";
import { createApp, deactivateUser, getUser, registerUser } from "./util";

describe("Organizations", () => {
  let app: INestApplication;
  let accessToken: string;

  const credentials = {
    email: "organizations-test@archesai.com",
    password: "password",
    username: "organizations-test",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;

    const usersService = app.get<UsersService>(UsersService);
    await usersService.setEmailVerifiedByEmail(credentials.email);
  });

  afterAll(async () => {
    await app.close();
  });

  it("should add a user to an organization when they create it", async () => {
    // Get user
    const user = await getUser(app, accessToken);

    // Get members of default organization
    const res = await request(app.getHttpServer())
      .get("/organizations/" + user.defaultOrgname + "/members")
      .set("Authorization", "Bearer " + accessToken);
    expect(res.status).toBe(200);
    expect(res.body.metadata.totalResults).toBe(1);

    // Delete the organization
    await request(app.getHttpServer())
      .delete("/organizations/" + user.defaultOrgname)
      .send()
      .set("Authorization", "Bearer " + accessToken);
    expect(res.status).toBe(200);

    // Deactivate the user
    await deactivateUser(app, accessToken);
  });
});
