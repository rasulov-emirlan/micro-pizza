package domain

import "time"

const (
	RoleOwner Role = iota
	RoleAdmin
	RoleModerator
	RoleDeliveryMan
	RoleUser

	ReadAllSortByID Sorting = iota
	ReadAllSortByFullNameASC
	ReadAllSortByFullNameDESC
	ReadAllSortByEmailASC
	ReadAllSortByEmailDESC

	AuthRefreshExp = time.Hour * 24
	AuthAccessExp  = time.Hour

	CodeLength = 6

	PasswordMinLength = 8
	PasswordMaxLength = 64

	RequestSignUpEmailTitle = `
	Micro-Pizzas sign up code
	`
	RequestSignUpEmailMessage = `
	Here is your code for signin up
	`

	RequestSignUpSMSTitle = `
	Micro-Pizzas sign up code
	`
	RequestSignUpSMSMessage = `
	Here is your code for signin up
	`

	RequestSignInEmailTitle = `
	Micro-Pizzas sign in code
	`
	RequestSignInEmailMessage = `
	Here is your code for signin in
	`

	RequestSignInSMSTitle = `
	Micro-Pizzas sign in code
	`
	RequestSignInSMSMessage = `
	Here is your code for signin in
	`
)
