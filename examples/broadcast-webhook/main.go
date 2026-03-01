// Package main demonstrates receiving OpenRouter Broadcast webhook traces.
//
// This example starts an HTTP server that listens for OTLP JSON trace payloads
// sent by OpenRouter's Broadcast feature configured with a Webhook destination.
//
// Usage:
//
//	go run examples/broadcast-webhook/main.go
//
// Then configure your OpenRouter Broadcast to send traces to http://<host>:8080/webhook
// and test with:
//
//	curl -X POST -H "Content-Type: application/json" \
//	  -d '{"resourceSpans":[]}' http://localhost:8080/webhook
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hra42/openrouter-go"
)

func main() {
	handler := openrouter.BroadcastWebhookHandler(func(traces []openrouter.BroadcastTrace) {
		for _, tr := range traces {
			fmt.Printf("--- Trace %s / Span %s ---\n", tr.TraceID, tr.SpanID)
			fmt.Printf("  Name:     %s\n", tr.SpanName)
			fmt.Printf("  Model:    %s\n", tr.ResponseModel)
			fmt.Printf("  Duration: %s\n", tr.Duration)
			fmt.Printf("  Tokens:   %d prompt + %d completion = %d total\n",
				tr.InputTokens, tr.OutputTokens, tr.TotalTokens)
			fmt.Printf("  Cost:     $%.6f\n", tr.TotalCost)
			if tr.UserID != "" {
				fmt.Printf("  User:     %s\n", tr.UserID)
			}
			if tr.SessionID != "" {
				fmt.Printf("  Session:  %s\n", tr.SessionID)
			}
			if len(tr.Metadata) > 0 {
				fmt.Printf("  Metadata: %v\n", tr.Metadata)
			}
			fmt.Println()
		}
	})

	http.Handle("/webhook", handler)

	fmt.Println("Broadcast webhook server listening on :8080/webhook")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
