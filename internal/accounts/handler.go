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

// ListAccounts handles GET /auth/accounts
func (h *Handler) ListAccounts(ctx echo.Context, params ListAccountsParams) error {
	// Handle pagination
	limit := 10
	offset := 0
	if params.Page.Number > 0 && params.Page.Size > 0 {
		offset = (params.Page.Number - 1) * params.Page.Size
		limit = params.Page.Size
	}

	listParams := ListAccountsParams{
		Page: PageQuery{
			Number: offset/limit + 1,
			Size:   limit,
		},
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

// GetAccount handles GET /auth/accounts/{id}
func (h *Handler) GetAccount(ctx echo.Context, id uuid.UUID) error {
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

// DeleteAccount handles DELETE /auth/accounts/{id}
func (h *Handler) DeleteAccount(ctx echo.Context, id uuid.UUID) error {
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

// CreateAccount handles POST /auth/accounts (registration)
func (h *Handler) CreateAccount(ctx echo.Context) error {
	var req CreateAccountJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err)
		return ctx.JSON(http.StatusBadRequest, Problem{
			Type:   "validation_failed",
			Title:  "Invalid Request",
			Status: http.StatusBadRequest,
			Detail: "Failed to parse request body",
		})
	}

	// Validate password strength
	validator := DefaultPasswordValidator()
	if err := validator.Validate(req.Password); err != nil {
		return ctx.JSON(http.StatusBadRequest, Problem{
			Type:   "validation_failed",
			Title:  "Weak Password",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		h.logger.Error("failed to hash password", "error", err)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to process registration",
		})
	}

	// Create account with local provider
	account := &Account{
		UserId:     uuid.New(), // This should be linked to an actual user
		ProviderId: Local,
		AccountId:  string(req.Email), // Convert types.Email to string
		Password:   hashedPassword,    // Use string, not pointer
	}

	_, err = h.service.Create(ctx.Request().Context(), account)
	if err != nil {
		if err == ErrDuplicateAccount {
			return ctx.JSON(http.StatusConflict, Problem{
				Type:   "account_exists",
				Title:  "Account Already Exists",
				Status: http.StatusConflict,
				Detail: "An account with this email already exists",
			})
		}
		h.logger.Error("failed to create account", "error", err)
		return ctx.JSON(http.StatusInternalServerError, Problem{
			Type:   "internal_error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "Failed to create account",
		})
	}

	// For registration, we should return a token response
	// This is a placeholder - you'll need to implement JWT generation
	tokenResponse := TokenResponse{
		AccessToken:  "placeholder_access_token",
		RefreshToken: "placeholder_refresh_token", // Use string, not pointer
		TokenType:    "bearer",
		ExpiresIn:    3600,
	}

	return ctx.JSON(http.StatusCreated, tokenResponse)
}
