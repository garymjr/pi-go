package rpc

import (
	"context"
	"reflect"

	"github.com/garymjr/pi-go/wire"
)

func (c *Client) Prompt(ctx context.Context, cmd wire.PromptCommand) error {
	_, err := c.send(ctx, "prompt", ensureCommandType(cmd, "prompt"))
	return err
}

func (c *Client) Steer(ctx context.Context, cmd wire.SteerCommand) error {
	_, err := c.send(ctx, "steer", ensureCommandType(cmd, "steer"))
	return err
}

func (c *Client) FollowUp(ctx context.Context, cmd wire.FollowUpCommand) error {
	_, err := c.send(ctx, "follow_up", ensureCommandType(cmd, "follow_up"))
	return err
}

func (c *Client) Abort(ctx context.Context) error {
	_, err := c.send(ctx, "abort", wire.AbortCommand{Type: "abort"})
	return err
}

func (c *Client) NewSession(ctx context.Context, cmd wire.NewSessionCommand) (*wire.NewSessionResponse, error) {
	resp, err := c.send(ctx, "new_session", ensureCommandType(cmd, "new_session"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.NewSessionResponse), nil
}

func (c *Client) GetState(ctx context.Context) (*wire.GetStateResponse, error) {
	resp, err := c.send(ctx, "get_state", wire.GetStateCommand{Type: "get_state"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetStateResponse), nil
}

func (c *Client) SetModel(ctx context.Context, cmd wire.SetModelCommand) (*wire.SetModelResponse, error) {
	resp, err := c.send(ctx, "set_model", ensureCommandType(cmd, "set_model"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.SetModelResponse), nil
}

func (c *Client) CycleModel(ctx context.Context) (*wire.CycleModelResponse, error) {
	resp, err := c.send(ctx, "cycle_model", wire.CycleModelCommand{Type: "cycle_model"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.CycleModelResponse), nil
}

func (c *Client) GetAvailableModels(ctx context.Context) (*wire.GetAvailableModelsResponse, error) {
	resp, err := c.send(ctx, "get_available_models", wire.GetAvailableModelsCommand{Type: "get_available_models"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetAvailableModelsResponse), nil
}

func (c *Client) SetThinkingLevel(ctx context.Context, cmd wire.SetThinkingLevelCommand) error {
	_, err := c.send(ctx, "set_thinking_level", ensureCommandType(cmd, "set_thinking_level"))
	return err
}

func (c *Client) CycleThinkingLevel(ctx context.Context) (*wire.CycleThinkingLevelResponse, error) {
	resp, err := c.send(ctx, "cycle_thinking_level", wire.CycleThinkingLevelCommand{Type: "cycle_thinking_level"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.CycleThinkingLevelResponse), nil
}

func (c *Client) SetSteeringMode(ctx context.Context, cmd wire.SetSteeringModeCommand) error {
	_, err := c.send(ctx, "set_steering_mode", ensureCommandType(cmd, "set_steering_mode"))
	return err
}

func (c *Client) SetFollowUpMode(ctx context.Context, cmd wire.SetFollowUpModeCommand) error {
	_, err := c.send(ctx, "set_follow_up_mode", ensureCommandType(cmd, "set_follow_up_mode"))
	return err
}

func (c *Client) Compact(ctx context.Context, cmd wire.CompactCommand) (*wire.CompactResponse, error) {
	resp, err := c.send(ctx, "compact", ensureCommandType(cmd, "compact"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.CompactResponse), nil
}

func (c *Client) SetAutoCompaction(ctx context.Context, cmd wire.SetAutoCompactionCommand) error {
	_, err := c.send(ctx, "set_auto_compaction", ensureCommandType(cmd, "set_auto_compaction"))
	return err
}

func (c *Client) SetAutoRetry(ctx context.Context, cmd wire.SetAutoRetryCommand) error {
	_, err := c.send(ctx, "set_auto_retry", ensureCommandType(cmd, "set_auto_retry"))
	return err
}

func (c *Client) AbortRetry(ctx context.Context) error {
	_, err := c.send(ctx, "abort_retry", wire.AbortRetryCommand{Type: "abort_retry"})
	return err
}

func (c *Client) Bash(ctx context.Context, cmd wire.BashCommand) (*wire.BashResponse, error) {
	resp, err := c.send(ctx, "bash", ensureCommandType(cmd, "bash"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.BashResponse), nil
}

func (c *Client) AbortBash(ctx context.Context) error {
	_, err := c.send(ctx, "abort_bash", wire.AbortBashCommand{Type: "abort_bash"})
	return err
}

func (c *Client) GetSessionStats(ctx context.Context) (*wire.GetSessionStatsResponse, error) {
	resp, err := c.send(ctx, "get_session_stats", wire.GetSessionStatsCommand{Type: "get_session_stats"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetSessionStatsResponse), nil
}

func (c *Client) ExportHTML(ctx context.Context, cmd wire.ExportHTMLCommand) (*wire.ExportHTMLResponse, error) {
	resp, err := c.send(ctx, "export_html", ensureCommandType(cmd, "export_html"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.ExportHTMLResponse), nil
}

func (c *Client) SwitchSession(ctx context.Context, cmd wire.SwitchSessionCommand) (*wire.SwitchSessionResponse, error) {
	resp, err := c.send(ctx, "switch_session", ensureCommandType(cmd, "switch_session"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.SwitchSessionResponse), nil
}

func (c *Client) Fork(ctx context.Context, cmd wire.ForkCommand) (*wire.ForkResponse, error) {
	resp, err := c.send(ctx, "fork", ensureCommandType(cmd, "fork"))
	if err != nil {
		return nil, err
	}
	return resp.(*wire.ForkResponse), nil
}

func (c *Client) GetForkMessages(ctx context.Context) (*wire.GetForkMessagesResponse, error) {
	resp, err := c.send(ctx, "get_fork_messages", wire.GetForkMessagesCommand{Type: "get_fork_messages"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetForkMessagesResponse), nil
}

func (c *Client) GetLastAssistantText(ctx context.Context) (*wire.GetLastAssistantTextResponse, error) {
	resp, err := c.send(ctx, "get_last_assistant_text", wire.GetLastAssistantTextCommand{Type: "get_last_assistant_text"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetLastAssistantTextResponse), nil
}

func (c *Client) SetSessionName(ctx context.Context, cmd wire.SetSessionNameCommand) error {
	_, err := c.send(ctx, "set_session_name", ensureCommandType(cmd, "set_session_name"))
	return err
}

func (c *Client) GetMessages(ctx context.Context) (*wire.GetMessagesResponse, error) {
	resp, err := c.send(ctx, "get_messages", wire.GetMessagesCommand{Type: "get_messages"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetMessagesResponse), nil
}

func (c *Client) GetCommands(ctx context.Context) (*wire.GetCommandsResponse, error) {
	resp, err := c.send(ctx, "get_commands", wire.GetCommandsCommand{Type: "get_commands"})
	if err != nil {
		return nil, err
	}
	return resp.(*wire.GetCommandsResponse), nil
}

func ensureCommandType[T any](cmd T, typ string) T {
	v := reflect.ValueOf(&cmd).Elem()
	if v.Kind() == reflect.Struct {
		field := v.FieldByName("Type")
		if field.IsValid() && field.CanSet() && field.Kind() == reflect.String && field.String() == "" {
			field.SetString(typ)
		}
	}
	return cmd
}
