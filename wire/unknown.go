package wire

import "encoding/json"

type UnknownFrame struct {
	Type string          `json:"type,omitempty"`
	Raw  json.RawMessage `json:"-"`
}

func (f UnknownFrame) frameType() string { return f.Type }

type UnknownEvent struct {
	Type string          `json:"type,omitempty"`
	Raw  json.RawMessage `json:"-"`
}

func (e UnknownEvent) frameType() string { return e.Type }
func (e UnknownEvent) eventType() string { return e.Type }

type UnknownUIRequest struct {
	Method string          `json:"method,omitempty"`
	Raw    json.RawMessage `json:"-"`
}

func (r UnknownUIRequest) frameType() string { return "extension_ui_request" }
func (r UnknownUIRequest) uiMethod() string  { return r.Method }
