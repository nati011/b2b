package handler

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
	"b2b.nati011.github.com/internal/core/application/user"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

var (
	ErrIdNotFound = errors.New("oopsy, distributor id not provided")
	ErrIdNotValid = errors.New("oopsy, distributor id not valid")
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
	util.WriteJSON(w, util.Envelope{"users": users_resp}, http.StatusAccepted)
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
		case user.ErrEmailNotValid,
			user.ErrPhoneNotValid,
			user.ErrPhoneOrEmailMandatory,
			user.ErrFirstNameMandatory:

			util.RequestErrorResponse(w, err)
			return
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"user": id}, http.StatusAccepted)
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
			util.WriteJSON(w, util.Envelope{"distributor": GetDistributorResponse{
				Id:          resp.Id,
				Name:        resp.Name,
				Tin:         resp.Tin,
				Latitude:    resp.Latitude,
				Longitude:   resp.Longitude,
				GeneralZone: resp.GeneralZone,
				Region:      resp.Region,
				Woreda:      resp.Woreda,
				Users:       users_resp.List,
			}}, http.StatusAccepted)
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
		util.WriteJSON(w, util.Envelope{"distributors": handler_resp}, http.StatusAccepted)
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
		util.WriteJSON(w, util.Envelope{"distributors": handler_resp}, http.StatusAccepted)
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
	util.WriteJSON(w, util.Envelope{"distributor": id}, http.StatusAccepted)
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
	util.WriteJSON(w, util.Envelope{"distributor": id}, http.StatusAccepted)
}
