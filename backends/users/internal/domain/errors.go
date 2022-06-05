package domain

import "errors"

var (
	ErrInvalidDependency = errors.New("domain: provided dependency is nil")

	ErrInvalidRequestSignUpInput = errors.New("domain: request sign up requires at least phone number or email")

	ErrInvalidToken       = errors.New("domain: provided jwt is invalid")
	ErrInvalidPhoneNumber = errors.New("domain: provided phone number is invalid")
	ErrInvalidFullName    = errors.New("domain: provided full name is invalid")
	ErrInvalidCode        = errors.New("domain: provided registration code is invalid")

	ErrOwnerCantBeRemoved = errors.New("domain: owner can't be deleted or updated")
	ErrNotAllowed         = errors.New("domain: not allowed")

	ErrPasswordIsNotSecure = errors.New("domain: password is not secure enough")

	ErrNoUsers = errors.New("domain: no users found")
)
