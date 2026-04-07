package wire

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestDecoderFraming(t *testing.T) {
	input := strings.NewReader("{\"type\":\"agent_start\"}\r\n{\"type\":\"agent_start\",\"note\":\"x\u2028y\u2029z\"}\n{\"type\":\"agent_start\"}")
	d := NewDecoder(input)
	for i := 0; i < 3; i++ {
		frame, err := d.Decode()
		if err != nil {
			t.Fatalf("decode %d: %v", i, err)
		}
		if _, ok := frame.(*AgentStartEvent); !ok {
			t.Fatalf("frame %d type = %T", i, frame)
		}
	}
}

func TestEncodeLine(t *testing.T) {
	line, err := EncodeLine(map[string]any{"type": "abort"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasSuffix(line, []byte("\n")) {
		t.Fatal("missing newline")
	}
}

func TestResponseDispatch(t *testing.T) {
	tests := []struct {
		line string
		want any
	}{
		{"{\"type\":\"response\",\"command\":\"prompt\",\"success\":true}", &PromptResponse{}},
		{"{\"type\":\"response\",\"command\":\"new_session\",\"success\":true,\"data\":{\"cancelled\":false}}", &NewSessionResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_state\",\"success\":true,\"data\":{\"thinkingLevel\":\"medium\",\"isStreaming\":false,\"isCompacting\":false,\"steeringMode\":\"all\",\"followUpMode\":\"one-at-a-time\",\"sessionId\":\"s\",\"autoCompactionEnabled\":true,\"messageCount\":0,\"pendingMessageCount\":0}}", &GetStateResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_model\",\"success\":true,\"data\":{\"id\":\"m\"}}", &SetModelResponse{}},
		{"{\"type\":\"response\",\"command\":\"cycle_model\",\"success\":true,\"data\":null}", &CycleModelResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_available_models\",\"success\":true,\"data\":{\"models\":[]}}", &GetAvailableModelsResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_thinking_level\",\"success\":true}", &SetThinkingLevelResponse{}},
		{"{\"type\":\"response\",\"command\":\"cycle_thinking_level\",\"success\":true,\"data\":{\"level\":\"high\"}}", &CycleThinkingLevelResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_steering_mode\",\"success\":true}", &SetSteeringModeResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_follow_up_mode\",\"success\":true}", &SetFollowUpModeResponse{}},
		{"{\"type\":\"response\",\"command\":\"compact\",\"success\":true,\"data\":{\"summary\":\"s\",\"firstKeptEntryId\":\"1\",\"tokensBefore\":1,\"details\":{}}}", &CompactResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_auto_compaction\",\"success\":true}", &SetAutoCompactionResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_auto_retry\",\"success\":true}", &SetAutoRetryResponse{}},
		{"{\"type\":\"response\",\"command\":\"abort_retry\",\"success\":true}", &AbortRetryResponse{}},
		{"{\"type\":\"response\",\"command\":\"bash\",\"success\":true,\"data\":{\"output\":\"x\",\"exitCode\":0,\"cancelled\":false,\"truncated\":false}}", &BashResponse{}},
		{"{\"type\":\"response\",\"command\":\"abort_bash\",\"success\":true}", &AbortBashResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_session_stats\",\"success\":true,\"data\":{\"sessionId\":\"s\",\"userMessages\":0,\"assistantMessages\":0,\"toolCalls\":0,\"toolResults\":0,\"totalMessages\":0,\"tokens\":{\"input\":0,\"output\":0,\"cacheRead\":0,\"cacheWrite\":0,\"total\":0},\"cost\":0}}", &GetSessionStatsResponse{}},
		{"{\"type\":\"response\",\"command\":\"export_html\",\"success\":true,\"data\":{\"path\":\"/tmp/x\"}}", &ExportHTMLResponse{}},
		{"{\"type\":\"response\",\"command\":\"switch_session\",\"success\":true,\"data\":{\"cancelled\":false}}", &SwitchSessionResponse{}},
		{"{\"type\":\"response\",\"command\":\"fork\",\"success\":true,\"data\":{\"text\":\"x\",\"cancelled\":false}}", &ForkResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_fork_messages\",\"success\":true,\"data\":{\"messages\":[{\"entryId\":\"1\",\"text\":\"x\"}]}}", &GetForkMessagesResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_last_assistant_text\",\"success\":true,\"data\":{\"text\":null}}", &GetLastAssistantTextResponse{}},
		{"{\"type\":\"response\",\"command\":\"set_session_name\",\"success\":true}", &SetSessionNameResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_messages\",\"success\":true,\"data\":{\"messages\":[{\"role\":\"user\",\"content\":\"x\",\"timestamp\":1}]}}", &GetMessagesResponse{}},
		{"{\"type\":\"response\",\"command\":\"get_commands\",\"success\":true,\"data\":{\"commands\":[]}}", &GetCommandsResponse{}},
	}
	for _, tt := range tests {
		frame, err := decodeFrame([]byte(tt.line))
		if err != nil {
			t.Fatalf("decode %T: %v", tt.want, err)
		}
		if reflect.TypeOf(frame) != reflect.TypeOf(tt.want) {
			t.Fatalf("got %T want %T", frame, tt.want)
		}
	}
}

func TestEventAndUIAndUnknownDispatch(t *testing.T) {
	eventLines := []struct {
		line string
		want any
	}{
		{"{\"type\":\"agent_start\"}", &AgentStartEvent{}},
		{"{\"type\":\"agent_end\",\"messages\":[]}", &AgentEndEvent{}},
		{"{\"type\":\"turn_start\"}", &TurnStartEvent{}},
		{"{\"type\":\"turn_end\",\"message\":{\"role\":\"user\",\"content\":\"x\",\"timestamp\":1},\"toolResults\":[]}", &TurnEndEvent{}},
		{"{\"type\":\"message_start\",\"message\":{\"role\":\"user\",\"content\":\"x\",\"timestamp\":1}}", &MessageStartEvent{}},
		{"{\"type\":\"message_update\",\"message\":{\"role\":\"assistant\",\"content\":[],\"api\":\"a\",\"provider\":\"p\",\"model\":\"m\",\"usage\":{\"input\":0,\"output\":0,\"cacheRead\":0,\"cacheWrite\":0,\"totalTokens\":0,\"cost\":{\"input\":0,\"output\":0,\"cacheRead\":0,\"cacheWrite\":0,\"total\":0}},\"stopReason\":\"stop\",\"timestamp\":1},\"assistantMessageEvent\":{\"type\":\"start\",\"partial\":{\"role\":\"assistant\",\"content\":[],\"api\":\"a\",\"provider\":\"p\",\"model\":\"m\",\"usage\":{\"input\":0,\"output\":0,\"cacheRead\":0,\"cacheWrite\":0,\"totalTokens\":0,\"cost\":{\"input\":0,\"output\":0,\"cacheRead\":0,\"cacheWrite\":0,\"total\":0}},\"stopReason\":\"stop\",\"timestamp\":1}}}", &MessageUpdateEvent{}},
		{"{\"type\":\"message_end\",\"message\":{\"role\":\"user\",\"content\":\"x\",\"timestamp\":1}}", &MessageEndEvent{}},
		{"{\"type\":\"tool_execution_start\",\"toolCallId\":\"1\",\"toolName\":\"bash\",\"args\":{}}", &ToolExecutionStartEvent{}},
		{"{\"type\":\"tool_execution_update\",\"toolCallId\":\"1\",\"toolName\":\"bash\",\"args\":{},\"partialResult\":{}}", &ToolExecutionUpdateEvent{}},
		{"{\"type\":\"tool_execution_end\",\"toolCallId\":\"1\",\"toolName\":\"bash\",\"result\":{},\"isError\":false}", &ToolExecutionEndEvent{}},
		{"{\"type\":\"queue_update\",\"steering\":[],\"followUp\":[]}", &QueueUpdateEvent{}},
		{"{\"type\":\"compaction_start\",\"reason\":\"manual\"}", &CompactionStartEvent{}},
		{"{\"type\":\"compaction_end\",\"reason\":\"manual\",\"result\":null,\"aborted\":false,\"willRetry\":false}", &CompactionEndEvent{}},
		{"{\"type\":\"auto_retry_start\",\"attempt\":1,\"maxAttempts\":3,\"delayMs\":1,\"errorMessage\":\"x\"}", &AutoRetryStartEvent{}},
		{"{\"type\":\"auto_retry_end\",\"success\":true,\"attempt\":1}", &AutoRetryEndEvent{}},
		{"{\"type\":\"extension_error\",\"extensionPath\":\"x\",\"event\":\"tool_call\",\"error\":\"boom\"}", &ExtensionErrorEvent{}},
	}
	for _, tt := range eventLines {
		frame, err := decodeFrame([]byte(tt.line))
		if err != nil {
			t.Fatal(err)
		}
		if reflect.TypeOf(frame) != reflect.TypeOf(tt.want) {
			t.Fatalf("got %T want %T", frame, tt.want)
		}
	}
	uiLines := []struct {
		line string
		want any
	}{
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"select\",\"title\":\"x\",\"options\":[\"a\"]}", &SelectUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"confirm\",\"title\":\"x\",\"message\":\"y\"}", &ConfirmUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"input\",\"title\":\"x\"}", &InputUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"editor\",\"title\":\"x\"}", &EditorUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"notify\",\"message\":\"x\"}", &NotifyUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"setStatus\",\"statusKey\":\"x\"}", &SetStatusUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"setWidget\",\"widgetKey\":\"x\"}", &SetWidgetUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"setTitle\",\"title\":\"x\"}", &SetTitleUIRequest{}},
		{"{\"type\":\"extension_ui_request\",\"id\":\"1\",\"method\":\"set_editor_text\",\"text\":\"x\"}", &SetEditorTextUIRequest{}},
	}
	for _, tt := range uiLines {
		frame, err := decodeFrame([]byte(tt.line))
		if err != nil {
			t.Fatal(err)
		}
		if reflect.TypeOf(frame) != reflect.TypeOf(tt.want) {
			t.Fatalf("got %T want %T", frame, tt.want)
		}
	}
	frame, err := decodeFrame([]byte(`{"type":"future_frame","x":1}`))
	if err != nil || reflect.TypeOf(frame) != reflect.TypeOf(UnknownFrame{}) {
		t.Fatalf("unknown frame: %T %v", frame, err)
	}
	frame, err = decodeFrame([]byte(`{"type":"extension_ui_request","id":"1","method":"future"}`))
	if err != nil || reflect.TypeOf(frame) != reflect.TypeOf(UnknownUIRequest{}) {
		t.Fatalf("unknown ui: %T %v", frame, err)
	}
}

func TestNestedUnionDispatch(t *testing.T) {
	resp, err := decodeFrame([]byte(`{"type":"response","command":"get_messages","success":true,"data":{"messages":[{"role":"user","content":"hello","timestamp":1},{"role":"assistant","content":[{"type":"text","text":"hi"}],"api":"a","provider":"p","model":"m","usage":{"input":0,"output":0,"cacheRead":0,"cacheWrite":0,"totalTokens":0,"cost":{"input":0,"output":0,"cacheRead":0,"cacheWrite":0,"total":0}},"stopReason":"stop","timestamp":1},{"role":"toolResult","toolCallId":"1","toolName":"bash","content":[{"type":"text","text":"ok"}],"isError":false,"timestamp":1},{"role":"bashExecution","command":"ls","output":"x","cancelled":false,"truncated":false,"timestamp":1}]}}`))
	if err != nil {
		t.Fatal(err)
	}
	messages := resp.(*GetMessagesResponse).Data.Messages
	if reflect.TypeOf(messages[0]) != reflect.TypeOf(&UserMessage{}) || reflect.TypeOf(messages[1]) != reflect.TypeOf(&AssistantMessage{}) || reflect.TypeOf(messages[2]) != reflect.TypeOf(&ToolResultMessage{}) || reflect.TypeOf(messages[3]) != reflect.TypeOf(&BashExecutionMessage{}) {
		t.Fatalf("unexpected message union types")
	}
	frame, err := decodeFrame([]byte(`{"type":"message_update","message":{"role":"assistant","content":[],"api":"a","provider":"p","model":"m","usage":{"input":0,"output":0,"cacheRead":0,"cacheWrite":0,"totalTokens":0,"cost":{"input":0,"output":0,"cacheRead":0,"cacheWrite":0,"total":0}},"stopReason":"stop","timestamp":1},"assistantMessageEvent":{"type":"toolcall_end","contentIndex":0,"toolCall":{"type":"toolCall","id":"1","name":"bash","arguments":{}},"partial":{"role":"assistant","content":[],"api":"a","provider":"p","model":"m","usage":{"input":0,"output":0,"cacheRead":0,"cacheWrite":0,"totalTokens":0,"cost":{"input":0,"output":0,"cacheRead":0,"cacheWrite":0,"total":0}},"stopReason":"stop","timestamp":1}}}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := frame.(*MessageUpdateEvent).AssistantMessageEvent.(*AssistantMessageEventToolCallEnd); !ok {
		t.Fatalf("assistant event type = %T", frame.(*MessageUpdateEvent).AssistantMessageEvent)
	}
}
