package distributor_subscription

import "strings"

func (d *DistributorSubscriptionService) validateCreate(req *CreatePlanRequest) error {
	if err := d.validateName(req.Name); err != nil {
		return err
	}

	if err := d.validateDesc(req.Description); err != nil {
		return err
	}

	if err := d.validatePrice(req.Price); err != nil {
		return err
	}

	if err := d.validateTerm(req.TermInMonth); err != nil {
		return err
	}
	return nil
}

func (d *DistributorSubscriptionService) validateName(name string) error {
	sliced := strings.Split(name, "")
	if len(sliced) == 0 {
		return ErrNameMandatory
	}
	return nil
}

func (d *DistributorSubscriptionService) validateDesc(desc string) error {
	sliced := strings.Split(desc, "")
	if len(sliced) == 0 {
		return ErrNameMandatory
	}
	return nil
}

func (d *DistributorSubscriptionService) validateTerm(term int) error {
	if term == 0 {
		return ErrTermCannotBeZero
	}
	if term < 0 {
		return ErrTermCannotBeNegative
	}
	return nil
}

func (d *DistributorSubscriptionService) validatePrice(price float64) error {
	if price == 0 {
		return ErrPriceCannotBeZero
	}
	if price < 0 {
		return ErrPriceCannotBeNegative
	}
	return nil
}
