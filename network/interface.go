package network

import (
	"fmt"
	"net"
	"os"
)

const DEFAULT_MTU = 1400


type Interface struct {
	isTAP bool
	name  string
	file  *os.File
}

func (i Interface) Name() string {
	return i.name
}

func (i *Interface) Write(data []byte) (n int, err error) {
	return i.file.Write(data)
}

func (i *Interface) Read(data []byte) (n int, err error) {
	return i.file.Read(data)
}

func CreateInterface(deviceType string, IPAddr string) (*Interface, error) {
	fmt.Println(deviceType)
	if deviceType != "tun" && deviceType != "tap" {
		return nil, fmt.Errorf("Unknown interface type: %s\n", deviceType)
	}
	iface := new(Interface)
	var err error

	if deviceType == "tun" {
		iface, err = newTUN()
		if err != nil {
			return nil, fmt.Errorf("Create new TUN interface %v err: %s", iface, err)
		}
	}

	if deviceType == "tap" {
		iface, err = newTAP()
		if err != nil {
			return nil, fmt.Errorf("Create new TAP interface %v err: %s", iface, err)
		}
	}

	err = UpInterface(iface.Name())

	if err != nil {
		return nil, fmt.Errorf("%s interface error: %s \n", deviceType, err)
	}
	err = AssignIpAddress(iface.Name(), IPAddr)
	return iface, err
}

func IPv4Destination(packet []byte) net.IP {
	return net.IPv4(packet[16], packet[17], packet[18], packet[19])
}
