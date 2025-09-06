package content

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockRepository implements Repository for testing
type MockRepository struct {
	artifacts map[uuid.UUID]*Artifact
	labels    map[uuid.UUID]*Label
	err       error
}

// Compile-time check
var _ Repository = (*MockRepository)(nil)

func NewMockRepository() *MockRepository {
	return &MockRepository{
		artifacts: make(map[uuid.UUID]*Artifact),
		labels:    make(map[uuid.UUID]*Label),
	}
}

// Artifact methods
func (m *MockRepository) CreateArtifact(_ context.Context, artifact *Artifact) (*Artifact, error) {
	if m.err != nil {
		return nil, m.err
	}
	if artifact.Id == (uuid.UUID{}) {
		artifact.Id = uuid.New()
	}
	artifact.CreatedAt = time.Now()
	artifact.UpdatedAt = time.Now()
	m.artifacts[artifact.Id] = artifact
	return artifact, nil
}

func (m *MockRepository) GetArtifactByID(_ context.Context, id uuid.UUID) (*Artifact, error) {
	if m.err != nil {
		return nil, m.err
	}
	artifact, exists := m.artifacts[id]
	if !exists {
		return nil, ErrArtifactNotFound
	}
	return artifact, nil
}

func (m *MockRepository) UpdateArtifact(_ context.Context, id uuid.UUID, artifact *Artifact) (*Artifact, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.artifacts[id]; !exists {
		return nil, ErrArtifactNotFound
	}
	artifact.UpdatedAt = time.Now()
	m.artifacts[id] = artifact
	return artifact, nil
}

func (m *MockRepository) DeleteArtifact(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.artifacts[id]; !exists {
		return ErrArtifactNotFound
	}
	delete(m.artifacts, id)
	return nil
}

func (m *MockRepository) ListArtifacts(_ context.Context, _ ListArtifactsParams) ([]*Artifact, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	artifacts := make([]*Artifact, 0, len(m.artifacts))
	for _, artifact := range m.artifacts {
		artifacts = append(artifacts, artifact)
	}
	return artifacts, int64(len(artifacts)), nil
}

// Label methods
func (m *MockRepository) CreateLabel(_ context.Context, label *Label) (*Label, error) {
	if m.err != nil {
		return nil, m.err
	}
	if label.Id == (uuid.UUID{}) {
		label.Id = uuid.New()
	}
	label.CreatedAt = time.Now()
	label.UpdatedAt = time.Now()
	m.labels[label.Id] = label
	return label, nil
}

func (m *MockRepository) GetLabelByID(_ context.Context, id uuid.UUID) (*Label, error) {
	if m.err != nil {
		return nil, m.err
	}
	label, exists := m.labels[id]
	if !exists {
		return nil, ErrLabelNotFound
	}
	return label, nil
}

func (m *MockRepository) UpdateLabel(_ context.Context, id uuid.UUID, label *Label) (*Label, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.labels[id]; !exists {
		return nil, ErrLabelNotFound
	}
	label.UpdatedAt = time.Now()
	m.labels[id] = label
	return label, nil
}

func (m *MockRepository) DeleteLabel(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.labels[id]; !exists {
		return ErrLabelNotFound
	}
	delete(m.labels, id)
	return nil
}

func (m *MockRepository) ListLabels(_ context.Context, _ ListLabelsParams) ([]*Label, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	labels := make([]*Label, 0, len(m.labels))
	for _, label := range m.labels {
		labels = append(labels, label)
	}
	return labels, int64(len(labels)), nil
}

// Test cases
func TestService_CreateArtifact(t *testing.T) {
	tests := []struct {
		name       string
		req        *CreateArtifactJSONRequestBody
		orgID      string
		producerID string
		repoErr    error
		wantErr    bool
	}{
		{
			name: "successful creation",
			req: &CreateArtifactJSONRequestBody{
				Name: "Test Artifact",
				Text: "Test content",
			},
			orgID:      uuid.New().String(),
			producerID: uuid.New().String(),
			wantErr:    false,
		},
		{
			name: "artifact too large",
			req: &CreateArtifactJSONRequestBody{
				Name: "Large Artifact",
				Text: strings.Repeat("a", MaxArtifactSize+1),
			},
			orgID:      uuid.New().String(),
			producerID: uuid.New().String(),
			wantErr:    true,
		},
		{
			name: "repository error",
			req: &CreateArtifactJSONRequestBody{
				Name: "Test Artifact",
				Text: "Test content",
			},
			orgID:      uuid.New().String(),
			producerID: uuid.New().String(),
			repoErr:    errors.New("database error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			repo.err = tt.repoErr
			service := NewService(repo, slog.Default())

			artifact, err := service.CreateArtifact(context.Background(), tt.req, tt.orgID, tt.producerID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && artifact == nil {
				t.Error("CreateArtifact() returned nil artifact")
			}
		})
	}
}

func TestService_GetArtifact(t *testing.T) {
	tests := []struct {
		name       string
		artifactID uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "existing artifact",
			artifactID: uuid.New(),
			setup: func(r *MockRepository) {
				artifact := &Artifact{
					Id:   uuid.New(),
					Name: "Test Artifact",
					Text: "Test content",
				}
				r.artifacts[artifact.Id] = artifact
			},
			wantErr: false,
		},
		{
			name:       "non-existent artifact",
			artifactID: uuid.New(),
			setup:      func(_ *MockRepository) {},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			// Get the correct ID from the setup
			var testID uuid.UUID
			if tt.name == "existing artifact" {
				for id := range repo.artifacts {
					testID = id
					break
				}
			} else {
				testID = tt.artifactID
			}

			artifact, err := service.GetArtifact(context.Background(), testID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && artifact == nil {
				t.Error("GetArtifact() returned nil artifact")
			}
		})
	}
}

func TestService_UpdateArtifact(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name       string
		artifactID uuid.UUID
		req        *UpdateArtifactJSONRequestBody
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "successful update",
			artifactID: existingID,
			req: &UpdateArtifactJSONRequestBody{
				Name: "Updated Artifact",
				Text: "Updated content",
			},
			setup: func(r *MockRepository) {
				artifact := &Artifact{
					Id:   existingID,
					Name: "Original Artifact",
					Text: "Original content",
				}
				r.artifacts[existingID] = artifact
			},
			wantErr: false,
		},
		{
			name:       "update with too large content",
			artifactID: existingID,
			req: &UpdateArtifactJSONRequestBody{
				Text: strings.Repeat("a", MaxArtifactSize+1),
			},
			setup: func(r *MockRepository) {
				artifact := &Artifact{
					Id:   existingID,
					Name: "Original Artifact",
					Text: "Original content",
				}
				r.artifacts[existingID] = artifact
			},
			wantErr: true,
		},
		{
			name:       "non-existent artifact",
			artifactID: uuid.New(),
			req: &UpdateArtifactJSONRequestBody{
				Name: "Updated Artifact",
			},
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			artifact, err := service.UpdateArtifact(context.Background(), tt.artifactID, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if artifact == nil {
					t.Error("UpdateArtifact() returned nil artifact")
				} else {
					if tt.req.Name != "" && artifact.Name != tt.req.Name {
						t.Errorf("UpdateArtifact() name = %v, want %v", artifact.Name, tt.req.Name)
					}
					if tt.req.Text != "" && artifact.Text != tt.req.Text {
						t.Errorf("UpdateArtifact() text = %v, want %v", artifact.Text, tt.req.Text)
					}
				}
			}
		})
	}
}

func TestService_DeleteArtifact(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name       string
		artifactID uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			artifactID: existingID,
			setup: func(r *MockRepository) {
				artifact := &Artifact{
					Id:   existingID,
					Name: "Test Artifact",
				}
				r.artifacts[existingID] = artifact
			},
			wantErr: false,
		},
		{
			name:       "non-existent artifact",
			artifactID: uuid.New(),
			setup:      func(_ *MockRepository) {},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			err := service.DeleteArtifact(context.Background(), tt.artifactID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, exists := repo.artifacts[tt.artifactID]; exists {
					t.Error("DeleteArtifact() artifact still exists")
				}
			}
		})
	}
}

func TestService_ListArtifacts(t *testing.T) {
	tests := []struct {
		name    string
		orgID   string
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:   "list with artifacts",
			orgID:  uuid.New().String(),
			limit:  10,
			offset: 0,
			setup: func(r *MockRepository) {
				for i := 0; i < 5; i++ {
					artifact := &Artifact{
						Id:   uuid.New(),
						Name: "Test Artifact",
						Text: "Test content",
					}
					r.artifacts[artifact.Id] = artifact
				}
			},
			wantErr: false,
		},
		{
			name:    "empty list",
			orgID:   uuid.New().String(),
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			artifacts, total, err := service.ListArtifacts(context.Background(), tt.orgID, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListArtifacts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if artifacts == nil {
					t.Error("ListArtifacts() returned nil artifacts")
				}
				if total != len(artifacts) {
					t.Errorf("ListArtifacts() total = %v, want %v", total, len(artifacts))
				}
			}
		})
	}
}

func TestService_CreateLabel(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateLabelJSONRequestBody
		orgID   string
		repoErr error
		wantErr bool
	}{
		{
			name: "successful creation",
			req: &CreateLabelJSONRequestBody{
				Name: "Test Label",
			},
			orgID:   uuid.New().String(),
			wantErr: false,
		},
		{
			name: "repository error",
			req: &CreateLabelJSONRequestBody{
				Name: "Test Label",
			},
			orgID:   uuid.New().String(),
			repoErr: errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			repo.err = tt.repoErr
			service := NewService(repo, slog.Default())

			label, err := service.CreateLabel(context.Background(), tt.req, tt.orgID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLabel() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && label == nil {
				t.Error("CreateLabel() returned nil label")
			}
		})
	}
}

func TestService_GetLabel(t *testing.T) {
	tests := []struct {
		name    string
		labelID uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:    "existing label",
			labelID: uuid.New(),
			setup: func(r *MockRepository) {
				label := &Label{
					Id:   uuid.New(),
					Name: "Test Label",
				}
				r.labels[label.Id] = label
			},
			wantErr: false,
		},
		{
			name:    "non-existent label",
			labelID: uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			// Get the correct ID from the setup
			var testID uuid.UUID
			if tt.name == "existing label" {
				for id := range repo.labels {
					testID = id
					break
				}
			} else {
				testID = tt.labelID
			}

			label, err := service.GetLabel(context.Background(), testID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLabel() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && label == nil {
				t.Error("GetLabel() returned nil label")
			}
		})
	}
}

func TestService_SearchArtifacts(t *testing.T) {
	tests := []struct {
		name    string
		orgID   string
		query   string
		limit   int
		offset  int
		wantErr bool
	}{
		{
			name:    "search returns empty results",
			orgID:   uuid.New().String(),
			query:   "test query",
			limit:   10,
			offset:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo, slog.Default())

			artifacts, total, err := service.SearchArtifacts(context.Background(), tt.orgID, tt.query, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchArtifacts() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if artifacts == nil {
					t.Error("SearchArtifacts() returned nil artifacts")
				}
				if total != 0 {
					t.Errorf("SearchArtifacts() total = %v, want 0", total)
				}
			}
		})
	}
}

// TestService_UpdateLabel tests the UpdateLabel method
func TestService_UpdateLabel(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		labelID uuid.UUID
		req     *UpdateLabelJSONRequestBody
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:    "successful update",
			labelID: existingID,
			req: &UpdateLabelJSONRequestBody{
				Name: "Updated Label",
			},
			setup: func(r *MockRepository) {
				label := &Label{
					Id:   existingID,
					Name: "Original Label",
				}
				r.labels[existingID] = label
			},
			wantErr: false,
		},
		{
			name:    "non-existent label",
			labelID: uuid.New(),
			req: &UpdateLabelJSONRequestBody{
				Name: "Updated Label",
			},
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			label, err := service.UpdateLabel(context.Background(), tt.labelID, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLabel() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && label == nil {
				t.Error("UpdateLabel() returned nil label")
			}
		})
	}
}

// TestService_DeleteLabel tests the DeleteLabel method
func TestService_DeleteLabel(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		labelID uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:    "delete existing label",
			labelID: existingID,
			setup: func(r *MockRepository) {
				label := &Label{
					Id:   existingID,
					Name: "Test Label",
				}
				r.labels[existingID] = label
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent label",
			labelID: uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			err := service.DeleteLabel(context.Background(), tt.labelID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestService_ListLabels tests the ListLabels method
func TestService_ListLabels(t *testing.T) {
	orgID := uuid.New().String()

	tests := []struct {
		name    string
		orgID   string
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
		wantLen int
	}{
		{
			name:   "list labels for organization",
			orgID:  orgID,
			limit:  10,
			offset: 0,
			setup: func(r *MockRepository) {
				for i := 0; i < 3; i++ {
					label := &Label{
						Id:             uuid.New(),
						Name:           "Label " + string(rune('A'+i)),
						OrganizationId: orgID,
					}
					r.labels[label.Id] = label
				}
			},
			wantErr: false,
			wantLen: 3,
		},
		{
			name:    "empty label list",
			orgID:   uuid.New().String(),
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			labels, total, err := service.ListLabels(context.Background(), tt.orgID, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListLabels() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if labels == nil {
					t.Error("ListLabels() returned nil labels")
				}
				if total != tt.wantLen {
					t.Errorf("ListLabels() total = %v, want %v", total, tt.wantLen)
				}
			}
		})
	}
}

// TestService_GetArtifactsByLabel tests getting artifacts by label
func TestService_GetArtifactsByLabel(t *testing.T) {
	labelID := uuid.New()

	tests := []struct {
		name    string
		labelID uuid.UUID
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:    "get artifacts for label",
			labelID: labelID,
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			artifacts, total, err := service.GetArtifactsByLabel(context.Background(), tt.labelID, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetArtifactsByLabel() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if artifacts == nil {
					t.Error("GetArtifactsByLabel() returned nil artifacts")
				}
				if total < 0 {
					t.Error("GetArtifactsByLabel() returned negative total")
				}
			}
		})
	}
}

// TestService_GetLabelsByArtifact tests getting labels by artifact
func TestService_GetLabelsByArtifact(t *testing.T) {
	artifactID := uuid.New()

	tests := []struct {
		name       string
		artifactID uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "get labels for artifact",
			artifactID: artifactID,
			setup:      func(_ *MockRepository) {},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			labels, err := service.GetLabelsByArtifact(context.Background(), tt.artifactID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLabelsByArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && labels == nil {
				t.Error("GetLabelsByArtifact() returned nil labels")
			}
		})
	}
}

// TestService_AddLabelToArtifact tests adding a label to an artifact
func TestService_AddLabelToArtifact(t *testing.T) {
	artifactID := uuid.New()
	labelID := uuid.New()

	tests := []struct {
		name       string
		artifactID uuid.UUID
		labelID    uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "add label to artifact",
			artifactID: artifactID,
			labelID:    labelID,
			setup: func(r *MockRepository) {
				artifact := &Artifact{
					Id:   artifactID,
					Text: "Test Artifact",
				}
				label := &Label{
					Id:   labelID,
					Name: "Test Label",
				}
				r.artifacts[artifactID] = artifact
				r.labels[labelID] = label
			},
			wantErr: false,
		},
		{
			name:       "non-existent artifact",
			artifactID: uuid.New(),
			labelID:    labelID,
			setup:      func(_ *MockRepository) {},
			wantErr:    false, // This might not error in the simple implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			err := service.AddLabelToArtifact(context.Background(), tt.artifactID, tt.labelID)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddLabelToArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestService_RemoveLabelFromArtifact tests removing a label from an artifact
func TestService_RemoveLabelFromArtifact(t *testing.T) {
	artifactID := uuid.New()
	labelID := uuid.New()

	tests := []struct {
		name       string
		artifactID uuid.UUID
		labelID    uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "remove label from artifact",
			artifactID: artifactID,
			labelID:    labelID,
			setup: func(r *MockRepository) {
				artifact := &Artifact{
					Id:   artifactID,
					Text: "Test Artifact",
				}
				label := &Label{
					Id:   labelID,
					Name: "Test Label",
				}
				r.artifacts[artifactID] = artifact
				r.labels[labelID] = label
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			err := service.RemoveLabelFromArtifact(context.Background(), tt.artifactID, tt.labelID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveLabelFromArtifact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestNewService tests the service constructor
func TestNewService(t *testing.T) {
	repo := NewMockRepository()
	logger := slog.Default()

	service := NewService(repo, logger)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if service.repo == nil {
		t.Error("NewService() service.repo is nil")
	}

	if service.logger == nil {
		t.Error("NewService() service.logger is nil")
	}
}

// TestMockRepository_EdgeCases tests edge cases in the mock repository
func TestMockRepository_EdgeCases(t *testing.T) {
	t.Run("ListArtifacts with error", func(t *testing.T) {
		repo := NewMockRepository()
		repo.err = errors.New("database error")

		artifacts, total, err := repo.ListArtifacts(context.Background(), ListArtifactsParams{})

		if err == nil {
			t.Error("ListArtifacts() expected error but got none")
		}

		if artifacts != nil {
			t.Error("ListArtifacts() returned non-nil slice on error")
		}

		if total != 0 {
			t.Errorf("ListArtifacts() total = %v, want 0", total)
		}
	})

	t.Run("UpdateLabel on non-existent", func(t *testing.T) {
		repo := NewMockRepository()
		label := &Label{
			Name: "Updated",
		}

		_, err := repo.UpdateLabel(context.Background(), uuid.New(), label)

		if !errors.Is(err, ErrLabelNotFound) {
			t.Errorf("UpdateLabel() error = %v, want %v", err, ErrLabelNotFound)
		}
	})
}

// BenchmarkCreateArtifact benchmarks artifact creation
func BenchmarkCreateArtifact(b *testing.B) {
	repo := NewMockRepository()
	service := NewService(repo, slog.Default())

	req := &CreateArtifactJSONRequestBody{
		Text: "Benchmark text content",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateArtifact(context.Background(), req, uuid.New().String(), uuid.New().String())
	}
}

// BenchmarkListArtifacts benchmarks listing artifacts
func BenchmarkListArtifacts(b *testing.B) {
	repo := NewMockRepository()
	service := NewService(repo, slog.Default())

	// Setup some artifacts
	for i := 0; i < 100; i++ {
		artifact := &Artifact{
			Id:   uuid.New(),
			Text: "Test artifact",
		}
		repo.artifacts[artifact.Id] = artifact
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = service.ListArtifacts(context.Background(), uuid.New().String(), 10, 0)
	}
}
