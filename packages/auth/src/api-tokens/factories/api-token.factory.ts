import { faker } from '@faker-js/faker'

import type { ApiTokenEntity } from '@archesai/schemas'

export const createRandomApiToken = (
  overrides?: Partial<ApiTokenEntity>
): ApiTokenEntity => {
  return {
    createdAt: faker.date.recent().toISOString(),
    id: faker.string.uuid(),
    key: faker.internet.password(),
    name: faker.lorem.words(3),
    organizationId: faker.string.uuid(),
    role: 'ADMIN',
    updatedAt: faker.date.recent().toISOString(),
    ...overrides
  }
}
