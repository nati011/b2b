package middleware

import (
	"errors"
	"log"
	"net/http"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/distributor_subscription"
)

var (
	ErrUserIdNotFound      = errors.New(" user not found")
	ErrSubscriptionExpired = errors.New(" subscription expired")
	ErrUnknown             = errors.New(" unknown error has occured")
)

type DistributorSubscription struct {
	distributor_subscription distributor_subscription.Prodvider
	distributor              distributor.Provider
}

func NewDistributorSubscription(
	ds distributor_subscription.Prodvider,
	d distributor.Provider) DistributorSubscription {
	return DistributorSubscription{
		distributor_subscription: ds,
		distributor:              d,
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
