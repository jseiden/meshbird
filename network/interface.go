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
type Ethertype [2]byte
type Tagging int

const (
	NotTagged    Tagging = 0
	Tagged       Tagging = 4
	DoubleTagged Tagging = 8
)

var (
	IPv4 = Ethertype{0x08, 0x00}
)

func (i Interface) Name() string {
	return i.name
}

func (i *Interface) Write(data []byte) (n int, err error) {
	return i.file.Write(data)
}

func (i *Interface) Read(data []byte) (n int, err error) {
	return i.file.Read(data)
}

func (i *Interface) IsTAP() bool {
	return i.isTAP
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
func IPv4Source(packet []byte) net.IP {
	return net.IPv4(packet[12], packet[13], packet[14], packet[15])
}

func L2Tagging(l2Frame []byte) Tagging {
	if l2Frame[12] == 0x81 && l2Frame[13] == 0x00 {
		return Tagged
	} else if l2Frame[12] == 0x88 && l2Frame[13] == 0xa8 {
		return DoubleTagged
	}
	return NotTagged
}

func L2Ethertype(l2Frame []byte) Ethertype {
	ethertypePos := 12 + L2Tagging(l2Frame)
	return Ethertype{l2Frame[ethertypePos], l2Frame[ethertypePos+1]}
}
func L2Payload(l2Frame []byte) []byte {
	return l2Frame[12+L2Tagging(l2Frame)+2:]
}

func IsIPv4(packet []byte) bool {
	return 4 == (packet[0] >> 4)
}
