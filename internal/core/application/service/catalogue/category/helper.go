package category

import "context"

func (c *CategoryService) validateName(ctx context.Context, name string) error {
	if name == "" {
		return ErrNameIsNotSupplied
	}

	//duplicate
	resp, err := c.GetAll(ctx)
	if err != ErrEmptyGetContent {
		return err
	}

	for _, i := range resp.List {
		if i.Name == name {
			return ErrDuplicateName
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
