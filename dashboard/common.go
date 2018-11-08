package dashboard

import (
	"fmt"
	"meta2/core"
	"sort"
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
	Id          string
	Status      DeviceWorkload
}

type DeviceWorkload struct {
	Name        string
	Performance string
	Status      string
}

type BaseDashboard struct {
	messageChannel  chan DashboardEvent
	responseChannel chan core.GlobalEvent
	config          map[string]string
	devices         map[string]DeviceWorkload
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
	d.devices = make(map[string]DeviceWorkload, 0)
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
			d.updateDevice(message.Id, message.Status)
		case RemoveDevice:
			d.removeDevice(message.Id)
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

func (d *BaseDashboard) getDeviceIds() []string {
	deviceIds := make([]string, len(d.devices))
	index := 0
	for deviceId, _ := range d.devices {
		deviceIds[index] = deviceId
		index++
	}

	sort.Strings(deviceIds)

	return deviceIds
}

func (d *BaseDashboard) removeDevice(id string) {
	d.deviceLock.Lock()
	delete(d.devices, id)
	d.deviceLock.Unlock()
}

func (d *BaseDashboard) updateDevice(id string, status DeviceWorkload) {
	d.deviceLock.Lock()
	d.devices[id] = status
	d.deviceLock.Unlock()
}

func (d *BaseDashboard) Terminate() {
	d.terminated = true;
	d.messageChannel <- DashboardEvent{MessageType: Terminate}
}

func (d *BaseDashboard) render() {
	println("Base renderer. To be replaced.")
}
