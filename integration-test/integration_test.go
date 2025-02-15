package integration_test

import (
	. "github.com/Eun/go-hit"
	"github.com/brianvoe/gofakeit/v7"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	host            = "app:8080"
	healthPath      = "http://" + host + "/health"
	defaultAttempts = 20

	basePath = "http://" + host + "/api"
)

func TestMain(m *testing.M) {
	err := healthCheck(defaultAttempts)
	if err != nil {
		log.Fatalf("integration tests: host %s is not available: %v", host, err)
	}

	log.Infof("integration tests: host %s is available", host)
	os.Exit(m.Run())
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Infof("integration tests: host %s is not available, attempts left: %d", host, attempts)
		time.Sleep(time.Second)
		attempts--
	}

	return err
}

func getValidAuthData(attempts int) (string, string, string) {
	var testUsername, testPassword, token string
	var err error
	for attempts > 0 {
		testUsername = "test_" + gofakeit.Username()
		testPassword = gofakeit.Password(true, true, true, true, false, 12) + "aA1!"
		token, err = getAuthToken(testUsername, testPassword)

		if err == nil {
			return testUsername, testPassword, token
		}

		log.Infof("integration tests: auth attempts left: %d", attempts)
		attempts--
		time.Sleep(time.Second)
	}
	log.Fatalf("integration tests: auth error: %v", err)
	return "", "", ""
}

func getAuthToken(username, password string) (string, error) {
	var token string
	var err error

	body := map[string]any{
		"username": username,
		"password": password,
	}
	err = Do(
		Post(basePath+"/auth"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(body),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".token").In(&token),
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
