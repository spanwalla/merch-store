package integration_test

import (
	. "github.com/Eun/go-hit"
	"github.com/brianvoe/gofakeit/v7"
	"net/http"
	"testing"
)

// HTTP POST: /auth
func TestAuth(t *testing.T) {
	testUsername := "test_" + gofakeit.Username()
	testPassword := gofakeit.Password(true, true, true, true, false, 12) + "aA1!"
	testWrongPassword := gofakeit.Password(true, false, false, true, false, 12) + "bB1!"
	testBadPassword := gofakeit.Password(true, false, true, false, false, 12)

	testCases := []struct {
		description      string
		body             map[string]interface{}
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description: "registration success",
			body: map[string]interface{}{
				"username": testUsername,
				"password": testPassword,
			},
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().JQ(".token").Len().GreaterThan(0),
		},
		{
			description: "authorization success",
			body: map[string]interface{}{
				"username": testUsername,
				"password": testPassword,
			},
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().JQ(".token").Len().GreaterThan(0),
		},
		{
			description: "authorization wrong password",
			body: map[string]interface{}{
				"username": testUsername,
				"password": testWrongPassword,
			},
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Equal("wrong password"),
		},
		{
			description: "registration bad password",
			body: map[string]interface{}{
				"username": "test_" + gofakeit.Username(),
				"password": testBadPassword,
			},
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description: "empty fields",
			body: map[string]interface{}{
				"username": "",
				"password": "",
			},
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description: "password wasnt provided",
			body: map[string]interface{}{
				"username": "test_" + gofakeit.Username(),
			},
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
	}

	for _, tc := range testCases {
		Test(t,
			Description(tc.description),
			Post(basePath+"/auth"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().JSON(tc.body),
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}
