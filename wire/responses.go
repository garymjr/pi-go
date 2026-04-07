package wire

import "encoding/json"

type Frame interface{ frameType() string }

type Response interface {
	Frame
	responseCommand() string
}

type responseBase struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Command string `json:"command"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (r responseBase) frameType() string       { return r.Type }
func (r responseBase) responseCommand() string { return r.Command }

type PromptResponse struct{ responseBase }
type SteerResponse struct{ responseBase }
type FollowUpResponse struct{ responseBase }
type AbortResponse struct{ responseBase }
type SetThinkingLevelResponse struct{ responseBase }
type SetSteeringModeResponse struct{ responseBase }
type SetFollowUpModeResponse struct{ responseBase }
type SetAutoCompactionResponse struct{ responseBase }
type SetAutoRetryResponse struct{ responseBase }
type AbortRetryResponse struct{ responseBase }
type AbortBashResponse struct{ responseBase }
type SetSessionNameResponse struct{ responseBase }

type NewSessionResponse struct {
	responseBase
	Data struct {
		Cancelled bool `json:"cancelled"`
	} `json:"data"`
}

type GetStateResponse struct {
	responseBase
	Data RpcSessionState `json:"data"`
}

type SetModelResponse struct {
	responseBase
	Data Model `json:"data"`
}

type CycleModelData struct {
	Model         Model  `json:"model"`
	ThinkingLevel string `json:"thinkingLevel"`
	IsScoped      bool   `json:"isScoped"`
}

type CycleModelResponse struct {
	responseBase
	Data *CycleModelData `json:"data"`
}

type GetAvailableModelsResponse struct {
	responseBase
	Data struct {
		Models []Model `json:"models"`
	} `json:"data"`
}

type CycleThinkingLevelResponse struct {
	responseBase
	Data *struct {
		Level string `json:"level"`
	} `json:"data"`
}

type CompactResponse struct {
	responseBase
	Data CompactionResult `json:"data"`
}

type BashResponse struct {
	responseBase
	Data BashResult `json:"data"`
}

type GetSessionStatsResponse struct {
	responseBase
	Data SessionStats `json:"data"`
}

type ExportHTMLResponse struct {
	responseBase
	Data struct {
		Path string `json:"path"`
	} `json:"data"`
}

type SwitchSessionResponse struct {
	responseBase
	Data struct {
		Cancelled bool `json:"cancelled"`
	} `json:"data"`
}

type ForkResponse struct {
	responseBase
	Data struct {
		Text      string `json:"text"`
		Cancelled bool   `json:"cancelled"`
	} `json:"data"`
}

type GetForkMessagesResponse struct {
	responseBase
	Data struct {
		Messages []struct {
			EntryID string `json:"entryId"`
			Text    string `json:"text"`
		} `json:"messages"`
	} `json:"data"`
}

type GetLastAssistantTextResponse struct {
	responseBase
	Data struct {
		Text *string `json:"text"`
	} `json:"data"`
}

type GetMessagesResponse struct {
	responseBase
	Data struct {
		Messages []AgentMessage `json:"messages"`
	} `json:"data"`
}

func (r *GetMessagesResponse) UnmarshalJSON(data []byte) error {
	type alias struct {
		responseBase
		Data struct {
			Messages []json.RawMessage `json:"messages"`
		} `json:"data"`
	}
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.responseBase = aux.responseBase
	r.Data.Messages = make([]AgentMessage, 0, len(aux.Data.Messages))
	for _, raw := range aux.Data.Messages {
		msg, err := DecodeAgentMessage(raw)
		if err != nil {
			return err
		}
		r.Data.Messages = append(r.Data.Messages, msg)
	}
	return nil
}

type GetCommandsResponse struct {
	responseBase
	Data struct {
		Commands []SlashCommand `json:"commands"`
	} `json:"data"`
}

type ErrorResponse struct{ responseBase }
