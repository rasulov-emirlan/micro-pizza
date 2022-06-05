package psql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/internal/domain"
)

func (r *Repository) Create(ctx context.Context, u domain.User) (domain.ID, error) {
	sql, args, err := sq.Insert("users").Columns(
		"full_name", "email", "phone_number", "password", "birth_date",
		"created_at",
	).Values(u.FullName, u.Email, u.PhoneNumber, u.Password, u.BirthDate, u.CreatedAt).
		Suffix("RETURNING \"id\"").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	id := domain.ID(0)
	if err := conn.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return id, err
	}

	insertAddress := sq.Insert("addresses").
		Columns("user_id", "country", "city", "street", "floor", "apartment", "instructions")
	for _, v := range u.Addresses {
		insertAddress = insertAddress.Values(id, v.CountryCode, v.City, v.Street, v.Floor, v.Apartment, v.Instructions)
	}

	sql, args, err = insertAddress.Suffix("RETURNING \"id\"").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return id, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return id, err
	}
	return id, tx.Commit(ctx)
}

// TODO: all the Read methods that return a single domain.User
// use almost the same query
// maybe we could somehow reuse the same query for all of them
// and just pass different args to it
func (r *Repository) Read(ctx context.Context, userID domain.ID) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
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

	batch := pgx.Batch{}

	// here we get our user
	batch.Queue(sql, args...)

	// here we get our user's addresses
	sql, args, err = sq.Select(
		"country", "city", "street", "floor", "apartment", "instructions",
	).From("addresses").
		Where("user_id = ?", userID).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.User{}, err
	}
	batch.Queue(sql, args...)

	// here we send our batch and read results
	res := conn.SendBatch(ctx, &batch)

	u := domain.User{}
	roles := []int{}
	if err := res.QueryRow().Scan(
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Password, &u.BirthDate,
		&roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}

	// now we read all addresses
	addresses := []domain.Address{}
	rows, err := res.Query()
	if err != nil {
		return u, err
	}
	for rows.Next() {
		address := domain.Address{}
		if err := rows.Scan(
			&address.CountryCode, &address.City, &address.Street, &address.Floor, &address.Apartment, &address.Instructions,
		); err != nil {
			return u, err
		}
	}

	u.Addresses = addresses
	return u, res.Close()
}

func (r *Repository) ReadByName(ctx context.Context, fullName string) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
		"ARRAY_AGGG(users_roles.role_id)",
		"created_at",
		"updated_at",
	).From("users").
		LeftJoin("users_roles").
		Where("full_name = ?", fullName).
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

	batch := pgx.Batch{}

	// here we get our user
	batch.Queue(sql, args...)

	// here we get our user's addresses
	sql, args, err = sq.Select(
		"country", "city", "street", "floor", "apartment", "instructions",
	).From("addresses").
		Where("user_id = (SELECT id FROM users WHERE full_name = ?)", fullName).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.User{}, err
	}
	batch.Queue(sql, args...)

	// here we send our batch and read results
	res := conn.SendBatch(ctx, &batch)

	u := domain.User{}
	roles := []int{}
	if err := res.QueryRow().Scan(
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Password, &u.BirthDate,
		&roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}

	// now we read all addresses
	addresses := []domain.Address{}
	rows, err := res.Query()
	if err != nil {
		return u, err
	}
	for rows.Next() {
		address := domain.Address{}
		if err := rows.Scan(
			&address.CountryCode, &address.City, &address.Street, &address.Floor, &address.Apartment, &address.Instructions,
		); err != nil {
			return u, err
		}
	}

	u.Addresses = addresses
	return u, res.Close()
}

func (r *Repository) ReadByEmail(ctx context.Context, email string) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
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

	batch := pgx.Batch{}

	// here we get our user
	batch.Queue(sql, args...)

	// here we get our user's addresses
	sql, args, err = sq.Select(
		"country", "city", "street", "floor", "apartment", "instructions",
	).From("addresses").
		Where("user_id = (SELECT id FROM users WHERE email = ?)", email).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.User{}, err
	}
	batch.Queue(sql, args...)

	// here we send our batch and read results
	res := conn.SendBatch(ctx, &batch)

	u := domain.User{}
	roles := []int{}
	if err := res.QueryRow().Scan(
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.Password, &u.BirthDate,
		&roles, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	for _, v := range roles {
		u.Roles = append(u.Roles, parseRoles(v))
	}

	// now we read all addresses
	addresses := []domain.Address{}
	rows, err := res.Query()
	if err != nil {
		return u, err
	}
	for rows.Next() {
		address := domain.Address{}
		if err := rows.Scan(
			&address.CountryCode, &address.City, &address.Street, &address.Floor, &address.Apartment, &address.Instructions,
		); err != nil {
			return u, err
		}
	}

	u.Addresses = addresses
	return u, res.Close()
}

func (r *Repository) ReadByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, error) {
	sql, args, err := sq.Select(
		"id", "full_name", "email", "password",
		"birth_date", "phone_number",
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
		&roles, &u.CreatedAt, &u.UpdatedAt,
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
