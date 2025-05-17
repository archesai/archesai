import { faker } from '@faker-js/faker'

import type { MemberEntity } from '@archesai/domain'

import { MEMBER_ENTITY_KEY, RoleTypes } from '@archesai/domain'

export function createRandomMember(
  overrides?: Partial<MemberEntity>
): MemberEntity {
  return {
    createdAt: faker.date.recent().toISOString(),
    id: faker.string.uuid(),
    invitationId: faker.string.uuid(),
    name: faker.person.fullName(),
    orgname: faker.internet.domainName(),
    role: faker.helpers.arrayElement(RoleTypes),
    slug: faker.lorem.slug(),
    type: MEMBER_ENTITY_KEY,
    updatedAt: faker.date.recent().toISOString(),
    userId: faker.string.uuid(),
    ...overrides
  }
}
