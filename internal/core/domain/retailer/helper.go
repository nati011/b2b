package retailer

import (
	"context"
	"regexp"
)

func (r *RetailerService) validateTin(ctx context.Context, tin string) error {
	// TIN must be 10 digits
	if matched, _ := regexp.MatchString(`^\d{10}$`, tin); !matched {
		return ErrInvalidTin
	}

	_, err := r.GetByParam(ctx, &GetByParamRequest{
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
