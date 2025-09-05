// Package http provides HTTP handlers for workflow operations
package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/archesai/archesai/internal/workflows/domain"
	"github.com/archesai/archesai/internal/workflows/generated/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// Handler handles HTTP requests for workflow operations
type Handler struct {
	service *domain.Service
	logger  *slog.Logger
}

// NewHandler creates a new workflow handler
func NewHandler(service *domain.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers workflow routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Pipeline routes
	g.POST("/pipelines", h.CreatePipeline)
	g.GET("/pipelines", h.FindManyPipelines)
	g.GET("/pipelines/:id", h.FindPipelineByID)
	g.PUT("/pipelines/:id", h.UpdatePipeline)
	g.DELETE("/pipelines/:id", h.DeletePipeline)

	// Run routes
	g.POST("/runs", h.CreateRun)
	g.GET("/runs", h.FindManyRuns)
	g.GET("/runs/:id", h.FindRunByID)
	g.POST("/runs/:id/start", h.StartRun)
	g.POST("/runs/:id/cancel", h.CancelRun)
	g.DELETE("/runs/:id", h.DeleteRun)

	// Tool routes
	g.POST("/tools", h.CreateTool)
	g.GET("/tools", h.FindManyTools)
	g.GET("/tools/:id", h.FindToolByID)
	g.PUT("/tools/:id", h.UpdateTool)
	g.DELETE("/tools/:id", h.DeleteTool)
}

// Pipeline handlers

// CreatePipeline creates a new pipeline
func (h *Handler) CreatePipeline(c echo.Context) error {
	var req domain.CreatePipelineRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get org ID from auth context
	orgID := orgPlaceholder

	pipeline, err := h.service.CreatePipeline(c.Request().Context(), &req, orgID)
	if err != nil {
		h.logger.Error("failed to create pipeline", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create pipeline",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": pipeline.PipelineEntity,
	})
}

// FindManyPipelines retrieves pipelines
func (h *Handler) FindManyPipelines(c echo.Context) error {
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

	pipelines, total, err := h.service.ListPipelines(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list pipelines", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve pipelines",
		})
	}

	data := make([]api.PipelineEntity, len(pipelines))
	for i, pipeline := range pipelines {
		data[i] = pipeline.PipelineEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindPipelineByID retrieves a pipeline by ID
func (h *Handler) FindPipelineByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid pipeline ID",
		})
	}

	pipeline, err := h.service.GetPipeline(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrPipelineNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Pipeline not found",
			})
		}
		h.logger.Error("failed to get pipeline", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve pipeline",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": pipeline.PipelineEntity,
	})
}

// UpdatePipeline updates a pipeline
func (h *Handler) UpdatePipeline(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid pipeline ID",
		})
	}

	var req domain.UpdatePipelineRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	pipeline, err := h.service.UpdatePipeline(c.Request().Context(), id, &req)
	if err != nil {
		if err == domain.ErrPipelineNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Pipeline not found",
			})
		}
		h.logger.Error("failed to update pipeline", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update pipeline",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": pipeline.PipelineEntity,
	})
}

// DeletePipeline deletes a pipeline
func (h *Handler) DeletePipeline(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid pipeline ID",
		})
	}

	err = h.service.DeletePipeline(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrPipelineNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Pipeline not found",
			})
		}
		h.logger.Error("failed to delete pipeline", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete pipeline",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Run handlers

// CreateRun creates a new run
func (h *Handler) CreateRun(c echo.Context) error {
	var req domain.CreateRunRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get org ID from auth context
	orgID := orgPlaceholder

	run, err := h.service.CreateRun(c.Request().Context(), &req, orgID)
	if err != nil {
		h.logger.Error("failed to create run", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create run",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": run.RunEntity,
	})
}

// FindManyRuns retrieves runs
func (h *Handler) FindManyRuns(c echo.Context) error {
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

	// Check if filtering by pipeline
	if pipelineID := c.QueryParam("pipeline_id"); pipelineID != "" {
		runs, total, err := h.service.ListRunsByPipeline(c.Request().Context(), pipelineID, limit, offset)
		if err != nil {
			h.logger.Error("failed to list runs by pipeline", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to retrieve runs",
			})
		}

		data := make([]api.RunEntity, len(runs))
		for i, run := range runs {
			data[i] = run.RunEntity
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"total": total,
			},
		})
	}

	runs, total, err := h.service.ListRuns(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list runs", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve runs",
		})
	}

	data := make([]api.RunEntity, len(runs))
	for i, run := range runs {
		data[i] = run.RunEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindRunByID retrieves a run by ID
func (h *Handler) FindRunByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid run ID",
		})
	}

	run, err := h.service.GetRun(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRunNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Run not found",
			})
		}
		h.logger.Error("failed to get run", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve run",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": run.RunEntity,
	})
}

// StartRun starts a run
func (h *Handler) StartRun(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid run ID",
		})
	}

	run, err := h.service.StartRun(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRunNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Run not found",
			})
		}
		h.logger.Error("failed to start run", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": run.RunEntity,
	})
}

// CancelRun cancels a run
func (h *Handler) CancelRun(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid run ID",
		})
	}

	run, err := h.service.CancelRun(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRunNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Run not found",
			})
		}
		h.logger.Error("failed to cancel run", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": run.RunEntity,
	})
}

// DeleteRun deletes a run
func (h *Handler) DeleteRun(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid run ID",
		})
	}

	err = h.service.DeleteRun(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrRunNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Run not found",
			})
		}
		h.logger.Error("failed to delete run", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete run",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Tool handlers

// CreateTool creates a new tool
func (h *Handler) CreateTool(c echo.Context) error {
	var req domain.CreateToolRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get org ID from auth context
	orgID := orgPlaceholder

	tool, err := h.service.CreateTool(c.Request().Context(), &req, orgID)
	if err != nil {
		h.logger.Error("failed to create tool", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create tool",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": tool.ToolEntity,
	})
}

// FindManyTools retrieves tools
func (h *Handler) FindManyTools(c echo.Context) error {
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

	tools, total, err := h.service.ListTools(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list tools", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve tools",
		})
	}

	data := make([]api.ToolEntity, len(tools))
	for i, tool := range tools {
		data[i] = tool.ToolEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindToolByID retrieves a tool by ID
func (h *Handler) FindToolByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid tool ID",
		})
	}

	tool, err := h.service.GetTool(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrToolNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Tool not found",
			})
		}
		h.logger.Error("failed to get tool", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve tool",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tool.ToolEntity,
	})
}

// UpdateTool updates a tool
func (h *Handler) UpdateTool(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid tool ID",
		})
	}

	var req domain.UpdateToolRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	tool, err := h.service.UpdateTool(c.Request().Context(), id, &req)
	if err != nil {
		if err == domain.ErrToolNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Tool not found",
			})
		}
		h.logger.Error("failed to update tool", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update tool",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": tool.ToolEntity,
	})
}

// DeleteTool deletes a tool
func (h *Handler) DeleteTool(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid tool ID",
		})
	}

	err = h.service.DeleteTool(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrToolNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Tool not found",
			})
		}
		h.logger.Error("failed to delete tool", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete tool",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
