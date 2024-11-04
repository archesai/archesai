import { ChatbotEntity } from "@/src/chatbots/entities/chatbot.entity";
import { MessageEntity } from "@/src/messages/entities/message.entity";
import { ThreadEntity } from "@/src/threads/entities/thread.entity";
import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { UsersService } from "../src/users/users.service";
import { createApp, getUser, registerUser } from "./util";

describe("Chatbots", () => {
  let app: INestApplication;
  let accessToken: string;
  let orgname: string;

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

  it("should be able to chat with a chatbot", async () => {
    // Create a chatbot
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/chatbots`)
      .send({
        description:
          "You are an educational teacher and you help people learn about the stories of Aesop's Fables.",
        llmBase: "GPT_3_5_TURBO_16_K",
        name: "My Test Chatbot",
      })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res.status).toBe(201);
    expect(res).toSatisfyApiSpec();
    const chatbot = res.body as ChatbotEntity;
    expect(chatbot.name).toBe("My Test Chatbot");

    // Create a thread
    const threadRes = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/chatbots/${chatbot.id}/threads`)
      .send({
        name: "Aesop's Fables",
      })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(threadRes.status).toBe(201);
    expect(threadRes).toSatisfyApiSpec();
    const thread = threadRes.body as ThreadEntity;

    // Chat with the chatbot
    const chatRes = await request(app.getHttpServer())
      .post(
        `/organizations/${orgname}/chatbots/${chatbot.id}/threads/${thread.id}/messages`
      )
      .send({
        question: "What is the story of the tortoise and the hare?",
      })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(chatRes.status).toBe(201);
    expect(chatRes).toSatisfyApiSpec();
    const message = chatRes.body as MessageEntity;
    expect(message.answer.length).toBeGreaterThan(0);

    // Delete the chatbot
    const deleteRes = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}/chatbots/${res.body.id}`)
      .set("Authorization", `Bearer ${accessToken}`);
    expect(deleteRes.status).toBe(200);
  });
});
