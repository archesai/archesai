import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { OrganizationsService } from "../src/organizations/organizations.service";
import { UsersService } from "../src/users/users.service";
import { createApp, sleep } from "./util";

describe("Agents", () => {
  let app: INestApplication;

  const credentials1 = {
    email: "agent1@archesai.com",
    password: "password",
    username: "agent1",
  };

  const credentials2 = {
    email: "agent2@archesai.com",
    password: "password",
    username: "agent2",
  };

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("should only let them read documents they own", async () => {
    // Create users
    let res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials1);
    expect(res.status).toBe(201);
    const agent1Token = res.body.apiToken;
    res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials2);
    expect(res.status).toBe(201);
    const agent2Token = res.body.apiToken;

    // Verify emails
    const usersService = app.get<UsersService>(UsersService);
    const organizationsService =
      app.get<OrganizationsService>(OrganizationsService);
    await usersService.setEmailVerifiedByEmail(credentials1.email);
    await usersService.setEmailVerifiedByEmail(credentials2.email);

    // Get orgnames
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + agent1Token);
    expect(res.status).toBe(200);
    const orgname1 = res.body.defaultOrgname;
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + agent2Token);
    expect(res.status).toBe(200);
    const orgname2 = res.body.defaultOrgname;

    // Add bunch of tokens
    await organizationsService.addCredits(orgname1, 100000);
    await organizationsService.addCredits(orgname2, 100000);

    // Index documents
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname1}/documents`)
      .send({
        name: "book1-test.pdf",
        url: "https://www.nraonlinetraining.org/documents/TEST_PDF_File.pdf",
      })
      .set("Authorization", "Bearer " + agent1Token);
    expect(res.status).toBe(201);
    const id1 = res.body.id;
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname2}/documents`)
      .send({
        name: "book2-test.pdf",
        url: "https://www.nraonlinetraining.org/documents/TEST_PDF_File.pdf",
      })
      .set("Authorization", "Bearer " + agent2Token);
    expect(res.status).toBe(201);
    const id2 = res.body.id;

    // Wait for documents to be processed
    let complete = false;
    for (let i = 0; i < 20; i++) {
      res = await request(app.getHttpServer())
        .get(`/organizations/${orgname1}/documents/${id1}`)
        .set("Authorization", "Bearer " + agent1Token);
      expect(res.status).toBe(200);
      if (res.body.job.status === "ERROR") {
        expect(true).toBe(false);
      }
      if (res.body.job.status === "COMPLETE") {
        complete = true;
        break;
      }
      await sleep(3000);
    }
    expect(complete).toBe(true);
    for (let i = 0; i < 20; i++) {
      res = await request(app.getHttpServer())
        .get(`/organizations/${orgname2}/documents/${id2}`)
        .set("Authorization", "Bearer " + agent2Token);
      expect(res.status).toBe(200);
      if (res.body.job.status === "ERROR") {
        expect(true).toBe(false);
      }
      if (res.body.job.status === "COMPLETE") {
        break;
      }
      await sleep(3000);
    }
    expect(complete).toBe(true);

    // Try to upload file and succeed
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname1}/chatbots`)
      .send({
        accessScope: "DOCUMENT",
        description:
          "You are an educational teacher and you help people learn about the stories of Aesop's Fables.",
        documentIds: ["doc1", "doc2"],
        llmBase: "GPT_3_5_TURBO_16_K",
        name: "My Test Agent",
      })
      .set("Authorization", "Bearer " + agent1Token);
    expect(res.status).toBe(404);
  });

  afterEach(async () => {
    await app.close();
  });
});
