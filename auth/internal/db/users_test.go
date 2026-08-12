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
	_, err := userModel.CreateUser(ctx, user)
	assert.NoError(t, err)
}
