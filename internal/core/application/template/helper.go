package template

import "context"

func (t *Template) ValidateName(ctx context.Context, n string) error {
	if n == "" {
		return ErrInvalidName
	}

	_, err := t.getByName(ctx, n)
	if err != ErrNameNotFound {
		return ErrDuplicateName
	}

	return nil
}

func ValidateHTML(h string) error {
	if h == "" {
		return ErrInvalidTemplate
	}
	return nil
}
