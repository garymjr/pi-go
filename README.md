# pi-go

Go RPC client for `pi --mode rpc` over stdin/stdout.

## Features

- Full v1 command surface in `rpc`
- Raw wire protocol access in `wire`
- Strict LF-only JSONL framing
- Request/response correlation with auto-generated request IDs
- Synchronous notification fanout for events and extension UI requests
- Forward-compatible unknown frame preservation

## Install

```bash
go get github.com/garymjr/pi-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/garymjr/pi-go/rpc"
    "github.com/garymjr/pi-go/wire"
)

func main() {
    ctx := context.Background()
    client := rpc.NewClient(rpc.Options{NoSession: true})
    if err := client.Start(ctx); err != nil {
        panic(err)
    }
    defer client.Close()

    client.OnNotification(func(n rpc.Notification) {
        evt, ok := n.(rpc.EventNotification)
        if !ok {
            return
        }
        if update, ok := evt.Event.(*wire.MessageUpdateEvent); ok {
            if delta, ok := update.AssistantMessageEvent.(*wire.AssistantMessageEventTextDelta); ok {
                fmt.Print(delta.Delta)
            }
        }
    })

    if _, err := client.PromptAndWait(ctx, wire.PromptCommand{Message: "Hello"}); err != nil {
        panic(err)
    }
}
```

## Packages

- `wire`: protocol structs, JSONL encoder, strict decoder, unknown variant preservation
- `rpc`: subprocess lifecycle, typed command methods, notifications, idle helpers

## Support Matrix

- Commands: all documented v1 commands implemented
- Stdout frames: responses, events, extension UI requests
- Extension UI responses: value, confirm, cancel
- Forward compatibility: unknown top-level frames, unknown event/UI/message/content variants preserved as raw JSON
