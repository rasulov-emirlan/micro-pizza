package domain

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/mail"
	"reflect"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo       Repository
	cache      Cache
	sms        SMSsender
	emailer    Emailer
	logger     Logger
	jwtManager JWTmanager
}

func NewService(
	r Repository,
	s SMSsender,
	m Emailer,
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
	if v := reflect.ValueOf(m); v.Kind() == reflect.Pointer &&
		v.IsNil() {
		return nil, ErrInvalidDependency
	}

	j.SetExp(AuthAccessExp, AuthRefreshExp)
	j.SetKey(jwtKey)

	return &service{
		repo:       r,
		sms:        s,
		emailer:    m,
		logger:     l,
		jwtManager: j,
	}, nil
}

func (s *service) Read(ctx context.Context, id ID) (User, error) {
	u, err := s.repo.Read(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("read(): could not read from db %w", err)
	}
	return u, nil
}

func (s *service) ReadByEmail(ctx context.Context, email string) (User, error) {
	u, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("readByEmail(): could not read from db %w", err)
	}
	return u, nil
}

func (s *service) ReadByPhoneNumber(ctx context.Context, phoneNumber string) (User, error) {
	u, err := s.repo.ReadByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return User{}, fmt.Errorf("readByPhoneNumber(): could not read from db %w", err)
	}
	return u, nil
}

func (s *service) ReadAll(ctx context.Context, cfg ReadAllInput) ([]User, error) {
	u, err := s.repo.ReadAll(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("readAll(): could not read from db %w", err)
	}
	return u, nil
}

func (s *service) RequestSignUp(ctx context.Context, inp RequestSignUpInput) error {
	code := make([]byte, CodeLength)
	if _, err := rand.Read(code); err != nil {
		return fmt.Errorf("requestSignUp(): %w", err)
	}

	switch {
	case utf8.RuneCountInString(inp.PhoneNumber) != 0:
		if err := s.sms.Send(
			inp.PhoneNumber,
			RequestSignUpSMSTitle,
			RequestSignUpSMSMessage+":"+string(code),
		); err != nil {
			return fmt.Errorf("requestSignUp(): could not send sms %w", err)
		}
		if err := s.cache.Store(inp.PhoneNumber, string(code)); err != nil {
			return fmt.Errorf("requestSignUp(): could not cache %w", err)
		}
	case utf8.RuneCountInString(inp.Email) != 0:
		if err := s.emailer.Send(
			inp.Email,
			RequestSignUpEmailTitle,
			RequestSignUpEmailMessage+":"+string(code),
		); err != nil {
			return fmt.Errorf("requestSignUp(): could not send email %w", err)
		}
		if err := s.cache.Store(inp.Email, string(code)); err != nil {
			return fmt.Errorf("requestSignUp(): could not cache %w", err)
		}
	default:
		return ErrInvalidRequestSignUpInput
	}
	return nil
}

func (s *service) SignUp(ctx context.Context, inp SignUpInput) (SignInOutput, error) {
	switch {
	case utf8.RuneCountInString(inp.PhoneNumber) != 0:
		code, err := s.cache.Get(inp.PhoneNumber)
		if err != nil {
			return SignInOutput{}, fmt.Errorf("signUp(): %w", err)
		}
		if code != inp.Code {
			return SignInOutput{}, fmt.Errorf("signUp(): %w", ErrInvalidCode)
		}
	case utf8.RuneCountInString(inp.Email) != 0:
		code, err := s.cache.Get(inp.Email)
		if err != nil {
			return SignInOutput{}, fmt.Errorf("signUp(): %w", err)
		}
		if code != inp.Code {
			return SignInOutput{}, fmt.Errorf("signUp(): %w", ErrInvalidCode)
		}
	default:
		return SignInOutput{}, ErrInvalidSignUpInput
	}

	var (
		u = User{
			FullName:    inp.FullName,
			Email:       inp.Email,
			PhoneNumber: inp.PhoneNumber,
			Addresses:   []Address{inp.Address},
			Roles:       []Role{RoleUser},
			BirthDate:   inp.BirthDate,
			CreatedAt:   time.Now().UTC(),
		}
		err error
	)
	u.ID, err = s.repo.Create(ctx, u)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("signUp(): %w", err)
	}

	claims, err := s.jwtManager.Generate(u.ID, u.Roles)
	if err != nil {
		return claims, fmt.Errorf("signUp(): could not generate jwt due to: %w", err)
	}
	return claims, nil
}

func (s *service) RequestSignIn(ctx context.Context, inp RequestSignInInput) error {
	code := make([]byte, CodeLength)
	if _, err := rand.Read(code); err != nil {
		return fmt.Errorf("requestSignIn(): %w", err)
	}

	switch {
	case utf8.RuneCountInString(inp.PhoneNumber) != 0:
		if err := s.sms.Send(
			inp.PhoneNumber,
			RequestSignInSMSTitle,
			RequestSignInSMSMessage+":"+string(code),
		); err != nil {
			return fmt.Errorf("requestSignIn(): could not send sms %w", err)
		}
		if err := s.cache.Store(inp.PhoneNumber, string(code)); err != nil {
			return fmt.Errorf("requestSignIn(): %w", err)
		}
	case utf8.RuneCountInString(inp.Email) != 0:
		if err := s.emailer.Send(
			inp.Email,
			RequestSignInEmailTitle,
			RequestSignInEmailMessage+":"+string(code),
		); err != nil {
			return fmt.Errorf("requestSignIn(): could not send email %w", err)
		}
		if err := s.cache.Store(inp.Email, string(code)); err != nil {
			return fmt.Errorf("requestSignIn(): %w", err)
		}
	default:
		return ErrInvalidSignInInput
	}
	return nil
}

func (s *service) SignIn(ctx context.Context, inp SignInInput) (SignInOutput, error) {
	var (
		u   User
		code string
		err error
	)
	switch {
	case utf8.RuneCountInString(inp.PhoneNumber) != 0:
		code, err = s.cache.Get(inp.PhoneNumber)
		if err != nil {
			return SignInOutput{}, fmt.Errorf("signIn(): could read from cache %w", err)
		}
		if code != inp.Code {
			return SignInOutput{}, fmt.Errorf("signIn(): %w", ErrInvalidCode)
		}
		u, err = s.repo.ReadByPhoneNumber(ctx, inp.PhoneNumber)
	case utf8.RuneCountInString(inp.Email) != 0:
		code, err = s.cache.Get(inp.Email)
		if err != nil {
			return SignInOutput{}, fmt.Errorf("signIn(): could read from cache %w", err)
		}
		if code != inp.Code {
			return SignInOutput{}, fmt.Errorf("signIn(): %w", ErrInvalidCode)
		}
		u, err = s.repo.ReadByEmail(ctx, inp.Email)
	default:
		return SignInOutput{}, ErrInvalidSignInInput
	}

	if err != nil {
		return SignInOutput{}, fmt.Errorf("signIn(): could not read from db %w", err)
	}

	claims, err := s.jwtManager.Generate(u.ID, u.Roles)
	if err != nil {
		return claims, fmt.Errorf("signIn(): could not generate jwt due to: %w", err)
	}
	return claims, nil
}

func (s *service) SignInEmailPassword(ctx context.Context, email, password string) (SignInOutput, error) {
	u, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return SignInOutput{}, fmt.Errorf("signInEmailPassword(): could not read from db %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return SignInOutput{}, fmt.Errorf("signInEmailPassword(): password is incorrect %w", err)
	}
	claims, err := s.jwtManager.Generate(u.ID, u.Roles)
	if err != nil {
		return claims, fmt.Errorf("signInEmailPassword(): %w", err)
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

func (s *service) AddRole(ctx context.Context, userID ID, role Role) error {
	if role == RoleOwner {
		return fmt.Errorf("addRole(): owners can be asigned only manualy %w", ErrInvalidRole)
	}
	return fmt.Errorf("addRole(): %w", s.repo.AddRole(ctx, userID, role))
}

func (s *service) RemoveRole(ctx context.Context, userID ID, role Role) error {
	u, err := s.repo.Read(ctx, userID)
	if err != nil {
		return fmt.Errorf("removeRole(): %w", err)
	}
	for _, v := range u.Roles {
		if v == RoleOwner {
			return fmt.Errorf("removeRole(): %w", ErrNotAllowed)
		}
	}
	return fmt.Errorf("removeRole(): %w", s.repo.RemoveRole(ctx, userID, role))
}

// Make sure that the one calling this function is the user itself
// If not then he has to be at least an admin. And keep in mind that nobody exept
// the owner can change owners fields
func (s *service) Update(ctx context.Context, changeset UpdateInput) error {
	// some users might not even have a password
	// so we do not force them to update it
	if changeset.Password != "" {
		// but if they have a password then force them to make a good one
		if l := utf8.RuneCountInString(changeset.Password); l > PasswordMaxLength || l < PasswordMinLength {
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
	if l := utf8.RuneCountInString(changeset.PhoneNumber); l < 6 || l > 30 {
		return fmt.Errorf("update(): %w", ErrInvalidPhoneNumber)
	}
	if l := utf8.RuneCountInString(changeset.FullName); l == 0 || l > 250 {
		return fmt.Errorf("update(): %w", ErrInvalidFullName)
	}
	return fmt.Errorf("update(): repo error: %w", s.repo.Update(ctx, changeset))
}

func (s *service) Delete(ctx context.Context, whomToDelete ID) error {
	u, err := s.repo.Read(ctx, whomToDelete)
	if err != nil {
		return fmt.Errorf("delete(): %w", err)
	}
	for _, v := range u.Roles {
		if v == RoleOwner {
			return fmt.Errorf("delete(): %w", ErrNotAllowed)
		}
	}
	return fmt.Errorf("delete(): repo error: %w", s.repo.Delete(ctx, whomToDelete))
}
