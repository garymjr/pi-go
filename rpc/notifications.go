package rpc

import "pi-go/wire"

type Notification interface{ notification() }

type EventNotification struct{ Event wire.Event }

func (EventNotification) notification() {}

type UIRequestNotification struct{ Request wire.UIRequest }

func (UIRequestNotification) notification() {}

type NotificationHandler func(Notification)
