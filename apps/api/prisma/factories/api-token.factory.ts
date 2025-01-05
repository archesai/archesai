import { ApiTokenEntity } from '@/src/api-tokens/entities/api-token.entity'
import { RoleTypeEnum } from '@/src/members/entities/member.entity'
import { faker } from '@faker-js/faker'

export function createRandomApiToken(
  overrides?: Partial<ApiTokenEntity>
): ApiTokenEntity {
  return new ApiTokenEntity({
    createdAt: faker.date.recent(),
    domains: faker.internet.domainName(),
    id: faker.string.uuid(),
    key: faker.internet.password(),
    name: faker.word.noun(),
    orgname: faker.internet.domainName(),
    role: RoleTypeEnum.ADMIN,
    updatedAt: faker.date.recent(),
    username: faker.internet.username(),
    ...overrides
  })
}
