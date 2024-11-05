import { EmailService } from "@/src/email/email.service";
import { ConfirmEmailVerificationDto } from "@/src/email-verification/dto/confirm-password-verification.dto";
import { INestApplication } from "@nestjs/common";
import request from "supertest";

import { createApp, getUser, registerUser } from "./util";

describe("Email Verification", () => {
  let app: INestApplication;
  let accessToken: string;
  let capturedToken: null | string = null;

  const userCredentials = {
    email: "email-verification-test@archesai.com",
    password: "password",
    username: "email-verification-test",
  };

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

    // Register user and admin
    accessToken = (await registerUser(app, userCredentials)).accessToken;
  });

  afterAll(async () => {
    await app.close();
  });

  it("should mark user as email-verified false by default", async () => {
    const user = await getUser(app, accessToken);
    expect(user.emailVerified).toBe(false);
  });

  it("should send an email verification link", async () => {
    const res = await request(app.getHttpServer())
      .post("/auth/email-verification/request")
      .set("Authorization", `Bearer ${accessToken}`);

    expect(res.status).toBe(201);
    expect(capturedToken).not.toBeNull();
  });

  it("should confirm email verification", async () => {
    await confirmEmailVerificationExpectStatus(
      {
        token: capturedToken!,
      },
      201
    );

    const user = await getUser(app, accessToken);
    expect(user.emailVerified).toBe(true);
  });

  it("should throw 400 if email token is invalid", async () => {
    await confirmEmailVerificationExpectStatus(
      {
        token: capturedToken!,
      },
      400
    );
  });

  const confirmEmailVerificationExpectStatus = async (
    confirmEmailVerificationDto: ConfirmEmailVerificationDto,
    status: number
  ) => {
    const res = await request(app.getHttpServer())
      .post("/auth/email-verification/confirm")
      .send(confirmEmailVerificationDto);
    expect(res).toSatisfyApiSpec();
    expect(res.status).toBe(status);
  };
});
