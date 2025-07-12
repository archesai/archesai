import { faker } from '@faker-js/faker'

import type { UserEntity } from '@archesai/schemas'

import { USER_ENTITY_KEY } from '@archesai/schemas'

export function createRandomUser(overrides?: Partial<UserEntity>): UserEntity {
  return {
    createdAt: faker.date.recent().toISOString(),
    deactivated: false,
    email: faker.internet.email(),
    emailVerified: new Date().toISOString(),
    id: faker.string.uuid(),
    image: faker.image.avatar(),
    name: faker.person.fullName(),
    orgname: faker.word.adverb(),
    slug: faker.lorem.slug(),
    type: USER_ENTITY_KEY,
    updatedAt: faker.date.recent().toISOString(),
    ...overrides
  }
}
