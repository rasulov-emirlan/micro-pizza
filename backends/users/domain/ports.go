package domain

import (
	"context"
	"time"
)

type (
	Service interface {
		SignUp(context.Context, SignUpInput) error

		RequestSignIn(ctx context.Context, phoneNumber string) error
		SignIn(ctx context.Context, phoneNumber string, code []byte) (SignInOutput, error)

		SignInEmail(ctx context.Context, email, password string) (SignInOutput, error)
		Refresh(ctx context.Context, refreshKey string) (SignInOutput, error)

		AddRole(ctx context.Context, adminID, userID ID, role Role) error
		RemoveRole(ctx context.Context, adminID, userID ID, role Role) error

		Update(ctx context.Context, whoIsUpdating ID, changeset UpdateInput) error
		Delete(ctx context.Context, whosDeleting, whomToDelete ID) error
	}

	SMSsender interface {
		Send(phoneNumber string, code []byte) error
	}

	Cache interface {
		Store(key, value string) error
		Get(key string) (value string, err error)
	}

	Repository interface {
		Create(context.Context, User) (ID, error)

		Read(context.Context, ID) (User, error)
		ReadByName(ctx context.Context, fullName string) (User, error)
		ReadByEmail(ctx context.Context, email string) (User, error)
		ReadByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)

		Update(ctx context.Context, changeset UpdateInput) error
		AddRole(context.Context, ID, Role) error

		RemoveRole(context.Context, ID, Role) error
		Delete(context.Context, ID) error
	}

	Logger interface {
		Infof(format string, args ...string)
		Errorf(format string, args ...string)
	}

	JWTmanager interface {
		SetKey(key []byte)
		SetExp(access, refresh time.Duration)

		Generate(userID ID, roles []Role) (SignInOutput, error)
		DecodeAccess(accessKey string) (AccessClaims, error)
		DecodeRefresh(refreshKey string) (RefreshClaims, error)
	}
)
