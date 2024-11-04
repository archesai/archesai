import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { createApp, deactivateUser, getUser, registerUser } from "./util";

describe("Auth", () => {
  let app: INestApplication;

  beforeAll(async () => {
    app = await createApp();
    await app.init();
  });

  it("Should protect private endpoints like /user", async () => {
    // Try to get /auth/me endpoint and fail
    const res = await request(app.getHttpServer()).get("/user");
    expect(res.status).toBe(401);

    // Create account
    const credentials = {
      email: "auth-test@archesai.com",
      password: "password",
      username: "auth-test",
    };

    const tokenDto = await registerUser(app, credentials);
    const accessToken = tokenDto.accessToken;

    await getUser(app, accessToken);

    await deactivateUser(app, accessToken);
  });

  afterAll(async () => {
    await app.close();
  });
});
