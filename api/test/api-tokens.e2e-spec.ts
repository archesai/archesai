import { ApiTokenEntity } from "@/src/api-tokens/entities/api-token.entity";
import { INestApplication } from "@nestjs/common";
import { RoleType } from "@prisma/client";
import request from "supertest";

import {
  createApp,
  getUser,
  registerUser,
  setEmailVerifiedByEmail,
} from "./util";

describe("Access Tokens", () => {
  let app: INestApplication;
  let accessToken: string;
  let orgname: string;
  const credentials = {
    email: "api-tokens-test@archesai.com",
    password: "password",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();
    accessToken = (await registerUser(app, credentials)).accessToken;
    orgname = (await getUser(app, accessToken)).defaultOrgname;
    await setEmailVerifiedByEmail(app, credentials.email);
  });

  afterAll(async () => {
    await app.close();
  });

  it("Should let users create and delete scoped api tokens", async () => {
    // Create USER token
    const apiToken = await createToken(orgname, accessToken, "USER");

    // Verify token exists
    await verifyTokenExists(orgname, accessToken);

    // Verify USER token role
    await verifyUserRole(apiToken.key, "USER");

    // Verify restricted actions with USER token
    await verifyRestrictedActions(orgname, apiToken.key);

    // Delete token
    await deleteToken(orgname, accessToken, apiToken.id);

    // Ensure token is no longer valid
    await verifyTokenRevocation(orgname, apiToken.key);

    // Create a token with restricted domain
    const badUserToken = await createToken(
      orgname,
      accessToken,
      "USER",
      "localhost"
    );

    // Verify restricted access due to bad domain
    await verifyRestrictedDomainAccess(orgname, badUserToken.key);
  });

  const createToken = async (
    orgname: string,
    accessToken: string,
    role: RoleType,
    domains = "*"
  ): Promise<ApiTokenEntity> => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/api-tokens`)
      .send({ domains, name: `${role}-token`, role })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(201);
    expect(res).toSatisfyApiSpec();
    return res.body;
  };

  const verifyTokenExists = async (orgname: string, accessToken: string) => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/api-tokens`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(200);
    expect(res.body.results.length).toBe(1);
    expect(res).toSatisfyApiSpec();
    return res.body.results[0].id;
  };

  const verifyUserRole = async (
    accessToken: string,
    expectedRole: RoleType
  ) => {
    const res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(200);
    expect(res).toSatisfyApiSpec();
    expect(res.body.memberships.length).toBe(1);
    expect(res.body.memberships[0].role).toBe(expectedRole);
  };

  const verifyRestrictedActions = async (
    orgname: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(403);
  };

  const deleteToken = async (orgname: string, accessToken: string, tokenId) => {
    const res = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}/api-tokens/${tokenId}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(200);
  };

  const verifyTokenRevocation = async (
    orgname: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(401);
  };

  const verifyRestrictedDomainAccess = async (
    orgname: string,
    accessToken: string
  ) => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(401);
  };
});
