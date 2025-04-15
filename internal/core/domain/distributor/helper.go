package distributor

import (
	"context"
	"regexp"
)

<<<<<<< HEAD
func (d *DistributorService) validateTin(ctx context.Context, tin string) error {
=======
func (r *DistributorService) validateTin(ctx context.Context, tin string) error {
>>>>>>> 8f0b9404 (init distributor refactor)
	// TIN must be 10 digits
	if matched, _ := regexp.MatchString(`^\d{10}$`, tin); !matched {
		return ErrInvalidTin
	}

<<<<<<< HEAD
	_, err := d.GetByParam(ctx, &GetByParamRequest{
=======
	_, err := r.GetByParam(ctx, &GetByParamRequest{
>>>>>>> 8f0b9404 (init distributor refactor)
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
<<<<<<< HEAD

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
=======
>>>>>>> 8f0b9404 (init distributor refactor)
