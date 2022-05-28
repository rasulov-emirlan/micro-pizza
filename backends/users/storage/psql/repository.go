package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/domain"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/storage/psql/migrations"
)

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(url string, withMigrations bool) (*repository, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		// if our db starts in another docker container simultaneously
		// with our app then we might have to wait for it to be ready
		time.Sleep(10 * time.Second)
		db, err = pgxpool.Connect(context.Background(), url)
		if err != nil {
			return nil, err
		}
	}
	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}
	if withMigrations {
		if err := migrations.Up(url); err != nil {
			db.Close()
			return nil, err
		}
	}
	return &repository{db}, nil
}

func (r *repository) Close() {
	r.conn.Close()
}

func parseRoles(role int) domain.Role {
	switch role {
	case 0:
		return domain.RoleOwner
	case 1:
		return domain.RoleAdmin
	case 2:
		return domain.RoleModerator
	case 3:
		return domain.RoleDeliveryMan
	default:
		return domain.RoleUser
	}
}
