package config

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PaginationBuilder struct {
	pagination *Pagination
}

func NewPaginationBuilder() *PaginationBuilder {
	return DefaultPaginationBuilder()
}

func DefaultPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{
		pagination: &Pagination{
			Limit:  100,
			Offset: 0,
		},
	}
}

func (b *PaginationBuilder) SetLimit(limit int) *PaginationBuilder {
	if limit > 0 {
		b.pagination.Limit = limit
	}
	return b
}

func (b *PaginationBuilder) SetOffset(offset int) *PaginationBuilder {
	if offset > 0 {
		b.pagination.Offset = offset
	}
	return b
}

func (b *PaginationBuilder) Build() *Pagination {
	return b.pagination
}
