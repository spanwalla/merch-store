package integration_test

import (
	. "github.com/Eun/go-hit"
	"github.com/spanwalla/merch-store/internal/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// HTTP GET: /info
func TestInfo(t *testing.T) {
	_, _, userToken := getValidAuthData(defaultAttempts)

	testCases := []struct {
		description      string
		authToken        string
		expectedStatus   IStep
		expectedResponse IStep
	}{
		{
			description:      "unauthorized",
			authToken:        "no-token",
			expectedStatus:   Expect().Status().Equal(http.StatusUnauthorized),
			expectedResponse: Expect().Body().JSON().JQ(".errors").Len().GreaterThan(0),
		},
		{
			description:      "success",
			authToken:        userToken,
			expectedStatus:   Expect().Status().Equal(http.StatusOK),
			expectedResponse: Expect().Body().JSON().Contains("coins", "inventory", "coinHistory"),
		},
	}

	for _, tc := range testCases {
		authHeader := Send().Headers("Authorization").Add("Bearer " + tc.authToken)

		Test(t,
			Description(tc.description),
			Get(basePath+"/info"),
			authHeader,
			tc.expectedStatus,
			tc.expectedResponse,
		)
	}
}

func TestBuyAndInfo(t *testing.T) {
	_, _, userToken := getValidAuthData(defaultAttempts)
	var response entity.UserReport

	MustDo(
		Description("buy item"),
		Get(basePath+"/buy/hoody"),
		Send().Headers("Authorization").Add("Bearer "+userToken),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Len().Equal(0),
	)

	MustDo(
		Description("get info"),
		Get(basePath+"/info"),
		Send().Headers("Authorization").Add("Bearer "+userToken),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().In(&response),
	)

	assert.Contains(t, response.Inventory, entity.Inventory{
		Type:     "hoody",
		Quantity: 1,
	})
}
