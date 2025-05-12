package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/welovemedia/ffmate/sev"
)

func TestWebsocketController(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := sev.New("test", "", "", 3000)
	controller := &WebsocketController{
		Prefix: "",
	}
	controller.Setup(s)

	server := httptest.NewServer(s.Gin())
	defer server.Close()

	// Convert http://... to ws://...
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/v1/ws"

	t.Run("websocket connection", func(t *testing.T) {
		// Connect to the server
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection: %v", err)
		}
		defer ws.Close()

		// Test connection is established
		if ws == nil {
			t.Error("expected websocket connection to be established")
		}
	})
}
