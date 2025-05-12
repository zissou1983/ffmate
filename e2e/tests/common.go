package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/welovemedia/ffmate/internal/dto"
)

type WebhookCall struct {
	Event dto.WebhookEvent
	Data  interface{}
}

func SetupWebhookServer(expectedCalls int) (string, chan WebhookCall) {
	calls := make(chan WebhookCall, expectedCalls)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		calls <- WebhookCall{
			Event: dto.WebhookEvent(payload["event"].(string)),
			Data:  payload["data"],
		}
		w.WriteHeader(http.StatusOK)
	}))

	return server.URL, calls
}

func RegisterWebhook(url string, event dto.WebhookEvent) error {
	_, err := resty.New().R().
		SetBody(&dto.NewWebhook{
			Event: event,
			Url:   url,
		}).
		Post("http://localhost:3000/api/v1/webhooks")

	if err != nil {
		return fmt.Errorf("failed to register webhook: %v", err)
	} else {
		fmt.Printf("Webhook '%s' registered\n", event)
	}
	return nil
}

func WaitForWebhook(calls chan WebhookCall, expectedEvent dto.WebhookEvent) {
	fmt.Printf("Waiting for webhook event '%s'...\n", expectedEvent)
	select {
	case call := <-calls:
		if call.Event != expectedEvent {
			fmt.Printf("Expected webhook event %s, got %s", expectedEvent, call.Event)
			os.Exit(1)
		} else {
			fmt.Printf("Webhook '%s' received\n", call.Event)
		}
	case <-time.After(time.Second * 5):
		fmt.Printf("Timeout waiting for webhook event %s", expectedEvent)
		os.Exit(1)

	}
}
