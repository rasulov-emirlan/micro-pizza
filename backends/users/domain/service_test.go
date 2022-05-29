package domain_test

import (
	"testing"

	"github.com/rasulov-emirlan/micro-pizzas/backends/users/domain"
)

// test domain/service.go

func setupService() {

}

func TestServiceSignUp(t *testing.T) {
	testCases := []struct {
		name    string
		input   domain.SignUpInput
		wantErr bool
	}{
		{
			name: "success",
			input: domain.SignUpInput{
				FullName:    "John Doe",
				PhoneNumber: "+1123-1234",
				Email:       "john@gmail.com",
				Address: domain.Address{
					Country:             "USA",
					City:                "New York",
					Street:              "123 Main Street",
					Floor:               1,
					AddressInstructions: "Live on a second floor",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// err := domain.SignUp(tc.input)
			// if tc.wantErr && err == nil {
			// 	t.Errorf("Expected error but got nil")
			// }
		})
	}
}
