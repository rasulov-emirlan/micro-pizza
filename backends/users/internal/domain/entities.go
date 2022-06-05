package domain

import "time"

type (
	ID   int64
	Role uint

	Address struct {
		CountryCode  string `json:"countryCode"`
		City         string `json:"city"`
		Street       string `json:"address"`
		Floor        *int   `json:"floor,omitempty"`
		Apartment    *int   `json:"apartment,omitempty"`
		Instructions string `json:"Instructions"`
	}

	User struct {
		ID ID `json:"id"`

		FullName string `json:"fullName"`
		Roles    []Role `json:"roles"`

		Email       string `json:"-"`
		PhoneNumber string `json:"phoneNumber"`
		Password    string `json:"-"`

		BirthDate time.Time `json:"birthDate"`

		Addresses []Address `json:"addresses"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
