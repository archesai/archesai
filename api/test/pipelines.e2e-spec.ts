import { CreatePipelineDto } from "@/src/pipelines/dto/create-pipeline.dto";
import { PipelineEntity } from "@/src/pipelines/entities/pipeline.entity";
import { INestApplication } from "@nestjs/common";
import request from "supertest";

import {
  createApp,
  getUser,
  registerUser,
  setEmailVerifiedByEmail,
} from "./util";

describe("Pipelines", () => {
  let app: INestApplication;
  let accessToken: string;
  let orgname: string;

  const credentials = {
    email: "pipelines-test@archesai.com",
    password: "password",
    username: "pipelines-test",
  };

  beforeAll(async () => {
    app = await createApp();
    await app.init();

    accessToken = (await registerUser(app, credentials)).accessToken;

    const user = await getUser(app, accessToken);
    orgname = user.defaultOrgname;
    await setEmailVerifiedByEmail(app, user.email);
  });

  afterAll(async () => {
    await app.close();
  });

  it("should create default pipeline on user creation", async () => {
    await getPipeline();
  });

  it("should create a new pipeline", async () => {
    const newPipeline: CreatePipelineDto = {
      description: "A new pipeline for testing purposes",
      name: "New Pipeline",
      pipelineTools: [],
    };

    await createPipeline(newPipeline);
  });

  it("should update an existing pipeline", async () => {
    // get original pipeline

    const pipeline = await getPipeline();

    const res = await request(app.getHttpServer())
      .patch(`/organizations/${orgname}/pipelines/${pipeline.id}`) // or .patch depending on your API
      .set("Authorization", `Bearer ${accessToken}`)
      .send({
        name: "Updated Pipeline",
        pipelineTools: [],
      });
    expect(res).toSatisfyApiSpec();
    expect(res.status).toBe(200); // or the appropriate success status code
    expect(res.body).toMatchObject({
      ...pipeline,
      name: "Updated Pipeline",
    }); // or adjust to match the expected response
  });

  const createPipeline = async (createPipelineDto: CreatePipelineDto) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/pipelines`)
      .set("Authorization", `Bearer ${accessToken}`)
      .send(createPipelineDto);

    expect(res).toSatisfyApiSpec();
    expect(res.status).toBe(201);
    return res.body as PipelineEntity;
  };

  const getPipeline = async () => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/pipelines`)
      .set("Authorization", `Bearer ${accessToken}`);
    console.log(JSON.stringify(res.body, null, 2));
    expect(res).toSatisfyApiSpec();
    expect(res.status).toBe(200);
    return res.body.results[0] as PipelineEntity;
  };
});
