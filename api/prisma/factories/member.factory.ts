import { MemberEntity } from '@/src/members/entities/member.entity'
import { faker } from '@faker-js/faker'

export function createRandomMember(
  overrides?: Partial<MemberEntity>
): MemberEntity {
  return new MemberEntity({
    createdAt: faker.date.recent(),
    id: faker.string.uuid(),
    inviteAccepted: faker.datatype.boolean(),
    inviteEmail: faker.internet.email(),
    orgname: faker.internet.domainName(),
    role: faker.helpers.arrayElement(['ADMIN', 'USER']),
    updatedAt: faker.date.recent(),
    username: faker.internet.username(),
    ...overrides
  })
}
