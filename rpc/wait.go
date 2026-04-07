package rpc

import (
	"context"

	"github.com/garymjr/pi-go/wire"
)

func (c *Client) WaitForIdle(ctx context.Context) error {
	_, err := c.CollectEventsUntilIdle(ctx)
	return err
}

func (c *Client) CollectEventsUntilIdle(ctx context.Context) ([]wire.Event, error) {
	events := []wire.Event{}
	done := make(chan struct{})
	unsub := c.OnNotification(func(n Notification) {
		evt, ok := n.(EventNotification)
		if !ok {
			return
		}
		events = append(events, evt.Event)
		if _, ok := evt.Event.(*wire.AgentEndEvent); ok {
			select {
			case <-done:
			default:
				close(done)
			}
		}
	})
	defer unsub()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return events, nil
	}
}

func (c *Client) PromptAndWait(ctx context.Context, cmd wire.PromptCommand) ([]wire.Event, error) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan struct{})
	var events []wire.Event
	var collectErr error
	go func() {
		events, collectErr = c.CollectEventsUntilIdle(ctx2)
		close(done)
	}()
	if err := c.Prompt(ctx, cmd); err != nil {
		cancel()
		<-done
		return nil, err
	}
	<-done
	return events, collectErr
}
