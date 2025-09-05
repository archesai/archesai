/**
 * Base types and entity definitions for the UI package
 * These types are used across UI components without dependency on the client package
 */

/**
 * Artifact entity type
 */
export interface ArtifactEntity extends BaseEntity {
  /** Artifact description */
  description?: string
  /** Artifact metadata */
  metadata?: Record<string, unknown>
  /** Artifact MIME type */
  mimeType: string
  /** Artifact name */
  name?: string
  /** Organization ID that owns this artifact */
  organizationId?: Uuid
  /** Artifact size in bytes */
  size?: number
  /** Artifact text content or URL */
  text?: null | string
}

/**
 * Base entity type matching the common fields from OpenAPI schemas
 * Based on the BaseEntity schema from api/components/schemas/BaseEntity.yaml
 */
export interface BaseEntity {
  /** The date this item was created */
  createdAt: string
  /** The ID of the item */
  id: Uuid
  /** The date this item was last updated */
  updatedAt: string
}

// Types that need to be defined for the UI components
export interface FilterCondition {
  field: string
  operator: string
  type: "condition"
  value: FilterValue
}

export interface FilterGroup {
  children: FilterNode[]
  operator: "and" | "or"
  type: "group"
}

/**
 * A recursive filter node that can be a condition or group
 */
export type FilterNode = FilterCondition | FilterGroup

export type FilterValue =
  | (boolean | null | number | string)[]
  | boolean
  | null
  | number
  | string

/**
 * Invitation entity type
 */
export interface InvitationEntity extends BaseEntity {
  /** The email address the invitation was sent to */
  email: string
  /** Expiration date */
  expiresAt: string
  /** Inviter user ID */
  invitedBy: Uuid
  /** Invitation metadata */
  metadata?: Record<string, unknown>
  /** The organization ID */
  organizationId: Uuid
  /** The role assigned to the invitation */
  role: string
  /** Invitation status */
  status: "accepted" | "cancelled" | "expired" | "pending"
}

/**
 * Label entity type
 */
export interface LabelEntity extends BaseEntity {
  /** Label color */
  color?: string
  /** Label description */
  description?: string
  /** Label metadata */
  metadata?: Record<string, unknown>
  /** Label name */
  name: string
  /** Organization ID that owns this label */
  organizationId: Uuid
}

/**
 * Member entity type (organization membership)
 */
export interface MemberEntity extends BaseEntity {
  /** Member metadata */
  metadata?: Record<string, unknown>
  /** The organization ID */
  organizationId: Uuid
  /** The member's role */
  role: string
  /** The user ID */
  userId: Uuid
}

/**
 * Organization entity type
 */
export interface OrganizationEntity extends BaseEntity {
  /** Organization description */
  description?: string
  /** Organization logo URL */
  logo?: string
  /** Organization metadata */
  metadata?: Record<string, unknown>
  /** The organization name */
  name: string
  /** Organization slug */
  slug?: string
}

/**
 * Pipeline entity type
 */
export interface PipelineEntity extends BaseEntity {
  /** Pipeline configuration */
  config?: Record<string, unknown>
  /** Pipeline description */
  description?: string
  /** Pipeline metadata */
  metadata?: Record<string, unknown>
  /** Pipeline name */
  name: string
  /** Organization ID that owns this pipeline */
  organizationId: Uuid
  /** Pipeline status */
  status: "active" | "archived" | "inactive"
}

/**
 * Run entity type
 */
export interface RunEntity extends BaseEntity {
  /** Run completed timestamp */
  completedAt?: string
  /** Run error message if failed */
  error?: string
  /** Run metadata */
  metadata?: Record<string, unknown>
  /** Run name */
  name?: string
  /** Pipeline ID this run belongs to */
  pipelineId: Uuid
  /** Run result data */
  result?: Record<string, unknown>
  /** Run started timestamp */
  startedAt?: string
  /** Run status */
  status: "cancelled" | "completed" | "failed" | "pending" | "running"
}

export interface SearchQuery {
  filter?: FilterNode
  page?: {
    number?: number
    size?: number
  }
  sort?: {
    field: string
    order: "asc" | "desc"
  }[]
}

/**
 * Tool entity type
 */
export interface ToolEntity extends BaseEntity {
  /** Tool configuration */
  config?: Record<string, unknown>
  /** Tool description */
  description?: string
  /** Tool metadata */
  metadata?: Record<string, unknown>
  /** Tool name */
  name: string
  /** Tool status */
  status: "active" | "deprecated" | "inactive"
  /** Tool type */
  type: string
}

/**
 * User entity type
 */
export interface UserEntity extends BaseEntity {
  /** The user's email address */
  email: string
  /** Whether the user's email has been verified */
  emailVerified: boolean
  /** The user's avatar image URL */
  image?: string
  /** The user's display name */
  name: string
}

/**
 * Universally Unique Identifier
 * @minLength 36
 */
export type Uuid = string
