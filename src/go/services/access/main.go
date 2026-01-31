package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"enterprise-llm-agent-platform/shared"
	"github.com/gin-gonic/gin"
)

const serviceName = "access"

type invokeRequest struct {
	Input    string            `json:"input"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type invokeResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	RequestID string            `json:"request_id"`
	Input     string            `json:"input"`
	Output    string            `json:"output"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type errorResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}

type openAIChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []openAIChoice `json:"choices"`
}

type openAIChoice struct {
	Index        int         `json:"index"`
	Delta        openAIDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type openAIDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

func buildServer() *gin.Engine {
	router := shared.NewServer(serviceName)
	protected := router.Group("/")
	protected.Use(authMiddleware())
	protected.POST("/invoke", func(c *gin.Context) {
		var req invokeRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil || req.Input == "" {
			requestID := shared.GetRequestID(c)
			c.JSON(http.StatusBadRequest, errorResponse{
				Status:    "error",
				Service:   serviceName,
				RequestID: requestID,
				Error:     "invalid_request",
				Message:   "input is required",
			})
			return
		}

		requestID := shared.GetRequestID(c)
		response := invokeResponse{
			Status:    "ok",
			Service:   serviceName,
			RequestID: requestID,
			Input:     req.Input,
			Output:    req.Input,
			Metadata:  req.Metadata,
		}

		if wantsSSE(c) {
			if err := streamInvoke(c, response); err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse{
					Status:    "error",
					Service:   serviceName,
					RequestID: requestID,
					Error:     "stream_failed",
					Message:   "failed to stream response",
				})
			}
			return
		}

		c.JSON(http.StatusOK, response)
	})
	return router
}

func wantsSSE(c *gin.Context) bool {
	if c.Query("stream") == "true" {
		return true
	}
	accept := c.GetHeader("Accept")
	return strings.Contains(accept, "text/event-stream")
}

func streamInvoke(c *gin.Context, response invokeResponse) error {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	created := time.Now().Unix()
	model := "access-stub"
	chunkID := response.RequestID
	first := openAIChunk{
		ID:      chunkID,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []openAIChoice{
			{
				Index: 0,
				Delta: openAIDelta{
					Content: response.Output,
				},
			},
		},
	}
	if err := writeSSEData(c, first); err != nil {
		return err
	}
	stop := "stop"
	last := openAIChunk{
		ID:      chunkID,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []openAIChoice{
			{
				Index:        0,
				Delta:        openAIDelta{},
				FinishReason: &stop,
			},
		},
	}
	if err := writeSSEData(c, last); err != nil {
		return err
	}
	return writeSSEDone(c)
}

func writeSSEData(c *gin.Context, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := c.Writer.WriteString("data: " + string(payload) + "\n\n"); err != nil {
		return err
	}
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

func writeSSEDone(c *gin.Context) error {
	if _, err := c.Writer.WriteString("data: [DONE]\n\n"); err != nil {
		return err
	}
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

func main() {
	settings := shared.LoadSettings(serviceName)
	router := buildServer()
	if err := shared.RunServer(router, settings); err != nil {
		panic(err)
	}
}
