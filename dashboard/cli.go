package dashboard

type CLIDashboard struct {
	BaseDashboard
}

func (d *CLIDashboard) Launch() (chan DashboardMessage, error) {
	d.messageChannel = make(chan DashboardMessage)
	go d.manageLifeCycle()
	return d.messageChannel, nil
}
