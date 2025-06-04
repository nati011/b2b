package middleware

import (
	"net/http"
	"strconv"

	"b2b.nati011.github.com/config"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
)

type PaginationMiddleware struct {
	pagination *config.Pagination
}

func NewPaginationMiddleware(pagination *config.Pagination) *PaginationMiddleware {
	return &PaginationMiddleware{
		pagination: pagination,
	}
}

func (p *PaginationMiddleware) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const ParamLimit = "limit"
		const ParamOffset = "offset"

		paramValues := r.URL.Query()
		paramLimitValue := paramValues.Get(ParamLimit)
		paramOffsetValue := paramValues.Get(ParamOffset)

		var limit int
		var offset int
		var err error

		// if (paramLimitValue == "") && (paramOffsetValue == "") {
		// 	next.ServeHTTP(w, r)
		// }

		if paramLimitValue != "" {
			limit, err = strconv.Atoi(paramLimitValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}

		if paramOffsetValue != "" {
			offset, err = strconv.Atoi(paramOffsetValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}
		p.pagination.Limit = limit
		p.pagination.Offset = offset

		next.ServeHTTP(w, r)
	})
}
