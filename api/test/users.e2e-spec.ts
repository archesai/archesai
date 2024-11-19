import { INestApplication } from "@nestjs/common";

import {
  createApp,
  deactivateUser,
  getOrganization,
  getUser,
  registerUser,
  setEmailVerifiedByEmail,
  testBaseControllerEndpoints,
} from "./util";

describe("Users", () => {
  let app: INestApplication;
  let accessToken: string;
  const baseRoute = "/user";

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

  testBaseControllerEndpoints(() => app, baseRoute, accessToken, {
    createCases: [],
    findAllCases: [
      {
        expectedResponse: [
          {
            email: "users-test-admin@example.com",
            id: expect.any(String),
            username: "testuser",
          },
        ],
        expectedStatus: 200,
        name: "should return all users",
      },
    ],
    findOneCases: [
      {
        expectedResponse: {
          email: "test@example.com",
          id: "user-id",
          username: "testuser",
        },
        expectedStatus: 200,
        id: "user-id",
        name: "should return a user by ID",
      },
    ],
    removeCases: [
      {
        expectedResponse: {},
        expectedStatus: 200,
        id: "user-id",
        name: "should remove a user",
      },
    ],
    updateCases: [
      {
        dto: { email: "updated@example.com" },
        expectedResponse: {
          email: "updated@example.com",
          id: "user-id",
          username: "testuser",
        },
        expectedStatus: 200,
        id: "user-id", // Replace with a valid ID
        name: "should update a user",
      },
      // Add more update test cases
    ],
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
