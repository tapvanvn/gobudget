package gobudget

import "github.com/tapvanvn/godbengine/engine"

var __eng *engine.Engine = nil
var __events map[EventName]*Event = map[EventName]*Event{}

func InitGoBudget(eng *engine.Engine) error {
	memdb := eng.GetMemPool()
	if memdb == nil {

		return ErrInvalidDBEngine
	}
	__eng = eng
	return nil
}

func AddEvent(event *Event) {

	__events[event.Name] = event
}
