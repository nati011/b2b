<<<<<<< HEAD
package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
=======
package handler

import (
	"encoding/json"
	"io"
	"net/http"
>>>>>>> 8f0b9404 (init distributor refactor)

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
<<<<<<< HEAD
	"b2b.nati011.github.com/internal/core/application/user"
=======
>>>>>>> 8f0b9404 (init distributor refactor)
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

<<<<<<< HEAD
var (
	ErrIdNotFound = errors.New("oopsy, distributor id not provided")
	ErrIdNotValid = errors.New("oopsy, distributor id not valid")
)
=======
type DistributorHandler struct {
	distributorService distributor.Provider
}
>>>>>>> 8f0b9404 (init distributor refactor)

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
<<<<<<< HEAD
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
	Users       []int  `json:"user"`
}

type GetAllDistributorResponse struct {
	List []GetDistributorResponse `json:"list"`
}

type GetDistributorByParamRequest struct {
	Name string `json:"name"`
	Tin  string `json:"tin"`
}

type UpdateDistributorRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Tin  string `json:"tin"`
}

type Distributor struct {
	service     distributor.Provider
	userService user.Provider
	middleware  util.AuthMiddleware
}

func InitDistributor() {
	handler.Register(new(Distributor))
}

func (d *Distributor) Init(applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	d.service = domainServices.DistributorService
	d.userService = applicationServices.UserService
	d.middleware = *applicationServices.AuthMiddleware
	return nil
}

func (d *Distributor) Routes(mux *http.ServeMux) {
	distributorHandler := http.HandlerFunc(d.GetDistributorHandler)
	mux.Handle("GET /api/v1/distributor", d.middleware.Authenticate(distributorHandler))
	mux.HandleFunc("POST /api/v1/distributor", d.CreateDistributorHandler)
	mux.HandleFunc("PUT /api/v1/distributor", d.UpdateDistributorHandler)

	mux.HandleFunc("POST /api/v1/distributor/{id}/user", d.CreateUserHandler)
	mux.HandleFunc("GET /api/v1/distributor/{id}/user", d.GetUserHandler)
}

func (de *Distributor) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	resp, err := de.service.Get(r.Context(), typedParamId)
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

	users_resp, err := de.service.GetAllUsers(r.Context(), resp.Id)
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
		Distributor_Id: typedParamId,
		FirstName:      requestBody.FirstName,
		LastName:       requestBody.LastName,
		Username:       requestBody.Username,
		Email:          requestBody.Email,
		Phone:          requestBody.Phone,
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
		resp, err := de.service.Get(r.Context(), typedParamId)
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
				Users:       users_resp.List,
			}})
		}
	} else if paramNameValue != "" || paramTinValue != "" {
		resp, err := de.service.GetByParam(r.Context(), &distributor.GetByParamRequest{
			Name: strings.Trim(paramNameValue, `"`),
			Tin:  strings.Trim(paramTinValue, `"`),
		})
		if err != nil {
			switch err {
			case distributor.ErrEmptyGetContent:
				util.RequestErrorResponse(w, err)
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
				Users:       users_resp.List,
			})
		}
		util.OperationSuccessResponse(w, util.Envelope{"distributors": handler_resp})
	} else {
		resp, err := de.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case distributor.ErrEmptyGetContent:

				util.RequestErrorResponse(w, err)
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
				Users:       users_resp.List,
			})
		}
		util.OperationSuccessResponse(w, util.Envelope{"distributors": handler_resp})
	}
}

func (de *Distributor) CreateDistributorHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
=======
	UserId    int    `json:"user_id"`
}

type BusinessLocation struct {
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
}

type UpdateBusinessRequest struct {
	Id            int              `json:"id"`
	DistributorId int              `json:"distributorId"`
	Name          string           `json:"name"`
	Tin           int              `json:"tin"`
	Location      BusinessLocation `json:"location"`
}

type RegisterDistributorResponse struct {
	Id      int    `json:"distributorId"`
	Message string `json:"message"`
}

type GetDistributorByParamRequest struct {
	Id    int
	Name  string
	Email string
}

func InitDistributor() {
	handler.Register(new(DistributorHandler))
}

func (d *DistributorHandler) Init(services *application_core.Container, domainService *domain_core.Container) error {
	d.distributorService = services.DistributorService
	return nil
}

func (d *DistributorHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/distributor/register", d.CreateHandler)
	mux.HandleFunc("GET /api/v1/retailer", d.GetHandler)
	mux.HandleFunc("PUT /api/v1/retailer", d.UpdateHandler)

}

func (h *DistributorHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
>>>>>>> 8f0b9404 (init distributor refactor)
		return
	}
	defer r.Body.Close()

	var requestBody CreateDistributorRequest
<<<<<<< HEAD
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
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
=======

	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}

	registerResponse, err := h.distributorService.Create(r.Context(), (*distributor.RegisterDistributorRequest)(&requestBody))

	if err != nil {
		switch err {
		case distributor.ErrUnknown:
			util.ServerErrorResponse(w, r, err)
			return
		default:
			util.RequestErrorResponse(w, r, err)
		}
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (h *DistributorHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateBusinessRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RequestErrorResponse(w, r, util.ErrInvalidRequestBody)
		return
	}

	registerResponse, err := h.distributorService.Update(r.Context(), (*distributor.UpdateRequest)(&req))

	if err != nil {
		switch err {
		case distributor.ErrUnknown:
			util.ServerErrorResponse(w, r, err)
			return
		default:
			util.RequestErrorResponse(w, r, err)
			return
		}
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (h *DistributorHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateBusinessInformation

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RequestErrorResponse(w, r, util.ErrInvalidRequestBody)
		return
	}

	registerResponse, err := h.distributorService.GetAll(r.Context())

	if err != nil {
		switch err {
		case distributor.ErrUnknown:
			util.ServerErrorResponse(w, r, err)
			return
		default:
			util.RequestErrorResponse(w, r, err)
			return
		}
	}

	json.NewEncoder(w).Encode(registerResponse)
>>>>>>> 8f0b9404 (init distributor refactor)
}
