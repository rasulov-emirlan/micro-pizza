package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/internal/domain"
	"github.com/rasulov-emirlan/micro-pizzas/backends/users/internal/domain/mocks"
)

func TestRequestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockCache(ctrl)
	mockRepo := mocks.NewMockRepository(ctrl)
	mockSMSsender := mocks.NewMockSMSsender(ctrl)
	mockEmailer := mocks.NewMockEmailer(ctrl)
	mockJWTmanager := mocks.NewMockJWTmanager(ctrl)
	mockJWTmanager.EXPECT().SetExp(gomock.Any(), gomock.Any()).Times(1)
	mockJWTmanager.EXPECT().SetKey(gomock.Any()).Times(1)
	mockLogger := mocks.NewMockLogger(ctrl)
	s, err := domain.NewService(
		mockRepo,
		mockSMSsender,
		mockEmailer,
		mockCache,
		mockLogger,
		mockJWTmanager,
		[]byte("secret"),
	)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		name string
		inp  domain.RequestSignUpInput
		err  error

		mockup func()
	}{
		{
			name: "success with email",
			inp: domain.RequestSignUpInput{
				Email:       "pizzas@gmail.com",
				PhoneNumber: "",
			},
			err: nil,
			mockup: func() {
				mockCache.EXPECT().Store("pizzas@gmail.com", gomock.AssignableToTypeOf(""))
				mockEmailer.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
		{
			name: "success with phone number",
			inp: domain.RequestSignUpInput{
				Email:       "",
				PhoneNumber: "+996702569123",
			},
			err: nil,
			mockup: func() {
				mockCache.EXPECT().Store(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockSMSsender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		{
			name: "fail with empty email and phone number",
			inp: domain.RequestSignUpInput{
				Email:       "",
				PhoneNumber: "",
			},
			err: domain.ErrInvalidRequestSignUpInput,
			mockup: func() {
				mockCache.EXPECT().Store(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockEmailer.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				mockSMSsender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockup()
			err := s.RequestSignUp(context.Background(), tc.inp)
			if !errors.Is(err, tc.err) {
				t.Errorf("got %v, want %v", err, tc.err)
			}
		})
	}
}
