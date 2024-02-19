package main

import (
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dundunlabs/gomon/app"
	"github.com/dundunlabs/gomon/fs"
	"github.com/fsnotify/fsnotify"
)

const (
	delay = 100 * time.Millisecond
)

func waitForExit() os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	return <-ch
}

func main() {
	app := app.New(os.Args[1:]...)
	app.Start()
	defer app.Stop()

	t := time.AfterFunc(math.MaxInt64, func() { app.Restart() })
	o := fs.NewObserver(func(e fs.Event) {
		switch e.Op {
		case fsnotify.Write:
			t.Stop()
			t.Reset(delay)
		}
	})
	defer o.Disconnect()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	o.Observe(wd)

	waitForExit()
}
