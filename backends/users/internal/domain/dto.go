package domain

import "time"

type (
	RequestSignUpInput struct {
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
	}

	SignUpInput struct {
		Code string `json:"code"`

		// which ever is not empty will be used
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`

		FullName  string    `json:"fullName"`
		BirthDate time.Time `json:"birthDate"`
		Address   Address   `json:"address"`
	}

	RequestSignInInput struct {
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
	}

	SignInInput struct {
		Code        string `json:"code"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
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
	}

	Sorting uint

	ReadAllInput struct {
		Limit       int     `json:"limit"`
		Offset      int     `json:"offset"`
		SortBy      Sorting `json:"sortBy"`
		CountryCode string  `json:"countryCode"`
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
