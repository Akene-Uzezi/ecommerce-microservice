package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var user = &User{
	Email:    "test@test.com",
	Password: "testpassword",
	Name:     "testuser",
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	_, err := UserModel.CreateUser(ctx, user)
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	_, err := UserModel.GetUserByEmail(ctx, user)
	assert.NoError(t, err)
}
