package middleware

import (
	"errors"
	"log"
	"net/http"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/distributor_subscription"
)

var (
	ErrUserIdNotFound      = errors.New(" user not found")
	ErrSubscriptionExpired = errors.New(" subscription expired")
	ErrUnknown             = errors.New(" unknown error has occured")
)

const (
	SUPERADMIN_ROLE_NAME = "superadmin"
	ADMIN_ROLE_NAME      = "admin"
)

type DistributorSubscription struct {
	distributor_subscription distributor_subscription.Prodvider
	distributor              distributor.Provider
	userService              user.Provider
	roleService              role.Provider
}

func NewDistributorSubscriptionMiddleware(
	distributor_subscription distributor_subscription.Prodvider,
	distributor distributor.Provider,
	userService user.Provider,
	roleService role.Provider) DistributorSubscription {
	return DistributorSubscription{
		distributor_subscription: distributor_subscription,
		distributor:              distributor,
		userService:              userService,
		roleService:              roleService,
	}
}

func (ds *DistributorSubscription) RequireSubscription(next http.Handler, options ...Option) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get userId from context
		userId, ok := r.Context().Value("userId").(int)
		if !ok {
			log.Printf("user not found in context or is not an integer")
			util.ServerErrorResponse(w, ErrUserIdNotFound)
			return
		}

		// check if user has admin or superadmin role - bypass subscription check
		assignedRoles, err := ds.userService.GetAllAssignedRoles(r.Context(), userId)
		if err != nil {
			switch err {
			case user.ErrNoRoleAssigned:
				// User has no roles, continue with subscription check
			default:
				log.Printf("failed to get assigned roles: %v", err)
				util.RequestErrorResponse(w, ErrUnknown)
				return
			}
		} else {
			// Check if user has admin or superadmin role
			for _, assignedRole := range assignedRoles.List {
				roleDetails, err := ds.roleService.Get(r.Context(), &role.GetRequest{Id: assignedRole.Id})
				if err != nil {
					log.Printf("failed to get role details for roleId %d: %v", assignedRole.Id, err)
					continue
				}
				if roleDetails.Name == SUPERADMIN_ROLE_NAME || roleDetails.Name == ADMIN_ROLE_NAME {
					log.Printf("User %d has %s role, bypassing subscription check", userId, roleDetails.Name)
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// get distributor from userId
		dist, err := ds.distributor.GetByUserId(r.Context(), userId)
		if err != nil {
			log.Printf("failed to check subscription err: %v", err)
			util.RequestErrorResponse(w, ErrUnknown)
			return
		}

		// check subscription status
		sub, err := ds.distributor_subscription.GetSubscriptionByDistributorId(r.Context(), dist.Id)
		if err != nil {
			log.Printf("failed to get subscription by distId err: %v", err)
		}
		if sub.Status != distributor_subscription.STATUS_SUBSCRIPTION_ACTIVE {
			util.RequestErrorResponse(w, ErrSubscriptionExpired)
			return
		}

		next.ServeHTTP(w, r)
	})
}
