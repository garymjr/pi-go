package wire

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder { return &Decoder{r: bufio.NewReader(r)} }

func (d *Decoder) Decode() (Frame, error) {
	line, err := d.r.ReadBytes('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && len(line) == 0 {
			return nil, io.EOF
		}
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	line = bytes.TrimSuffix(line, []byte("\n"))
	line = bytes.TrimSuffix(line, []byte("\r"))
	if len(bytes.TrimSpace(line)) == 0 {
		if errors.Is(err, io.EOF) {
			return nil, io.EOF
		}
		return d.Decode()
	}
	return decodeFrame(line)
}

func decodeFrame(line []byte) (Frame, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(line, &envelope); err != nil {
		return nil, fmt.Errorf("invalid json frame: %w", err)
	}
	switch envelope.Type {
	case "response":
		return decodeResponse(line)
	case "extension_ui_request":
		return decodeUIRequest(line)
	case "agent_start":
		var v AgentStartEvent
		return &v, json.Unmarshal(line, &v)
	case "agent_end":
		var v AgentEndEvent
		return &v, json.Unmarshal(line, &v)
	case "turn_start":
		var v TurnStartEvent
		return &v, json.Unmarshal(line, &v)
	case "turn_end":
		var v TurnEndEvent
		return &v, json.Unmarshal(line, &v)
	case "message_start":
		var v MessageStartEvent
		return &v, json.Unmarshal(line, &v)
	case "message_update":
		var v MessageUpdateEvent
		return &v, json.Unmarshal(line, &v)
	case "message_end":
		var v MessageEndEvent
		return &v, json.Unmarshal(line, &v)
	case "tool_execution_start":
		var v ToolExecutionStartEvent
		return &v, json.Unmarshal(line, &v)
	case "tool_execution_update":
		var v ToolExecutionUpdateEvent
		return &v, json.Unmarshal(line, &v)
	case "tool_execution_end":
		var v ToolExecutionEndEvent
		return &v, json.Unmarshal(line, &v)
	case "queue_update":
		var v QueueUpdateEvent
		return &v, json.Unmarshal(line, &v)
	case "compaction_start":
		var v CompactionStartEvent
		return &v, json.Unmarshal(line, &v)
	case "compaction_end":
		var v CompactionEndEvent
		return &v, json.Unmarshal(line, &v)
	case "auto_retry_start":
		var v AutoRetryStartEvent
		return &v, json.Unmarshal(line, &v)
	case "auto_retry_end":
		var v AutoRetryEndEvent
		return &v, json.Unmarshal(line, &v)
	case "extension_error":
		var v ExtensionErrorEvent
		return &v, json.Unmarshal(line, &v)
	default:
		return UnknownFrame{Type: envelope.Type, Raw: append(json.RawMessage(nil), line...)}, nil
	}
}

func decodeResponse(line []byte) (Frame, error) {
	var envelope struct {
		Type    string `json:"type"`
		Command string `json:"command"`
		Success bool   `json:"success"`
	}
	if err := json.Unmarshal(line, &envelope); err != nil {
		return nil, err
	}
	if envelope.Command == "" {
		return nil, fmt.Errorf("response missing command")
	}
	if !envelope.Success {
		var v ErrorResponse
		if err := json.Unmarshal(line, &v); err != nil {
			return nil, err
		}
		return &v, nil
	}
	switch envelope.Command {
	case "prompt":
		var v PromptResponse
		return &v, json.Unmarshal(line, &v)
	case "steer":
		var v SteerResponse
		return &v, json.Unmarshal(line, &v)
	case "follow_up":
		var v FollowUpResponse
		return &v, json.Unmarshal(line, &v)
	case "abort":
		var v AbortResponse
		return &v, json.Unmarshal(line, &v)
	case "new_session":
		var v NewSessionResponse
		return &v, json.Unmarshal(line, &v)
	case "get_state":
		var v GetStateResponse
		return &v, json.Unmarshal(line, &v)
	case "set_model":
		var v SetModelResponse
		return &v, json.Unmarshal(line, &v)
	case "cycle_model":
		var v CycleModelResponse
		return &v, json.Unmarshal(line, &v)
	case "get_available_models":
		var v GetAvailableModelsResponse
		return &v, json.Unmarshal(line, &v)
	case "set_thinking_level":
		var v SetThinkingLevelResponse
		return &v, json.Unmarshal(line, &v)
	case "cycle_thinking_level":
		var v CycleThinkingLevelResponse
		return &v, json.Unmarshal(line, &v)
	case "set_steering_mode":
		var v SetSteeringModeResponse
		return &v, json.Unmarshal(line, &v)
	case "set_follow_up_mode":
		var v SetFollowUpModeResponse
		return &v, json.Unmarshal(line, &v)
	case "compact":
		var v CompactResponse
		return &v, json.Unmarshal(line, &v)
	case "set_auto_compaction":
		var v SetAutoCompactionResponse
		return &v, json.Unmarshal(line, &v)
	case "set_auto_retry":
		var v SetAutoRetryResponse
		return &v, json.Unmarshal(line, &v)
	case "abort_retry":
		var v AbortRetryResponse
		return &v, json.Unmarshal(line, &v)
	case "bash":
		var v BashResponse
		return &v, json.Unmarshal(line, &v)
	case "abort_bash":
		var v AbortBashResponse
		return &v, json.Unmarshal(line, &v)
	case "get_session_stats":
		var v GetSessionStatsResponse
		return &v, json.Unmarshal(line, &v)
	case "export_html":
		var v ExportHTMLResponse
		return &v, json.Unmarshal(line, &v)
	case "switch_session":
		var v SwitchSessionResponse
		return &v, json.Unmarshal(line, &v)
	case "fork":
		var v ForkResponse
		return &v, json.Unmarshal(line, &v)
	case "get_fork_messages":
		var v GetForkMessagesResponse
		return &v, json.Unmarshal(line, &v)
	case "get_last_assistant_text":
		var v GetLastAssistantTextResponse
		return &v, json.Unmarshal(line, &v)
	case "set_session_name":
		var v SetSessionNameResponse
		return &v, json.Unmarshal(line, &v)
	case "get_messages":
		var v GetMessagesResponse
		return &v, json.Unmarshal(line, &v)
	case "get_commands":
		var v GetCommandsResponse
		return &v, json.Unmarshal(line, &v)
	default:
		var v ErrorResponse
		if err := json.Unmarshal(line, &v); err != nil {
			return nil, err
		}
		return &v, nil
	}
}

func decodeUIRequest(line []byte) (Frame, error) {
	var envelope struct {
		Type   string `json:"type"`
		Method string `json:"method"`
	}
	if err := json.Unmarshal(line, &envelope); err != nil {
		return nil, err
	}
	if envelope.Method == "" {
		return nil, fmt.Errorf("extension_ui_request missing method")
	}
	switch envelope.Method {
	case "select":
		var v SelectUIRequest
		return &v, json.Unmarshal(line, &v)
	case "confirm":
		var v ConfirmUIRequest
		return &v, json.Unmarshal(line, &v)
	case "input":
		var v InputUIRequest
		return &v, json.Unmarshal(line, &v)
	case "editor":
		var v EditorUIRequest
		return &v, json.Unmarshal(line, &v)
	case "notify":
		var v NotifyUIRequest
		return &v, json.Unmarshal(line, &v)
	case "setStatus":
		var v SetStatusUIRequest
		return &v, json.Unmarshal(line, &v)
	case "setWidget":
		var v SetWidgetUIRequest
		return &v, json.Unmarshal(line, &v)
	case "setTitle":
		var v SetTitleUIRequest
		return &v, json.Unmarshal(line, &v)
	case "set_editor_text":
		var v SetEditorTextUIRequest
		return &v, json.Unmarshal(line, &v)
	default:
		return UnknownUIRequest{Method: envelope.Method, Raw: append(json.RawMessage(nil), line...)}, nil
	}
}
