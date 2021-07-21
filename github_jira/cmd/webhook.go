package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// For non-crucial messaging that doesn't represent program failure on webhook fail.
func sendWebhookMessage(url, message string) {
	msg, err := json.Marshal(map[string]string{"text": message})
	if err != nil {
		fmt.Printf("Unable to send log to webhook: %+v\n", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(msg))
	if err != nil {
		fmt.Printf("Unable to send log to webhook: %+v\n", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending log %s\n", err)
	}
	if resp.StatusCode >= 400 {
		fmt.Printf("Sending log failed with status code %d\n", resp.StatusCode)
	}
}
