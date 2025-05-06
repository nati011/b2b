package configurable_product

import (
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"b2b.nati011.github.com/internal/core/domain/product"
	port "b2b.nati011.github.com/internal/port/domain/configurable_product"
	"github.com/buckket/go-blurhash"
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

func generateBlurHash(images []string) ([]port.Image, error) {
	var imageWithBlurHash []port.Image
	for _, value := range images {
		res, err := http.Get(value)
		if err != nil {
			return []port.Image{}, ErrUnknown
		}
		imageFile := res.Body
		loadedImage, _, err := image.Decode(imageFile)
		str, _ := blurhash.Encode(4, 3, loadedImage)
		if err != nil {
			return []port.Image{}, ErrUnknown
		}
		image := port.Image{
			ImageUrl: value,
			BlurHash: str,
		}

		log.Printf("Blurhash %v:", str)
		imageWithBlurHash = append(imageWithBlurHash, image)
	}
	return imageWithBlurHash, nil
}
