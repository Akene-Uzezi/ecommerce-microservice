package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	payload, _ := json.Marshal(&map[string]any{
		"email":    "test@test.com",
		"password": "testpwd",
		"name":     "testuser",
	})
	req, _ := http.NewRequest("POST", "/api/v1/create_user", bytes.NewBuffer(payload))
	client := &http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)
	t.Logf("response: %v", res)
}
