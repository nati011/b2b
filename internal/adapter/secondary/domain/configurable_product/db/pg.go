package configurable_product

import (
	"context"
	"database/sql"
	"log"

	port "b2b.nati011.github.com/internal/port/domain/configurable_product"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_configurable_products_by_id($1);"

	err := p.Pool.QueryRowContext(ctx, query, id).Scan(
		&response.Id,
		&response.Name,
		&response.Desc,
		&response.ExternalId,
		&response.IsAvailable)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}

	// get images
	var productImages []string
	query = "SELECT * FROM public.get_images_by_cp_Id($1);"
	rows, err := p.Pool.QueryContext(ctx, query, id)
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

	// get attribute-values
	var productAttruteValue = []map[string]string{}
	query = "SELECT * FROM public.get_all_configurable_product_attributes($1)"
	rows, err = p.Pool.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	defer rows.Close()

	for rows.Next() {
		var attributeId int
		if err := rows.Scan(&attributeId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetResponse{}, err
		}

		// get product attribute value
		query = "SELECT * FROM public.get_attributes_values_by_attribute_id($1)"
		rows, err = p.Pool.QueryContext(ctx, query, attributeId)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return port.GetResponse{}, port.ErrSysUnknown
			}
		}

		for rows.Next() {
			var attributeName string
			var attributeValue string
			if err := rows.Scan(&attributeName, &attributeValue); err != nil {
				log.Printf("unable to scan row: %q", err)
				return port.GetResponse{}, err
			}

			productAttruteValue = append(productAttruteValue, map[string]string{
				attributeName: attributeValue,
			})
		}

	}
	response.Attributes = productAttruteValue

	// get member products
	var member_productIds []int
	query = "SELECT * FROM public.get_all_configurable_product_members($1);"
	rows, err = p.Pool.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var productId int
		if err := rows.Scan(&productId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetResponse{}, err
		}
		member_productIds = append(member_productIds, productId)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetResponse{}, port.ErrSysUnknown
	}
	response.Products = member_productIds
	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	// get all Id
	var cp_Ids []int
	query := "SELECT * FROM public.get_all_configurable_products();"
	rows, err := p.Pool.QueryContext(ctx, query)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
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
		cp_Ids = append(cp_Ids, productId)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	// get member products
	var cps []port.GetResponse
	for _, i := range cp_Ids {
		resp, err := p.Get(ctx, i)
		if err != nil {
			switch err {
			default:
				return port.GetAllResponse{}, err
			}
		}
		cps = append(cps, resp)
	}

	if len(cps) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: cps,
	}, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	// get all Id
	var cp_Ids []int
	query := "SELECT * FROM public.get_all_configurable_products_by_name($1);"
	rows, err := p.Pool.QueryContext(ctx, query, name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
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
		cp_Ids = append(cp_Ids, productId)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	// get member products
	var cps []port.GetResponse
	for _, i := range cp_Ids {
		resp, err := p.Get(ctx, i)
		if err != nil {
			switch err {
			default:
				return port.GetAllResponse{}, err
			}
		}
		cps = append(cps, resp)
	}

	if len(cps) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: cps,
	}, nil

}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	// get all Id
	var cp_Ids []int
	query := "SELECT * FROM public.get_all_configurable_products_by_ext_id($1);"
	rows, err := p.Pool.QueryContext(ctx, query, extId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
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
		cp_Ids = append(cp_Ids, productId)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	// get member products
	var cps []port.GetResponse
	for _, i := range cp_Ids {
		resp, err := p.Get(ctx, i)
		if err != nil {
			switch err {
			default:
				return port.GetAllResponse{}, err
			}
		}
		cps = append(cps, resp)
	}

	if len(cps) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: cps,
	}, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var configurable_product_id int
	query := "SELECT * FROM public.create_configurable_product($1, $2, $3);"

	err := p.Pool.QueryRowContext(ctx, query,
		req.Name,
		req.Desc,
		req.ExternalId,
	).Scan(&configurable_product_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	// add images
	for _, i := range req.Images {
		query = "SELECT * FROM public.add_image_to_configurable_product($1, $2, $3);"

		_, err = p.Pool.QueryContext(ctx, query,
			i,
			"",
			configurable_product_id,
		)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}
	// create attributes
	for _, v := range req.AttributeKeys {
		for ak, av := range v {
			print(av)
			// get attribute id by name and productId
			var attribute_ids []int

			query = "SELECT * FROM public.get_attribute_id_by_name($1)"
			rows, err := p.Pool.QueryContext(ctx, query, ak)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
				default:
					return 0, port.ErrSysUnknown
				}
			}
			defer rows.Close()

			for rows.Next() {
				var attribute_id int
				if err := rows.Scan(&attribute_id); err != nil {
					log.Printf("unable to scan row: %q", err)
					return 0, err
				}
				attribute_ids = append(attribute_ids, attribute_id)
			}

			for _, i := range attribute_ids {
				query = "SELECT * FROM public.add_attribute_to_configurable_product($1, $2);"
				_, err = p.Pool.QueryContext(ctx, query,
					i,
					configurable_product_id,
				)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
					default:
						return 0, port.ErrSysUnknown
					}
				}
			}
		}

	}

	// add configurable product members
	for _, i := range req.Products {
		query = "SELECT * FROM public.add_product_to_configurable_product($1, $2);"

		_, err = p.Pool.QueryContext(ctx, query,
			configurable_product_id,
			i,
		)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}

	return configurable_product_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_name($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.Name).Scan(&resourceId)
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
	query := "SELECT * FROM public.update_configurable_product_desc($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.Desc).Scan(&resourceId)
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

func (p *Postgres) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_externalId($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.ExternalId).Scan(&resourceId)
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

func (p *Postgres) UpdateIsAvailableStatus(ctx context.Context, req *port.UpdateIsAvailableStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_isAvailable_status($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.Status).Scan(&resourceId)
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

func (p *Postgres) UpdateProducts(ctx context.Context, req *port.UpdateProductRequest) error {
	//find diff

	//remove all products
	query := "SELECT * FROM public.remove_all_configurable_product_members($1);"

	_, err := p.Pool.QueryContext(ctx, query, req.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	//add product
	for _, i := range req.ProductIds {
		query = "SELECT * FROM public.add_product_to_configurable_product($1, $2);"

		_, err = p.Pool.QueryContext(ctx, query,
			req.Id,
			i,
		)
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

func (p *Postgres) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	//remove images
	query := "SELECT * FROM public.remove_all_configurable_product_images($1);"

	_, err := p.Pool.QueryContext(ctx, query, req.Id)
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
		query := "SELECT * FROM public.add_image_to_configurable_product($1, $2, $3);"

		_, err := p.Pool.QueryContext(ctx, query, i, "", req.Id)
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

func (p *Postgres) UpdateAttributes(ctx context.Context, req *port.UpdateAttributes) error {
	//remove attributes
	query := "SELECT * FROM public.remove_all_configurable_product_attributes($1);"

	_, err := p.Pool.QueryContext(ctx, query, req.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	//attach new attributes
	for _, i := range req.AttributeKeys {
		query := "SELECT * FROM public.add_attribute_to_configurable_product($1, $2);"

		for key := range i {
			_, err := p.Pool.QueryContext(ctx, query, key, req.Id)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return port.ErrSysNoRows
				default:
					return port.ErrSysUnknown
				}
			}
		}

	}
	return nil
}
