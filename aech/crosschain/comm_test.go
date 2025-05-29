package crosschain

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// Helper function to capture log output
func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf) // Redirect log output to buffer
	f()                  // Execute the function that logs
	log.SetOutput(os.Stderr) // Reset log output to default
	return buf.String()
}

func TestSendMessage(t *testing.T) {
	msg := CrossChainMessage{
		FromPlaneID: "TestPlaneFrom",
		ToPlaneID:   "TestPlaneTo",
		BlockHash:   "testBlockHash123",
		Payload:     "Test payload for send",
		Timestamp:   time.Now(),
	}

	logOutput := captureLogOutput(func() {
		err := SendMessage(msg)
		if err != nil {
			t.Errorf("SendMessage returned an unexpected error: %v", err)
		}
	})

	// Check if the log output for SendMessage contains expected parts
	if !strings.Contains(logOutput, "Crosschain Mock Send:") {
		t.Errorf("Log output from SendMessage does not contain 'Crosschain Mock Send:'. Log: %s", logOutput)
	}
	if !strings.Contains(logOutput, "from Plane TestPlaneFrom to Plane TestPlaneTo") {
		t.Errorf("Log output from SendMessage does not contain correct plane IDs. Log: %s", logOutput)
	}
	if !strings.Contains(logOutput, "regarding BlockHash testBlockHash123") {
		t.Errorf("Log output from SendMessage does not contain correct BlockHash. Log: %s", logOutput)
	}
	if !strings.Contains(logOutput, "Payload: Test payload for send") {
		t.Errorf("Log output from SendMessage does not contain correct Payload. Log: %s", logOutput)
	}

	// Also check if the ReceiveMessage log is present due to the loopback call in SendMessage
	if !strings.Contains(logOutput, "Crosschain Mock Receive:") {
		t.Errorf("Log output from SendMessage (via loopback) does not contain 'Crosschain Mock Receive:'. Log: %s", logOutput)
	}
    if !strings.Contains(logOutput, "Message for Plane TestPlaneTo from Plane TestPlaneFrom") {
        t.Errorf("Log output from ReceiveMessage (via loopback) does not contain correct plane IDs. Log: %s", logOutput)
    }
}

func TestReceiveMessage(t *testing.T) {
	msg := CrossChainMessage{
		FromPlaneID: "AnotherPlaneFrom",
		ToPlaneID:   "AnotherPlaneTo",
		BlockHash:   "anotherBlockHash456",
		Payload:     "Test payload for receive",
		Timestamp:   time.Now(),
	}

	logOutput := captureLogOutput(func() {
		ReceiveMessage(msg)
	})

	// Check if the log output for ReceiveMessage contains expected parts
	if !strings.Contains(logOutput, "Crosschain Mock Receive:") {
		t.Errorf("Log output from ReceiveMessage does not contain 'Crosschain Mock Receive:'. Log: %s", logOutput)
	}
	if !strings.Contains(logOutput, "Message for Plane AnotherPlaneTo from Plane AnotherPlaneFrom") {
		t.Errorf("Log output from ReceiveMessage does not contain correct plane IDs. Log: %s", logOutput)
	}
	if !strings.Contains(logOutput, "regarding BlockHash anotherBlockHash456") {
		t.Errorf("Log output from ReceiveMessage does not contain correct BlockHash. Log: %s", logOutput)
	}
	if !strings.Contains(logOutput, "Payload: Test payload for receive") {
		t.Errorf("Log output from ReceiveMessage does not contain correct Payload. Log: %s", logOutput)
	}
}
