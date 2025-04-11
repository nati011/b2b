package payment_option

import "b2b.nati011.github.com/internal/core/application/payment_partner"

func paymentResponseMapper(in payment_partner.GetResponse) GetResponse {
	return GetResponse(in)
}

func allPaymentResponseMapper(in payment_partner.GetAllResponse) GetAllResponse {
	var response GetAllResponse
	for _, i := range response.List {
		response.List = append(response.List, GetResponse{
			Id:               i.Id,
			Name:             i.Name,
			Icon:             i.Icon,
			Status:           i.Status,
			Init_payment_url: i.Init_payment_url,
		})
	}
	return GetAllResponse(response)
}
