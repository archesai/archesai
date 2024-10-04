import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { OrganizationsService } from "../src/organizations/organizations.service";
import { createApp } from "./util";

describe("Access Tokens", () => {
  let app: INestApplication;

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("Should let users create and delete scoped access tokens", async () => {
    // Create account
    const credentials = {
      email: "tokens@archesai.com",
      password: "password",
      username: "tokens",
    };

    const organizationsService =
      app.get<OrganizationsService>(OrganizationsService);

    let res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials);
    expect(res.status).toBe(201);
    const authToken = res.body.apiToken;

    // Get /user endpoint
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + authToken);
    expect(res.status).toBe(200);
    expect(res.body.memberships.length).toBe(1);
    expect(res.body.memberships[0].role).toBe("ADMIN");
    const orgname = res.body.memberships[0].orgname;

    // Create USER access token without API access and fail
    res = await request(app.getHttpServer())
      .post("/organizations/" + orgname + "/tokens")
      .send({ role: "USER" })
      .set("Authorization", "Bearer " + authToken);
    expect(res.status).toBe(403);
    await organizationsService.setPlan(orgname, "API");

    // Create USER access token
    res = await request(app.getHttpServer())
      .post("/organizations/" + orgname + "/tokens")
      .send({ domains: "*", name: "my-token", role: "USER" })
      .set("Authorization", "Bearer " + authToken);

    expect(res.status).toBe(201);
    const userToken = res.body.apiToken;

    // Should see tokens
    res = await request(app.getHttpServer())
      .get("/organizations/" + orgname + "/tokens")
      .set("Authorization", "Bearer " + authToken);
    expect(res.status).toBe(200);
    expect(res.body.length).toBe(1);
    const id = res.body[0].id;

    // Get /user endpoint and see USER
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + userToken);
    expect(res.status).toBe(200);
    expect(res.body.memberships.length).toBe(1);
    expect(res.body.memberships[0].role).toBe("USER");

    // Should be restricted
    res = await request(app.getHttpServer())
      .delete("/organizations/" + orgname)
      .set("Authorization", "Bearer " + userToken);
    expect(res.status).toBe(403);

    // Delete token
    res = await request(app.getHttpServer())
      .delete("/organizations/" + orgname + "/tokens/" + id)
      .set("Authorization", "Bearer " + authToken);
    expect(res.status).toBe(200);

    // Should no longer have access
    res = await request(app.getHttpServer())
      .delete("/organizations/" + orgname)
      .set("Authorization", "Bearer " + userToken);
    expect(res.status).toBe(401);

    // Create bad token
    res = await request(app.getHttpServer())
      .post("/organizations/" + orgname + "/tokens")
      .send({ domains: "localhost", name: "localhost-token", role: "USER" })
      .set("Authorization", "Bearer " + authToken);

    expect(res.status).toBe(201);
    const badUserToken = res.body.apiToken;

    // Should be restricted due to bad domain
    res = await request(app.getHttpServer())
      .get("/organizations/" + orgname)
      .set("Authorization", "Bearer " + badUserToken);
    expect(res.status).toBe(401);
  });

  afterEach(async () => {
    await app.close();
  });
});
