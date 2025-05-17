import { faker } from '@faker-js/faker'

import type { ApiTokenEntity } from '@archesai/domain'

import { API_TOKEN_ENTITY_KEY } from '@archesai/domain'

export const createRandomApiToken = (
  overrides?: Partial<ApiTokenEntity>
): ApiTokenEntity => {
  return {
    createdAt: faker.date.recent().toISOString(),
    id: faker.string.uuid(),
    key: faker.internet.password(),
    name: faker.word.noun(),
    orgname: faker.internet.domainName(),
    role: 'ADMIN',
    slug: faker.string.alpha(),
    type: API_TOKEN_ENTITY_KEY,
    updatedAt: faker.date.recent().toISOString(),
    ...overrides
  }
}
