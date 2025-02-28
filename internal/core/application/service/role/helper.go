package role

import "context"

func (r *RoleProvider) validateName(ctx context.Context, name string) error {
	//empty name
	if name == "" {
		return ErrEmptyName
	}
	//duplicate name
	res, err := r.Get(ctx, &GetRequest{Name: name})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return nil
		default:
			return err
		}
	}
	if res.Id != 0 {
		return ErrDuplicateName
	}
	return nil
}

func (r *RoleProvider) validateId(ctx context.Context, id int) error {
	//check if id is non_zero
	if id == 0 {
		return ErrIdNotFound
	}
	//check if id exists
	resp, err := r.Get(ctx, &GetRequest{Id: id})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return ErrIdNotFound
		default:
			return err
		}
	}
	if resp.Id == 0 {
		return ErrIdNotFound
	}
	return nil
}
