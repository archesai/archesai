import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser, setEmailVerified } from '#utils/helpers'

describe('Tools', () => {
  let app: HttpInstance
  let accessToken: string
  let orgname: string

  const credentials = {
    email: 'tools-test@archesai.com',
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

  it('should create default tools on user creation', async () => {
    const res = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/tools`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(200)
  })

  it('should create a new tool', async () => {
    const newTool = {
      description: 'A new tool for testing purposes',
      inputType: 'TEXT',
      name: 'New Tool',
      outputType: 'TEXT',
      toolBase: 'extract-text'
    }

    await createTool(newTool)
  })

  it('should update an existing tool', async () => {
    const newTool = {
      description: 'A new tool for testing purposes',
      inputType: 'TEXT',
      name: 'New Tool',
      outputType: 'TEXT',
      toolBase: 'extract-text'
    }

    const tool = await createTool(newTool)

    // request to update tool
    const res = await request(app.getHttpServer())
      .patch(`/organizations/${orgname}/tools/${tool.id}`) // or .patch depending on your API
      .set('Authorization', `Bearer ${accessToken}`)
      .send({
        name: 'Updated Tool'
      })
    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(200) // or the appropriate success status code
    expect(res.body).toMatchObject({
      ...tool,
      name: 'Updated Tool'
    }) // or adjust to match the expected response
  })

  const createTool = async (tool: CreateToolDto) => {
    const res = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/tools`)
      .set('Authorization', `Bearer ${accessToken}`)
      .send(tool)
    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(201)
    return res.body as ToolEntity
  }
})
