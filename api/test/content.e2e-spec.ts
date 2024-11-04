import { INestApplication } from "@nestjs/common";
import * as fs from "fs";
import request from "supertest";

import {
  createApp,
  getUser,
  registerUser,
  setEmailVerifiedByEmail,
} from "./util";

describe("Content", () => {
  let app: INestApplication;
  let accessToken: string;
  let orgname: string;

  const credentials = {
    email: "content-test@archesai.com",
    password: "password",
    username: "content-test",
  };

  const filePath = `book-${new Date().valueOf().toString()}.pdf`;

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;

    const user = await getUser(app, accessToken);
    orgname = user.defaultOrgname;
    await setEmailVerifiedByEmail(app, credentials.email);
  });

  afterAll(async () => {
    await app.close();
  });

  it("Should process uploaded content", async () => {
    // Upload the file
    const readUrl = await uploadFile(orgname, accessToken, filePath);

    // Add credits and attempt to upload again
    await expectContentUpload(orgname, accessToken, readUrl, 201);

    // Poll for document processing completion
    // await waitForDocumentProcessing(orgname, accessToken, contentId);
  });

  // Helper function to upload a document and assert status
  const expectContentUpload = async (
    orgname,
    accessToken,
    readUrl,
    expectedStatus
  ) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/content`)
      .send({ name: "book.pdf", pipelineId: "", url: readUrl })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(res).toSatisfyApiSpec();
    expect(res.status).toBe(expectedStatus);
    return res.body.id;
  };

  // Helper function to upload a file and return the read URL
  const uploadFile = async (orgname, accessToken, filePath) => {
    const fileRes = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/storage/write`)
      .send({ path: filePath })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(fileRes.status).toBe(201);
    expect(fileRes).toSatisfyApiSpec();
    const { write } = fileRes.body;

    const fileData = fs.readFileSync("./test/testdata/book.pdf");
    const uploadRes = await fetch(write, {
      body: fileData,
      method: "PUT",
    });
    expect(uploadRes.status).toBe(200);

    const readRes = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/storage/read`)
      .send({ path: filePath })
      .set("Authorization", `Bearer ${accessToken}`);
    expect(readRes.status).toBe(201);
    return readRes.body.read;
  };

  // Helper function to poll for document processing status
  // const waitForDocumentProcessing = async (orgname, accessToken, contentId) => {
  //   let complete = false;

  //   for (let i = 0; i < 20; i++) {
  //     const res = await request(app.getHttpServer())
  //       .get(`/organizations/${orgname}/content/${contentId}`)
  //       .set("Authorization", `Bearer ${accessToken}`);
  //     expect(res.status).toBe(200);

  //     if (res.body.job.status === "ERROR") {
  //       throw new Error("Document processing failed");
  //     }
  //     if (res.body.job.status === "COMPLETE") {
  //       complete = true;
  //       break;
  //     }

  //     await sleep(3000);
  //   }

  //   expect(complete).toBe(true);
  // };
});
