package fs

import (
	"github.com/fsnotify/fsnotify"
)

func NewObserver(handle Handler) *Observer {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				panic(err)
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				handle(Event(ev))
			}
		}
	}()

	return &Observer{w}
}

type Observer struct {
	w *fsnotify.Watcher
}

func (o *Observer) Observe(name string) error {
	return o.w.Add(name)
}

func (o *Observer) Unobserve(name string) error {
	return o.w.Remove(name)
}

func (o *Observer) Disconnect() error {
	return o.w.Close()
}
