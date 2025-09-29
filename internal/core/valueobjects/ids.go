package valueobjects

import (
	"github.com/google/uuid"
)

// ID value objects for domain entities

type OrganizationID uuid.UUID
type UserID uuid.UUID
type LabelID uuid.UUID
type PipelineID uuid.UUID
type PipelineStepID uuid.UUID
type RunID uuid.UUID
type ToolID uuid.UUID
type ArtifactID uuid.UUID
type InvitationID uuid.UUID
type MemberID uuid.UUID
type AccountID uuid.UUID
type APIKeyID uuid.UUID
type SessionID uuid.UUID

// Status value objects

type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusRunning   RunStatus = "running"
	RunStatusCompleted RunStatus = "completed"
	RunStatusFailed    RunStatus = "failed"
	RunStatusCancelled RunStatus = "cancelled"
)

// MIME type value object

type MIMEType string
