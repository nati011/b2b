package product

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/product"
)

type Postgres struct {
	db         *sql.DB
	Pagination *config.Pagination
}

func NewPostgres(DB *sql.DB, pagination *config.Pagination) port.DB {
	return &Postgres{
		db:         DB,
		Pagination: pagination,
	}
}

func (p *Postgres) GetAllStockLedger(ctx context.Context) (port.GetStockLedgerResponse, error) {
	var response port.GetStockLedgerResponse
	var responseBase port.GetStockLedgerBaseResponse
	query := "SELECT * FROM public.get_stock_ledger_entries();"
	dest := []any{
		&responseBase.Id,
		&responseBase.Quantity,
		&responseBase.Product_id,
		&responseBase.Operation,
		&responseBase.CreatedOn,
	}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(nil, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return port.GetStockLedgerResponse{}, port_commons.ErrSysNoRows
		default:
			return port.GetStockLedgerResponse{}, err
		}
	}
	for _, res := range result {
		val := port.GetStockLedgerBaseResponse{
			Id:         int(res[0].(int64)),
			Quantity:   int(res[1].(int64)),
			Product_id: int(res[2].(int64)),
			Operation:  res[3].(string),
			CreatedOn:  res[4].(time.Time),
		}
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) GetStockLedger(ctx context.Context, id int) (port.GetStockLedgerResponse, error) {
	var response port.GetStockLedgerResponse
	var responseBase port.GetStockLedgerBaseResponse
	query := "SELECT * FROM public.get_stock_ledger_entries_by_product_id($1);"
	dest := []any{
		&responseBase.Id,
		&responseBase.Quantity,
		&responseBase.Product_id,
		&responseBase.Operation,
		&responseBase.CreatedOn,
	}
	args := []any{&id}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return port.GetStockLedgerResponse{}, port_commons.ErrSysNoRows
		default:
			return port.GetStockLedgerResponse{}, err
		}
	}
	for _, res := range result {
		val := port.GetStockLedgerBaseResponse{
			Id:         int(res[0].(int64)),
			Quantity:   int(res[1].(int64)),
			Product_id: int(res[2].(int64)),
			Operation:  res[3].(string),
			CreatedOn:  res[4].(time.Time),
		}
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_products_by_id($1);"
	args := []any{&id}
	result := []any{
		&response.Id,
		&response.Name,
		&response.Desc,
		&response.ExternalID,
		&response.IsActive,
		&response.DistributorId,
		&response.Stock,
		&response.AvailableStock,
		&response.ReservedStock,
		&response.Price,
	}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}

	//
	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Desc = *result[2].(*string)
	response.ExternalID = *result[3].(*string)
	response.IsActive = *result[4].(*bool)
	response.DistributorId = *result[5].(*int)
	response.Stock = *result[6].(*int)
	response.AvailableStock = *result[7].(*int)
	response.ReservedStock = *result[8].(*int)
	response.Price = *result[9].(*float64)

	// images
	//--------------------
	var imageResponse []port.Image
	var imageResponseBase port.Image
	query = "SELECT * FROM public.get_images_by_productId($1);"

	productImageArgs := []any{id}
	imagesDest := []any{
		&imageResponseBase.ImageUrl,
		&imageResponseBase.BlurHash,
	}

	imagesResult, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
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
		image := port.Image{
			ImageUrl: i[0].(string),
			BlurHash: i[1].(string),
		}
		imageResponse = append(imageResponse, image)
	}
	response.Images = imageResponse
	//--------------------

	// categories
	//--------------------
	var productCategories []int
	var productCategoryBase int
	query = "SELECT * FROM public.get_categories_by_productId($1);"

	categoryArgs := []any{id}
	categoryDest := []any{
		&productCategoryBase,
	}

	categoryResult, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetResponse{}, err
		}
	}

	for _, i := range categoryResult {
		productCategories = append(productCategories, int(i[0].(int64)))
	}
	response.CategoryId = productCategories

	// attribute-values
	//--------------------
	type avProductReq struct {
		AttributeKey   string
		AttributeValue string
	}
	var productAttruteValue = map[string]string{}
	var attributeValueResponseBase avProductReq

	query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
	avArgs := []any{id}
	avDest := []any{
		&attributeValueResponseBase.AttributeKey,
		&attributeValueResponseBase.AttributeValue,
	}

	avResult, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
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
		productAttruteValue[a[0].(string)] = a[1].(string)

	}
	response.Attributes = productAttruteValue
	//--------------------

	return response, nil
}

func (p *Postgres) Search(ctx context.Context, req *port.SearchRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	var totalCount int64
	query := "SELECT * FROM public.get_all_products_paginated($1, $2, $3, $4, $5);"
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
		&totalCount,
	}
	args := []any{
		p.Pagination.Limit,
		p.Pagination.Offset,
		req.Name,
		req.PriceMin,
		req.PriceMax}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetAllResponse{}, err
		}
	}
	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}

		totalCount = res[10].(int64)
		log.Print(totalCount)
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	response.TotalCount = totalCount
	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	var totalCount int64
	query := "SELECT * FROM public.get_all_products();"
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
	}
	args := []any{}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetAllResponse{}, err
		}
	}
	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	response.TotalCount = totalCount
	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, req *port.GetByNameRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_products_by_name($1);"
	args := []any{&req.Name}
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
	}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetAllResponse{}, err
		}
	}
	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, req *port.GetByExternalIdRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_products_by_externalId($1);"
	args := []any{&req.ExternalId}
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
	}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetAllResponse{}, err
		}
	}

	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) GetByDistributorId(ctx context.Context, req *port.GetByDistributorIdRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_products_by_distributorId($1);"
	args := []any{&req.DistributorId}
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
	}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}
	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) GetByCategory(ctx context.Context, req *port.GetByCategoryRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_products_by_categoryIds($1);"
	args := []any{&req.CategoryId}
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
	}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}
	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) GetByPriceRange(ctx context.Context, req *port.GetByPriceRangeRequest) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	// get product ids
	query := "SELECT * FROM public.get_products_by_price_range($1, $2);"
	args := []any{&req.PriceMin, &req.PriceMax}
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Desc,
		&responseBase.ExternalID,
		&responseBase.IsActive,
		&responseBase.DistributorId,
		&responseBase.Stock,
		&responseBase.AvailableStock,
		&responseBase.ReservedStock,
		&responseBase.Price,
	}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return port.GetAllResponse{}, err
		}
	}
	for _, res := range result {
		v, _ := strconv.ParseFloat(res[9].(string), 64)
		val := port.GetResponse{
			Id:             int(res[0].(int64)),
			Name:           res[1].(string),
			Desc:           res[2].(string),
			ExternalID:     res[3].(string),
			IsActive:       res[4].(bool),
			DistributorId:  int(res[5].(int64)),
			Stock:          int(res[6].(int64)),
			AvailableStock: int(res[7].(int64)),
			ReservedStock:  int(res[8].(int64)),
			Price:          v,
		}
		// images
		//--------------------
		var imageResponse []port.Image
		var imageResponseBase port.Image
		query = "SELECT * FROM public.get_images_by_productId($1);"

		productImageArgs := []any{val.Id}
		imagesDest := []any{
			&imageResponseBase.ImageUrl,
			&imageResponseBase.BlurHash,
		}

		imagesResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			image := port.Image{
				ImageUrl: i[0].(string),
				BlurHash: i[1].(string),
			}
			imageResponse = append(imageResponse, image)
		}
		val.Images = imageResponse

		// categories
		var productCategories []int
		var productCategoryBase int
		query = "SELECT * FROM public.get_categories_by_productId($1);"

		categoryArgs := []any{val.Id}
		categoryDest := []any{
			&productCategoryBase,
		}

		categoryResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithMultiRowResultSet(categoryArgs, categoryDest),
		).DoMultiQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return port.GetAllResponse{}, err
			}
		}

		for _, i := range categoryResult {
			productCategories = append(productCategories, int(i[0].(int64)))
		}
		val.CategoryId = productCategories

		type avProductReq struct {
			AttributeKey   string
			AttributeValue string
		}
		var productAttruteValue = map[string]string{}
		var attributeValueResponseBase avProductReq

		query = "SELECT * FROM public.get_attributes_values_by_productId($1)"
		avArgs := []any{val.Id}
		avDest := []any{
			&attributeValueResponseBase.AttributeKey,
			&attributeValueResponseBase.AttributeValue,
		}

		avResult, err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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
			productAttruteValue[a[0].(string)] = a[1].(string)

		}
		val.Attributes = productAttruteValue
		response.List = append(response.List, val)
	}
	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var product_id int
	query := "SELECT * FROM public.create_product($1, $2, $3, $4, $5);"
	args := []any{
		req.Name,
		req.Desc,
		req.ExternalID,
		req.DistributorId,
		req.Price,
	}
	result := []any{
		&product_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	// create images
	for _, i := range req.Images {
		query := "SELECT * FROM public.add_image_to_product($1, $2, $3);"
		productImageArgs := []any{
			i.ImageUrl,
			i.BlurHash,
			product_id,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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

	// create category
	for _, i := range req.CategoryId {
		query := "SELECT * FROM public.add_category_to_product($1, $2);"

		productImageArgs := []any{
			product_id,
			i,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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

	// product attributes
	for key, value := range req.Attributes {
		var attributeId int

		query = "SELECT * FROM public.create_product_attribute($1);"
		productImageArgs := []any{
			key,
		}
		pAResult := []any{&attributeId}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(productImageArgs, pAResult),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, err
			}
		}

		query = "SELECT * FROM public.create_product_attribute_value($1, $2, $3);"
		productAVArgs := []any{
			value,
			product_id,
			attributeId,
		}
		err = query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(productAVArgs, nil),
		).DoSingleQuery()
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return 0, err
			}
		}
	}
	return product_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_name($1, $2);"
	args := []any{req.Id, req.Name}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateExternalID(ctx context.Context, req *port.UpdateExternalIDRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_external_Id($1, $2);"
	args := []any{req.Id, req.ExternalId}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
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
	query := "SELECT * FROM public.update_product_desc($1, $2);"

	args := []any{req.Id, req.Desc}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateActiveStatus(ctx context.Context, req *port.UpdateActiveStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_status($1, $2);"

	args := []any{req.Id, req.Status}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdatePrice(ctx context.Context, req *port.UpdatePriceRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_product_price($1, $2);"

	args := []any{req.Id, req.Price}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateImages(ctx context.Context, req *port.UpdateImagesRequest) error {
	//remove images
	query := "SELECT * FROM public.remove_all_product_images($1);"

	args := []any{req.Id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//attach new images
	for _, i := range req.Images {
		query := "SELECT * FROM public.add_image_to_product($1, $2, $3);"
		productImageArgs := []any{
			i.ImageUrl,
			i.BlurHash,
			req.Id,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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

func (p *Postgres) UpdateCategoryId(ctx context.Context, req *port.UpdateCategoryIdRequest) error {
	// remove all categories
	query := "SELECT * FROM public.remove_all_categories_from_product($1);"

	args := []any{req.Id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	// attach new categories to product
	for _, i := range req.CategoryId {
		query := "SELECT * FROM public.add_category_to_product($1, $2);"

		productImageArgs := []any{
			req.Id,
			i,
		}
		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.db),
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

const (
	STOCK_OPERATION_GOODS_RECEIVING  = "GOODS_RECEIVING"
	STOCK_OPERATION_DEPLETION        = "DEPLETION"
	STOCK_OPERATION_RESERVE          = "RESERVE"
	STOCK_OPERATION_FREE_RESERVATION = "FREE_RESERVATION"
)

func (p *Postgres) GoodsReceiving(ctx context.Context, req *port.GoodsReceivingRequest) error {
	var stock int
	query := "SELECT * FROM public.get_stock_by_productId($1);"
	productStockArgs := []any{
		req.Id,
	}
	pResult := []any{&stock}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productStockArgs, pResult),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//update stock
	query = "SELECT * FROM public.update_product_stock($1, $2);"

	updatedAmount := stock + req.Amount
	productUpdateStockArgs := []any{
		req.Id,
		updatedAmount,
	}
	pUResult := []any{&stock}
	err = query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productUpdateStockArgs, pUResult),
	).DoSingleQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return err
		}
	}

	//store to stock ledger
	query = "SELECT * FROM public.stock_operation_ledger_entry($1, $2, $3, $4);"
	productUpdateLedgerArgs := []any{
		req.Amount,
		req.Id,
		STOCK_OPERATION_GOODS_RECEIVING,
		0,
	}
	pULesult := []any{&stock}
	err = query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productUpdateLedgerArgs, pULesult),
	).DoSingleQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return err
		}
	}
	return nil
}

func (p *Postgres) Dispatch(ctx context.Context, req *port.DispatchRequest) error {
	var stock int
	query := "SELECT * FROM public.get_stock_by_productId($1);"
	productStockArgs := []any{
		req.Id,
	}
	pResult := []any{&stock}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productStockArgs, pResult),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//update stock
	query = "SELECT * FROM public.update_product_stock($1, $2);"

	updatedAmount := stock - req.Amount
	productUpdateStockArgs := []any{
		req.Id,
		updatedAmount,
	}
	pUResult := []any{&stock}
	err = query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productUpdateStockArgs, pUResult),
	).DoSingleQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return err
		}
	}

	//store to stock ledger
	query = "SELECT * FROM public.stock_operation_ledger_entry($1, $2, $3, $4);"
	productUpdateLedgerArgs := []any{
		req.Amount,
		req.Id,
		STOCK_OPERATION_DEPLETION,
		0,
	}
	pULesult := []any{&stock}
	err = query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productUpdateLedgerArgs, pULesult),
	).DoSingleQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return err
		}
	}
	return nil
}

func (p *Postgres) Reserve(ctx context.Context, req *port.ReserveRequest) error {
	query := "SELECT * FROM public.reserve_product_stock($1, $2);"
	productStockArgs := []any{
		req.Id,
		req.Amount,
	}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productStockArgs, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//store to stock ledger
	query = "SELECT * FROM public.stock_operation_ledger_entry($1, $2, $3, $4);"
	productUpdateLedgerArgs := []any{
		req.Amount,
		req.Id,
		STOCK_OPERATION_RESERVE,
		0,
	}
	err = query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productUpdateLedgerArgs, nil),
	).DoSingleQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return err
		}
	}

	return nil
}

func (p *Postgres) FreeReservation(ctx context.Context, req *port.FreeReservedRequest) error {
	query := "SELECT * FROM public.free_reserved_product_stock($1, $2);"
	productStockArgs := []any{
		req.Id,
		req.Amount,
	}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productStockArgs, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	//store to stock ledger
	query = "SELECT * FROM public.stock_operation_ledger_entry($1, $2, $3, $4);"
	productUpdateLedgerArgs := []any{
		req.Amount,
		req.Id,
		STOCK_OPERATION_FREE_RESERVATION,
		0,
	}
	err = query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(productUpdateLedgerArgs, nil),
	).DoSingleQuery()
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return err
		}
	}
	return nil
}
