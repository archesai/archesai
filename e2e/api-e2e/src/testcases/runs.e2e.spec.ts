import _jestOpenAPI from 'jest-openapi'
import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser, setEmailVerified } from '#utils/helpers'

describe('Runs', () => {
  let app: HttpInstance
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
    const badInputs = [
      {
        inputs: [],
        pipelineId: '1',
        runType: 'TOOL_RUN'
      },
      {
        inputs: [],
        pipelineId: '1',
        runType: 'TOOL_RUN'
      },
      {
        inputs: [],
        pipelineId: '1',
        runType: 'TOOL_RUN'
      },
      {
        inputs: [],
        pipelineId: '1',
        runType: 'PIPELINE_RUN'
      },
      {
        inputs: [],
        pipelineId: '1',
        runType: 'PIPELINE_RUN'
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
    const summarizerTool = tools.find((tool) => tool.name === 'Summarize')

    expect(summarizerTool).toBeDefined()

    const createRunDto = {
      inputs: [],
      pipelineId: summarizerTool.id,
      runType: 'TOOL_RUN'
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
    return res.body.results as {
      id: string
      name: string
    }[]
  }
})
