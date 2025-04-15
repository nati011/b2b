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

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Name, &response.Desc, &response.ExternalID, &response.IsActive, &response.DistributorId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}

	// get price
	query = "SELECT * FROM public.get_price_by_productId($1);"
	var price float64
	err = p.db.QueryRowContext(ctx, query, id).Scan(&price)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	response.Price = price

	// get images
	var productImages []string
	query = "SELECT * FROM public.get_images_by_productId($1);"
	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var productImage string
		if err := rows.Scan(&productImage); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetResponse{}, err
		}
		productImages = append(productImages, productImage)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetResponse{}, port.ErrSysUnknown
	}
	response.Images = productImages

	// get categories
	var productCategories []int
	query = "SELECT * FROM public.get_categories_by_productId($1);"
	rows, err = p.db.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId int
		if err := rows.Scan(&categoryId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetResponse{}, err
		}
		productCategories = append(productCategories, categoryId)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetResponse{}, port.ErrSysUnknown
	}
	response.CategoryId = productCategories

	// get stock
	query = "SELECT * FROM public.get_stock_by_productId($1);"
	var stock int
	err = p.db.QueryRowContext(ctx, query, id).Scan(&stock)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	response.Stock = stock

	// get attribute-values
	var productAttruteValue = map[string]string{}
	query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
	rows, err = p.db.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	defer rows.Close()

	for rows.Next() {
		var attributeName string
		var attributeValue string
		if err := rows.Scan(&attributeName, &attributeValue); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetResponse{}, err
		}
		productAttruteValue[attributeName] = attributeValue
	}
	response.Attributes = productAttruteValue

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	query := "SELECT * FROM public.get_all_products();"
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
		if err := rows.Scan(&product.Id, &product.Name, &product.Desc, &product.ExternalID, &product.IsActive, &product.DistributorId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}

		// get price
		query = "SELECT * FROM public.get_price_by_productId($1);"
		var price float64
		err = p.db.QueryRowContext(ctx, query, product.Id).Scan(&price)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return port.GetAllResponse{}, port.ErrSysNoRows
			default:
				return port.GetAllResponse{}, port.ErrSysUnknown
			}
		}
		product.Price = price

		// get images
		var productImages []string
		query = "SELECT * FROM public.get_images_by_productId($1);"
		rows, err := p.db.QueryContext(ctx, query, product.Id)
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
			var productImage string
			if err := rows.Scan(&productImage); err != nil {
				log.Printf("unable to scan row: %q", err)
				return port.GetAllResponse{}, err
			}
			productImages = append(productImages, productImage)
		}

		if err := rows.Err(); err != nil {
			log.Printf("error occurred during rows iteration: %q", err)
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
		product.Images = productImages

		// get categories
		var productCategories []int
		query = "SELECT * FROM public.get_categories_by_productId($1);"
		rows, err = p.db.QueryContext(ctx, query, product.Id)
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
			var categoryId int
			if err := rows.Scan(&categoryId); err != nil {
				log.Printf("unable to scan row: %q", err)
				return port.GetAllResponse{}, err
			}
			productCategories = append(productCategories, categoryId)
		}

		if err := rows.Err(); err != nil {
			log.Printf("error occurred during rows iteration: %q", err)
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
		product.CategoryId = productCategories

		// get stock
		query = "SELECT * FROM public.get_stock_by_productId($1);"
		var stock int
		err = p.db.QueryRowContext(ctx, query, product.Id).Scan(&stock)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return port.GetAllResponse{}, port.ErrSysNoRows
			default:
				return port.GetAllResponse{}, port.ErrSysUnknown
			}
		}
		product.Stock = stock

		// get attribute-values
		var productAttruteValue = map[string]string{}
		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		rows, err = p.db.QueryContext(ctx, query, product.Id)
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
			var attributeName string
			var attributeValue string
			if err := rows.Scan(&attributeName, &attributeValue); err != nil {
				log.Printf("unable to scan row: %q", err)
				return port.GetAllResponse{}, err
			}
			productAttruteValue[attributeName] = attributeValue
		}
		product.Attributes = productAttruteValue
		response.List = append(response.List, product)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}
	if len(response.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
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
		var product port.GetResponse
		if err := rows.Scan(&product.Id, &product.Name, &product.Desc, &product.ExternalID, &product.IsActive, &product.DistributorId); err != nil {
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
		var product port.GetResponse
		if err := rows.Scan(&product.Id, &product.Name, &product.Desc, &product.ExternalID, &product.IsActive, &product.DistributorId); err != nil {
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
		var product port.GetResponse
		if err := rows.Scan(&product.Id, &product.Name, &product.Desc, &product.ExternalID, &product.IsActive, &product.DistributorId); err != nil {
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

func (p *Postgres) GetByCategory(ctx context.Context, req *port.GetByCategoryRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	// get product ids
	var productIds []int
	for _, i := range req.CategoryId {
		query := "SELECT * FROM public.get_products_by_categoryId($1);"
		rows, err := p.db.QueryContext(ctx, query, i)
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
			var productId int
			if err := rows.Scan(&productId); err != nil {
				log.Printf("unable to scan row: %q", err)
				return port.GetAllResponse{}, err
			}
			productIds = append(productIds, productId)
		}

		if err := rows.Err(); err != nil {
			log.Printf("error occurred during rows iteration: %q", err)
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
	}

	//get products
	for _, i := range productIds {
		product_resp, err := p.Get(ctx, i)
		if err != nil {
			log.Printf("error occurred during products iteration: %q", err)
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
		response.List = append(response.List, product_resp)
	}
	return response, nil
}

func (p *Postgres) GetByPriceRange(ctx context.Context, req *port.GetByPriceRangeRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	// get product ids
	var productIds []int
	query := "SELECT * FROM public.get_products_by_price_range($1, $2);"
	rows, err := p.db.QueryContext(ctx, query, req.PriceMin, req.PriceMax)
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
		var productId int
		if err := rows.Scan(&productId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		productIds = append(productIds, productId)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	//get products
	for _, i := range productIds {
		product_resp, err := p.Get(ctx, i)
		if err != nil {
			log.Printf("error occurred during products iteration: %q", err)
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
		response.List = append(response.List, product_resp)
	}
	return response, nil
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
	query = "SELECT * FROM public.add_price_to_product($1, $2);"

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
	query = "SELECT * FROM public.create_product_stock($1, $2);"

	_, err = p.db.QueryContext(ctx, query, 0, product_id)
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
		var attributeId int

		query = "SELECT * FROM public.create_product_attribute($1);"

		err = p.db.QueryRowContext(ctx, query, key).Scan(&attributeId)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return 0, port.ErrSysNoRows
			default:
				return 0, port.ErrSysUnknown
			}
		}

		query = "SELECT * FROM public.create_product_attribute_value($1, $2, $3);"

		_, err = p.db.QueryContext(ctx, query, value, product_id, attributeId)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
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
	query := "SELECT * FROM public.update_product_price($1, $2);"

	_, err := p.db.QueryContext(ctx, query, req.Id, req.Price)
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
	_, err = p.db.QueryContext(ctx, query, req.Amount, req.Id, STOCK_OPERATION_GOODS_RECEIVING, 0)
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
	_, err = p.db.QueryContext(ctx, query, req.Amount, req.Id, STOCK_OPERATION_DEPLETION, 0)
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
