package resource

import "context"

func (r *ResourceProvider) validateName(ctx context.Context, name string) error {
	//empty name
	if name == "" {
		return ErrEmptyName
	}
	// Check for duplicate name
	if _, err := r.GetByName(ctx, name); err == nil {
		return ErrDuplicateName
	} else if err != nil && err != ErrNameNotFound {
		return err
	}
	return nil
}

func (r *ResourceProvider) validateAction(ctx context.Context, action string) error {
	//empty action
	if action == "" {
		return ErrEmptyAction
	}
	return nil
}

func (r *ResourceProvider) validateId(ctx context.Context, id int) error {
	_, err := r.Get(ctx, id)
	return err
}
