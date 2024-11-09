import { ContentService } from "@/src/content/content.service";
import { faker } from "@faker-js/faker";
import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { UsersService } from "../src/users/users.service";
import { createApp, getUser, registerUser } from "./util";

describe("Labels", () => {
  let app: INestApplication;
  let accessToken: string;
  let orgname: string;
  let labelId: string;

  const credentials = {
    email: "chatbots-test@archesai.com",
    password: "password",
    username: "chatbots-test",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;

    const usersService = app.get<UsersService>(UsersService);
    await usersService.setEmailVerifiedByEmail(credentials.email);

    const user = await getUser(app, accessToken);
    orgname = user.defaultOrgname;

    const contentService = app.get(ContentService);
    for (let i = 0; i < 100; i++) {
      await contentService.create(orgname, {
        name: faker.lorem.words(),
        text: faker.lorem.sentence(),
      });
    }
  });

  afterAll(async () => {
    await app.close();
  });

  it("should be able to create a label", async () => {
    // Create a label
    const labelRes = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/labels`)
      .send({
        name: "Aesop's Fables",
      })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(labelRes.status).toBe(201);
    expect(labelRes).toSatisfyApiSpec();
    labelId = labelRes.body.id;
  });

  it("should be able to get content with a label", async () => {
    // Get messages
    const contentRes = await request(app.getHttpServer())
      .get(
        `/organizations/${orgname}/content?filters=${JSON.stringify([{ field: "labelId", operator: "equals", value: labelId }])}`
      )
      .set("Authorization", `Bearer ${accessToken}`);
    expect(contentRes.status).toBe(200);
    expect(contentRes).toSatisfyApiSpec();
    const content = contentRes.body;
    expect(content.results.length).toBeGreaterThan(0);
  });

  it("should be able to delete a label", async () => {
    // Delete the chatbot
    const deleteRes = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}/labels/${labelId}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(deleteRes.status).toBe(200);
  });
});
