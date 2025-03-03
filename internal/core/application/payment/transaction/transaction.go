package transaction

import "time"

type GetTransactionResponse struct {
	Date   time.Time
	Amount int64
}
