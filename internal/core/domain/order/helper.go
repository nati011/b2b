package order

import (
	"context"

	product "b2b.nati011.github.com/internal/core/domain/product"
)

func (o *OrderService) validate_retailerId(ctx context.Context, id int) error {
	if id == 0 {
		return ErrRetailerIdNotSupplied
	}
	_, err := o.ProductService.Get(ctx, id)
	if err != nil {
		switch err {
		case product.ErrIdNotFound:
			return ErrRetailerIdNotSupplied
		default:
			return ErrUnknown
		}
	}
	return nil
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
			case ErrIdNotFound:
				return ErrItemMemberProductNotFound
			default:
				return ErrUnknown
			}
		}
		//validate provided qty exists
		if prod_resp.Stock < i.Quantity {
			return ErrItemMemberProductQuantityNotFound
		}
	}

	return nil
}
