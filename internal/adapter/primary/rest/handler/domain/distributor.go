package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/user"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

var (
	ErrIdNotFound = errors.New(" distributor id not provided")
	ErrIdNotValid = errors.New(" distributor id not valid")
)

type CreateDistributorRequest struct {
	Tin         string `json:"tin"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
}

type CreateDistributorUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetDistributorResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Tin         string `json:"tin"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
	IsActive    bool   `json:"is_active"`
	Users       []int  `json:"user"`
	Verdict     string `json:"verdict"`
}

type GetAllDistributorResponse struct {
	List       []GetDistributorResponse `json:"distributors"`
	TotalCount int64                    `json:"total_count"`
}

type GetDistributorByParamRequest struct {
	Name    string `json:"name"`
	Tin     string `json:"tin"`
	Status  string `json:"status"`
	Verdict string `json:"verdict"`
}

type UpdateDistributorRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Tin  string `json:"tin"`
}

type Distributor struct {
	service        distributor.Provider
	userService    user.Provider
	authMiddleware middleware.Auth
}

func InitDistributor() {
	handler.Register(new(Distributor))

	handler.RegisterResource("/api/v1/distributor")
	handler.RegisterResource("/api/v1/distributor/{param}/user")
	handler.RegisterResource("/api/v1/distributor/{param}/status")
	handler.RegisterResource("/api/v1/distributor/{param}/onboarding_review")
}

func (d *Distributor) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	d.service = domainServices.DistributorService
	d.userService = applicationServices.UserService
	d.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (d *Distributor) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/distributor", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.GetDistributorHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/distributor", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.CreateDistributorHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PUT /api/v1/distributor", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.UpdateDistributorHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/distributor/{id}/user", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.CreateUserHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/distributor/user", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.GetDistributorUsersFromContext)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/distributor/{id}/user", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.GetUserHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/distributor/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.StatusHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/distributor/{id}/onboarding_review", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.OnboardingApprovalHandler)).ServeHTTP(w, r)
	})
}

func (de *Distributor) GetDistributorUsersFromContext(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("userId").(int)

	resp, err := de.service.GetByUserId(r.Context(), id)
	if err != nil {
		switch err {
		case distributor.ErrIdNotFound:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}

	users_resp, err := de.service.GetAllUserDetail(r.Context(), resp.Id)
	if err != nil {
		switch err {
		case distributor.ErrIdNotFound:
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"users": users_resp})

}

func (de *Distributor) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	resp, err := de.service.GetByUserId(r.Context(), typedParamId)
	if err != nil {
		switch err {
		case distributor.ErrIdNotFound:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}

	users_resp, err := de.service.GetAllUserDetail(r.Context(), resp.Id)
	if err != nil {
		switch err {
		case distributor.ErrIdNotFound:
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"users": users_resp})
}

func (de *Distributor) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateDistributorUserRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := de.service.CreateUser(r.Context(), &distributor.CreateUserRequest{
		DistributorId: typedParamId,
		FirstName:     requestBody.FirstName,
		LastName:      requestBody.LastName,
		Username:      requestBody.Phone,
		Email:         requestBody.Email,
		Phone:         requestBody.Phone,
	})
	if err != nil {
		switch err {
		default:
			util.RequestErrorResponse(w, err)
			return
		case user.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"user": id})
}

func (de *Distributor) GetDistributorHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	const ParamTin = "tin"
	const ParamStatus = "status"
	const ParamApprovalStatus = "verdict"

	paramValues := r.URL.Query()
	paramNameValue := paramValues.Get(ParamName)
	paramTinValue := paramValues.Get(ParamTin)
	paramStatusValue := paramValues.Get(ParamStatus)
	paramApprovalStatusValue := paramValues.Get(ParamApprovalStatus)

	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := de.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case distributor.ErrIdNotFound:
				util.OperationSuccessResponse(w, util.Envelope{"distributor": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		// get all users
		users_resp, err := de.service.GetAllUsers(r.Context(), resp.Id)
		if err != nil {
			switch err {
			case distributor.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}

		if len(users_resp.List) != 0 {
			util.OperationSuccessResponse(w, util.Envelope{"distributor": GetDistributorResponse{
				Id:          resp.Id,
				Name:        resp.Name,
				Tin:         resp.Tin,
				Latitude:    resp.Latitude,
				Longitude:   resp.Longitude,
				GeneralZone: resp.GeneralZone,
				Region:      resp.Region,
				Woreda:      resp.Woreda,
				IsActive:    resp.IsActive,
				Users:       users_resp.List,
				Verdict:     resp.Verdict,
			}})
		}
	} else if paramNameValue != "" || paramTinValue != "" || paramStatusValue != "" || paramApprovalStatusValue != "" {
		resp, err := de.service.GetByParam(r.Context(), &distributor.GetByParamRequest{
			Name:           strings.Trim(paramNameValue, `"`),
			Tin:            strings.Trim(paramTinValue, `"`),
			Status:         strings.Trim(paramStatusValue, `"`),
			ApprovalStatus: strings.Trim(paramApprovalStatusValue, `"`),
		})
		if err != nil {
			switch err {
			case distributor.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, util.Envelope{"distributor": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		handler_resp := GetAllDistributorResponse{}
		for _, i := range resp.List {
			// get all users
			users_resp, err := de.service.GetAllUsers(r.Context(), i.Id)
			if err != nil {
				switch err {
				case distributor.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}

			handler_resp.List = append(handler_resp.List, GetDistributorResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
				Verdict:     i.Verdict,
				Users:       users_resp.List,
			})

		}
		handler_resp.TotalCount = resp.TotalCount
		util.OperationSuccessResponse(w, handler_resp)
	} else {
		resp, err := de.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case distributor.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, util.Envelope{"distributor": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		handler_resp := GetAllDistributorResponse{}
		for _, i := range resp.List {
			// get all users
			users_resp, err := de.service.GetAllUsers(r.Context(), i.Id)
			if err != nil {
				switch err {
				case distributor.ErrEmptyGetContent:
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}

			handler_resp.List = append(handler_resp.List, GetDistributorResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
				Users:       users_resp.List,
				Verdict:     i.Verdict,
			})
		}
		handler_resp.TotalCount = resp.TotalCount
		util.OperationSuccessResponse(w, handler_resp)
	}
}

func (de *Distributor) CreateDistributorHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateDistributorRequest

	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	requestBody.Username = requestBody.Phone
	id, err := de.service.Create(r.Context(), (*distributor.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case distributor.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"distributor": id})
}

func (de *Distributor) UpdateDistributorHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateDistributorRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := de.service.Update(r.Context(), (*distributor.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case distributor.ErrIdNotFound:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"distributor": id})

}

const (
	ACTIVATE_DISTRIBUTOR_COMMAND   = "activate"
	DEACTIVATE_DISTRIBUTOR_COMMAND = "deactivate"
)

func (d *Distributor) StatusHandler(w http.ResponseWriter, r *http.Request) {
	const ParamCommand = "command"
	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)

	if paramCommandValue != "" {
		typedParamId, err := util.GetPathParam(r, 4)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case ACTIVATE_DISTRIBUTOR_COMMAND:
			err = d.service.Activate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case distributor.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"detail": "distributor successfully activated"})
		case DEACTIVATE_DISTRIBUTOR_COMMAND:
			err = d.service.Dectivate(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case distributor.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"detail": "distributor successfully deactivated"})
		default:
			util.RequestErrorResponse(w, ErrUnknownProductCommand)
		}
	}
}

const (
	APPROVE_DISTRIBUTOR_COMMAND = "approve"
	REJECT_DISTRIBUTOR_COMMAND  = "reject"
)

func (d *Distributor) OnboardingApprovalHandler(w http.ResponseWriter, r *http.Request) {
	const ParamCommand = "command"
	const ParamComment = "comment"
	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	paramCommentValue := paramValues.Get(ParamComment)

	if paramCommandValue != "" {
		typedParamId, err := util.GetPathParam(r, 4)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case APPROVE_DISTRIBUTOR_COMMAND:
			err = d.service.ApproveOnboardingRequest(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case distributor.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"detail": "distributor successfully approved"})
		case REJECT_DISTRIBUTOR_COMMAND:
			err = d.service.RejectOnboardingRequest(r.Context(), typedParamId, paramCommentValue)
			if err != nil {
				switch err {
				case distributor.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
			util.OperationSuccessResponse(w, util.Envelope{"detail": "distributor successfully rejected"})
		default:
			util.RequestErrorResponse(w, ErrUnknownProductCommand)
		}
	}
}
