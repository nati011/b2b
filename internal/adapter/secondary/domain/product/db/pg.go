package product

import (
	"context"
	"database/sql"
	"log"

	port "b2b.nati011.github.com/internal/port/domain/product"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		db: DB,
	}
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_products_by_id($1);"

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Desc, &response.ExternalID, &response.IsActive, &response.DistributorId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_invoices();"
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var product port.GetResponse
		if err := rows.Scan(&product.Id, &product.Desc, &product.ExternalID, &product.IsActive, &product.DistributorId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, product)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, req *port.GetByNameRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_products_by_name($1);"
	rows, err := p.db.QueryContext(ctx, query, req.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var invoice port.GetResponse
		if err := rows.Scan(&invoice.Id, &invoice.Desc, &invoice.ExternalID, &invoice.IsActive, &invoice.DistributorId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, req *port.GetByExternalIdRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_products_by_externalId($1);"
	rows, err := p.db.QueryContext(ctx, query, req.ExternalId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var invoice port.GetResponse
		if err := rows.Scan(&invoice.Id, &invoice.Desc, &invoice.ExternalID, &invoice.IsActive, &invoice.DistributorId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByDistributorId(ctx context.Context, req *port.GetByDistributorIdRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_products_by_distributorId($1);"
	rows, err := p.db.QueryContext(ctx, query, req.DistributorId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var invoice port.GetResponse
		if err := rows.Scan(&invoice.Id, &invoice.Desc, &invoice.ExternalID, &invoice.IsActive, &invoice.DistributorId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByCategory(ctx context.Context, req *port.GetByCategoryRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByPriceRange(ctx context.Context, req *port.GetByPriceRangeRequest) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var product_id int
	query := "SELECT * FROM public.create_product($1, $2, $3, $4);"

	err := p.db.QueryRowContext(ctx, query,
		req.Name,
		req.Desc,
		req.ExternalID,
		req.DistributorId,
	).Scan(&product_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	// create images
	for _, i := range req.Images {
		query := "SELECT * FROM public.add_image_to_product($1, $2, $3);"

		_, err := p.db.QueryContext(ctx, query, i, "", product_id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return 0, port.ErrSysNoRows
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}

	// create price
	query = "SELECT * FROM public.update_product_price($1, $2);"

	_, err = p.db.QueryContext(ctx, query, product_id, req.Price)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	// create category
	for _, i := range req.CategoryId {
		query := "SELECT * FROM public.add_category_to_product($1, $2);"

		_, err := p.db.QueryContext(ctx, query, product_id, i)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return 0, port.ErrSysNoRows
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}

	// create stock
	query = "SELECT * FROM public.update_product_stock($1, $2);"

	_, err = p.db.QueryContext(ctx, query, product_id, 0)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	// product attributes
	for key, value := range req.Attributes {
		query = "SELECT * FROM public.create_product_attribute($1);"
		var attributeId int
		err = p.db.QueryRowContext(ctx, query, key).Scan(&attributeId)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return 0, port.ErrSysNoRows
			default:
				return 0, port.ErrSysUnknown
			}
		}

		query = "SELECT * FROM public.create_product_attribute_value($1);"

		_, err = p.db.QueryContext(ctx, query, value, product_id, attributeId)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return 0, port.ErrSysNoRows
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}

	return product_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_name($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Name).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}

func (p *Postgres) UpdateExternalID(ctx context.Context, req *port.UpdateExternalIDRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_external_Id($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.ExternalId).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}

func (p *Postgres) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_desc($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Desc).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}

func (p *Postgres) UpdateActiveStatus(ctx context.Context, req *port.UpdateActiveStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_status($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Status).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}

func (p *Postgres) UpdatePrice(ctx context.Context, req *port.UpdatePriceRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_price($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Price).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}

func (p *Postgres) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	//remove images
	query := "SELECT * FROM public.remove_all_product_images($1);"

	_, err := p.db.QueryContext(ctx, query, req.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	//attach new images
	for _, i := range req.Images {
		query := "SELECT * FROM public.add_image_to_product($1, $2, $3);"

		_, err := p.db.QueryContext(ctx, query, i, "", req.Id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return port.ErrSysNoRows
			default:
				return port.ErrSysUnknown
			}
		}
	}
	return nil
}

func (p *Postgres) UpdateCategoryId(ctx context.Context, req *port.UpdateCategoryIdRequest) error {
	// remove all categories
	query := "SELECT * FROM public.remove_all_categories_from_product($1);"

	_, err := p.db.QueryContext(ctx, query, req.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	// attach new categories to product
	for _, i := range req.CategoryId {
		query := "SELECT * FROM public.add_category_to_product($1, $2);"

		_, err := p.db.QueryContext(ctx, query, req.Id, i)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return port.ErrSysNoRows
			default:
				return port.ErrSysUnknown
			}
		}
	}
	return nil
}

const (
	STOCK_OPERATION_GOODS_RECEIVING = "GOODS_RECEIVING"
	STOCK_OPERATION_DEPLETION       = "DEPLETION"
)

func (p *Postgres) GoodsReceiving(ctx context.Context, req *port.GoodsReceivingRequest) error {
	//get stock
	query := "SELECT * FROM public.get_stock_by_productId($1);"
	var stock int
	err := p.db.QueryRowContext(ctx, query, req.Id).Scan(&stock)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	//update stock
	query = "SELECT * FROM public.update_product_stock($1, $2);"
	updatedAmount := stock + req.Amount
	_, err = p.db.QueryContext(ctx, query, req.Id, updatedAmount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	//store to stock ledger
	query = "SELECT * FROM public.stock_operation_ledger_entry($1, $2, $3, $4);"
	_, err = p.db.QueryContext(ctx, query, req.Amount, req.Id, STOCK_OPERATION_GOODS_RECEIVING, "")
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (p *Postgres) Dispatch(ctx context.Context, req *port.DispatchRequest) error {
	//get stock
	query := "SELECT * FROM public.get_stock_by_productId($1);"
	var stock int
	err := p.db.QueryRowContext(ctx, query, req.Id).Scan(&stock)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	//update stock
	query = "SELECT * FROM public.update_product_stock($1, $2);"
	updatedAmount := stock - req.Amount
	_, err = p.db.QueryContext(ctx, query, req.Id, updatedAmount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	//store to stock ledger
	query = "SELECT * FROM public.stock_operation_ledger_entry($1, $2, $3, $4);"
	_, err = p.db.QueryContext(ctx, query, req.Amount, req.Id, STOCK_OPERATION_DEPLETION, "")
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}
