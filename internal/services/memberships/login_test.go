package memberships

import (
	"testing"

	"github.com/amiulam/music-catalog/internal/configs"
	"github.com/amiulam/music-catalog/internal/models/memberships"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_service_Login(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	type args struct {
		request memberships.LoginRequest
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				request: memberships.LoginRequest{
					Email:    "test@mail.com",
					Password: "password",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, "", int64(0)).Return(&memberships.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "test@mail.com",
					Password: "$2a$10$03nzIVJ8O36iifVU2sI6hOv5519k0929vcCRuJogheXAAbiWr58CG",
					Username: "testusername",
				}, nil)
			},
		},
		{
			name: "fail when get user",
			args: args{
				request: memberships.LoginRequest{
					Email:    "test@mail.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, "", int64(0)).Return(nil, assert.AnError)
			},
		},
		{
			name: "invalid credentials",
			args: args{
				request: memberships.LoginRequest{
					Email:    "test@mail.com",
					Password: "password",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, "", int64(0)).Return(&memberships.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "test@mail.com",
					Password: "wrong password",
					Username: "testusername",
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &service{
				cfg: &configs.Config{
					Service: configs.Service{
						SecretJWT: "abc",
					},
				},
				membershipRepo: mockRepo,
			}

			got, err := s.Login(tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotEmpty(t, got)
			} else {
				assert.Empty(t, got)
			}
		})
	}
}
