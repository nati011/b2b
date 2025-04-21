package template

import "context"

func (t *Template) ValidateName(ctx context.Context, n string) error {
	if n == "" {
		return ErrInvalidName
	}
	if _, err := t.GetByName(ctx, n); err != nil {
		switch err {
		case ErrNameNotFound:
			return nil
		default:
			return err
		}
	}
	return ErrDuplicateName
}

func ValidateHTML(h string) error {
	if h == "" {
		return ErrInvalidTemplate
	}
	return nil
}
