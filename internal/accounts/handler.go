package accounts

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AccountService defines the service interface
type AccountService interface {
	Create(ctx context.Context, account *Account) (*Account, error)
	Get(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByProviderID(ctx context.Context, provider string, providerAccountID string) (*Account, error)
	List(ctx context.Context, params ListAccountsParams) ([]*Account, int64, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error)
	Update(ctx context.Context, id uuid.UUID, account *Account) (*Account, error)
	Delete(ctx context.Context, id uuid.UUID) error
	LinkAccount(ctx context.Context, userID uuid.UUID, account *Account) (*Account, error)
	UnlinkAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
}

// Handler handles HTTP requests for accounts
type Handler struct {
	service AccountService
	logger  *slog.Logger
}

// NewHandler creates a new account handler
func NewHandler(service AccountService, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// AccountsFindMany handles GET /auth/accounts
func (h *Handler) AccountsFindMany(ctx echo.Context, params AccountsFindManyParams) error {
	// Convert AccountsFindManyParams to ListAccountsParams
	listParams := ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}

	// Handle pagination
	if params.Page.Number > 0 && params.Page.Size > 0 {
		listParams.Offset = (params.Page.Number - 1) * params.Page.Size
		listParams.Limit = params.Page.Size
	}

	accounts, total, err := h.service.List(ctx.Request().Context(), listParams)
	if err != nil {
		h.logger.Error("failed to list accounts", "error", err)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to retrieve accounts",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":  accounts,
		"total": total,
		"page":  params.Page,
	})
}

// AccountsGetOne handles GET /auth/accounts/{id}
func (h *Handler) AccountsGetOne(ctx echo.Context, id uuid.UUID) error {
	account, err := h.service.Get(ctx.Request().Context(), id)
	if err != nil {
		if err == ErrAccountNotFound {
			return ctx.JSON(http.StatusNotFound, Problem{
				Type:   "not_found",
				Title:  "Account Not Found",
				Status: http.StatusNotFound,
				Detail: "The requested account does not exist",
			})
		}
		h.logger.Error("failed to get account", "error", err, "id", id)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to retrieve account",
		})
	}

	return ctx.JSON(http.StatusOK, account)
}

// AccountsDelete handles DELETE /auth/accounts/{id}
func (h *Handler) AccountsDelete(ctx echo.Context, id uuid.UUID) error {
	err := h.service.Delete(ctx.Request().Context(), id)
	if err != nil {
		if err == ErrAccountNotFound {
			return ctx.JSON(http.StatusNotFound, Problem{
				Type:   "not_found",
				Title:  "Account Not Found",
				Status: http.StatusNotFound,
				Detail: "The requested account does not exist",
			})
		}
		h.logger.Error("failed to delete account", "error", err, "id", id)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to delete account",
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
