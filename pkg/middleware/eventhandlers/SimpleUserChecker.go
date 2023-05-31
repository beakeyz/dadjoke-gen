package eventhandlers

import (
	"fmt"

	"github.com/beakeyz/dadjoke-gen/pkg/bus/event"
)

func SimpleUserChecker(event event.UserCreateEvent) event.UserCreateEvent {
	// TODO
	event.SetCanceled(true)
	fmt.Println("funnie")

	return event
}
