package db

import (
	"context"
	"database/sql"
	"log"

	port "b2b.nati011.github.com/internal/port/application/user"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_users_by_id($1);"

	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&response.Id,
		&response.FirstName,
		&response.LastName,
		&response.Email,
		&response.Phone,
		&response.Username,
		&response.DOB,
		&response.IsActive,
		&response.ExternalId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func (p *Postgres) GetByEmail(ctx context.Context, email string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_users_by_email($1);"
	rows, err := p.db.QueryContext(ctx, query, email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var user port.GetResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.DOB, &user.IsActive, &user.ExternalId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByPhone(ctx context.Context, phone string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_users_by_phone($1);"
	rows, err := p.db.QueryContext(ctx, query, phone)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var user port.GetResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.DOB, &user.IsActive, &user.ExternalId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByUsername(ctx context.Context, username string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_users_by_username($1);"
	rows, err := p.db.QueryContext(ctx, query, username)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var user port.GetResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.DOB, &user.IsActive, &user.ExternalId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByActiveStatus(ctx context.Context, status bool, pagination *port.Pagination) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	limit := 10
	offset := pagination.Offset

	if pagination.Limit != 0 {
		limit = pagination.Limit
	}

	query := "SELECT * FROM public.get_users_by_active_status($1,$2,$3);"
	rows, err := p.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var user port.GetResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.DOB, &user.IsActive, &user.ExternalId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_users_by_external_id($1);"
	rows, err := p.db.QueryContext(ctx, query, extId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var user port.GetResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.DOB, &user.IsActive, &user.ExternalId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetUserProvider(ctx context.Context, id int) (port.GetUserProviderResponse, error) {
	var response port.GetUserProviderResponse

	query := "SELECT * FROM public.get_user_provider($1);"
	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetUserProviderResponse{}, port.ErrSysNoRows
		default:
			return port.GetUserProviderResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var userProvider port.UserProvider
		if err := rows.Scan(&userProvider.UserId, &userProvider.ProviderId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetUserProviderResponse{}, err
		}
		response.List = append(response.List, userProvider)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetUserProviderResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context, pagination *port.Pagination) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	limit := 10
	offset := pagination.Offset

	if pagination.Limit != 0 {
		limit = pagination.Limit
	}
	query := "SELECT * FROM public.get_all_users($1,$2);"
	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var user port.GetResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.DOB, &user.IsActive, &user.ExternalId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	//activate by default
	var resourceId int
	query := "SELECT * FROM public.create_user($1, $2, $3, $4, $5, $6, $7);"

	err := p.db.QueryRowContext(ctx, query,
		req.FirstName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
		req.IsActive,
		req.ExternalId,
	).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) CreateAndActivate(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_and_activate_user($1, $2, $3, $4, $5, $6, $7);"

	err := p.db.QueryRowContext(ctx, query,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
		req.ExternalId,
	).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) CreateUserProvider(ctx context.Context, req *port.CreateUserProviderRequest) error {
	var resourceId any
	query := "SELECT * FROM public.create_user_provider($1, $2);"

	err := p.db.QueryRowContext(ctx, query,
		req.UserId,
		req.ProviderId,
	).Scan(&resourceId)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	log.Print(resourceId)

	return nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := "SELECT * FROM public.delete_user($1);"

	err := p.db.QueryRowContext(ctx, query, id).Err()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (p *Postgres) UpdateFirstName(ctx context.Context, req *port.UpdateFirstNameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_FirstName($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.FirstName).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdateEmail(ctx context.Context, req *port.UpdateEmailRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_email($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Email).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdateDOB(ctx context.Context, req *port.UpdateDOBRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_dob($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.DOB).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdateIsActiveStatus(ctx context.Context, req *port.UpdateIsActiveRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_is_active_status($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.IsActive).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdatePhone(ctx context.Context, req *port.UpdatePhoneRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_phone($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Phone).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdateUsername(ctx context.Context, req *port.UpdateUsernameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_name($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Username).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) AssignRole(ctx context.Context, id int, role_id int) error {
	query := "SELECT * FROM public.add_role_to_user($1, $2);"

	_, err := p.db.QueryContext(ctx, query, id, role_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (p *Postgres) RemoveAssignedRole(ctx context.Context, id int, role_id int) error {
	query := "SELECT * FROM public.remove_role_from_user($1, $2);"

	err := p.db.QueryRowContext(ctx, query, id, role_id).Err()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (p *Postgres) GetAllAssignedRole(ctx context.Context, id int) (port.GetAllAssignedRoleResponse, error) {
	var response port.GetAllAssignedRoleResponse

	query := "SELECT * FROM public.get_all_role_by_user($1);"
	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllAssignedRoleResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllAssignedRoleResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var role port.GetAssignedRoleResponse
		if err := rows.Scan(&role.Id); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllAssignedRoleResponse{}, err
		}
		response.List = append(response.List, role)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllAssignedRoleResponse{}, port.ErrSysUnknown
	}

	return response, nil
}
