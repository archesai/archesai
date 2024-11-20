import { INestApplication } from "@nestjs/common";

import {
  createApp,
  deactivateUser,
  getOrganization,
  getUser,
  registerUser,
  setEmailVerified,
} from "./util";

describe("Users", () => {
  let app: INestApplication;
  let accessToken: string;

  const credentials = {
    email: "users-test-admin@archesai.com",
    password: "password",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;
  });

  afterAll(async () => {
    await app.close();
  });

  it("should be defined", () => {
    expect(app).toBeDefined();
  });

  it("should create an internal user on first API call", async () => {
    // Verify user data
    const user = await getUser(app, accessToken);
    await setEmailVerified(app, user.id);

    // Verify organization data
    await getOrganization(app, user.defaultOrgname, accessToken);

    // Deactivate the user
    await deactivateUser(app, accessToken);
  });
});
