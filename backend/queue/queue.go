package dispatch

type Dispatcher interface {
	Push(Event) error
	Run()
}

type EventDispatcher struct {
	Opts     Options
	Queue    chan models.Notification
	Finished bool
}

type Options struct {
	MaxWorkers   int
	MaxQueueSize int
}

func NewEventDispatcher(opts Options) Dispatcher {
	return EventDispatcher{
		Opts:     opts,
		Queue:    make(chan Event, opts.MaxQueueSize),
		Finished: false,
	}
}
