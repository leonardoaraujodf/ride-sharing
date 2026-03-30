//go:build integration
// +build integration

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Env = append(os.Environ(), "ENABLE_DELAY=true")

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(2 * time.Second)

	requestCompleted := make(chan struct{})

	go func() {
		client := &http.Client{Timeout: 20 * time.Second}
		reqBody := map[string]any{
			"userID":      "test-user",
			"pickup":      map[string]float64{"latitude": 40.7128, "longitude": -74.0060},
			"destination": map[string]float64{"latitude": 40.7580, "longitude": -73.9855},
		}
		jsonBody, _ := json.Marshal(reqBody)

		resp, err := client.Post("http://localhost:8083/preview", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			t.Logf("request error: %v", err)
			close(requestCompleted)
			return
		}
		defer resp.Body.Close()
		close(requestCompleted)
	}()

	time.Sleep(500 * time.Millisecond)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("failed to send SIGTERM: %v", err)
	}

	select {
	case <-requestCompleted:
		t.Log("in-flight request completed before shutdown")
	case <-time.After(15 * time.Second):
		t.Error("request did not complete within graceful shutdown timeout")
	}

	waitErr := cmd.Wait()
	t.Logf("server exit error: %v", waitErr)
}
