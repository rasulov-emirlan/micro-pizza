package domain

import "time"

type (
	ID   int64
	Role uint
)

const (
	RoleOwner Role = iota
	RoleAdmin
	RoleModerator
	RoleDeliveryMan
	RoleUser
)

type (
	Address struct {
		Country             string `json:"country"`
		City                string `json:"city"`
		Street              string `json:"address"`
		HomeNumber          int    `json:"homeNumber"`
		AddressInstructions string `json:"addressInstructions"`
	}

	User struct {
		ID ID `json:"id"`

		FullName string `json:"fullName"`
		Roles    []Role `json:"roles"`

		Email       string `json:"-"`
		PhoneNumber string `json:"phoneNumber"`
		Password    string `json:"-"`

		BirthDate time.Time `json:"birthDate"`

		Address Address `json:"address"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
