package dashboard

import (
	ui "github.com/gizak/termui"
	"meta2/core"
	"sort"
	"sync"
)

type CLIDashboard struct {
	BaseDashboard
}

func (d *CLIDashboard) render() {
	ui.Init()
	defer ui.Close()
	table := ui.NewTable()
	table.FgColor = ui.ColorWhite
	table.BgColor = ui.ColorDefault
	table.Height = 25

	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(6, 0, table)),
	)

	ui.Handle("/timer/1s", func(e ui.Event) {
		if (d.terminated) {
			ui.StopLoop()
		} else {
			d.deviceLock.Lock()
			rows := make([][]string, len(d.devices)+1)

			// Header row
			rows[0] = []string{"Device", "Workload", "Performance", "Status"}
			//table.BgColors[0] = ui.ColorGreen

			deviceIds := make([]string, 0)
			for deviceId, _ := range d.devices {
				deviceIds = append(deviceIds, deviceId)
			}

			sort.Strings(deviceIds)

			for i, identifier := range deviceIds {
				status := d.devices[identifier]
				rows[i+1] = []string{
					identifier,
					status.WorkloadName,
					status.WorkloadPerformance,
					status.WorkloadStatus,
				}
			}

			d.deviceLock.Unlock()

			table.Rows = rows
			table.Analysis()
			table.SetSize()
			table.BgColors[0] = ui.ColorGreen

			ui.Body.Align()
			ui.Render(ui.Body)
		}
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		d.responseChannel <- core.GlobalEvent{EventType: core.Shutdown}
	})

	ui.Loop()
}

func (d *CLIDashboard) Launch(group *sync.WaitGroup, globalEvents chan core.GlobalEvent) (chan DashboardEvent, error) {
	go d.render()
	return d.BaseDashboard.Launch(group, globalEvents)
}
