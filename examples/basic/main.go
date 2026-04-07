package main

import (
	"context"
	"fmt"
	"os"

	"github.com/garymjr/pi-go/rpc"
	"github.com/garymjr/pi-go/wire"
)

func main() {
	ctx := context.Background()
	client := rpc.NewClient(rpc.Options{NoSession: true})
	if err := client.Start(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer client.Close()

	client.OnNotification(func(n rpc.Notification) {
		evt, ok := n.(rpc.EventNotification)
		if !ok {
			return
		}
		update, ok := evt.Event.(*wire.MessageUpdateEvent)
		if !ok {
			return
		}
		delta, ok := update.AssistantMessageEvent.(*wire.AssistantMessageEventTextDelta)
		if ok {
			fmt.Print(delta.Delta)
		}
	})

	if _, err := client.PromptAndWait(ctx, wire.PromptCommand{Message: "Hello"}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println()
}
