package invoice

func validateStatus(status string) error {
	if status == "" {
		return ErrSysStatusNotSupplied
	}
	return nil
}
