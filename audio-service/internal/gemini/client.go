package gemini

import (
    "fmt"
    "net/url"
    "github.com/gorilla/websocket"
)

type GeminiClient struct {
    conn *websocket.Conn
}

func (gc *GeminiClient) Close() error {
    return gc.conn.Close()
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
    u := url.URL{Scheme: "wss", Host: "generativelanguage.googleapis.com", Path: "/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent", RawQuery: "key=" + apiKey}
    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("dial: %w", err)
    }
    return &GeminiClient{conn: c}, nil
}