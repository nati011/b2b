package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	application_handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/user"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/retailer"
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
	UserId    int    `json:"user_id"`
}

type GetDistributorResponse struct {
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

type GetAllDistributorResponse struct {
	List []GetResponse `json:"list"`
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
}

func InitDistributor() {
	handler.Register(new(Distributor))
}

func (r *Distributor) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = applicationServices.DistributorService
	r.userService = applicationServices.UserService
	return nil
}

func (d *Distributor) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/distributor", d.GetDistributorHandler)
	mux.HandleFunc("POST /api/v1/retailer", d.CreateDistributorHandler)
	mux.HandleFunc("PUT /api/v1/retailer", d.UpdateDistributorHandler)
}

func (re *Distributor) GetDistributorHandler(w http.ResponseWriter, r *http.Request) {
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
			util.RequestErrorResponse(w, r, err)
			return
		}
		resp, err := re.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case retailer.ErrIdNotFound:
				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}
		// get all users
		users_resp, err := re.service.GetAllUsers(r.Context(), resp.Id)
		if err != nil {
			switch err {
			case retailer.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}
		// ASSEMPTION: retailer has one user ERGO users_resp.List[0]
		// there isn't a case where a retailer doesnot have a user agent ERGO users_resp.List[0] cannot throw an exception
		resp_user, err := re.userService.GetByParam(r.Context(), &user.GetByParam{
			ID: users_resp.List[0],
		})
		if err != nil {
			switch err {
			case user.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, r, err)
			}
		}
		if len(users_resp.List) != 0 {
			util.WriteJSON(w, util.Envelope{"retailer": GetResponse{
				Id:          resp.Id,
				Name:        resp.Name,
				Tin:         resp.Tin,
				Latitude:    resp.Latitude,
				Longitude:   resp.Longitude,
				GeneralZone: resp.GeneralZone,
				Region:      resp.Region,
				Woreda:      resp.Woreda,
				Users:       (application_handler.GetUserResponse)(resp_user.List[0]),
			}}, http.StatusAccepted)
		}
	} else if paramNameValue != "" || paramTinValue != "" {
		resp, err := re.service.GetByParam(r.Context(), &distributor.GetByParamRequest{
			Name: paramNameValue,
			Tin:  paramTinValue,
		})
		if err != nil {
			switch err {
			case retailer.ErrEmptyGetContent:

				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}
		handler_resp := GetAllResponse{}
		for _, i := range resp.List {
			// get all users
			users_resp, err := re.service.GetAllUsers(r.Context(), i.Id)
			if err != nil {
				switch err {
				case retailer.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, r, err)
				}
			}
			// ASSEMPTION: retailer has one user ERGO users_resp.List[0]
			// there isn't a case where a retailer doesnot have a user agent ERGO users_resp.List[0] cannot throw an exception
			resp_user, err := re.userService.GetByParam(r.Context(), &user.GetByParam{
				ID: users_resp.List[0],
			})
			if err != nil {
				switch err {
				case user.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, r, err)
				}
			}
			handler_resp.List = append(handler_resp.List, GetResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				Users:       (application_handler.GetUserResponse)(resp_user.List[0]),
			})
		}
		util.WriteJSON(w, util.Envelope{"retailers": handler_resp}, http.StatusAccepted)
	} else {
		resp, err := re.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case retailer.ErrEmptyGetContent:

				util.RequestErrorResponse(w, r, err)
				return
			default:
				util.ServerErrorResponse(w, r, err)
				return
			}
		}
		handler_resp := GetAllResponse{}
		for _, i := range resp.List {
			// get all users
			users_resp, err := re.service.GetAllUsers(r.Context(), i.Id)
			if err != nil {
				switch err {
				case retailer.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, r, err)
				}
			}
			// ASSEMPTION: retailer has one user ERGO users_resp.List[0]
			// there isn't a case where a retailer doesnot have a user agent ERGO users_resp.List[0] cannot throw an exception
			resp_user, err := re.userService.GetByParam(r.Context(), &user.GetByParam{
				ID: users_resp.List[0],
			})
			if err != nil {
				switch err {
				case user.ErrIdNotFound:
				default:
					util.ServerErrorResponse(w, r, err)
				}
			}
			handler_resp.List = append(handler_resp.List, GetResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				Users:       (application_handler.GetUserResponse)(resp_user.List[0]),
			})
		}
		util.WriteJSON(w, util.Envelope{"retailers": handler_resp}, http.StatusAccepted)
	}
}

func (p *Distributor) CreateDistributorHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	id, err := p.service.Create(r.Context(), (*distributor.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case retailer.ErrIdNotFound,
			retailer.ErrDuplicateTin,
			retailer.ErrInvalidTin:

			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"retailer": id}, http.StatusAccepted)
}

func (re *Distributor) UpdateDistributorHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	var requestBody UpdateRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, r, err)
		return
	}
	id, err := re.service.Update(r.Context(), (*distributor.UpdateRequest)(&requestBody))
	if err != nil {
		switch err {
		case retailer.ErrIdNotFound:
			util.RequestErrorResponse(w, r, err)
			return
		default:
			util.ServerErrorResponse(w, r, err)
			return
		}
	}
	util.WriteJSON(w, util.Envelope{"retailer": id}, http.StatusAccepted)
}
