package common

import (
	"fmt"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network"
	"strconv"
	"net"
)

type InterfaceService struct {
	BaseService

	ln       *LocalNode
	instance *network.Interface
	netTable *NetTable
	logger   log.Logger
}

func (is *InterfaceService) Name() string {
	return "iface"
}

func (is *InterfaceService) Init(ln *LocalNode) (err error) {
	is.logger = log.L(is.Name())
	is.ln = ln
	is.netTable = ln.NetTable()
	netSize, _ := ln.State().Secret().Net.Mask.Size()
	IPAddress := fmt.Sprintf("%s/%s", ln.State().PrivateIP(), strconv.Itoa(netSize))
	if is.instance, err = network.CreateInterface(is.ln.Config().DeviceType, IPAddress); err != nil {
		return err
	}

	if err = network.SetMTU(is.instance.Name(), 1400); err != nil {
		is.logger.Warning("unable to set mtu, %v", err)
	}
	return nil
}

func (is *InterfaceService) Run() error {
	for {
		buf := make([]byte, 1500)
		n, err := is.instance.Read(buf)
		var dst net.IP
		var src net.IP
		if err != nil {
			is.logger.Error("error on read from interface, %v", err)
			return err
		}
		packet := buf[:n]
		if is.instance.IsTAP() {
			ethertype := network.L2Ethertype(buf)
			if ethertype == network.IPv4 {
				packet := network.L2Payload(buf)
				if network.IsIPv4(packet) {
					dst = network.IPv4Destination(packet)
					src = network.IPv4Source(packet)
				}
			}
		} else {
			dst = network.IPv4Destination(packet)
			src = network.IPv4Source(packet)

		}

		is.netTable.SendPacket(src, dst, packet)
		is.logger.Debug("successfully been read %d bytes", n)
	}
	return nil
}

func (is *InterfaceService) WritePacket(packet []byte) error {
	is.logger.Debug("ready to write %d bytes", len(packet))
	if _, err := is.instance.Write(packet); err != nil {
		return err
	}
	return nil
}
