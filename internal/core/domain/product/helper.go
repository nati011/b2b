package product

import "context"

func (p *ProductService) validateName(ctx context.Context, name string) error {
	if name == "" {
		return ErrNameNotSupplied
	}
	//duplicate check
	_, err := p.GetByParam(ctx, &GetByParamRequest{
		Name: name,
	})
	if err != ErrEmptyGetContent {
		return ErrNameDuplicate
	}
	return nil
}

func validateDesc(desc string) error {
	if desc == "" {
		return ErrDescNotSupplied
	}
	return nil
}

func validateImages(images []string) error {
	count := 0
	for range images {
		count++
	}
	if count < 2 {
		return ErrImagesMustBeAtleastTwo
	}
	return nil
}

func update_validatePrice(price int) error {
	if price < 0 {
		return ErrPriceCannotBeNegative
	}
	return nil
}

func create_validatePrice(price int) error {
	if price == 0 {
		return ErrPriceNotSupplied
	}
	return nil
}

func validateAttributes(attributes map[string]string) error {
	for _, value := range attributes {
		if value == "" {
			return ErrAttributeValuesCannotBeEmpty
		}
	}
	return nil
}
