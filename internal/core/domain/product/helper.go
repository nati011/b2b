package product

import (
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	port "b2b.nati011.github.com/internal/port/domain/product"
	"github.com/buckket/go-blurhash"
)

func (p *ProductService) validateName(ctx context.Context, name string) error {
	if name == "" {
		return ErrNameNotSupplied
	}
	//duplicate check
	_, err := p.GetByParam(ctx, &GetByParamRequest{
		Name: name,
	})
	if err != ErrEmptyGetContent && err != ErrUnknown {
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

func (p *ProductService) validateDistributor(ctx context.Context, id int) error {
	if id == 0 {
		return ErrDistributorIdMandatory
	}
	dist, err := p.DistrbutorService.Get(ctx, id)
	if err == distributor.ErrIdNotFound {
		return ErrDistributorNotFound
	}
	if !dist.IsActive {
		return ErrDistributorInactive
	}
	return nil
}

func generateBlurHash(images []string) ([]port.Image, error) {
	var imageWithBlurHash []port.Image
	for _, value := range images {
		res, err := http.Get(value)
		if err != nil {
			log.Printf("failed to genrate bluhash err: %v", err)
			return []port.Image{}, ErrUnknown
		}
		imageFile := res.Body
		loadedImage, _, err := image.Decode(imageFile)
		if err != nil {
			log.Printf("failed to genrate bluhash err: %v", err)
			return []port.Image{}, ErrUnknown
		}
		str, err := blurhash.Encode(4, 3, loadedImage)
		if err != nil {
			log.Printf("failed to genrate bluhash err: %v", err)
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
