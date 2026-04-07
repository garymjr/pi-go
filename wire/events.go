package wire

import "encoding/json"

type Event interface {
	Frame
	eventType() string
}

type eventBase struct {
	Type string `json:"type"`
}

func (e eventBase) frameType() string { return e.Type }
func (e eventBase) eventType() string { return e.Type }

type AgentStartEvent struct{ eventBase }

type AgentEndEvent struct {
	eventBase
	Messages []AgentMessage `json:"messages"`
}

func (e *AgentEndEvent) UnmarshalJSON(data []byte) error {
	type alias struct {
		eventBase
		Messages []json.RawMessage `json:"messages"`
	}
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.eventBase = aux.eventBase
	e.Messages = make([]AgentMessage, 0, len(aux.Messages))
	for _, raw := range aux.Messages {
		msg, err := DecodeAgentMessage(raw)
		if err != nil {
			return err
		}
		e.Messages = append(e.Messages, msg)
	}
	return nil
}

type TurnStartEvent struct{ eventBase }

type TurnEndEvent struct {
	eventBase
	Message     AgentMessage        `json:"message"`
	ToolResults []ToolResultMessage `json:"toolResults"`
}

func (e *TurnEndEvent) UnmarshalJSON(data []byte) error {
	type alias struct {
		eventBase
		Message     json.RawMessage     `json:"message"`
		ToolResults []ToolResultMessage `json:"toolResults"`
	}
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	msg, err := DecodeAgentMessage(aux.Message)
	if err != nil {
		return err
	}
	e.eventBase = aux.eventBase
	e.Message = msg
	e.ToolResults = aux.ToolResults
	return nil
}

type MessageStartEvent struct {
	eventBase
	Message AgentMessage `json:"message"`
}

func (e *MessageStartEvent) UnmarshalJSON(data []byte) error {
	type alias struct {
		eventBase
		Message json.RawMessage `json:"message"`
	}
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	msg, err := DecodeAgentMessage(aux.Message)
	if err != nil {
		return err
	}
	e.eventBase = aux.eventBase
	e.Message = msg
	return nil
}

type MessageEndEvent struct {
	eventBase
	Message AgentMessage `json:"message"`
}

func (e *MessageEndEvent) UnmarshalJSON(data []byte) error {
	type alias struct {
		eventBase
		Message json.RawMessage `json:"message"`
	}
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	msg, err := DecodeAgentMessage(aux.Message)
	if err != nil {
		return err
	}
	e.eventBase = aux.eventBase
	e.Message = msg
	return nil
}

type AssistantMessageEvent interface{ assistantMessageEventType() string }

type assistantMessageEventBase struct {
	Type string `json:"type"`
}

func (e assistantMessageEventBase) assistantMessageEventType() string { return e.Type }

type AssistantMessageEventStart struct {
	assistantMessageEventBase
	Partial AssistantMessage `json:"partial"`
}

type AssistantMessageEventTextStart struct {
	assistantMessageEventBase
	ContentIndex int              `json:"contentIndex"`
	Partial      AssistantMessage `json:"partial"`
}

type AssistantMessageEventTextDelta struct {
	assistantMessageEventBase
	ContentIndex int              `json:"contentIndex"`
	Delta        string           `json:"delta"`
	Partial      AssistantMessage `json:"partial"`
}

type AssistantMessageEventTextEnd struct {
	assistantMessageEventBase
	ContentIndex int              `json:"contentIndex"`
	Content      string           `json:"content"`
	Partial      AssistantMessage `json:"partial"`
}

type AssistantMessageEventThinkingStart = AssistantMessageEventTextStart
type AssistantMessageEventThinkingDelta = AssistantMessageEventTextDelta
type AssistantMessageEventThinkingEnd = AssistantMessageEventTextEnd

type AssistantMessageEventToolCallStart = AssistantMessageEventTextStart

type AssistantMessageEventToolCallDelta struct {
	assistantMessageEventBase
	ContentIndex int              `json:"contentIndex"`
	Delta        string           `json:"delta"`
	Partial      AssistantMessage `json:"partial"`
}

type AssistantMessageEventToolCallEnd struct {
	assistantMessageEventBase
	ContentIndex int              `json:"contentIndex"`
	ToolCall     ToolCallContent  `json:"toolCall"`
	Partial      AssistantMessage `json:"partial"`
}

type AssistantMessageEventDone struct {
	assistantMessageEventBase
	Reason  string           `json:"reason"`
	Message AssistantMessage `json:"message"`
}

type AssistantMessageEventError struct {
	assistantMessageEventBase
	Reason string           `json:"reason"`
	Error  AssistantMessage `json:"error"`
}

type UnknownAssistantMessageEvent struct {
	Type string          `json:"type,omitempty"`
	Raw  json.RawMessage `json:"-"`
}

func (e UnknownAssistantMessageEvent) assistantMessageEventType() string { return e.Type }

type MessageUpdateEvent struct {
	eventBase
	Message               AgentMessage          `json:"message"`
	AssistantMessageEvent AssistantMessageEvent `json:"assistantMessageEvent"`
}

func (e *MessageUpdateEvent) UnmarshalJSON(data []byte) error {
	type alias struct {
		eventBase
		Message               json.RawMessage `json:"message"`
		AssistantMessageEvent json.RawMessage `json:"assistantMessageEvent"`
	}
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	msg, err := DecodeAgentMessage(aux.Message)
	if err != nil {
		return err
	}
	delta, err := DecodeAssistantMessageEvent(aux.AssistantMessageEvent)
	if err != nil {
		return err
	}
	e.eventBase = aux.eventBase
	e.Message = msg
	e.AssistantMessageEvent = delta
	return nil
}

func DecodeAssistantMessageEvent(data json.RawMessage) (AssistantMessageEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	switch envelope.Type {
	case "start":
		var v AssistantMessageEventStart
		return &v, json.Unmarshal(data, &v)
	case "text_start":
		var v AssistantMessageEventTextStart
		return &v, json.Unmarshal(data, &v)
	case "text_delta":
		var v AssistantMessageEventTextDelta
		return &v, json.Unmarshal(data, &v)
	case "text_end":
		var v AssistantMessageEventTextEnd
		return &v, json.Unmarshal(data, &v)
	case "thinking_start":
		var v AssistantMessageEventThinkingStart
		return &v, json.Unmarshal(data, &v)
	case "thinking_delta":
		var v AssistantMessageEventThinkingDelta
		return &v, json.Unmarshal(data, &v)
	case "thinking_end":
		var v AssistantMessageEventThinkingEnd
		return &v, json.Unmarshal(data, &v)
	case "toolcall_start":
		var v AssistantMessageEventToolCallStart
		return &v, json.Unmarshal(data, &v)
	case "toolcall_delta":
		var v AssistantMessageEventToolCallDelta
		return &v, json.Unmarshal(data, &v)
	case "toolcall_end":
		var v AssistantMessageEventToolCallEnd
		return &v, json.Unmarshal(data, &v)
	case "done":
		var v AssistantMessageEventDone
		return &v, json.Unmarshal(data, &v)
	case "error":
		var v AssistantMessageEventError
		return &v, json.Unmarshal(data, &v)
	default:
		return UnknownAssistantMessageEvent{Type: envelope.Type, Raw: append(json.RawMessage(nil), data...)}, nil
	}
}

type ToolExecutionStartEvent struct {
	eventBase
	ToolCallID string          `json:"toolCallId"`
	ToolName   string          `json:"toolName"`
	Args       json.RawMessage `json:"args"`
}

type ToolExecutionUpdateEvent struct {
	eventBase
	ToolCallID    string          `json:"toolCallId"`
	ToolName      string          `json:"toolName"`
	Args          json.RawMessage `json:"args"`
	PartialResult json.RawMessage `json:"partialResult"`
}

type ToolExecutionEndEvent struct {
	eventBase
	ToolCallID string          `json:"toolCallId"`
	ToolName   string          `json:"toolName"`
	Result     json.RawMessage `json:"result"`
	IsError    bool            `json:"isError"`
}

type QueueUpdateEvent struct {
	eventBase
	Steering []string `json:"steering"`
	FollowUp []string `json:"followUp"`
}

type CompactionStartEvent struct {
	eventBase
	Reason string `json:"reason"`
}

type CompactionEndEvent struct {
	eventBase
	Reason       string            `json:"reason"`
	Result       *CompactionResult `json:"result"`
	Aborted      bool              `json:"aborted"`
	WillRetry    bool              `json:"willRetry"`
	ErrorMessage string            `json:"errorMessage,omitempty"`
}

type AutoRetryStartEvent struct {
	eventBase
	Attempt      int    `json:"attempt"`
	MaxAttempts  int    `json:"maxAttempts"`
	DelayMS      int    `json:"delayMs"`
	ErrorMessage string `json:"errorMessage"`
}

type AutoRetryEndEvent struct {
	eventBase
	Success    bool   `json:"success"`
	Attempt    int    `json:"attempt"`
	FinalError string `json:"finalError,omitempty"`
}

type ExtensionErrorEvent struct {
	eventBase
	ExtensionPath string `json:"extensionPath"`
	Event         string `json:"event"`
	Error         string `json:"error"`
}
