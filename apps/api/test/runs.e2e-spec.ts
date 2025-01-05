import { CreateRunDto } from '@/src/runs/dto/create-run.dto'
import { ToolEntity } from '@/src/tools/entities/tool.entity'
import { INestApplication } from '@nestjs/common'
import request from 'supertest'

import { createApp, getUser, registerUser, setEmailVerified } from './util'
import { RunTypeEnum } from '@/src/runs/entities/run.entity'

describe('Runs', () => {
  let app: INestApplication
  let accessToken: string
  let orgname: string

  const credentials = {
    email: 'runs-test@archesai.com',
    password: 'password'
  }

  beforeAll(async () => {
    app = await createApp()
    await app.init()

    accessToken = (await registerUser(app, credentials)).accessToken

    const user = await getUser(app, accessToken)
    orgname = user.defaultOrgname
    await setEmailVerified(app, user.id)
  })

  afterAll(async () => {
    await app.close()
  })

  it('CREATE - should block bad requests to create a run', async () => {
    const badInputs: CreateRunDto[] = [
      {
        contentIds: ['1'],
        runType: RunTypeEnum.TOOL_RUN
      },
      {
        runType: RunTypeEnum.TOOL_RUN,
        text: 'This is the text to use as input for the run.',
        contentIds: []
      },
      {
        pipelineId: '1',
        runType: RunTypeEnum.TOOL_RUN,
        url: 'https://example.com',
        contentIds: []
      },
      {
        contentIds: ['1'],
        runType: RunTypeEnum.PIPELINE_RUN
      },
      {
        runType: RunTypeEnum.PIPELINE_RUN,
        contentIds: []
      }
    ]

    for (const badRun of badInputs) {
      const res = await request(app.getHttpServer())
        .post(`/organizations/${orgname}/runs`)
        .set('Authorization', `Bearer ${accessToken}`)
        .send(badRun)
      expect(res.status).toBe(400)
    }
  })

  it('CREATE - should create a new run with a tool', async () => {
    const tools = await getDefaultTools()
    const summarizerTool = tools.find((tool) => tool.name === 'Summarize')!

    const createRunDto: CreateRunDto = {
      runType: RunTypeEnum.TOOL_RUN,
      text: 'This is the text to use as input for the run.',
      toolId: summarizerTool.id,
      contentIds: []
    }
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/runs`)
      .set('Authorization', `Bearer ${accessToken}`)
      .send(createRunDto)
    expect(res.status).toBe(201)
    expect(res).toSatisfyApiSpec()
  })

  const getDefaultTools = async () => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/tools`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(200)
    expect(res).toSatisfyApiSpec()
    return res.body.results as ToolEntity[]
  }
})
