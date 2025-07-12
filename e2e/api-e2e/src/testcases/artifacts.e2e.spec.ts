import fs from 'fs'

import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

import { getUser, registerUser, setEmailVerified } from '#utils/helpers'

describe('Artifacts', () => {
  let app: HttpInstance
  let accessToken: string
  let organizationId: string
  let artifactId: string

  const credentials = {
    email: 'artifacts-test@archesai.com',
    password: 'password'
  }

  const filePath = `book-${new Date().valueOf().toString()}.pdf`

  beforeAll(async () => {
    app = await createApp()
    await app.init()

    accessToken = (await registerUser(app, credentials)).accessToken

    const user = await getUser(app, accessToken)
    organizationId = user.defaultOrgname
    await setEmailVerified(app, user.id)
  })

  afterAll(async () => {
    await app.close()
  })

  it('CREATE - should be able to create artifacts', async () => {
    // Upload the file
    const readUrl = await uploadFile(organizationId, accessToken, filePath)

    const res = await request(app.getHttpServer())
      .post(`/organizations/${organizationId}/artifacts`)
      .send({ name: 'book.pdf', url: readUrl })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res).toSatisfyApiSpec()
    expect(res.status).toBe(201)
    artifactId = res.body.id
  })

  it('UPDATE - should be able to update artifacts name', async () => {
    const res = await request(app.getHttpServer())
      .patch(`/organizations/${organizationId}/artifacts/${artifactId}`)
      .send({ name: 'new-book.pdf' })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res).toSatisfyApiSpec()
    expect(res.body.name).toBe('new-book.pdf')
    expect(res.status).toBe(200)
  })

  it('UPDATE - should fail if you try to create with bad labels', async () => {
    const res = await request(app.getHttpServer())
      .patch(`/organizations/${organizationId}/artifacts/${artifactId}`)
      .send({ labels: ['label1', 'label2'] })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(404)
  })

  it('UPDATE - should be able to update artifacts labels', async () => {
    const label = await request(app.getHttpServer())
      .post(`/organizations/${organizationId}/labels`)
      .send({ name: 'artifacts-test-label' })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(label.status).toBe(201)
    expect(label).toSatisfyApiSpec()

    const res = await request(app.getHttpServer())
      .patch(`/organizations/${organizationId}/artifacts/${artifactId}`)
      .send({ labels: [label.body.name] })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(res.status).toBe(200)
    expect(res).toSatisfyApiSpec()

    const getRes = await request(app.getHttpServer())
      .get(`/organizations/${organizationId}/artifacts/${artifactId}`)
      .set('Authorization', `Bearer ${accessToken}`)
    expect(getRes.status).toBe(200)
    expect(getRes.body.labels.length).toBe(1)
    expect(getRes.body.labels[0].id).toBe(label.body.id)
    expect(getRes.body.labels[0].name).toBe(label.body.name)
    expect(getRes).toSatisfyApiSpec()
  })

  // Helper function to get a write url, upload a file, and get a read url
  const uploadFile = async (
    organizationId: string,
    accessToken: string,
    filePath: string
  ) => {
    const fileRes = await request(app.getHttpServer())
      .post(`/organizations/${organizationId}/storage/write`)
      .send({ path: filePath })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(fileRes.status).toBe(201)
    expect(fileRes).toSatisfyApiSpec()
    const { write } = fileRes.body

    const fileData = fs.readFileSync('./test/testdata/book.pdf')
    const uploadRes = await fetch(write, {
      body: fileData,
      method: 'PUT'
    })
    expect(uploadRes.status).toBe(200)

    const readRes = await request(app.getHttpServer())
      .post(`/organizations/${organizationId}/storage/read`)
      .send({ path: filePath })
      .set('Authorization', `Bearer ${accessToken}`)
    expect(readRes.status).toBe(201)
    return readRes.body.read
  }
})
