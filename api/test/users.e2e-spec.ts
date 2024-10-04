import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { createApp } from "./util";

describe("User Onboard", () => {
  let app: INestApplication;
  let token: string;

  const credentials = {
    email: "onboard@archesai.com",
    password: "password",
    username: "admin",
  };

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("should create an internal user on first api call", async () => {
    // Create user
    let res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials);
    expect(res.status).toBe(201);
    token = res.body.apiToken;

    // Check user data
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(200);
    expect(res.body.email).toBe("onboard@archesai.com");
    expect(res.body.defaultOrgname).toBeTruthy();
    const orgname = res.body.defaultOrgname;

    // check organization data
    res = await request(app.getHttpServer())
      .get("/organizations/" + orgname)
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
