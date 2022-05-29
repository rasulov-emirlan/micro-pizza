package psql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/domain"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/storage/psql/migrations"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(url string, withMigrations bool) (*Repository, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
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
	return &Repository{db}, nil
}

func (r *Repository) Close() {
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

func getRoleID(role domain.Role) int {
	// TODO: change this code for some
	// actual db calls maybe?
	switch role {
	case domain.RoleOwner:
		return 1
	case domain.RoleAdmin:
		return 2
	case domain.RoleModerator:
		return 3
	case domain.RoleDeliveryMan:
		return 4
	default:
		return 5
	}
}
