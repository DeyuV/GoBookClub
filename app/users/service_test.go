package users

import (
	"context"
	"errors"
	"testing"
)

type mockedRepository struct {
	users map[uint]*UserResponse
}

func (m *mockedRepository) VerifyUser(ctx context.Context, userId uint) error {
	_, has := m.users[userId]
	if !has {
		return errors.New("user not found")
	}
	return nil
}

func (m *mockedRepository) GetUserById(ctx context.Context, userId uint) (*UserResponse, error) {
	user, has := m.users[userId]
	if !has {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockedRepository) Create(ctx context.Context, FirstName, LastName, Email, PasswordHash string, Enabled bool, Role string) (*UserResponse, error) {

	for _, user := range m.users {
		if user.Email == Email {
			return nil, errors.New("user already exists")
		}
	}

	usersCount := uint(len(m.users))

	result := UserResponse{
		ID:        usersCount,
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
		Password:  PasswordHash,
		Enabled:   Enabled,
		Role:      Role,
	}
	m.users[usersCount] = &result
	return &result, nil
}

func newMockedRepository() Repository {
	result := mockedRepository{
		users: make(map[uint]*UserResponse),
	}
	result.Create(context.Background(), "Andrei", "Varodi", "andrei.varodi@gmail.com", "test", true, "USER")
	result.Create(context.Background(), "Andrei", "Mic", "andrei.mic@gmail.com", "test", true, "USER")
	return &result
}

func TestCreateServiceUser(t *testing.T) {
	const Password = "123456"
	repo := newMockedRepository()
	svc := NewService(repo)
	newUser, err := svc.CreateUser(context.Background(), "Bogdan", "Dinu", "badu@badu.ro", Password)
	if err != nil {
		t.Fatalf("error creating user : %#v", err)
	}

	if len(newUser.Password) == 0 || Password == newUser.Password {
		t.Fatal("expecting hashed password")
	}

	t.Logf("user was created : %#v", newUser)
}
