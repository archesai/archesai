import { INestApplication } from "@nestjs/common";
import * as fs from "fs";
import request from "supertest";

import { OrganizationsService } from "../src/organizations/organizations.service";
import { UsersService } from "../src/users/users.service";
import { createApp, sleep } from "./util";

describe("Upload Document", () => {
  let app: INestApplication;
  let token: string;

  const credentials = {
    email: "upload-doc-test@archesai.com",
    password: "password",
    username: "upload-doc-test",
  };

  beforeEach(async () => {
    app = await createApp();
    await app.init();
  });

  it("Should process uploaded files", async () => {
    // Create user
    let res = await request(app.getHttpServer())
      .post("/auth/register")
      .send(credentials);
    expect(res.status).toBe(201);
    token = res.body.apiToken;

    // Verify email
    const usersService = app.get<UsersService>(UsersService);
    await usersService.setEmailVerifiedByEmail(credentials.email);

    // Get organization name
    res = await request(app.getHttpServer())
      .get("/user")
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(200);

    const orgname = res.body.defaultOrgname;

    const filePath = `book-${new Date().valueOf().toString()}.pdf`;

    // Get write urls
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/storage/write`)
      .send({
        path: filePath,
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(201);
    const { write } = res.body;

    // Write file
    const file = fs.readFileSync("./test/testdata/book.pdf");
    const response = await fetch(write, {
      body: file,
      method: "PUT",
    });
    expect(response.status).toBe(200);

    // Get read url
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/storage/read`)
      .send({
        path: filePath,
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(201);
    const { read } = res.body;

    // Remove credits
    const organizationsService =
      app.get<OrganizationsService>(OrganizationsService);
    await organizationsService.removeCredits(orgname, 10000);
    await organizationsService.setPlan(orgname, "FREE");

    // Try to upload file and fail bc not enough credits
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/documents`)
      .send({
        documentUrl: read,
        name: "book.pdf",
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(403);

    // Add credits
    await organizationsService.addCredits(orgname, 500000);

    // Try to upload file and succeed
    res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/documents`)
      .send({
        documentUrl: read,
        name: "book.pdf",
      })
      .set("Authorization", "Bearer " + token);
    expect(res.status).toBe(201);
    const { id } = res.body;

    // Try to upload file and succeed
    let complete = false;
    for (let i = 0; i < 20; i++) {
      res = await request(app.getHttpServer())
        .get(`/organizations/${orgname}/documents/${id}`)
        .set("Authorization", "Bearer " + token);
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
  });

  afterEach(async () => {
    await app.close();
  });
});
