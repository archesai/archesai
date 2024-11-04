import { INestApplication } from "@nestjs/common";

import {
  createApp,
  deactivateUser,
  getOrganization,
  getUser,
  registerUser,
  setEmailVerifiedByEmail,
} from "./util";

describe("Users", () => {
  let app: INestApplication;
  let accessToken: string;

  const credentials = {
    email: "users-test-admin@archesai.com",
    password: "password",
    username: "users-test-admin",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;
  });

  afterAll(async () => {
    await app.close();
  });

  it("should create an internal user on first API call", async () => {
    // Verify user data
    const user = await getUser(app, accessToken);
    await setEmailVerifiedByEmail(app, user.email);

    // Verify organization data
    await getOrganization(app, user.defaultOrgname, accessToken);

    // Deactivate the user
    await deactivateUser(app, accessToken);
  });
});
