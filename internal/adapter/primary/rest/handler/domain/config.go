package domain

import (
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/config"
)

type Config struct {
	authMiddleware middleware.Auth
	service        config.Provider
}

func InitConfig() {
	handler.Register(new(Config))

	handler.RegisterResource("/api/v1/config")
}

func (c *Config) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	c.service = domainService.ConfigService
	c.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (c *Config) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/config", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireAuthentication(http.HandlerFunc(c.GetOrderExpiryConfig)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/config", func(w http.ResponseWriter, r *http.Request) {
		c.authMiddleware.RequireAuthentication(http.HandlerFunc(c.SetOrderExpiryConfig)).ServeHTTP(w, r)
	})
}

func (c *Config) GetOrderExpiryConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := c.service.GetOrderExpiryConfig(r.Context())
	if err != nil {
		switch err {
		case config.ErrUnknown:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, resp)
}

func (c *Config) SetOrderExpiryConfig(w http.ResponseWriter, r *http.Request) {
	const ParamOrderExpiryConfig = "order_expiry_duration"

	paramValues := r.URL.Query()
	paramExpiryDurationValue := paramValues.Get(ParamOrderExpiryConfig)
	if paramExpiryDurationValue != "" {
		typedParamExpirationDurationInMinutes, err := strconv.Atoi(paramExpiryDurationValue)
		if err != nil {
			switch err {
			case config.ErrUnknown:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}

		err = c.service.SetOrderExpiryConfig(r.Context(), &config.SetOrderExpiryRequest{
			ExpiryDurationInMinues: typedParamExpirationDurationInMinutes,
		})
		if err != nil {
			switch err {
			case config.ErrUnknown:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
	}
}
