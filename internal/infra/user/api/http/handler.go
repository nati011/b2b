package user

import (
	roleDomain "marketplace/internal/infra/authz/role/domain"
	"marketplace/internal/infra/user/domain"
	userservice "marketplace/internal/infra/user/service"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/pagination"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	ExternalID  string `json:"external_id"`
	Email       string `json:"email,omitempty"`        // Optional: at least one of email or phone_number required
	PhoneNumber string `json:"phone_number,omitempty"` // Optional: at least one of email or phone_number required
	Name        string `json:"name"`
	UserType    string `json:"user_type"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	ExternalID *string `json:"external_id,omitempty"` // Optional: can be updated later (nil = preserve, "" = clear, "value" = set)
	Email      string  `json:"email"`
	Name       string  `json:"name"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID                string             `json:"id"`
	ExternalID        string             `json:"external_id"`
	Email             string             `json:"email,omitempty"`
	PhoneNumber       string             `json:"phone_number,omitempty"`
	Name              string             `json:"name"`
	Status            string             `json:"status"`
	UserType          string             `json:"user_type"`
	Roles             []UserRoleResponse `json:"roles"`
	RegistrationToken string             `json:"registration_token,omitempty"` // Only present for self-service users on creation
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

// UserRoleResponse summarizes a role assigned to a user.
type UserRoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UserListResponse represents a paginated list of users.
// It is a specialized view of the generic pagination.PageResult type.
type UserListResponse = pagination.PageResult[UserResponse]

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// UserRoleResponseEnvelope wraps role list responses.
type UserRoleResponseEnvelope struct {
	Roles []UserRoleResponse `json:"roles"`
}

// ToUserResponse converts a domain User to a UserResponse DTO
// roles parameter is optional - if nil, roles will be empty in response
func ToUserResponse(user *domain.User, roles []roleDomain.Role) UserResponse {
	resp := UserResponse{
		ID:         user.ID,
		ExternalID: user.ExternalID,
		Name:       user.Name,
		Status:     string(user.Status),
		UserType:   string(user.UserType),
		Roles:      ToUserRoleResponses(roles),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if !user.Email.IsEmpty() {
		resp.Email = user.Email.String()
	}

	if !user.PhoneNumber.IsEmpty() {
		resp.PhoneNumber = user.PhoneNumber.String()
	}

	return resp
}

// AssignRoleRequest represents the payload to assign a role to a user.
type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

// ToUserRoleResponses converts domain roles to DTOs.
func ToUserRoleResponses(roles []roleDomain.Role) []UserRoleResponse {
	if len(roles) == 0 {
		return []UserRoleResponse{}
	}
	result := make([]UserRoleResponse, len(roles))
	for i, role := range roles {
		result[i] = UserRoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		}
	}
	return result
}

type UserHandler struct {
	userService       *userservice.Service
	permissionChecker userservice.PermissionChecker
}

func NewUserHandler(userService *userservice.Service, permissionChecker userservice.PermissionChecker) *UserHandler {
	return &UserHandler{
		userService:       userService,
		permissionChecker: permissionChecker,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if err := h.validateCreateUserRequest(w, req); err != nil {
		return
	}

	userType := domain.UserTypeFromString(req.UserType)
	if err := h.authorizeOfficerCreation(w, r, userType); err != nil {
		return
	}

	email, phoneNumber, err := h.parseUserIdentifiers(w, r, req)
	if err != nil {
		return
	}

	ctx := r.Context()
	result, err := h.userService.Create(ctx, req.ExternalID, email, phoneNumber, req.Name, userType)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the response
	roles, _ := h.userService.ListRolesForUser(r.Context(), result.User.ID)
	response := ToUserResponse(result.User, roles)
	response.RegistrationToken = result.RegistrationToken
	httputil.JSON(w, http.StatusCreated, response)
}

// validateCreateUserRequest validates the create user request.
func (h *UserHandler) validateCreateUserRequest(w http.ResponseWriter, req CreateUserRequest) error {
	if req.Email == "" && req.PhoneNumber == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("either email or phone_number must be provided"))
		return errors.New("validation failed")
	}
	return nil
}

// authorizeOfficerCreation checks authorization for creating officer users.
func (h *UserHandler) authorizeOfficerCreation(w http.ResponseWriter, r *http.Request, userType domain.UserType) error {
	if userType != domain.UserTypeOfficer {
		return nil
	}

	authenticatedUser := httputil.UserFromContext(r.Context())
	if authenticatedUser == nil {
		httputil.Error(w, http.StatusUnauthorized, errors.New("authentication required to create officer users"))
		return errors.New("authorization failed")
	}

	ctx := r.Context()
	hasPerm, err := h.permissionChecker.HasPermission(ctx, authenticatedUser.ID, ResourceUsers, ActionCreate)
	if err != nil || !hasPerm {
		httputil.Error(w, http.StatusForbidden, errors.New("insufficient permissions to create officer users"))
		return errors.New("authorization failed")
	}

	return nil
}

// parseUserIdentifiers parses email and phone number from the request.
func (h *UserHandler) parseUserIdentifiers(w http.ResponseWriter, r *http.Request, req CreateUserRequest) (domain.Email, domain.PhoneNumber, error) {
	var email domain.Email
	if req.Email != "" {
		var err error
		email, err = domain.NewEmail(req.Email)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return email, domain.PhoneNumber{}, err
		}
	}

	var phoneNumber domain.PhoneNumber
	if req.PhoneNumber != "" {
		validator := func(value string) error {
			return h.userService.ValidatePhoneNumber(r.Context(), value)
		}
		var err error
		phoneNumber, err = domain.NewPhoneNumber(req.PhoneNumber, validator)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return email, phoneNumber, err
		}
	}

	return email, phoneNumber, nil
}

// GetMe handles GET /users/me
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx := r.Context()
	authenticatedUser := httputil.UserFromContext(ctx)
	if authenticatedUser == nil {
		httputil.Error(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	// Load roles for authenticated user
	roles, _ := h.userService.ListRolesForUser(r.Context(), authenticatedUser.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(authenticatedUser, roles))
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	id := h.extractID(r.URL.Path)

	ctx := r.Context()
	u, err := h.userService.Get(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	id := h.extractID(r.URL.Path)

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	email, err := domain.NewEmail(req.Email)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	u, err := h.userService.Update(ctx, id, req.ExternalID, email, req.Name)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	id := h.extractID(r.URL.Path)

	ctx := r.Context()
	err := h.userService.Delete(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	pageReq := pagination.FromRequest(r)

	ctx := r.Context()
	pageResult, err := h.userService.FindPage(ctx, pageReq)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response := UserListResponse{
		Items:      make([]UserResponse, len(pageResult.Items)),
		Total:      pageResult.Total,
		Page:       pageResult.Page,
		Limit:      pageResult.Limit,
		TotalPages: pageResult.TotalPages,
		HasNext:    pageResult.HasNext,
		HasPrev:    pageResult.HasPrev,
	}

	for i, u := range pageResult.Items {
		// Load roles for each user
		roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
		response.Items[i] = ToUserResponse(u, roles)
	}

	httputil.JSON(w, http.StatusOK, response)
}

// ActivateUser handles POST /users/:id/activate
func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	id := h.extractID(r.URL.Path)

	ctx := r.Context()
	u, err := h.userService.Activate(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// DeactivateUser handles POST /users/:id/deactivate
func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	id := h.extractID(r.URL.Path)

	ctx := r.Context()
	u, err := h.userService.Deactivate(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// SuspendUser handles POST /users/:id/suspend
func (h *UserHandler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	id := h.extractID(r.URL.Path)

	ctx := r.Context()
	u, err := h.userService.Suspend(ctx, id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// AssignRole handles POST /users/:id/roles requests.
func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	userID := h.extractID(r.URL.Path)
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("user identifier is required"))
		return
	}

	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(req.RoleID) == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("role_id is required"))
		return
	}

	u, err := h.userService.AssignRole(r.Context(), userID, req.RoleID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// RevokeRole handles DELETE /users/:id/roles/:roleId.
func (h *UserHandler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodDelete) {
		return
	}

	userID := h.extractID(r.URL.Path)
	roleID := h.extractRoleID(r.URL.Path)
	if userID == "" || roleID == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("user and role identifiers are required"))
		return
	}

	u, err := h.userService.RevokeRole(r.Context(), userID, roleID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	// Load roles for the user
	roles, _ := h.userService.ListRolesForUser(r.Context(), u.ID)
	httputil.JSON(w, http.StatusOK, ToUserResponse(u, roles))
}

// ListUserRoles handles GET /users/:id/roles.
func (h *UserHandler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodGet) {
		return
	}

	userID := h.extractID(r.URL.Path)
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("user identifier is required"))
		return
	}

	roles, err := h.userService.ListRolesForUser(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, UserRoleResponseEnvelope{
		Roles: ToUserRoleResponses(roles),
	})
}

func (h *UserHandler) writeError(w http.ResponseWriter, err error) {
	status := h.mapError(err)
	httputil.Error(w, status, err)
}

func (h *UserHandler) mapError(err error) int {
	if errors.Is(err, userservice.ErrUserNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, userservice.ErrEmailAlreadyExists) || errors.Is(err, userservice.ErrPhoneNumberAlreadyExists) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func (h *UserHandler) extractID(path string) string {
	if httputil.ExtractPathSegment(path, 0) != ResourceUsers {
		return ""
	}
	return httputil.ExtractPathSegment(path, 1)
}

func (h *UserHandler) extractRoleID(path string) string {
	if httputil.ExtractPathSegment(path, 0) != ResourceUsers || httputil.ExtractPathSegment(path, 2) != PathSegmentRoles {
		return ""
	}
	return httputil.ExtractPathSegment(path, 3)
}
