package wire

type Command interface{ commandType() string }

type commandBase struct{}

func (commandBase) commandType() string { return "" }

type PromptCommand struct {
	ID                string         `json:"id,omitempty"`
	Type              string         `json:"type"`
	Message           string         `json:"message"`
	Images            []ImageContent `json:"images,omitempty"`
	StreamingBehavior string         `json:"streamingBehavior,omitempty"`
}

func (PromptCommand) commandType() string { return "prompt" }

type SteerCommand struct {
	ID      string         `json:"id,omitempty"`
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Images  []ImageContent `json:"images,omitempty"`
}

func (SteerCommand) commandType() string { return "steer" }

type FollowUpCommand struct {
	ID      string         `json:"id,omitempty"`
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Images  []ImageContent `json:"images,omitempty"`
}

func (FollowUpCommand) commandType() string { return "follow_up" }

type AbortCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (AbortCommand) commandType() string { return "abort" }

type NewSessionCommand struct {
	ID            string `json:"id,omitempty"`
	Type          string `json:"type"`
	ParentSession string `json:"parentSession,omitempty"`
}

func (NewSessionCommand) commandType() string { return "new_session" }

type GetStateCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetStateCommand) commandType() string { return "get_state" }

type SetModelCommand struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Provider string `json:"provider"`
	ModelID  string `json:"modelId"`
}

func (SetModelCommand) commandType() string { return "set_model" }

type CycleModelCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (CycleModelCommand) commandType() string { return "cycle_model" }

type GetAvailableModelsCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetAvailableModelsCommand) commandType() string { return "get_available_models" }

type SetThinkingLevelCommand struct {
	ID    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Level string `json:"level"`
}

func (SetThinkingLevelCommand) commandType() string { return "set_thinking_level" }

type CycleThinkingLevelCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (CycleThinkingLevelCommand) commandType() string { return "cycle_thinking_level" }

type SetSteeringModeCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
	Mode string `json:"mode"`
}

func (SetSteeringModeCommand) commandType() string { return "set_steering_mode" }

type SetFollowUpModeCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
	Mode string `json:"mode"`
}

func (SetFollowUpModeCommand) commandType() string { return "set_follow_up_mode" }

type CompactCommand struct {
	ID                 string `json:"id,omitempty"`
	Type               string `json:"type"`
	CustomInstructions string `json:"customInstructions,omitempty"`
}

func (CompactCommand) commandType() string { return "compact" }

type SetAutoCompactionCommand struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func (SetAutoCompactionCommand) commandType() string { return "set_auto_compaction" }

type SetAutoRetryCommand struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func (SetAutoRetryCommand) commandType() string { return "set_auto_retry" }

type AbortRetryCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (AbortRetryCommand) commandType() string { return "abort_retry" }

type BashCommand struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Command string `json:"command"`
}

func (BashCommand) commandType() string { return "bash" }

type AbortBashCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (AbortBashCommand) commandType() string { return "abort_bash" }

type GetSessionStatsCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetSessionStatsCommand) commandType() string { return "get_session_stats" }

type ExportHTMLCommand struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type"`
	OutputPath string `json:"outputPath,omitempty"`
}

func (ExportHTMLCommand) commandType() string { return "export_html" }

type SwitchSessionCommand struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	SessionPath string `json:"sessionPath"`
}

func (SwitchSessionCommand) commandType() string { return "switch_session" }

type ForkCommand struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	EntryID string `json:"entryId"`
}

func (ForkCommand) commandType() string { return "fork" }

type GetForkMessagesCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetForkMessagesCommand) commandType() string { return "get_fork_messages" }

type GetLastAssistantTextCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetLastAssistantTextCommand) commandType() string { return "get_last_assistant_text" }

type SetSessionNameCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (SetSessionNameCommand) commandType() string { return "set_session_name" }

type GetMessagesCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetMessagesCommand) commandType() string { return "get_messages" }

type GetCommandsCommand struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type"`
}

func (GetCommandsCommand) commandType() string { return "get_commands" }

type ExtensionUIResponseValue struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (ExtensionUIResponseValue) commandType() string { return "extension_ui_response" }

type ExtensionUIResponseConfirm struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Confirmed bool   `json:"confirmed"`
}

func (ExtensionUIResponseConfirm) commandType() string { return "extension_ui_response" }

type ExtensionUIResponseCancel struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Cancelled bool   `json:"cancelled"`
}

func (ExtensionUIResponseCancel) commandType() string { return "extension_ui_response" }
