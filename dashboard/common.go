package dashboard

import (
	"fmt"
	"meta2/core"
	"sync"
)

type MessageType int

const (
	UpdateDevice MessageType = iota
	RemoveDevice
	Terminate
)

type DashboardEvent struct {
	MessageType MessageType
	Identifier  string
	Status      DeviceWorkloadStatus
}

type DeviceWorkloadStatus struct {
	WorkloadName        string
	WorkloadPerformance string
	WorkloadStatus      string
}

type BaseDashboard struct {
	messageChannel  chan DashboardEvent
	responseChannel chan core.GlobalEvent
	config          map[string]string
	devices         map[string]DeviceWorkloadStatus
	deviceLock      *sync.Mutex
	initialized     bool
	terminated      bool
}

type Dashboard interface {
	Init(config map[string]string)
	Launch(group *sync.WaitGroup, globalEvents chan core.GlobalEvent) (chan DashboardEvent, error)
	Terminate()
}

func (d *BaseDashboard) Init(config map[string]string) {
	d.devices = make(map[string]DeviceWorkloadStatus, 0)
	d.deviceLock = &sync.Mutex{}
	d.initialized = true
}

func (d *BaseDashboard) Launch(group *sync.WaitGroup, globalEvents chan core.GlobalEvent) (chan DashboardEvent, error) {
	if (!d.initialized) {
		d.Init(nil)
	}

	d.messageChannel = make(chan DashboardEvent)
	d.responseChannel = globalEvents
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
	delete(d.devices, identifier)
	d.deviceLock.Unlock()
}

func (d *BaseDashboard) updateDevice(identifier string, status DeviceWorkloadStatus) {
	d.deviceLock.Lock()
	d.devices[identifier] = status
	d.deviceLock.Unlock()
}

func (d *BaseDashboard) Terminate() {
	d.terminated = true;
	d.messageChannel <- DashboardEvent{MessageType: Terminate}
}

func (d *BaseDashboard) render() {
	println("Base renderer. To be replaced.")
}
