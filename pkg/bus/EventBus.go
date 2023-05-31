package bus

import (
	"reflect"

	"github.com/beakeyz/dadjoke-gen/pkg/bus/event"
	"github.com/beakeyz/dadjoke-gen/pkg/logger"
)

var mainLog = *logger.New("Event Bus", 93, false)

type eventBus struct {
	Calls []eventCall
}

type eventCall struct {
	EventType reflect.Type
	Method    reflect.Value
}

func (self *eventBus) Register(method interface{}) {

	fun := reflect.TypeOf(method)
	event := fun.In(0)

	var call eventCall = eventCall{
		EventType: event,
		Method:    reflect.ValueOf(method),
	}
	self.Calls = append(self.Calls, call)

	mainLog.Out("Registered event")

}

func (self *eventBus) CallEvents(event event.Event) {
	for _, call := range self.Calls {
		if call.EventType == reflect.TypeOf(event) {
			args := make([]reflect.Value, 1)
			args[0] = reflect.ValueOf(event)
			out := call.Method.Call(args)
			event = out[0].Interface()
		}
	}
}

var EVENTBUS eventBus = eventBus{}
