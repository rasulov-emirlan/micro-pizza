package domain

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type service struct {
	repo       Repository
	cache      Cache
	sms        SMSsender
	logger     Logger
	jwtManager JWTmanager
}

func NewService(
	r Repository,
	s SMSsender,
	c Cache,
	l Logger,
	j JWTmanager,
	jwtKey []byte,
) (Service, error) {
	if v := reflect.ValueOf(r); v.Kind() == reflect.Pointer &&
		reflect.ValueOf(r).IsNil() {
		return nil, ErrInvalidDependency
	}
	if v := reflect.ValueOf(c); v.Kind() == reflect.Pointer &&
		v.IsNil() {
		return nil, ErrInvalidDependency
	}
	if v := reflect.ValueOf(s); v.Kind() == reflect.Pointer &&
		v.IsNil() {
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

	// errors here can be ignored
	j.SetExp(time.Minute*20, time.Hour*24)
	j.SetKey(jwtKey)

	return &service{
		repo:       r,
		sms:        s,
		logger:     l,
		jwtManager: j,
	}, nil
}

func (s *service) SignUp(ctx context.Context, inp SignUpInput) error {
	u := User{
		FullName:    inp.FullName,
		Email:       inp.Email,
		PhoneNumber: inp.PhoneNumber,
		Address:     inp.Address,
		Roles:       []Role{RoleUser},
		BirthDate:   inp.BirthData,
		CreatedAt:   time.Now().UTC(),
	}
	_, err := s.repo.Create(ctx, u)
	if err != nil {
		return fmt.Errorf("signUp(): %w", err)
	}
	return nil
}

func (s *service) RequestSignIn(ctx context.Context, phoneNumber string) error {
	code := make([]byte, 4)
	if _, err := rand.Read(code); err != nil {
		return fmt.Errorf("requestSignIn(): %w", err)
	}
	if err := s.sms.Send(phoneNumber, code); err != nil {
		return fmt.Errorf("requestSignIn: %w", err)
	}
	return s.cache.Store(phoneNumber, string(code))
}

func (s *service) SignIn(ctx context.Context, phoneNumber string, code []byte) (SignInOutput, error) {
	if len(code) != 4 {
		return SignInOutput{}, errors.New("signIn(): incorrect code")
	}
	realCode, err := s.cache.Get(phoneNumber)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("signIn(): %w", err)
	}
	if realCode != string(code) {
		return SignInOutput{}, errors.New("signIn(): incorrect code")
	}
	u, err := s.repo.ReadByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("signIn(): %w", err)
	}
	claims, err := s.jwtManager.Generate(u.ID, u.Roles)
	if err != nil {
		return claims, fmt.Errorf("signIn(): %w", err)
	}
	return claims, nil
}

func (s *service) SignInEmail(ctx context.Context, email, password string) (SignInOutput, error) {
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

func (s *service) Update(ctx context.Context, whoIsUpdating ID, changeset UpdateInput) error {
	panic("not implemented")
}

func (s *service) Delete(ctx context.Context, whosDeleting, whomToDelete ID) error {
	panic("not implemented")
}
