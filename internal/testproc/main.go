package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var scenario string
	var mode string
	flag.StringVar(&scenario, "scenario", "startup_success", "test scenario")
	flag.StringVar(&mode, "mode", "", "ignored mode")
	flag.Parse()
	_ = mode

	switch scenario {
	case "startup_exit":
		_, _ = os.Stderr.WriteString("startup failed\n")
		os.Exit(2)
	case "invalid_json":
		time.Sleep(700 * time.Millisecond)
		fmt.Fprintln(os.Stdout, "{bad")
		time.Sleep(5 * time.Second)
		return
	}

	r := bufio.NewReader(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	var saved []map[string]any
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return
		}
		var msg map[string]any
		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}
		typ, _ := msg["type"].(string)
		id, _ := msg["id"].(string)

		switch scenario {
		case "get_state":
			if typ == "get_state" {
				_ = enc.Encode(map[string]any{
					"type":    "response",
					"id":      id,
					"command": "get_state",
					"success": true,
					"data": map[string]any{
						"model": nil, "thinkingLevel": "medium", "isStreaming": false, "isCompacting": false,
						"steeringMode": "all", "followUpMode": "one-at-a-time", "sessionId": "s1",
						"autoCompactionEnabled": true, "messageCount": 0, "pendingMessageCount": 0,
					},
				})
			}
		case "concurrent":
			saved = append(saved, msg)
			if len(saved) == 2 {
				_ = enc.Encode(map[string]any{"type": "response", "id": saved[1]["id"], "command": saved[1]["type"], "success": true})
				_ = enc.Encode(map[string]any{"type": "response", "id": saved[0]["id"], "command": saved[0]["type"], "success": true})
			}
		case "concurrent_large":
			saved = append(saved, msg)
			if len(saved) == 2 {
				for _, item := range saved {
					_ = enc.Encode(map[string]any{"type": "response", "id": item["id"], "command": item["type"], "success": true})
				}
			}
		case "notification_get_state":
			if typ == "prompt" {
				_ = enc.Encode(map[string]any{"type": "response", "id": id, "command": "prompt", "success": true})
				_ = enc.Encode(map[string]any{"type": "agent_start"})
				continue
			}
			if typ == "get_state" {
				_ = enc.Encode(map[string]any{
					"type":    "response",
					"id":      id,
					"command": "get_state",
					"success": true,
					"data": map[string]any{
						"model": nil, "thinkingLevel": "medium", "isStreaming": false, "isCompacting": false,
						"steeringMode": "all", "followUpMode": "one-at-a-time", "sessionId": "s1",
						"autoCompactionEnabled": true, "messageCount": 0, "pendingMessageCount": 0,
					},
				})
				_ = enc.Encode(map[string]any{"type": "agent_end", "messages": []any{}})
			}
		case "prompt_stream", "ui_roundtrip":
			if typ == "prompt" {
				message := map[string]any{"role": "assistant", "content": []any{map[string]any{"type": "text", "text": "hi"}}, "api": "x", "provider": "p", "model": "m", "usage": map[string]any{"input": 1, "output": 1, "cacheRead": 0, "cacheWrite": 0, "totalTokens": 2, "cost": map[string]any{"input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0, "total": 0}}, "stopReason": "stop", "timestamp": 1}
				_ = enc.Encode(map[string]any{"type": "response", "id": id, "command": "prompt", "success": true})
				_ = enc.Encode(map[string]any{"type": "agent_start"})
				_ = enc.Encode(map[string]any{
					"type":                  "message_update",
					"message":               message,
					"assistantMessageEvent": map[string]any{"type": "text_delta", "contentIndex": 0, "delta": "hi", "partial": message},
				})
				if scenario == "ui_roundtrip" {
					_ = enc.Encode(map[string]any{"type": "extension_ui_request", "id": "ui1", "method": "input", "title": "Value"})
				} else {
					_ = enc.Encode(map[string]any{"type": "agent_end", "messages": []any{}})
				}
			}
			if scenario == "ui_roundtrip" && typ == "extension_ui_response" {
				_ = enc.Encode(map[string]any{"type": "agent_end", "messages": []any{map[string]any{"role": "user", "content": "ok", "timestamp": 1}}})
			}
		case "command_failure":
			_ = enc.Encode(map[string]any{"type": "response", "id": id, "command": typ, "success": false, "error": "boom"})
		case "process_exit_pending":
			os.Exit(3)
		default:
			if typ == "extension_ui_response" {
				continue
			}
			_ = enc.Encode(map[string]any{"type": "response", "id": id, "command": typ, "success": true})
			if typ == "prompt" {
				time.Sleep(10 * time.Millisecond)
				_ = enc.Encode(map[string]any{"type": "agent_end", "messages": []any{}})
			}
		}
	}
}
