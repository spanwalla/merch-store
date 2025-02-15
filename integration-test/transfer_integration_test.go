package integration_test

import (
	. "github.com/Eun/go-hit"
	"net/http"
	"testing"
)

// HTTP POST: /sendCoin
func TestTransfer(t *testing.T) {
	_, _, firstToken := getValidAuthData(defaultAttempts)
	secondUsername, _, secondToken := getValidAuthData(defaultAttempts)

	testCases := []struct {
		description      string
		body             map[string]interface{}
		authToken        string
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description: "unauthorized",
			body: map[string]interface{}{
				"toUser": secondUsername,
				"amount": 88,
			},
			authToken:        "",
			expectedStatus:   Expect().Status().Equal(http.StatusUnauthorized),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description: "success",
			body: map[string]interface{}{
				"toUser": secondUsername,
				"amount": 13,
			},
			authToken:        firstToken,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().String().Len().Equal(0),
		},
		{
			description: "wrong username",
			body: map[string]interface{}{
				"toUser": secondUsername + "$(!kf",
				"amount": 49,
			},
			authToken:        firstToken,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description: "self transfer",
			body: map[string]interface{}{
				"toUser": secondUsername,
				"amount": 10,
			},
			authToken:        secondToken,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description: "negative amount",
			body: map[string]interface{}{
				"toUser": secondUsername,
				"amount": -10,
			},
			authToken:        firstToken,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description: "not all fields are set",
			body: map[string]interface{}{
				"toUser": secondUsername,
			},
			authToken:        firstToken,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
	}

	for _, tc := range testCases {
		Test(t,
			Description(tc.description),
			Post(basePath+"/sendCoin"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Headers("Authorization").Add("Bearer "+tc.authToken),
			Send().Body().JSON(tc.body),
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}
