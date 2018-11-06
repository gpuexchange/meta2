package dashboard

import (
	"fmt"
	"sync"
)

type MessageType int

const (
	UpdateDevice MessageType = iota
	RemoveDevice
	Terminate
)

type DashboardMessage struct {
	MessageType MessageType
	Identifier  string
	Status      DeviceWorkloadStatus
}

type DeviceWorkloadStatus struct {
	WorkloadName        string
	WorkloadPerformance string
	WorkloadStatus      string
}

type DashboardDevice struct {
	identifier int
	status     DeviceWorkloadStatus
}

type BaseDashboard struct {
	messageChannel chan DashboardMessage
	config         map[string]string
	devices        []DashboardDevice
	deviceLock     *sync.Mutex
	terminated     bool
}

type Dashboard interface {
	Init(config map[string]string)
	Launch(group *sync.WaitGroup) (chan DashboardMessage, error)
	Terminate()
}

func (d *BaseDashboard) Init(config map[string]string) {
	d.devices = make([]DashboardDevice, 0)
	d.deviceLock = &sync.Mutex{}
}

func (d *BaseDashboard) Launch(group *sync.WaitGroup) (chan DashboardMessage, error) {
	d.messageChannel = make(chan DashboardMessage)
	go d.manageLifeCycle(group)
	return d.messageChannel, nil
}

func (d *BaseDashboard) manageLifeCycle(group *sync.WaitGroup) {

	if group != nil {
		group.Add(1)
	}

	channel := d.messageChannel
	for {
		message := <-channel

		if (d.terminated && message.MessageType != Terminate) {
			continue;
		}

		switch message.MessageType {
		case UpdateDevice:
			d.updateDevice(message.Identifier, message.Status)
		case RemoveDevice:
			d.removeDevice(message.Identifier)
		case Terminate:
			if group != nil {
				group.Done()
			}
			break
		default:
			fmt.Printf("Unknown message type %s", string(message.MessageType))
		}
	}
}

func (d *BaseDashboard) removeDevice(identifier string) {
	d.deviceLock.Lock()
	fmt.Printf("Removing device %s\n", identifier)
	d.deviceLock.Unlock()
}

func (d *BaseDashboard) updateDevice(identifier string, status DeviceWorkloadStatus) {
	d.deviceLock.Lock()
	fmt.Printf("Updating device %s with Status %s\n", identifier, status)
	d.deviceLock.Unlock()
}

func (d *BaseDashboard) Terminate() {
	d.terminated = true;
	d.messageChannel <- DashboardMessage{MessageType: Terminate}
}
