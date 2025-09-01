// AUTH TYPES
export const AuthTypes = ['email', 'oauth', 'oidc', 'webauthn'] as const
export type AuthType = (typeof AuthTypes)[number]

// CONTENT BASE TYPES
export const ContentBaseTypes = ['AUDIO', 'IMAGE', 'TEXT', 'VIDEO'] as const
export type ContentBaseType = (typeof ContentBaseTypes)[number]

// PLAN TYPES
export const PlanTypes = [
  'BASIC',
  'FREE',
  'PREMIUM',
  'STANDARD',
  'UNLIMITED'
] as const
export type PlanType = (typeof PlanTypes)[number]

// PROVIDER TYPES
export const ProviderTypes = [
  'API_KEY',
  'FIREBASE',
  'LOCAL',
  'TWITTER'
] as const
export type ProviderType = (typeof ProviderTypes)[number]

// ROLE TYPES
export const RoleTypes = ['admin', 'owner', 'member'] as const
export type RoleType = (typeof RoleTypes)[number]

// STATUS TYPES
export const StatusTypes = [
  'COMPLETED',
  'FAILED',
  'PROCESSING',
  'QUEUED'
] as const
export type StatusType = (typeof StatusTypes)[number]
