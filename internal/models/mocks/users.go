package mocks

import (
	"time"

	"gosnipit.ricci2511.dev/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	// emulate a duplicate email error
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "mocked@example.com" && password == "mocked1234" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	if id == 1 {
		return &models.User{
			ID:      1,
			Name:    "Mocky McMockface",
			Email:   "mocked@example.com",
			Created: time.Now(),
		}, nil
	}

	return nil, models.ErrNoRecord
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 && currentPassword == "mocked1234" && newPassword == "mocked5678" {
		return nil
	}

	return models.ErrInvalidCredentials
}
