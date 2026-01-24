package role

import (
	"context"
	permissiondomain "marketplace/internal/infra/authz/permission/domain"
	roledomain "marketplace/internal/infra/authz/role/domain"
	roleservice "marketplace/internal/infra/authz/role/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/pagination"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PermissionDTO represents a permission in API responses.
type PermissionDTO struct {
	ID       string `json:"id"`
	Resource string `json:"resource_code"`
	Action   string `json:"action"`
}

// CreateRoleRequest represents the payload required to create a role.
type CreateRoleRequest struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

// UpdateRoleRequest represents the payload required to update a role.
type UpdateRoleRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

// RoleResponse represents a role returned to API clients.
type RoleResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permissions []PermissionDTO `json:"permissions"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// RoleListResponse represents a paginated list of roles.
// It is a specialized view of the generic pagination.PageResult type.
type RoleListResponse = pagination.PageResult[RoleResponse]

// ErrorResponse mirrors the generic error payload used in other modules.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// ToRoleResponse converts a domain role and permissions into a response DTO.
func ToRoleResponse(role *roledomain.Role, permissions []*permissiondomain.Permission) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: toPermissionDTOs(permissions),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// ToRoleListResponse converts a paginated domain result into a DTO response.
// Note: This function requires the handler to load permissions separately for each role.
func (h *RoleHandler) ToRoleListResponse(ctx context.Context, result pagination.PageResult[*roledomain.Role]) (RoleListResponse, error) {
	items := make([]RoleResponse, len(result.Items))
	for i, role := range result.Items {
		permissions, err := h.roleService.GetPermissions(ctx, role.ID)
		if err != nil {
			return RoleListResponse{}, err
		}
		items[i] = ToRoleResponse(role, permissions)
	}
	return pagination.PageResult[RoleResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}, nil
}

// toPermissionDTOs converts domain permissions to DTO permissions.
func toPermissionDTOs(perms []*permissiondomain.Permission) []PermissionDTO {
	result := make([]PermissionDTO, len(perms))
	for i, perm := range perms {
		result[i] = PermissionDTO{
			ID:       perm.ID,
			Resource: perm.Resource.Code,
			Action:   perm.Action,
		}
	}
	return result
}

// RoleHandler exposes HTTP endpoints for managing roles.
type RoleHandler struct {
	roleService *roleservice.RoleService
}

// NewRoleHandler constructs a RoleHandler.
func NewRoleHandler(roleService *roleservice.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// CreateRole handles POST /roles requests.
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}
	id := strings.TrimSpace(req.ID)
	if id == "" {
		id = uuid.NewString()
	}

	ctx := r.Context()
	role, err := h.roleService.Create(ctx, id, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		h.writeError(w, h.mapError(err), err)
		return
	}

	permissions, err := h.roleService.GetPermissions(ctx, role.ID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToRoleResponse(role, permissions))
}

// GetRole handles GET /roles/:id requests.
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id := h.extractID(r.URL.Path)
	if id == "" {
		h.writeError(w, http.StatusBadRequest, errMissingIdentifier)
		return
	}

	ctx := r.Context()
	role, err := h.roleService.Get(ctx, id)
	if err != nil {
		if errors.Is(err, roleservice.ErrRoleNotFound) {
			h.writeError(w, http.StatusNotFound, err)
			return
		}
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	permissions, err := h.roleService.GetPermissions(ctx, role.ID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToRoleResponse(role, permissions))
}

// UpdateRole handles PUT /roles/:id requests.
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id := h.extractID(r.URL.Path)
	if id == "" {
		h.writeError(w, http.StatusBadRequest, errMissingIdentifier)
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	role, err := h.roleService.Update(ctx, id, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		h.writeError(w, h.mapError(err), err)
		return
	}

	permissions, err := h.roleService.GetPermissions(ctx, role.ID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ToRoleResponse(role, permissions))
}

// DeleteRole handles DELETE /roles/:id requests.
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	id := h.extractID(r.URL.Path)
	if id == "" {
		h.writeError(w, http.StatusBadRequest, errMissingIdentifier)
		return
	}

	ctx := r.Context()
	if err := h.roleService.Delete(ctx, id); err != nil {
		if errors.Is(err, roleservice.ErrRoleNotFound) {
			h.writeError(w, http.StatusNotFound, err)
			return
		}
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListRoles handles GET /roles requests with pagination.
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	pageReq := pagination.FromRequest(r)

	ctx := r.Context()
	result, err := h.roleService.List(ctx, pageReq)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := h.ToRoleListResponse(ctx, result)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, response)
}

func (h *RoleHandler) writeError(w http.ResponseWriter, status int, err error) {
	httputil.JSON(w, status, ErrorResponse{
		Error:   err.Error(),
		Message: err.Error(),
	})
}

func (h *RoleHandler) mapError(err error) int {
	switch err {
	case roleservice.ErrRoleAlreadyExists, roleservice.ErrRoleNameConflict:
		return http.StatusConflict
	case roleservice.ErrRoleNotFound:
		return http.StatusNotFound
	case roleservice.ErrInvalidPermissionIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (h *RoleHandler) extractID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == ResourceRoles {
		return parts[1]
	}
	return ""
}

var errMissingIdentifier = errors.New("resource identifier is required")
