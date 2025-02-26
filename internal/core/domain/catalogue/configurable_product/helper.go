package configurable_product

import "context"

func (p *ConfigurableProductService) validateName(ctx context.Context, name string) error {
	if name == "" {
		return ErrNameIsNotSupplied
	}

	//validate uniqueness
	_, err := p.GetByParam(ctx, &GetByParamRequest{
		Name: name,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
		default:
			return ErrUnknown
		}
	}
	return nil
}

func validateDesc(desc string) error {
	if desc == "" {
		return ErrDescIsNotSupplied
	}
	return nil
}

func (p *ConfigurableProductService) validateAttributekeys(ctx context.Context, attributeKeys []string) error {
	//altease one attribute
	count := 0
	for range attributeKeys {
		count++
	}
	if count < 1 {
		return ErrAttributeKeysMustBeAtleastOne
	}
	//validate attributes exist in all products
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

func (p *ConfigurableProductService) validateProducts(ctx context.Context, products []int) error {
	count := 0
	for range products {
		count++
	}
	if count < 1 {
		return ErrProductsMustBeAtleastOne
	}
	return nil
}
