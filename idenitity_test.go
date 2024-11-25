package identity

import (
	"log"
	"testing"
)

func (TestClientContext) CheckExistedUser(username string) bool {
	return true
}

type TestClientContext struct {
}

type TestClient struct {
	Client
}

func TestRegister(t *testing.T) {
	request := UserRegisterRequest{
		Username: "test_user01",
		Password: "test123456",
		Email:    "test@example.com",
	}

	client := TestClient{
		Client{
			Context: TestClientContext{},
		},
	}

	registerErr := client.Register(&request)

	if registerErr != nil {
		log.Fatal(registerErr)
	}
}
