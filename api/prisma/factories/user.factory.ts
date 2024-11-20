import { UserEntity } from "@/src/users/entities/user.entity";
import { faker } from "@faker-js/faker";

export function createRandomUser(overrides?: Partial<UserEntity>): UserEntity {
  return new UserEntity({
    authProviders: [],
    createdAt: faker.date.recent(),
    deactivated: false,
    defaultOrgname: faker.word.adverb(),
    email: faker.internet.email(),
    emailVerified: true,
    firstName: faker.person.firstName(),
    id: faker.string.uuid(),
    lastName: faker.person.lastName(),
    memberships: [],
    password: faker.internet.password(),
    photoUrl: faker.image.avatar(),
    refreshToken: faker.internet.password(),
    updatedAt: faker.date.recent(),
    username: faker.internet.username(),
    ...overrides,
  });
}
