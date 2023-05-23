package mocks

import "gosnipit.ricci2511.dev/internal/models"

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
	if email == "mocked@example.com" && password == "pass" {
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
