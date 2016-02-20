// +build freebsd

package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func newTAP() (iface *Interface, err error) {
	iface, err = interfaceOpen("tap", "")
	if err != nil {
		return nil, err
	}
	iface.isTAP = true
	return iface, nil
}

func newTUN() (iface *Interface, err error) {
	iface, err = interfaceOpen("tun", "")
	if err != nil {
		return nil, err
	}
	return iface, nil
}

func interfaceOpen(ifType string) (*Interface, error) {
	var err error
	if ifType != "tun" && ifType != "tap" {
		return nil, fmt.Errorf("unknown interface type: %s", ifType)
	}
	iface := new(Interface)
	for i := 0; i < 10; i++ {
		ifPath := fmt.Sprintf("/dev/%s%d", ifType, i)
		fmt.Println(ifPath)
		iface.file, err = os.OpenFile(ifPath, os.O_RDWR, 0)
		if err != nil {
			continue
		}
		iface.name = fmt.Sprintf("%s%d", ifType, i)
		break
	}
	return iface, err
}

func AssignIpAddress(iface string, IpAddr string) error {
	ip, ipnet, _ := net.ParseCIDR(IpAddr)
	err := exec.Command("ifconfig", iface, "inet", ip.To4().String(), ip.To4().String(), "netmask", fmt.Sprintf("0x%s", ipnet.Mask.String())).Run()
	if err != nil {
		return fmt.Errorf("assign ip %s to %s err: %s", IpAddr, iface, err)
	}
	err = exec.Command("route", "add", IpAddr, ip.To4().String()).Run()
	if err != nil {
		return fmt.Errorf("Failed to add route %s", err)
	}
	return nil
}

func UpInterface(iface string) error {
	err := exec.Command("ifconfig", iface, "up").Run()
	if err != nil {
		return fmt.Errorf("up interface %s err: %s", iface, err)
	}
	return err
}

func SetMTU(iface string, mtu int) error {
	err := exec.Command("ifconfig", iface, "mtu", strconv.Itoa(mtu)).Run()
	if err != nil {
		return fmt.Errorf("Can't set MTU %s to %s err: %s", iface, strconv.Itoa(mtu), err)
	}
	return nil
}
