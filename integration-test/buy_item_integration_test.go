package integration_test

import (
	. "github.com/Eun/go-hit"
	"net/http"
	"testing"
)

// HTTP GET: /buy/:item
func TestBuyItem(t *testing.T) {
	testCases := []struct {
		description      string
		item             string
		withAuth         bool
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description:      "unauthorized",
			item:             "pen",
			withAuth:         false,
			expectedStatus:   Expect().Status().Equal(http.StatusUnauthorized),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description:      "success",
			item:             "pink-hoody",
			withAuth:         true,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().String().Len().Equal(0),
		},
		{
			description:      "wrong item name",
			item:             "ho1di2d",
			withAuth:         true,
			expectedStatus:   Expect().Status().Equal(http.StatusBadRequest),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
	}

	_, _, testToken := getValidAuthData(defaultAttempts)

	for _, tc := range testCases {
		authToken := "none"
		if tc.withAuth {
			authToken = testToken
		}
		authHeader := Send().Headers("Authorization").Add("Bearer " + authToken)

		Test(t,
			Description(tc.description),
			Get(basePath+"/buy/"+tc.item),
			authHeader,
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}
