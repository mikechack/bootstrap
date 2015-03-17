package models

import (
	"errors"
	"log"
	"sync"
)

type DeviceInfo struct {
	ConnectorUuid   string `json:"connectorUuid,omitempty"`
	OrganizationId  string `json:"organizationId,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	DeviceType      string `json:"deviceType,omitempty"`
	DeviceId        string `json:"deviceId,omitempty"`
	Ipaddress       string `json:"ipaddress,omitempty"`
	SoftwareVersion string `json:"softwareVersion,omitempty"`
	OsVersion       string `json:"osVersion,omitempty"`
}

type devices map[string]DeviceInfo

type byreguuid map[string]string

type registrar struct {
	lock            sync.RWMutex
	deviceInventory devices
	lookup          byreguuid
}

var Registrar registrar = registrar{deviceInventory: devices{}, lookup: byreguuid{}}

func AddDevice(device DeviceInfo) error {
	return Registrar.addDevice(device)
}

func GetDeviceBySerialNumber(serialNumber string) (dev DeviceInfo, err error) {
	return Registrar.getDeviceBySerialNumber(serialNumber)
}

func GetDeviceByReguuid(reguuid string) (dev DeviceInfo, err error) {
	return Registrar.getDeviceByReguuid(reguuid)
}

func MapDeviceByReguuid(reguuid, serialnumber string) error {
	return Registrar.mapDeviceByReguuid(reguuid, serialnumber)
}

func (r *registrar) mapDeviceByReguuid(requuid, serialnumber string) (err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.deviceInventory[serialnumber]; !ok {
		err = errors.New("Can't find device")
		log.Printf("Error %s\n", err.Error())
	} else {
		r.lookup[requuid] = serialnumber
	}
	return err

}

func (r *registrar) getDeviceByReguuid(reguuid string) (dev DeviceInfo, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if deviceid, ok := r.lookup[reguuid]; ok {
		log.Printf("Registrar: Found device by reguuid = %s\n", reguuid)
		dev = r.deviceInventory[deviceid]
	} else {
		err = errors.New("Device not mapped to requuid")
		log.Printf("Error %s\n", err.Error())
	}

	return dev, err
}

func (r *registrar) addDevice(device DeviceInfo) error {
	var err error
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.deviceInventory[device.SerialNumber]; ok {
		err = errors.New("Adding Duplicate Device")
		log.Printf("Error %s\n", err.Error())
	} else {
		err = nil
		r.deviceInventory[device.SerialNumber] = device
		log.Printf("Registrar: Added Device %s\n", device.SerialNumber)
	}

	return err

}

func (r *registrar) getDeviceBySerialNumber(serialNumber string) (dev DeviceInfo, err error) {
	var ok bool
	r.lock.Lock()
	defer r.lock.Unlock()
	if dev, ok = r.deviceInventory[serialNumber]; ok {

	} else {
		err = errors.New("Can't find device by serial number -" + serialNumber)
		log.Printf("Error %s\n", err.Error())
	}

	return dev, err

}

func init() {

	dev := DeviceInfo{
		SerialNumber:    "12341234",
		DeviceType:      "vTS",
		DeviceId:        "12341234",
		Ipaddress:       "10.10.10.20",
		SoftwareVersion: "1.0",
		OsVersion:       "CentOS-7",
	}

	AddDevice(dev)
}
