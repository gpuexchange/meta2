package dashboard

import (
	"fmt"
	"sync"
)

type CLIDashboard struct {
	BaseDashboard
}

func (d *CLIDashboard) render() {
	d.deviceLock.Lock()
	fmt.Println("Rendering status")
	for identifier, status := range d.devices {
		fmt.Printf("Device: %s Status: %s", identifier, status)
		fmt.Println()
	}

	fmt.Println()
	d.deviceLock.Unlock()
}

func (d *CLIDashboard) Launch(group *sync.WaitGroup) (chan DashboardMessage, error) {
	go func() {
		for ; !d.terminated; {
			d.render()
			// time.Sleep(time.Duration(1) * time.Second)
		}
	}()
	return d.BaseDashboard.Launch(group)
}
