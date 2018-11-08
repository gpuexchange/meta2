package dashboard

import (
	ui "github.com/gizak/termui"
	"meta2/core"
	"sync"
)

type CLIDashboard struct {
	BaseDashboard
}

func (d *CLIDashboard) uiRenderer() func() {
	table := ui.NewTable()

	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(6, 0, table)),
	)

	return func() {
		d.deviceLock.Lock()
		deviceCount := len(d.devices)
		rows := make([][]string, deviceCount+1)

		// Header row
		rows[0] = []string{"Device", "Workload", "Performance", "Status"}

		deviceIds := d.getDeviceIds()

		for index, id := range deviceIds {
			workload := d.devices[id]
			rows[index+1] = []string{
				id,
				workload.Name,
				workload.Performance,
				workload.Status,
			}
		}

		d.deviceLock.Unlock()

		table.Rows = rows
		table.Height = deviceCount*2 + 3
		table.Analysis()
		table.SetSize()
		table.FgColors[0] = ui.ColorGreen

		ui.Body.Align()
		ui.Render(ui.Body)
	}
}

func (d *CLIDashboard) render() {
	ui.Init()
	defer ui.Close()

	render := d.uiRenderer()

	ui.Handle("/timer/1s", func(e ui.Event) {
		if (d.terminated) {
			ui.StopLoop()
		} else {
			render()
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
