package order

import "context"

func (o *OrderService) validate_retailerId(ctx context.Context, id int) error {
	if id == 0 {
		return ErrRetailerIdNotSupplied
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
	}
	//validate items exist
	//validate provided qty exists
	return nil
}
