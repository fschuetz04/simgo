package simgo

// Awaitable represents any awaitable thing like an event.
type Awaitable interface {
	// Processed must return true when the event is processed.
	Processed() bool

	// Aborted must return true when the event is aborted.
	Aborted() bool

	// AddHandler must add the given normal handler. When the event is
	// processed, the normal handler must be called.
	AddHandler(handler Handler)

	// AddAbortHandler must add the given abort handler. When the event is
	// aborted, the abort handler must be called.
	AddAbortHandler(handler Handler)
}
