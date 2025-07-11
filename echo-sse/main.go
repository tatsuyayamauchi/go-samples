package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// see: https://github.com/langgenius/dify-docs/blob/main/ja-jp/openapi_chat.json
type DifyRequest struct {
	Inputs         map[string]any `json:"inputs"`
	Query          string         `json:"query"`
	ResponseMode   string         `json:"response_mode"`
	ConversationID string         `json:"conversation_id"`
	User           string         `json:"user"`
	Files          []any          `json:"files"`
}

type DifyEvent struct {
	Event          string `json:"event"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskID         string `json:"task_id"`
	ID             string `json:"id"`
	Answer         string `json:"answer"`
	Metadata       any    `json:"metadata"`
}

type SSEClient struct {
	apiKey  string
	baseURL string
}

func NewSSEClient(apiKey, baseURL string) *SSEClient {
	return &SSEClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (c *SSEClient) StreamChat(query string, conversationID string, callback func(DifyEvent)) error {
	reqBody := DifyRequest{
		Inputs:         make(map[string]any),
		Query:          query,
		ResponseMode:   "streaming",
		ConversationID: conversationID,
		User:           "go-client",
		Files:          []any{},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v1/chat-messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{
		Timeout: 0,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// 空行はイベントの区切り
		if line == "" {
			if buffer.Len() > 0 {
				if err := c.processEventData(buffer.String(), callback); err != nil {
					log.Printf("Error processing event: %v", err)
				}
				buffer.Reset()
			}
			continue
		}

		// "data: " で始まる行を処理
		if after, ok := strings.CutPrefix(line, "data: "); ok {
			data := after
			buffer.WriteString(data)
		}
	}

	if buffer.Len() > 0 {
		if err := c.processEventData(buffer.String(), callback); err != nil {
			log.Printf("Error processing final event: %v", err)
		}
	}

	return scanner.Err()
}

func (c *SSEClient) processEventData(data string, callback func(DifyEvent)) error {
	var event DifyEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	callback(event)
	return nil
}

func chatHandler(client *SSEClient) func(echo.Context) error {
	type ChatRequest struct {
		Query          string `json:"query"`
		ConversationID string `json:"conversation_id"`
	}

	return func(c echo.Context) error {
		var req ChatRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		c.Response().Header().Set("Content-Type", "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")

		callback := func(event DifyEvent) {
			eventData, _ := json.Marshal(event)
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", string(eventData))
			c.Response().Flush()
		}

		if err := client.StreamChat(req.Query, req.ConversationID, callback); err != nil {
			log.Printf("Error streaming chat: %v", err)
			fmt.Fprintf(c.Response().Writer, "event: error\ndata: %s\n\n", err.Error())
			return nil
		}

		fmt.Fprintf(c.Response().Writer, "event: done\ndata: {}\n\n")
		return nil
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	apiKey := os.Getenv("DIFY_API_KEY")
	baseURL := os.Getenv("DIFY_API_ENDPOINT")
	if baseURL == "" {
		baseURL = "https://api.dify.ai"
	}

	client := NewSSEClient(apiKey, baseURL)

	e.POST("/chat", chatHandler(client))

	e.Logger.Fatal(e.Start(":8080"))
}
