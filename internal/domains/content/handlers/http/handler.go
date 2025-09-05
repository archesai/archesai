// Package http provides HTTP handlers for content management.
package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/archesai/archesai/internal/domains/content/core"
	"github.com/archesai/archesai/internal/domains/content/generated/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// Handler handles HTTP requests for content operations
type Handler struct {
	service *core.Service
	logger  *slog.Logger
}

// NewHandler creates a new content handler
func NewHandler(service *core.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers content routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Artifact routes
	g.POST("/artifacts", h.CreateArtifact)
	g.GET("/artifacts", h.FindManyArtifacts)
	g.GET("/artifacts/:id", h.FindArtifactByID)
	g.PUT("/artifacts/:id", h.UpdateArtifact)
	g.DELETE("/artifacts/:id", h.DeleteArtifact)

	// Label routes
	g.POST("/labels", h.CreateLabel)
	g.GET("/labels", h.FindManyLabels)
	g.GET("/labels/:id", h.FindLabelByID)
	g.PUT("/labels/:id", h.UpdateLabel)
	g.DELETE("/labels/:id", h.DeleteLabel)

	// Artifact-Label relationship routes
	g.POST("/artifacts/:id/labels/:labelId", h.AddLabelToArtifact)
	g.DELETE("/artifacts/:id/labels/:labelId", h.RemoveLabelFromArtifact)
	g.GET("/labels/:id/artifacts", h.GetArtifactsByLabel)
	g.GET("/artifacts/:id/labels", h.GetLabelsByArtifact)
}

// Artifact handlers

// CreateArtifact creates a new artifact
func (h *Handler) CreateArtifact(c echo.Context) error {
	var req core.CreateArtifactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get org ID and producer ID from auth context
	orgID := orgPlaceholder
	producerID := "user-placeholder"

	artifact, err := h.service.CreateArtifact(c.Request().Context(), &req, orgID, producerID)
	if err != nil {
		if err == core.ErrArtifactTooLarge {
			return c.JSON(http.StatusRequestEntityTooLarge, map[string]interface{}{
				"error": "Artifact exceeds maximum size",
			})
		}
		h.logger.Error("failed to create artifact", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create artifact",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": artifact.ArtifactEntity,
	})
}

// FindManyArtifacts retrieves artifacts
func (h *Handler) FindManyArtifacts(c echo.Context) error {
	limit := 50
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	// TODO: Get org ID from auth context
	orgID := orgPlaceholder

	// Check if this is a search query
	if query := c.QueryParam("search"); query != "" {
		artifacts, total, err := h.service.SearchArtifacts(c.Request().Context(), orgID, query, limit, offset)
		if err != nil {
			h.logger.Error("failed to search artifacts", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to search artifacts",
			})
		}

		data := make([]api.ArtifactEntity, len(artifacts))
		for i, artifact := range artifacts {
			data[i] = artifact.ArtifactEntity
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"total": total,
			},
		})
	}

	artifacts, total, err := h.service.ListArtifacts(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list artifacts", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve artifacts",
		})
	}

	data := make([]api.ArtifactEntity, len(artifacts))
	for i, artifact := range artifacts {
		data[i] = artifact.ArtifactEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindArtifactByID retrieves an artifact by ID
func (h *Handler) FindArtifactByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid artifact ID",
		})
	}

	artifact, err := h.service.GetArtifact(c.Request().Context(), id)
	if err != nil {
		if err == core.ErrArtifactNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Artifact not found",
			})
		}
		h.logger.Error("failed to get artifact", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve artifact",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": artifact.ArtifactEntity,
	})
}

// UpdateArtifact updates an artifact
func (h *Handler) UpdateArtifact(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid artifact ID",
		})
	}

	var req core.UpdateArtifactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	artifact, err := h.service.UpdateArtifact(c.Request().Context(), id, &req)
	if err != nil {
		if err == core.ErrArtifactNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Artifact not found",
			})
		}
		if err == core.ErrArtifactTooLarge {
			return c.JSON(http.StatusRequestEntityTooLarge, map[string]interface{}{
				"error": "Artifact exceeds maximum size",
			})
		}
		h.logger.Error("failed to update artifact", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update artifact",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": artifact.ArtifactEntity,
	})
}

// DeleteArtifact deletes an artifact
func (h *Handler) DeleteArtifact(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid artifact ID",
		})
	}

	err = h.service.DeleteArtifact(c.Request().Context(), id)
	if err != nil {
		if err == core.ErrArtifactNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Artifact not found",
			})
		}
		h.logger.Error("failed to delete artifact", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete artifact",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Label handlers - implemented similarly to artifacts
// FindManyLabels retrieves labels

// CreateLabel creates a new label
func (h *Handler) CreateLabel(c echo.Context) error {
	var req core.CreateLabelRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get org ID from auth context
	orgID := orgPlaceholder

	label, err := h.service.CreateLabel(c.Request().Context(), &req, orgID)
	if err != nil {
		if err == core.ErrLabelExists {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": "Label already exists",
			})
		}
		h.logger.Error("failed to create label", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create label",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": label.LabelEntity,
	})
}

// FindManyLabels retrieves labels
func (h *Handler) FindManyLabels(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// FindLabelByID retrieves a label by ID
func (h *Handler) FindLabelByID(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// UpdateLabel updates a label
func (h *Handler) UpdateLabel(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// DeleteLabel deletes a label
func (h *Handler) DeleteLabel(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// AddLabelToArtifact adds a label to an artifact
func (h *Handler) AddLabelToArtifact(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// RemoveLabelFromArtifact removes a label from an artifact
func (h *Handler) RemoveLabelFromArtifact(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// GetArtifactsByLabel retrieves artifacts by label
func (h *Handler) GetArtifactsByLabel(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}

// GetLabelsByArtifact retrieves labels for an artifact
func (h *Handler) GetLabelsByArtifact(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"error": "Not implemented yet",
	})
}
