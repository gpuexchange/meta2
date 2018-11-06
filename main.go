package main

import (
	"fmt"
	"meta2/dashboard"
	"time"
)

func main() {
	d := dashboard.CLIDashboard{}
	d.Init(nil)
	channel, _ := d.Launch()
	channel <- dashboard.DashboardMessage{
		dashboard.UpdateDevice, "d1", dashboard.DeviceWorkloadStatus{
			"workload1",
			"ok",
			"OK",
		},
	}
	channel <- dashboard.DashboardMessage{
		dashboard.UpdateDevice, "d2", dashboard.DeviceWorkloadStatus{
		},
	}

	time.Sleep(10000)
	fmt.Println("Good bye")
}
