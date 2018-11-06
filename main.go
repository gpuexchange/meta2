package main

import (
	"fmt"
	"math/rand"
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

	cleanupChannel := make(chan int)
	dashboardEventChannel, _ := d.Launch(waitGroup)

	go func() {
		<-interuptChannel
		fmt.Println("Received an interupt.")
		// Cleaning up
		d.Terminate()
		// Close the dashboardEventChannel
		close(cleanupChannel)
	}()

	// Demo dashboard
	for i := 0; i < 3; i++ {
		go func(deviceIndex int) {
			identifier := fmt.Sprintf("dev/d%d", deviceIndex)
			status := dashboard.DeviceWorkloadStatus{
				"docker/x",
				"ok",
				"OK",
			}

			for i := 0; i < 100; i += rand.Intn(10) + 20 {
				status.WorkloadPerformance = fmt.Sprintf("%d Units", i)
				dashboardEventChannel <- dashboard.DashboardMessage{
					dashboard.UpdateDevice, identifier, status,
				}

				time.Sleep(time.Duration(500) * time.Millisecond)
			}

		}(i)
	}

	<-cleanupChannel

	waitGroup.Wait()
}
