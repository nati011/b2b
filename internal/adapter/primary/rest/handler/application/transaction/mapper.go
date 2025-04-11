package transaction

import transaction "b2b.nati011.github.com/internal/core/application/transaction"

func allTransactionMapper(in transaction.GetAllResponse) GetAllResponse {
	var response GetAllResponse
	for _, i := range in.List {
		response.List = append(response.List, GetResponse(i))
	}
	return response
}

func transactionMapper(in transaction.GetResponse) GetResponse {
	return GetResponse(in)
}
