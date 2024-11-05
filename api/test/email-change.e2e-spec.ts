import { EmailService } from "@/src/email/email.service";
import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { ConfirmEmailChangeDto } from "../src/email-change/dto/confirm-email-change.dto";
import { RequestEmailChangeDto } from "../src/email-change/dto/request-email-change.dto";
import { createApp, getUser, registerUser } from "./util";

describe("Email Change", () => {
  let app: INestApplication;
  let accessToken: string;
  let capturedToken: null | string = null;

  const credentials = {
    email: "email-change-test@archesai.com",
    password: "password",
    username: "email-change-test",
  };

  const newEmail = "new-email@archesai.com";

  beforeAll(async () => {
    app = await createApp();

    // Mock EmailService
    const emailService = app.get(EmailService);
    jest.spyOn(emailService, "sendMail").mockImplementation(({ html }) => {
      const tokenMatch = (html as string).match(/token=([a-zA-Z0-9]+)/);
      if (tokenMatch) {
        capturedToken = tokenMatch[1];
      }
      return Promise.resolve();
    });

    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;
  });

  afterAll(async () => {
    await app.close();
  });

  it("should request an email change", async () => {
    const requestEmailChangeDto: RequestEmailChangeDto = {
      email: newEmail,
    };
    const res = await request(app.getHttpServer())
      .post("/auth/email-change/request")
      .set("Authorization", `Bearer ${accessToken}`)
      .send(requestEmailChangeDto);

    expect(res.status).toBe(201);
  });

  it("should throw an error if email change token is invalid", async () => {
    await confirmEmailChangeExpectStatus(
      {
        token: "invalid",
      },
      400
    );
  });

  it("should confirm the email change", async () => {
    await confirmEmailChangeExpectStatus(
      {
        token: capturedToken,
      },
      201
    );
  });

  it("should update the email", async () => {
    const user = await getUser(app, accessToken);
    expect(user.email).toBe(newEmail);
  });

  const confirmEmailChangeExpectStatus = async (
    confirmEmailChangeDto: ConfirmEmailChangeDto,
    status: number
  ) => {
    const res = await request(app.getHttpServer())
      .post("/auth/email-change/confirm")
      .send(confirmEmailChangeDto);
    expect(res).toSatisfyApiSpec();
    expect(res.status).toBe(status);
  };
});
