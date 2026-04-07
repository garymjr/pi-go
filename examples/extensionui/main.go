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
		switch v := n.(type) {
		case rpc.UIRequestNotification:
			switch req := v.Request.(type) {
			case *wire.SelectUIRequest:
				_ = client.SendUIResponse(ctx, wire.ExtensionUIResponseValue{Type: "extension_ui_response", ID: req.ID, Value: req.Options[0]})
			case *wire.ConfirmUIRequest:
				_ = client.SendUIResponse(ctx, wire.ExtensionUIResponseConfirm{Type: "extension_ui_response", ID: req.ID, Confirmed: true})
			case *wire.InputUIRequest:
				_ = client.SendUIResponse(ctx, wire.ExtensionUIResponseValue{Type: "extension_ui_response", ID: req.ID, Value: "value"})
			case *wire.EditorUIRequest:
				_ = client.SendUIResponse(ctx, wire.ExtensionUIResponseValue{Type: "extension_ui_response", ID: req.ID, Value: req.Prefill})
			case *wire.NotifyUIRequest, *wire.SetStatusUIRequest, *wire.SetWidgetUIRequest, *wire.SetTitleUIRequest, *wire.SetEditorTextUIRequest:
				fmt.Printf("ui notification: %T\n", req)
			}
		}
	})

	if err := client.Prompt(ctx, wire.PromptCommand{Message: "/my-extension-command"}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := client.WaitForIdle(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
