package http

import (
	"marketplace/internal/infra/authz/permission/domain"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/pagination"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreatePermissionRequest represents the payload to create a permission.
type CreatePermissionRequest struct {
	ID          string `json:"id"`
	Resource    string `json:"resource_code"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// UpdatePermissionRequest represents the payload to update a permission.
type UpdatePermissionRequest struct {
	Resource    string `json:"resource_code"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// PermissionResponse represents a permission returned to clients.
type PermissionResponse struct {
	ID              string    `json:"id"`
	Resource        string    `json:"resource_code"`
	ResourceService string    `json:"resource_service"`
	Action          string    `json:"action"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PermissionListResponse is a paginated collection of permissions.
// It is a specialized view of the generic pagination.PageResult type.
type PermissionListResponse = pagination.PageResult[PermissionResponse]

// ErrorResponse mirrors the error payload shared by other modules.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// ToPermissionResponse converts a domain permission into a DTO.
func ToPermissionResponse(perm *domain.Permission) PermissionResponse {
	return PermissionResponse{
		ID:              perm.ID,
		Resource:        perm.Resource.Code,
		ResourceService: perm.Resource.Service,
		Action:          perm.Action,
		Description:     perm.Description,
		CreatedAt:       perm.CreatedAt,
		UpdatedAt:       perm.UpdatedAt,
	}
}

// ToPermissionListResponse converts a paginated domain result into a DTO.
func ToPermissionListResponse(result pagination.PageResult[*domain.Permission]) PermissionListResponse {
	items := make([]PermissionResponse, len(result.Items))
	for i, perm := range result.Items {
		items[i] = ToPermissionResponse(perm)
	}
	return pagination.PageResult[PermissionResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}
}

// PermissionHandler exposes HTTP endpoints for managing permissions.
type PermissionHandler struct {
	service *permissionservice.PermissionService
}

func NewPermissionHandler(service *permissionservice.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

// CreatePermission handles POST /permissions requests.
func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		id = uuid.NewString()
	}

	ctx := r.Context()
	perm, err := h.service.Create(ctx, id, req.Resource, req.Action, req.Description)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, permissionservice.ErrPermissionAlreadyExists) {
			status = http.StatusConflict
		}
		h.writeError(w, status, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToPermissionResponse(perm))
}

// GetPermission handles GET /permissions/:id.
func (h *PermissionHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id := h.extractID(r.URL.Path)
	if id == "" {
		h.writeError(w, http.StatusBadRequest, errMissingIdentifier)
		return
	}

	ctx := r.Context()
	perm, err := h.service.Get(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToPermissionResponse(perm))
}

// UpdatePermission handles PUT /permissions/:id.
func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id := h.extractID(r.URL.Path)
	if id == "" {
		h.writeError(w, http.StatusBadRequest, errMissingIdentifier)
		return
	}

	var req UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	perm, err := h.service.Update(ctx, id, req.Resource, req.Action, req.Description)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToPermissionResponse(perm))
}

// DeletePermission handles DELETE /permissions/:id.
func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	id := h.extractID(r.URL.Path)
	if id == "" {
		h.writeError(w, http.StatusBadRequest, errMissingIdentifier)
		return
	}

	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPermissions handles GET /permissions.
func (h *PermissionHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	pageReq := pagination.FromRequest(r)

	ctx := r.Context()
	result, err := h.service.List(ctx, pageReq)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToPermissionListResponse(result))
}

func (h *PermissionHandler) writeError(w http.ResponseWriter, status int, err error) {
	httputil.Error(w, status, err)
}

func (h *PermissionHandler) extractID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == ResourcePermissions {
		return parts[1]
	}
	return ""
}

var errMissingIdentifier = errors.New("resource identifier is required")
