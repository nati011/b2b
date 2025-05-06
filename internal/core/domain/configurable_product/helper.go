package configurable_product

import (
	"context"

	"b2b.nati011.github.com/internal/core/domain/product"
)

func (p *ConfigurableProductService) validateName(ctx context.Context, name string) error {
	if name == "" {
		return ErrNameIsNotSupplied
	}

	//validate uniqueness
	resp, err := p.GetByParam(ctx, &GetByParamRequest{
		Name: name,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
		default:
			return ErrUnknown
		}
	}
	wantLen := 0
	if len(resp.List) != wantLen {
		return ErrNameDuplicate
	}
	return nil
}

func validateDesc(desc string) error {
	if desc == "" {
		return ErrDescIsNotSupplied
	}
	return nil
}

func (p *ConfigurableProductService) validateAttributekeys(ctx context.Context, attributeKeys []string, products []int) error {
	//altease one attribute
	count := 0
	for range attributeKeys {
		count++
	}
	if count < 1 {
		return ErrAttributeKeysMustBeAtleastOne
	}
	//validate attributes exist in all products
	notFound := true
	for _, i := range attributeKeys {
		for _, j := range products {
			resp, err := p.ProductService.Get(ctx, j)
			if err != nil {
				continue
			}
			for k := range resp.Attributes {
				if k == i {
					notFound = false
					break
				}
			}
		}
	}
	if notFound {
		return ErrAttributeKeysDoNotExistInProduct
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

func (p *ConfigurableProductService) validateProducts(ctx context.Context, products []int) error {
	count := 0
	for range products {
		count++
	}
	if count < 1 {
		return ErrProductsMustBeAtleastOne
	}

	//check if products exist
	for _, i := range products {
		_, err := p.ProductService.Get(ctx, i)
		if err != nil {
			switch err {
			case product.ErrIdNotFound:
				return ErrProductNotFound
			}
		}
	}

	return nil
}
