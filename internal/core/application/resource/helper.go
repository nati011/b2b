package resource

import "context"

func (r *ResourceProvider) validateName(ctx context.Context, name string) error {
	//empty name
	if name == "" {
		return ErrEmptyName
	}
	//duplicate name
	resp, err := r.Get(ctx, &GetRequest{Name: name})
	if err != nil {
		return err
	}
	if resp.Id != 0 {
		return ErrDuplicateName
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
	//check if id is non_zero
	if id == 0 {
		return ErrIdNotFound
	}
	//check if id exists
	resp, err := r.Get(ctx, &GetRequest{Id: id})
	if err != nil {
		return err
	}
	if resp.Id == 0 {
		return ErrIdNotFound
	}
	return nil
}
