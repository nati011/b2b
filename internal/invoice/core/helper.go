package invoice

func validateStatus(status string) error {
	if status == "" {
		return ErrSysStatusNotSupplied
	}
	return nil
}

func validateOrderId(id int) error {
	if id == 0 {
		return ErrSysOrderIdNotSupplied
	}
	return nil
}

func validateSubTotal(subTotal float64) error {
	if subTotal == 0 {
		return ErrSysSubTotalMandatory
	}
	return nil
}
