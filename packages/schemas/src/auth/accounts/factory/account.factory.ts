import type { AccountEntity } from '#auth/accounts/entities/account.entity'

export const accountFactory = (
  overrides: Partial<AccountEntity>
): AccountEntity => {
  return {
    accessToken: null,
    accessTokenExpiresAt: null,
    accountId: '',
    createdAt: new Date().toISOString(),
    id: crypto.randomUUID(),
    idToken: null,
    password: null,
    providerId: '',
    refreshToken: null,
    refreshTokenExpiresAt: null,
    scope: null,
    updatedAt: new Date().toISOString(),
    userId: '',
    ...overrides
  }
}
