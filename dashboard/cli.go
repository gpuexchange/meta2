package dashboard

import "sync"

type CLIDashboard struct {
	BaseDashboard
}

func (d *CLIDashboard) Launch(group *sync.WaitGroup) (chan DashboardMessage, error) {
	d.messageChannel = make(chan DashboardMessage)
	go d.manageLifeCycle(group)
	return d.messageChannel, nil
}
