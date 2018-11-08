package main

import (
	"fmt"
	"math/rand"
	"meta2/core"
	"meta2/dashboard"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	waitGroup := &sync.WaitGroup{}

	d := dashboard.CLIDashboard{}
	d.Init(nil)

	interuptChannel := make(chan os.Signal, 1)
	signal.Notify(interuptChannel, os.Interrupt, syscall.SIGTERM)

	globalEventChannel := make(chan core.GlobalEvent)

	dashboardEventChannel, _ := d.Launch(waitGroup, globalEventChannel)

	cleanupChannel := make(chan int)

	cleanupFunction := func() {
		// Terminate all dashboards
		d.Terminate()
		// Close the dashboardEventChannel
		close(cleanupChannel)
	}

	go func() {
		select {
		case e := <-globalEventChannel:
			if e.EventType == core.Shutdown {
				cleanupFunction()
			}
		case <-interuptChannel:
			cleanupFunction()
		}
	}()

	// Demo dashboard
	for i := 0; i < 3; i++ {
		go func(deviceIndex int) {
			id := fmt.Sprintf("dev/d%d", deviceIndex)
			status := dashboard.DeviceWorkload{
				"docker/x",
				"ok",
				"OK",
			}

			for i := 0; i < 100; i += rand.Intn(10) + 20 {
				status.Performance = fmt.Sprintf("%d Units", i)
				dashboardEventChannel <- dashboard.DashboardEvent{
					dashboard.UpdateDevice, id, status,
				}

				time.Sleep(time.Duration(500) * time.Millisecond)
			}

		}(i)
	}

	<-cleanupChannel

	waitGroup.Wait()

	fmt.Println("Good bye")
}
