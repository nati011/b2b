package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

type DistributorHandler struct {
	distributorService distributor.Provider
}

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
		return
	}
	defer r.Body.Close()

	var requestBody CreateDistributorRequest

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
}
