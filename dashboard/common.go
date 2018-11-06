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
}

type Dashboard interface {
	Init(config map[string]string)
	Launch(incoming chan DashboardMessage) error
}

func (d *BaseDashboard) Init(config map[string]string) {
	d.devices = make([]DashboardDevice, 0)
	d.deviceLock = &sync.Mutex{}
}

func (d *BaseDashboard) manageLifeCycle() {
	channel := d.messageChannel
	for {
		message := <-channel
		switch message.MessageType {
		case UpdateDevice:
			d.updateDevice(message.Identifier, message.Status)
		case RemoveDevice:
			d.removeDevice(message.Identifier)
		case Terminate:
			break
		default:
			fmt.Println("Unknown message type %d", message.MessageType)
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
