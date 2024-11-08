import { PaginatedDto } from "@/src/common/paginated.dto";
import { MessageEntity } from "@/src/messages/entities/message.entity";
import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { UsersService } from "../src/users/users.service";
import { createApp, getUser, registerUser } from "./util";

describe("Threads", () => {
  let app: INestApplication;
  let accessToken: string;
  let orgname: string;
  let threadId: string;

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
  });

  afterAll(async () => {
    await app.close();
  });

  it("should be able to create a thread", async () => {
    // Create a thread
    const threadRes = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/threads`)
      .send({
        name: "Aesop's Fables",
      })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(threadRes.status).toBe(201);
    expect(threadRes).toSatisfyApiSpec();
    threadId = threadRes.body.id;
  });

  it("should be able to create a message in a thread", async () => {
    // Chat with the chatbot
    const chatRes = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/threads/${threadId}/messages`)
      .send({
        question: "What is the story of the tortoise and the hare?",
      })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(chatRes.status).toBe(201);
    expect(chatRes).toSatisfyApiSpec();
    const message = chatRes.body as MessageEntity;
    expect(message.answer.length).toBeGreaterThan(0);
  });

  it("should be able to get messages in a thread", async () => {
    // Get messages
    const messagesRes = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/threads/${threadId}/messages`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(messagesRes.status).toBe(200);
    expect(messagesRes).toSatisfyApiSpec();
    const messages = messagesRes.body as PaginatedDto<MessageEntity>;
    expect(messages.results.length).toBeGreaterThan(0);
  });

  it("should be able to delete a thread", async () => {
    // Delete the chatbot
    const deleteRes = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}/threads/${threadId}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(deleteRes.status).toBe(200);
  });
});
