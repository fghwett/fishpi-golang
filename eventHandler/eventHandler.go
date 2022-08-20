package eventHandler

import "fishpi/logger"

type eventHandler struct {
	name    string
	methods map[EventType][]func(interface{})

	logger logger.Logger
}

func NewEventHandler(name string, logger logger.Logger) *eventHandler {
	return &eventHandler{
		name:    name,
		methods: make(map[EventType][]func(interface{})),

		logger: logger,
	}
}

func (eh *eventHandler) Pub(event EventType, data interface{}) {
	methods, ok := eh.methods[event]
	if !ok {
		eh.logger.Logf("EventHandler %s: no methods for event %s\n", eh.name, event)
		return
	}
	for _, method := range methods {
		go method(data)
	}
}

func (eh *eventHandler) Sub(event EventType, method func(interface{})) {
	eh.methods[event] = append(eh.methods[event], method)
}
