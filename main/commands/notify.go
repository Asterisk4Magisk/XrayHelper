package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

type NotifyCommand struct{}

func (this *NotifyCommand) Execute(args []string) error {
	err := builds.LoadConfig()
	if err != nil {
		return err
	}
	if len(args) != 2 {
		return errors.New("usage: notify DIR FILE").WithPrefix("notify").WithPathObj(*this)
	}
	err = startNotify(args[0], args[1])
	if err != nil {
		return errors.New("start notify failed, ", startNotify(args[0], args[1])).WithPrefix("notify").WithPathObj(*this)
	}
	return nil
}

func startNotify(dir string, file string) error {
	// TODO
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			return
		}
	}(watcher)

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Print("%s %s\n", event.Name, event.Op)
			default:
				return
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}
	<-done
	return nil
}
