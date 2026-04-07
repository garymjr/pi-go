package wire

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ContentPart interface{ contentType() string }

type TextContent struct {
	Type          string `json:"type"`
	Text          string `json:"text"`
	TextSignature string `json:"textSignature,omitempty"`
}

func (TextContent) contentType() string { return "text" }

type ImageContent struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	MimeType string `json:"mimeType"`
}

func (ImageContent) contentType() string { return "image" }

type ThinkingContent struct {
	Type              string `json:"type"`
	Thinking          string `json:"thinking"`
	ThinkingSignature string `json:"thinkingSignature,omitempty"`
	Redacted          bool   `json:"redacted,omitempty"`
}

func (ThinkingContent) contentType() string { return "thinking" }

type ToolCallContent struct {
	Type             string          `json:"type"`
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Arguments        json.RawMessage `json:"arguments"`
	ThoughtSignature string          `json:"thoughtSignature,omitempty"`
}

func (ToolCallContent) contentType() string { return "toolCall" }

type UnknownContent struct {
	Type string          `json:"type,omitempty"`
	Raw  json.RawMessage `json:"-"`
}

func (u UnknownContent) contentType() string { return u.Type }

type UserContent struct {
	String *string
	Parts  []ContentPart
}

func (c UserContent) MarshalJSON() ([]byte, error) {
	if c.String != nil {
		return json.Marshal(*c.String)
	}
	return marshalContentParts(c.Parts)
}

func (c *UserContent) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*c = UserContent{}
		return nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return err
		}
		*c = UserContent{String: &s}
		return nil
	}
	parts, err := decodeContentParts(trimmed)
	if err != nil {
		return err
	}
	*c = UserContent{Parts: parts}
	return nil
}

type AgentMessage interface{ messageRole() string }

type UserMessage struct {
	Role      string      `json:"role"`
	Content   UserContent `json:"content"`
	Timestamp int64       `json:"timestamp"`
}

func (UserMessage) messageRole() string { return "user" }

type UsageCost struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheRead  float64 `json:"cacheRead"`
	CacheWrite float64 `json:"cacheWrite"`
	Total      float64 `json:"total"`
}

type Usage struct {
	Input       int64     `json:"input"`
	Output      int64     `json:"output"`
	CacheRead   int64     `json:"cacheRead"`
	CacheWrite  int64     `json:"cacheWrite"`
	TotalTokens int64     `json:"totalTokens"`
	Cost        UsageCost `json:"cost"`
}

type AssistantMessage struct {
	Role         string          `json:"role"`
	Content      []ContentPart   `json:"content"`
	API          string          `json:"api"`
	Provider     string          `json:"provider"`
	Model        string          `json:"model"`
	ResponseID   string          `json:"responseId,omitempty"`
	Usage        Usage           `json:"usage"`
	StopReason   string          `json:"stopReason"`
	ErrorMessage string          `json:"errorMessage,omitempty"`
	Timestamp    int64           `json:"timestamp"`
	Raw          json.RawMessage `json:"-"`
}

func (m AssistantMessage) messageRole() string { return "assistant" }

func (m AssistantMessage) MarshalJSON() ([]byte, error) {
	type alias AssistantMessage
	return json.Marshal(struct {
		alias
		Content json.RawMessage `json:"content"`
	}{alias: alias(m), Content: mustMarshalContentParts(m.Content)})
}

func (m *AssistantMessage) UnmarshalJSON(data []byte) error {
	type alias AssistantMessage
	aux := struct {
		alias
		Content json.RawMessage `json:"content"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	parts, err := decodeContentParts(aux.Content)
	if err != nil {
		return err
	}
	*m = AssistantMessage(aux.alias)
	m.Content = parts
	m.Raw = append(json.RawMessage(nil), data...)
	return nil
}

type ToolResultMessage struct {
	Role       string          `json:"role"`
	ToolCallID string          `json:"toolCallId"`
	ToolName   string          `json:"toolName"`
	Content    []ContentPart   `json:"content"`
	Details    json.RawMessage `json:"details,omitempty"`
	IsError    bool            `json:"isError"`
	Timestamp  int64           `json:"timestamp"`
}

func (ToolResultMessage) messageRole() string { return "toolResult" }

func (m ToolResultMessage) MarshalJSON() ([]byte, error) {
	type alias ToolResultMessage
	return json.Marshal(struct {
		alias
		Content json.RawMessage `json:"content"`
	}{alias: alias(m), Content: mustMarshalContentParts(m.Content)})
}

func (m *ToolResultMessage) UnmarshalJSON(data []byte) error {
	type alias ToolResultMessage
	aux := struct {
		alias
		Content json.RawMessage `json:"content"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	parts, err := decodeContentParts(aux.Content)
	if err != nil {
		return err
	}
	*m = ToolResultMessage(aux.alias)
	m.Content = parts
	return nil
}

type BashExecutionMessage struct {
	Role               string `json:"role"`
	Command            string `json:"command"`
	Output             string `json:"output"`
	ExitCode           *int   `json:"exitCode"`
	Cancelled          bool   `json:"cancelled"`
	Truncated          bool   `json:"truncated"`
	FullOutputPath     string `json:"fullOutputPath,omitempty"`
	Timestamp          int64  `json:"timestamp"`
	ExcludeFromContext bool   `json:"excludeFromContext,omitempty"`
}

func (BashExecutionMessage) messageRole() string { return "bashExecution" }

type CustomMessage struct {
	Role       string          `json:"role"`
	CustomType string          `json:"customType"`
	Content    json.RawMessage `json:"content"`
	Display    bool            `json:"display"`
	Details    json.RawMessage `json:"details,omitempty"`
	Timestamp  int64           `json:"timestamp"`
}

func (CustomMessage) messageRole() string { return "custom" }

type BranchSummaryMessage struct {
	Role      string `json:"role"`
	Summary   string `json:"summary"`
	FromID    string `json:"fromId"`
	Timestamp int64  `json:"timestamp"`
}

func (BranchSummaryMessage) messageRole() string { return "branchSummary" }

type CompactionSummaryMessage struct {
	Role         string `json:"role"`
	Summary      string `json:"summary"`
	TokensBefore int64  `json:"tokensBefore"`
	Timestamp    int64  `json:"timestamp"`
}

func (CompactionSummaryMessage) messageRole() string { return "compactionSummary" }

type UnknownMessage struct {
	Role string          `json:"role,omitempty"`
	Raw  json.RawMessage `json:"-"`
}

func (m UnknownMessage) messageRole() string { return m.Role }

func DecodeAgentMessage(data json.RawMessage) (AgentMessage, error) {
	var envelope struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	switch envelope.Role {
	case "user":
		var v UserMessage
		return &v, json.Unmarshal(data, &v)
	case "assistant":
		var v AssistantMessage
		return &v, json.Unmarshal(data, &v)
	case "toolResult":
		var v ToolResultMessage
		return &v, json.Unmarshal(data, &v)
	case "bashExecution":
		var v BashExecutionMessage
		return &v, json.Unmarshal(data, &v)
	case "custom":
		var v CustomMessage
		return &v, json.Unmarshal(data, &v)
	case "branchSummary":
		var v BranchSummaryMessage
		return &v, json.Unmarshal(data, &v)
	case "compactionSummary":
		var v CompactionSummaryMessage
		return &v, json.Unmarshal(data, &v)
	default:
		return &UnknownMessage{Role: envelope.Role, Raw: append(json.RawMessage(nil), data...)}, nil
	}
}

func decodeContentParts(data json.RawMessage) ([]ContentPart, error) {
	if len(bytes.TrimSpace(data)) == 0 || bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil, nil
	}
	var raws []json.RawMessage
	if err := json.Unmarshal(data, &raws); err != nil {
		return nil, err
	}
	parts := make([]ContentPart, 0, len(raws))
	for _, raw := range raws {
		part, err := decodeContentPart(raw)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return parts, nil
}

func decodeContentPart(data json.RawMessage) (ContentPart, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	switch envelope.Type {
	case "text":
		var v TextContent
		return v, json.Unmarshal(data, &v)
	case "image":
		var v ImageContent
		return v, json.Unmarshal(data, &v)
	case "thinking":
		var v ThinkingContent
		return v, json.Unmarshal(data, &v)
	case "toolCall":
		var v ToolCallContent
		return v, json.Unmarshal(data, &v)
	default:
		return UnknownContent{Type: envelope.Type, Raw: append(json.RawMessage(nil), data...)}, nil
	}
}

func marshalContentParts(parts []ContentPart) ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, part := range parts {
		if i > 0 {
			buf.WriteByte(',')
		}
		var raw []byte
		var err error
		switch v := part.(type) {
		case UnknownContent:
			raw = v.Raw
		case *UnknownContent:
			raw = v.Raw
		default:
			raw, err = json.Marshal(part)
		}
		if err != nil {
			return nil, err
		}
		buf.Write(raw)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func mustMarshalContentParts(parts []ContentPart) json.RawMessage {
	data, err := marshalContentParts(parts)
	if err != nil {
		panic(fmt.Sprintf("marshal content parts: %v", err))
	}
	return data
}
