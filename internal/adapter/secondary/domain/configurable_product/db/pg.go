package configurable_product

import (
	"context"
	"database/sql"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
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
	args := []any{&id}
	result := []any{
		&response.Id,
		&response.Name,
		&response.Desc,
		&response.ExternalId,
		&response.IsAvailable}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}

	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Desc = *result[2].(*string)
	response.ExternalId = *result[3].(*string)
	response.IsAvailable = *result[4].(*bool)

	// images
	//--------------------
	var imageResponse []string
	var imageResponseBase string
	query = "SELECT * FROM public.get_images_by_cp_Id($1);"

	productImageArgs := []any{id}
	imagesDest := []any{
		&imageResponseBase,
	}
	imagesResult, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(productImageArgs, imagesDest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetResponse{}, err
		}
	}
	for _, i := range imagesResult {
		imageResponse = append(imageResponse, i[0].(string))
	}
	response.Images = imageResponse
	//--------------------

	// attribute-values
	//--------------------
	type avProductReq struct {
		AttributeKey   string
		AttributeValue string
	}
	var productAttrutes = []string{}
	var attributeValueResponseBase avProductReq

	query = "SELECT * FROM public.get_all_configurable_product_attributes_values($1)"
	avArgs := []any{id}
	avDest := []any{
		&attributeValueResponseBase.AttributeKey,
		&attributeValueResponseBase.AttributeValue,
	}

	avResult, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(avArgs, avDest),
	).DoMultiQuery()

	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetResponse{}, err
		}
	}
	for _, a := range avResult {
		productAttrutes = append(productAttrutes, a[0].(string))
	}
	response.Attributes = productAttrutes
	//--------------------

	// member products
	//--------------------
	var memberProductIds []int
	var memberProductBase int
	query = "SELECT * FROM public.get_all_configurable_product_members($1);"
	memberProductArgs := []any{id}
	memberProductDest := []any{
		&memberProductBase,
	}

	memberProductResult, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(memberProductArgs, memberProductDest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetResponse{}, err
		}
	}
	for _, i := range memberProductResult {
		memberProductIds = append(memberProductIds, int(i[0].(int64)))
	}
	response.Products = memberProductIds
	//--------------------

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_all_configurable_products();"
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalId,
		&responseBase.IsAvailable}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(nil, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		val := port.GetResponse{
			Id:          *res[0].(*int),
			Name:        *res[1].(*string),
			Desc:        *res[2].(*string),
			ExternalId:  *res[3].(*string),
			IsAvailable: *res[4].(*bool),
		}
		// images
		//--------------------
		var imageResponse []string
		var imageResponseBase string
		query = "SELECT * FROM public.get_images_by_cp_Id($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase,
		}
		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(productImageArgs, imagesDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, i := range imagesResult {
			imageResponse = append(imageResponse, i[0].(string))
		}
		val.Images = imageResponse
		//--------------------

		// attribute-values
		//--------------------
		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttrutes = []string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_all_configurable_product_attributes_values($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(avArgs, avDest),
		).DoMultiQuery()

		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, a := range avResult {
			productAttrutes = append(productAttrutes, a[0].(string))
		}
		val.Attributes = productAttrutes
		//--------------------

		// member products
		//--------------------
		var memberProductIds []int
		var memberProductBase int
		query = "SELECT * FROM public.get_all_configurable_product_members($1);"
		memberProductArgs := []any{val.Id}
		memberProductDest := []any{
			&memberProductBase,
		}

		memberProductResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(memberProductArgs, memberProductDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, i := range memberProductResult {
			memberProductIds = append(memberProductIds, int(i[0].(int64)))
		}
		val.Products = memberProductIds
		response.List = append(response.List, val)
	}

	if len(response.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_all_configurable_products_by_name($1);"
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalId,
		&responseBase.IsAvailable}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(nil, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		val := port.GetResponse{
			Id:          *res[0].(*int),
			Name:        *res[1].(*string),
			Desc:        *res[2].(*string),
			ExternalId:  *res[3].(*string),
			IsAvailable: *res[4].(*bool),
		}
		// images
		//--------------------
		var imageResponse []string
		var imageResponseBase string
		query = "SELECT * FROM public.get_images_by_cp_Id($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase,
		}
		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(productImageArgs, imagesDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, i := range imagesResult {
			imageResponse = append(imageResponse, i[0].(string))
		}
		val.Images = imageResponse
		//--------------------

		// attribute-values
		//--------------------
		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttrutes = []string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_all_configurable_product_attributes_values($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(avArgs, avDest),
		).DoMultiQuery()

		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, a := range avResult {
			productAttrutes = append(productAttrutes, a[0].(string))
		}
		val.Attributes = productAttrutes
		//--------------------

		// member products
		//--------------------
		var memberProductIds []int
		var memberProductBase int
		query = "SELECT * FROM public.get_all_configurable_product_members($1);"
		memberProductArgs := []any{val.Id}
		memberProductDest := []any{
			&memberProductBase,
		}

		memberProductResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(memberProductArgs, memberProductDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, i := range memberProductResult {
			memberProductIds = append(memberProductIds, int(i[0].(int64)))
		}
		val.Products = memberProductIds
		response.List = append(response.List, val)
	}

	if len(response.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_all_configurable_products_by_ext_id($1);"
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalId,
		&responseBase.IsAvailable}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(nil, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		val := port.GetResponse{
			Id:          *res[0].(*int),
			Name:        *res[1].(*string),
			Desc:        *res[2].(*string),
			ExternalId:  *res[3].(*string),
			IsAvailable: *res[4].(*bool),
		}
		// images
		//--------------------
		var imageResponse []string
		var imageResponseBase string
		query = "SELECT * FROM public.get_images_by_cp_Id($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase,
		}
		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(productImageArgs, imagesDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, i := range imagesResult {
			imageResponse = append(imageResponse, i[0].(string))
		}
		val.Images = imageResponse
		//--------------------

		// attribute-values
		//--------------------
		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttrutes = []string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_all_configurable_product_attributes_values($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(avArgs, avDest),
		).DoMultiQuery()

		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, a := range avResult {
			productAttrutes = append(productAttrutes, a[0].(string))
		}
		val.Attributes = productAttrutes
		//--------------------

		// member products
		//--------------------
		var memberProductIds []int
		var memberProductBase int
		query = "SELECT * FROM public.get_all_configurable_product_members($1);"
		memberProductArgs := []any{val.Id}
		memberProductDest := []any{
			&memberProductBase,
		}

		memberProductResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(memberProductArgs, memberProductDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}
		for _, i := range memberProductResult {
			memberProductIds = append(memberProductIds, int(i[0].(int64)))
		}
		val.Products = memberProductIds
		response.List = append(response.List, val)
	}

	if len(response.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var configurable_product_id int
	query := "SELECT * FROM public.create_configurable_product($1, $2, $3);"

	args := []any{
		req.Name,
		req.Desc,
		req.ExternalId}
	result := []any{
		&configurable_product_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	// add images
	for _, i := range req.Images {
		query = "SELECT * FROM public.add_image_to_configurable_product($1, $2, $3);"
		productImageArgs := []any{
			i,
			"",
			configurable_product_id,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(productImageArgs, nil),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, err
			}
		}
	}

	// create attributes
	var attributeIds []int
	for _, v := range req.AttributeKeys {
		// get attribute id by name and productId
		var attributeId int

		query = "SELECT * FROM public.get_attribute_id_by_name($1)"
		productAttrArgs := []any{
			v,
		}
		pAResult := []any{&attributeId}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(productAttrArgs, pAResult),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, err
			}
		}
		attributeIds = append(attributeIds, attributeId)
	}

	for _, i := range attributeIds {
		query = "SELECT * FROM public.add_attribute_to_configurable_product($1, $2);"
		cpAttrArgs := []any{
			i,
			configurable_product_id,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(nil, cpAttrArgs),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, err
			}
		}
	}

	// add configurable product members
	for _, i := range req.Products {
		query = "SELECT * FROM public.add_product_to_configurable_product($1, $2);"
		cpAttrArgs := []any{
			configurable_product_id,
			i,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(cpAttrArgs, nil),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, err
			}
		}
	}

	return configurable_product_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_name($1, $2);"
	args := []any{
		req.Id,
		req.Name}
	result := []any{&resourceId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_desc($1, $2);"
	args := []any{
		req.Id,
		req.Desc}
	result := []any{&resourceId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_externalId($1, $2);"
	args := []any{
		req.Id,
		req.ExternalId}
	result := []any{&resourceId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateIsAvailableStatus(ctx context.Context, req *port.UpdateIsAvailableStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_configurable_product_isAvailable_status($1, $2);"

	args := []any{
		req.Id,
		req.Status}
	result := []any{&resourceId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateProducts(ctx context.Context, req *port.UpdateProductRequest) error {
	//find diff

	//remove all products
	query := "SELECT * FROM public.remove_all_configurable_product_members($1);"
	args := []any{req.Id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//add product
	for _, i := range req.ProductIds {
		query = "SELECT * FROM public.add_product_to_configurable_product($1, $2);"
		productArgs := []any{
			req.Id,
			i,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(productArgs, nil),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return err
			}
		}
	}
	return nil
}

func (p *Postgres) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	//remove images
	query := "SELECT * FROM public.remove_all_configurable_product_images($1);"

	args := []any{req.Id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//attach new images
	for _, i := range req.Images {
		query := "SELECT * FROM public.add_image_to_configurable_product($1, $2, $3);"
		productImageArgs := []any{
			i,
			"",
			req.Id,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(productImageArgs, nil),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return err
			}
		}
	}
	return nil
}

func (p *Postgres) UpdateAttributes(ctx context.Context, req *port.UpdateAttributes) error {
	//remove attributes
	query := "SELECT * FROM public.remove_all_configurable_product_attributes($1);"
	args := []any{req.Id}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//attach new attributes
	for _, i := range req.AttributeKeys {
		query := "SELECT * FROM public.add_attribute_to_configurable_product($1, $2);"

		for key := range i {
			productImageArgs := []any{
				key,
				req.Id,
			}
			err := query_handler.NewQuery(
				query_handler.WithCtx(ctx),
				query_handler.WithDB(p.Pool),
				query_handler.WithQuery(query),
				query_handler.WithSingleRowResultSet(productImageArgs, nil),
			).DoSingleQuery()
			if err != nil {
				switch err {
				case port_commons.ErrSysNoRows:
				default:
					return err
				}
			}
		}

	}
	return nil
}
