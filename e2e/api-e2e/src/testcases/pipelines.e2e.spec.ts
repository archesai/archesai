import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser, setEmailVerified } from '#utils/helpers'

describe('Pipelines', () => {
  let app: HttpInstance
  let accessToken: string
  let orgname: string

  const credentials = {
    email: 'pipelines-test@archesai.com',
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

  it('should create default pipeline on user creation', async () => {
    await getPipeline()
  })

  it('should create a new pipeline', async () => {
    const newPipeline = {
      description: 'A new pipeline for testing purposes',
      name: 'New Pipeline',
      steps: []
    }

    await createPipeline(newPipeline)
  })

  it('should update an existing pipeline', async () => {
    // get original pipeline

    const pipeline = await getPipeline()

    const res = await request(app.getHttpServer())
      .patch(`/organizations/${orgname}/pipelines/${pipeline.id}`) // or .patch depending on your API
      .set('Authorization', `Bearer ${accessToken}`)
      .send({
        name: 'Updated Pipeline',
        steps: []
      })
    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(200) // or the appropriate success status code
    expect(res.body).toMatchObject({
      ...pipeline,
      name: 'Updated Pipeline'
    }) // or adjust to match the expected response
  })

  const createPipeline = async (
    createPipelineRequest: CreatePipelineRequest
  ) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/pipelines`)
      .set('Authorization', `Bearer ${accessToken}`)
      .send(createPipelineRequest)

    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(201)
    return res.body
  }

  const getPipeline = async () => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/pipelines`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(200)
    return res.body.results[0]
  }
})
