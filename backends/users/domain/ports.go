package domain

import "context"

type (
	Service interface {
		SignUp(context.Context, SignUpInput) error
		SignIn(context.Context, SignInInput) (SignInOutput, error)
		Refresh(ctx context.Context, refreshKey string) (SignInOutput, error)

		AddRole(ctx context.Context, adminID, userID ID, role Role) error
		RemoveRole(ctx context.Context, adminID, userID ID, role Role) error

		Update(ctx context.Context, changeset UpdateInput) error
		Delte(ctx context.Context, userID ID) error
	}

	Repository interface {
		Create(context.Context, SignInInput) (ID, error)
		Read(context.Context, ID) (User, error)
		ReadByName(ctx context.Context, fullName string) (User, error)
		ReadByEmail(ctx context.Context, email string) (User, error)
		Update(ctx context.Context, changeset UpdateInput) error
		Delete(context.Context, ID) error
	}

	Logger interface {
		Infof(format string, args ...string)
	}

	JWTmanager interface {
		Generate(SignInInput) (SignInOutput, error)
		DecodeRefresh(refreshKey string) (RefreshClaims, error)
	}
)
