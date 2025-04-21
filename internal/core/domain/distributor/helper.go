package distributor

import (
	"context"
	"regexp"
)

func (d *DistributorService) validateTin(ctx context.Context, tin string) error {
	// TIN must be 10 digits
	if matched, _ := regexp.MatchString(`^\d{10}$`, tin); !matched {
		return ErrInvalidTin
	}

	_, err := d.GetByParam(ctx, &GetByParamRequest{
		Tin: tin,
	})
	if err != ErrEmptyGetContent {
		switch err {
		case nil:
			return ErrDuplicateTin
		default:
			return err
		}
	}
	return nil
}

func (d *DistributorService) validateDistributor(ctx context.Context, id int) error {
	//distributor with id must exist
	_, err := d.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}
	return nil
}
