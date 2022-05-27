package domain

import (
	"context"
	"reflect"
)

type service struct {
	repo       Repository
	logger     Logger
	jwtManager JWTmanager
}

func NewService(r Repository, l Logger, j JWTmanager) (Service, error) {
	if v := reflect.ValueOf(r); v.Kind() == reflect.Pointer &&
		reflect.ValueOf(r).IsNil() {
		return nil, ErrInvalidDependency
	}
	if v := reflect.ValueOf(l); v.Kind() == reflect.Pointer &&
		v.IsNil() {
		return nil, ErrInvalidDependency
	}
	if v := reflect.ValueOf(j); v.Kind() == reflect.Pointer &&
		v.IsNil() {
		return nil, ErrInvalidDependency
	}
	return &service{
		repo:       r,
		logger:     l,
		jwtManager: j,
	}, nil
}

func (s *service) SignUp(context.Context, SignUpInput) error {
	panic("not implemented")
}

func (s *service) SignIn(context.Context, SignInInput) (SignInOutput, error) {
	panic("not implemented")
}

func (s *service) Refresh(ctx context.Context, refreshKey string) (SignInOutput, error) {
	panic("not implemented")
}

func (s *service) AddRole(ctx context.Context, adminID, userID ID, role Role) error {
	panic("not implemented")
}

func (s *service) RemoveRole(ctx context.Context, adminID, userID ID, role Role) error {
	panic("not implemented")
}

func (s *service) Update(ctx context.Context, changeset UpdateInput) error {
	panic("not implemented")
}

func (s *service) Delte(ctx context.Context, userID ID) error {
	panic("not implemented")
}
