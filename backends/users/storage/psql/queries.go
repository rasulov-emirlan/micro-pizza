package psql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/domain"
)

func (r *Repository) Create(ctx context.Context, u domain.User) (domain.ID, error) {
	sql, args, err := sq.Insert("users").Columns(
		"full_name", "email", "phone_number", "password", "birth_date",
		"country", "city", "street", "floor",
		"address_instructions", "created_at",
	).Values(u.FullName, u.Email, u.PhoneNumber, u.Password, u.BirthDate,
		u.Address.Country, u.Address.City, u.Address.Street,
		u.Address.Floor, u.Address.AddressInstructions,
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

// TODO: all the Read methods that return a single domain.User
// use almost the same query
// maybe we could somehow reuse the same query for all of them
// and just pass different args to it
func (r *Repository) Read(ctx context.Context, userID domain.ID) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
		"country", "city", "street", "floor", "address_instructions",
		"ARRAY_AGGG(users_roles.role_id)",
		"created_at",
		"updated_at",
	).From("users").
		LeftJoin("users_roles").
		Where("id = ?", userID).
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
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Password, &u.BirthDate,
		&u.Address.Country, &u.Address.City, &u.Address.Street, &u.Address.Floor,
		&u.Address.AddressInstructions, &roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}
	return u, nil
}

func (r *Repository) ReadByName(ctx context.Context, fullName string) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
		"country", "city", "street", "floor", "address_instructions",
		"ARRAY_AGGG(users_roles.role_id)",
		"created_at",
		"updated_at",
	).From("users").
		LeftJoin("users_roles").
		Where("name = ?", fullName).
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
		&u.Address.Country, &u.Address.City, &u.Address.Street, &u.Address.Floor,
		&u.Address.AddressInstructions, &roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}
	return u, nil
}

func (r *Repository) ReadByEmail(ctx context.Context, email string) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
		"country", "city", "street", "floor", "address_instructions",
		"ARRAY_AGGG(users_roles.role_id)",
		"created_at",
		"updated_at",
	).From("users").
		LeftJoin("users_roles").
		Where("email = ?", email).
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
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Password, &u.BirthDate,
		&u.Address.Country, &u.Address.City, &u.Address.Street, &u.Address.Floor,
		&u.Address.AddressInstructions, &roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}
	return u, nil
}

func (r *Repository) ReadByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
		"country", "city", "street", "floor", "address_instructions",
		"ARRAY_AGGG(users_roles.role_id)",
		"created_at",
		"updated_at",
	).From("users").
		LeftJoin("users_roles").
		Where(sq.Eq{"phone_number": phoneNumber}).
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
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Password, &u.BirthDate,
		&u.Address.Country, &u.Address.City, &u.Address.Street, &u.Address.Floor,
		&u.Address.AddressInstructions, &roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}
	return u, nil
}

func (r *Repository) Update(ctx context.Context, changeset domain.UpdateInput) error {
	sql, args, err := sq.
		Update("users").
		Set("full_name", changeset.FullName).
		Set("email", changeset.Email).
		Set("phone_number", changeset.PhoneNumber).
		Set("password", changeset.Password).
		Set("country", changeset.Address.Country).
		Set("city", changeset.Address.City).
		Set("street", changeset.Address.Street).
		Set("floor", changeset.Address.Floor).
		Set("address_instructions", changeset.Address.AddressInstructions).
		Set("updated_at", time.Now().UTC()).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) AddRole(ctx context.Context, userID domain.ID, inp domain.Role) error {
	sql, args, err := sq.Insert("users_roles").
		Columns("user_id", "role_id").
		Values(userID, sq.Expr("")).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) RemoveRole(ctx context.Context, userID domain.ID, role domain.Role) error {
	sql, args, err := sq.Delete("users_roles").
		Where(sq.Eq{"user_id": userID, "role_id": getRoleID(role)}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) Delete(ctx context.Context, userID domain.ID) error {
	sql, args, err := sq.Delete("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	return err
}
