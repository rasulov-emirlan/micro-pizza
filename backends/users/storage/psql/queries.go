package psql

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/domain"
)

func (r *repository) Create(ctx context.Context, u domain.User) (domain.ID, error) {
	sql, args, err := sq.Insert("users").Columns(
		"full_name", "email", "password", "birth_date",
		"country", "city", "street", "home_number",
		"address_instructions", "created_at",
	).Values(u.FullName, u.Email, u.Password, u.BirthDate,
		u.Address.Country, u.Address.City, u.Address.Street,
		u.Address.HomeNumber, u.Address.AddressInstructions,
	).Suffix("RETURNING \"id\"").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	id := domain.ID(0)
	if err := conn.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return id, err
	}
	return id, err
}

func (r *repository) Read(ctx context.Context, userID domain.ID) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date",
		"country", "city", "street", "home_number", "address_instructions",
		"ARRAY_AGGG(users_roles.role_id)",
		"created_at",
		"updated_at",
	).From("users").
		LeftJoin("users_roles").
		Where("name = ?", userID).
		GroupBy("id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.User{}, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer conn.Release()

	u := domain.User{}
	roles := []int{}
	if err := conn.QueryRow(ctx, sql, args).Scan(
		&u.ID, &u.FullName, &u.Email, &u.Password, &u.BirthDate,
		&u.Address.Country, &u.Address.City, &u.Address.Street, &u.Address.HomeNumber,
		&u.Address.AddressInstructions, &roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}
	return u, nil
}

func (r *repository) ReadByName(ctx context.Context, fullName string) (domain.User, error) {
	panic("not implemented")
}
func (r *repository) ReadByEmail(ctx context.Context, email string) (domain.User, error) {
	panic("not implemented")
}
func (r *repository) ReadByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, error) {
	panic("not implemented")
}
func (r *repository) Update(ctx context.Context, changeset domain.UpdateInput) error {
	panic("not implemented")
}
func (r *repository) AddRole(context.Context, domain.ID, domain.Role) error {
	panic("not implemented")
}
func (r *repository) RemoveRole(context.Context, domain.ID, domain.Role) error {
	panic("not implemented")
}
func (r *repository) Delete(context.Context, domain.ID) error {
	panic("not implemented")
}
