package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"test-backend-1-X1ag/internal/app"
	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/logger"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.TestMode)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	application, err := app.New(context.Background(), cfg, logger.NewTestLogger())
	if err != nil {
		t.Fatalf("create app for tests: %v", err)
	}
	t.Cleanup(application.Close)

	server := httptest.NewServer(application.Router)
	t.Cleanup(server.Close)

	return server
}

func doRequest(t *testing.T, client *http.Client, serverURL, method, path, token string, body any) (*http.Response, []byte) {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		rawBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(rawBody)
	}

	req, err := http.NewRequest(method, serverURL+path, bodyReader)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("perform request %s %s: %v", method, path, err)
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	return resp, respBody
}

func nextDateForWeekdays(daysOfWeek []int, slotStart string) string {
	now := time.Now().UTC()

	startClock, err := time.Parse("15:04", slotStart)
	if err != nil {
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	}

	for i := 0; i < 14; i++ {
		candidate := now.AddDate(0, 0, i)
		weekday := int(candidate.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		for _, day := range daysOfWeek {
			if weekday != day {
				continue
			}

			candidateStart := time.Date(
				candidate.Year(),
				candidate.Month(),
				candidate.Day(),
				startClock.Hour(),
				startClock.Minute(),
				0,
				0,
				time.UTC,
			)
			if candidateStart.After(now) {
				return candidate.Format("2006-01-02")
			}
		}
	}

	return now.AddDate(0, 0, 1).Format("2006-01-02")
}
