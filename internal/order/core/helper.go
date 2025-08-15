package order

import (
	"context"
	"log"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

func (o *OrderService) validate_retailerId(ctx context.Context, id int) error {
	if id == 0 {
		return ErrRetailerIdNotSupplied
	}
	_, err := o.RetailerService.Get(ctx, id)
	if err != nil {
		switch err {
		case retailer.ErrIdNotFound:
			return ErrRetailerIdNotFound
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (o *OrderService) getUserRetailer(ctx context.Context, userId int) (retailer.GetResponse, error) {
	if userId == 0 {
		return retailer.GetResponse{}, ErrRetailerIdNotFound
	}
	retailer_resp, err := o.RetailerService.GetByUserId(ctx, userId)
	if err != nil {
		log.Printf("failed to get retailers %v", err)
		return retailer.GetResponse{}, ErrRetailerIdNotFound
	}

	return retailer_resp, nil
}

func (o *OrderService) getUserDistributor(ctx context.Context, userId int) (distributor.GetResponse, error) {
	if userId == 0 {
		return distributor.GetResponse{}, ErrRetailerIdNotFound
	}
	distributor_resp, err := o.DistributorService.GetByUserId(ctx, userId)
	if err != nil {
		log.Printf("failed to get distributor %v", err)
		return distributor.GetResponse{}, ErrRetailerIdNotFound
	}

	return distributor_resp, nil
}

func (o *OrderService) validate_items(ctx context.Context, items []Item) error {
	//atleast one item in list
	wantItemLen := 1
	itemCount := 0
	for range items {
		itemCount++
	}
	if itemCount < wantItemLen {
		return ErrAtleastOneOrderItemNeeded
	}

	//validate both products and qty mandatory
	for _, i := range items {
		if i.ProductId == 0 || i.Quantity == 0 {
			return ErrItemMemberProductIdOrQuantityEmpty
		}
		//validate items exist
		prod_resp, err := o.ProductService.Get(ctx, i.ProductId)
		if err != nil {
			switch err {
			case product.ErrIdNotFound:
				return ErrItemMemberProductNotFound
			default:
				return ErrUnknown
			}
		}
		//validate provided qty exists
		if prod_resp.AvailableStock < i.Quantity {
			return ErrItemMemberProductQuantityNotFound
		}
	}

	return nil
}
