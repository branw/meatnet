package logger

import (
	"fmt"
	"log/slog"
	"meatnet/pkg/meatnet"
	"os"
	"os/signal"
	"reflect"
	"tinygo.org/x/bluetooth"
)

func runMeatNet() error {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return err
	}

	service := meatnet.NewService(bluetooth.DefaultAdapter)

	events := make(chan meatnet.Event, 1)
	go func() {
		for {
			event := service.ReadEvent()
			events <- event
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

loop:
	for {
		select {
		case <-interrupt:
			break loop

		case event := <-events:
			switch event := event.(type) {
			case meatnet.ProbeAddedEvent:
				fmt.Printf("probe added %+v\n", *event.Probe)
			case meatnet.ProbeUpdatedEvent:
				fmt.Printf("probe updated %+v\n", *event.Probe)
			case meatnet.DeviceAddedEvent:
				fmt.Printf("device added %+v\n", *event.Device)
			case meatnet.DeviceUpdatedEvent:
				fmt.Printf("device updated %+v\n", *event.Device)

			default:
				fmt.Printf("skipping event %+v %+v\n", reflect.TypeOf(event), event)
			}
		}
	}

	err = service.Stop()
	if err != nil {
		fmt.Printf("service stop failed: %+v\n", err)
	}
	return nil
}

func main() {
	if err := runMeatNet(); err != nil {
		slog.Error("fatal error",
			slog.Any("err", err))
	}
}
