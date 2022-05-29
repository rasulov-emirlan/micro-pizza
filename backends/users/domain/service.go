package domain

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo       Repository
	cache      Cache
	sms        SMSsender
	logger     Logger
	jwtManager JWTmanager
}

const (
	authRefreshExp = time.Hour * 24
	authAccessExp  = time.Hour
)

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

	j.SetExp(authAccessExp, authRefreshExp)
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
	u, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("signInEmail(): %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return SignInOutput{}, fmt.Errorf("signInemail(): password is incorrect %w", err)
	}
	claims, err := s.jwtManager.Generate(u.ID, u.Roles)
	if err != nil {
		return claims, fmt.Errorf("signInemail(): %w", err)
	}
	return claims, nil
}

func (s *service) Refresh(ctx context.Context, refreshKey string) (SignInOutput, error) {
	refClaims, err := s.jwtManager.DecodeRefresh(refreshKey)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("refresh(): %w", err)
	}
	u, err := s.repo.Read(ctx, refClaims.UserID)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("refresh(): %w", err)
	}
	claims, err := s.jwtManager.Generate(u.ID, u.Roles)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("refresh(): %w", err)
	}
	return claims, nil
}

// TODO: not sure if we need role validation here
// could be easier to do it in our transport layer
// since every request will have jwt with roles

func (s *service) AddRole(ctx context.Context, adminID, userID ID, role Role) error {
	admin, err := s.repo.Read(ctx, adminID)
	if err != nil {
		return fmt.Errorf("addRole(): %w", err)
	}
	isAllowed := false
	for _, v := range admin.Roles {
		if v == RoleAdmin || v == RoleOwner {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("addRole(): %w", ErrNotAllowed)
	}
	return fmt.Errorf("addRole(): %w", s.repo.AddRole(ctx, userID, role))
}

func (s *service) RemoveRole(ctx context.Context, adminID, userID ID, role Role) error {
	u, err := s.repo.Read(ctx, userID)
	if err != nil {
		return fmt.Errorf("removeRole(): %w", err)
	}
	for _, v := range u.Roles {
		if v == RoleOwner {
			return fmt.Errorf("removeRole(): %w", ErrOwnerCantBeRemoved)
		}
	}
	admin, err := s.repo.Read(ctx, adminID)
	if err != nil {
		return fmt.Errorf("removeRole(): %w", err)
	}
	isAllowed := false
	for _, v := range admin.Roles {
		if v == RoleAdmin || v == RoleOwner {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("removeRole(): %w", ErrNotAllowed)
	}
	return fmt.Errorf("removeRole(): %w", s.repo.RemoveRole(ctx, userID, role))
}

func (s *service) Update(ctx context.Context, whoIsUpdating ID, changeset UpdateInput) error {
	if whoIsUpdating != changeset.ID {
		return fmt.Errorf("update(): %w", ErrNotAllowed)
	}

	// some users might not even have a password
	// so we do not force them to update it
	if changeset.Password != "" {
		// but if they have a password then force them to make a good one
		if l := len(changeset.Password); l > 100 || l < 8 {
			return ErrPasswordIsNotSecure
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(changeset.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("update(): error while hashing password: %w", err)
		}
		changeset.Password = string(hash)
	}

	if _, err := mail.ParseAddress(changeset.Email); err != nil {
		return fmt.Errorf("update(): invalid email: %w", err)
	}
	// TODO: add a better phone number validation
	if l := len(changeset.PhoneNumber); l < 6 || l > 30 {
		return fmt.Errorf("update(): %w", ErrInvalidPhoneNumber)
	}
	if l := len(changeset.FullName); l == 0 || l > 250 {
		return fmt.Errorf("update(): %w", ErrInvalidFullName)
	}
	return fmt.Errorf("update(): repo error: %w", s.repo.Update(ctx, changeset))
}

func (s *service) Delete(ctx context.Context, whosDeleting, whomToDelete ID) error {
	if whosDeleting != whomToDelete {
		u, err := s.repo.Read(ctx, whomToDelete)
		if err != nil {
			return fmt.Errorf("delete(): %w", err)
		}
		for _, v := range u.Roles {
			if v == RoleOwner {
				return fmt.Errorf("delete(): %w", ErrOwnerCantBeRemoved)
			}
		}
		admin, err := s.repo.Read(ctx, whosDeleting)
		if err != nil {
			return fmt.Errorf("delete(): %w", err)
		}
		isAllowed := false
		for _, v := range admin.Roles {
			if v == RoleAdmin || v == RoleOwner {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return fmt.Errorf("delete(): %w", ErrNotAllowed)
		}
	}
	return fmt.Errorf("delete(): repo error: %w", s.repo.Delete(ctx, whomToDelete))
}
