package domain

import "time"

type (
	SignUpInput struct {
		// FullName is not required
		FullName    string    `json:"fullName"`
		PhoneNumber string    `json:"phoneNumber"`
		BirthData   time.Time `json:"birthDate"`

		// Email is not required
		Email   string  `json:"email"`
		Address Address `json:"address"`
	}

	SignInInput struct {
		PhoneNumber string `json:"phoneNumber"`

		// If PhoneNumber is provided then these fields are not required
		// But they can be used to login using email and password
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInOutput struct {
		AccessKey  string `json:"accessKey"`
		RefreshKey string `json:"refreshKey"`
	}

	UpdateInput struct {
		ID ID `json:"id"`

		// Fields that are empty fill not be updated
		FullName    string `json:"fullName"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Address
	}

	// Structs bellow are for JWTmanager's use

	AccessClaims struct {
		UserID ID     `json:"userID"`
		Roles  []Role `json:"roles"`
	}
	RefreshClaims struct {
		UserID ID `json:"userID"`
	}
)
