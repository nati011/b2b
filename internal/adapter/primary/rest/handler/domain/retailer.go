package domain

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/user"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

type CreateRetailerRequest struct {
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
	Password  string `json:"password"`
}

type GetRetailerResponse struct {
	Id          int                                 `json:"id"`
	Name        string                              `json:"name"`
	Tin         string                              `json:"tin"`
	Latitude    string                              `json:"latitude"`
	Longitude   string                              `json:"longitude"`
	GeneralZone string                              `json:"general_zone"`
	Region      string                              `json:"region"`
	Woreda      string                              `json:"woreda"`
	Users       application_handler.GetUserResponse `json:"user"`
}

type GetAllRetailerResponse struct {
	List       []GetRetailerResponse `json:"list"`
	TotalCount int64                 `json:"total_count"`
}

type GetRetailerByParamRequest struct {
	Name string `json:"name"`
	Tin  string `json:"tin"`
}

type UpdateRetailerRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Tin  string `json:"tin"`
}

type Retailer struct {
	authMiddleware middleware.Auth
	service        retailer.Provider
	userService    user.Provider
}

func InitRetailer() {
	handler.Register(new(Retailer))

	handler.RegisterResource("/api/v1/retailer")
	handler.RegisterResource("/api/v1/retailer/{id}/user")
}

func (r *Retailer) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	r.service = domainServices.RetailerService
	r.userService = applicationServices.UserService
	r.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (re *Retailer) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/retailer", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.GetHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/retailer", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireNoAuthentication(http.HandlerFunc(re.CreateHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PUT /api/v1/retailer", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.UpdateHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/retailer/{id}/user", func(w http.ResponseWriter, r *http.Request) {
		re.authMiddleware.RequireAuthentication(http.HandlerFunc(re.GetUserHandler)).ServeHTTP(w, r)
	})
}

func (re *Retailer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	resp, err := re.service.Get(r.Context(), typedParamId)
	if err != nil {
		switch err {
		case retailer.ErrIdNotFound:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}

	users_resp, err := re.service.GetAllUsers(r.Context(), resp.Id)
	if err != nil {
		switch err {
		case retailer.ErrIdNotFound:
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"users": users_resp}, http.StatusAccepted)
}

func (re *Retailer) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	const ParamTin = "tin"

	paramValues := r.URL.Query()
	paramNameValue := paramValues.Get(ParamName)
	paramTinValue := paramValues.Get(ParamTin)

	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := re.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case retailer.ErrIdNotFound:
				util.OperationSuccessResponse(w, util.Envelope{"retailer": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		// get all users
		users_resp, err := re.service.GetAllUsers(r.Context(), resp.Id)
		if err != nil {
			switch err {
			case retailer.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		// ASSUMPTION: retailer has one user ERGO users_resp.List[0]
		// there isn't a case where a retailer doesnot have a user agent ERGO users_resp.List[0] cannot throw an exception
		resp_user, err := re.userService.Get(r.Context(), users_resp.List[0])
		if err != nil {
			switch err {
			case user.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}

		if len(users_resp.List) != 0 {
			util.WriteJSON(w, util.Envelope{"retailer": GetRetailerResponse{
				Id:          resp.Id,
				Name:        resp.Name,
				Tin:         resp.Tin,
				Latitude:    resp.Latitude,
				Longitude:   resp.Longitude,
				GeneralZone: resp.GeneralZone,
				Region:      resp.Region,
				Woreda:      resp.Woreda,
				Users:       (application_handler.GetUserResponse)(resp_user),
			}}, http.StatusAccepted)
		}
	} else if paramNameValue != "" || paramTinValue != "" {
		resp, err := re.service.GetByParam(r.Context(), &retailer.GetByParamRequest{
			Name: strings.Trim(paramNameValue, `"`),
			Tin:  strings.Trim(paramTinValue, `"`),
		})
		if err != nil {
			switch err {
			case retailer.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, util.Envelope{"retailer": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		handler_resp := GetAllRetailerResponse{}
		for _, i := range resp.List {
			// get all users
			users_resp, err := re.service.GetAllUsers(r.Context(), i.Id)
			if err != nil {
				switch err {
				case retailer.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}
			// ASSUMPTION: retailer has one user ERGO users_resp.List[0]
			// there isn't a case where a retailer doesnot have a user agent ERGO users_resp.List[0] cannot throw an exception
			resp_user, err := re.userService.Get(r.Context(), users_resp.List[0])
			if err != nil {
				switch err {
				case user.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}
			handler_resp.List = append(handler_resp.List, GetRetailerResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				Users:       (application_handler.GetUserResponse)(resp_user),
			})
		}
		util.WriteJSON(w, util.Envelope{"retailers": handler_resp}, http.StatusAccepted)
	} else {
		resp, err := re.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case retailer.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, util.Envelope{"retailer": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		handler_resp := GetAllRetailerResponse{}
		for _, i := range resp.List {
			// get all users
			users_resp, err := re.service.GetAllUsers(r.Context(), i.Id)
			if err != nil {
				switch err {
				case retailer.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}
			// ASSUMPTION: retailer has one user ERGO users_resp.List[0]
			// there isn't a case where a retailer doesnot have a user agent ERGO users_resp.List[0] cannot throw an exception
			resp_user, err := re.userService.Get(r.Context(), users_resp.List[0])
			if err != nil {
				switch err {
				case user.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, err)
					return
				}
			}
			handler_resp.List = append(handler_resp.List, GetRetailerResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				Users:       (application_handler.GetUserResponse)(resp_user),
			})
		}
		handler_resp.TotalCount = resp.TotalCount
		util.WriteJSON(w, util.Envelope{"retailers": handler_resp}, http.StatusAccepted)
	}
}

func (p *Retailer) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateRetailerRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	requestBody.Username = requestBody.Phone
	id, err := p.service.Create(r.Context(), (*retailer.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case retailer.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"retailer": id})
}

func (re *Retailer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateRetailerRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := re.service.Update(r.Context(), (*retailer.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case retailer.ErrIdNotFound:
			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"retailer": id}, http.StatusAccepted)
}
