import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { createApp } from "./util";

describe("Auth", () => {
  let app: INestApplication;

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("Should protect private endpoints like /user", async () => {
    // Try to get /auth/me endpoint and fail
    let res = await request(app.getHttpServer()).get("/user");
    expect(res.status).toBe(401);

    // Create account
    const credentials = {
      email: "admin@archesai.com",
      password: "password",
      username: "admin",
    };

    res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials);
    expect(res.status).toBe(201);
    const token = res.body.apiToken;

    // Get /auth/me endpoint and pass
    res = await request(app.getHttpServer())
      .get("/user")
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
